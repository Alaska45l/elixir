package admin

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"elixir/backend/internal/audit"
	"elixir/backend/internal/httpx"
	"elixir/backend/internal/media"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	Pool     *pgxpool.Pool
	Auth     AuthService
	Sessions SessionManager
	Media    *media.StorageService
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	ok, err := h.Auth.Login(r.Context(), r.RemoteAddr, req.Username, req.Password)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo iniciar sesión")
		return
	}
	if !ok {
		httpx.Error(w, r, http.StatusUnauthorized, "credenciales inválidas")
		return
	}
	h.Sessions.SetCookie(w, req.Username)
	httpx.WriteJSON(w, http.StatusOK, meResponse{Username: req.Username})
}

func (h Handler) Logout(w http.ResponseWriter, r *http.Request) {
	h.Sessions.ClearCookie(w)
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h Handler) Me(w http.ResponseWriter, r *http.Request) {
	username, err := h.Sessions.Username(r)
	if err != nil {
		httpx.Error(w, r, http.StatusUnauthorized, "sesión requerida")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, meResponse{Username: username})
}

func (h Handler) Metrics(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"total_products": 0, "total_orders": 0, "paid_revenue_cents": 0, "pending_orders": 0, "low_stock_count": 0, "recent_orders": []any{}})
		return
	}
	var totalProducts, totalOrders, pendingOrders, lowStockCount int
	var revenue int64
	threshold := h.lowStockThreshold(r)
	_ = h.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM products WHERE active=true`).Scan(&totalProducts)
	_ = h.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM orders`).Scan(&totalOrders)
	_ = h.Pool.QueryRow(r.Context(), `SELECT COALESCE(SUM(total_ars_cents),0) FROM orders WHERE status='paid'`).Scan(&revenue)
	_ = h.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM orders WHERE status='pending'`).Scan(&pendingOrders)
	_ = h.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM product_variants v JOIN products p ON p.id=v.product_id WHERE v.stock <= $1 AND v.active=true AND p.active=true`, threshold).Scan(&lowStockCount)
	recent := h.recentOrders(r)
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"total_products": totalProducts, "total_orders": totalOrders, "paid_revenue_cents": revenue, "pending_orders": pendingOrders, "low_stock_count": lowStockCount, "recent_orders": recent})
}

func (h Handler) AdminProducts(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": []any{}})
		return
	}
	rows, err := h.Pool.Query(r.Context(), `
		WITH variant_stats AS (
			SELECT product_id,
				MIN(price_ars_cents) FILTER (WHERE active=true) AS min_price_ars_cents,
				SUM(stock) FILTER (WHERE active=true) AS total_stock,
				COUNT(*) FILTER (WHERE active=true) AS variant_count
			FROM product_variants
			GROUP BY product_id
		)
		SELECT p.id, p.slug, p.name, COALESCE(p.tagline,''), p.featured, p.active, p.display_order,
			COALESCE(vs.min_price_ars_cents, 0),
			COALESCE(vs.total_stock, 0),
			COALESCE(vs.variant_count, 0)
		FROM products p
		LEFT JOIN variant_stats vs ON vs.product_id=p.id
		ORDER BY p.display_order ASC, p.created_at DESC
		LIMIT 200`)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudieron listar productos")
		return
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var id, slug, name, tagline string
		var featured, active bool
		var order, totalStock, variantCount int
		var minPrice int64
		if err := rows.Scan(&id, &slug, &name, &tagline, &featured, &active, &order, &minPrice, &totalStock, &variantCount); err != nil {
			httpx.Error(w, r, http.StatusInternalServerError, "no se pudieron leer productos")
			return
		}
		items = append(items, map[string]any{"id": id, "slug": slug, "name": name, "tagline": tagline, "featured": featured, "active": active, "display_order": order, "min_price_ars_cents": minPrice, "total_stock": totalStock, "variant_count": variantCount})
	}
	if err := rows.Err(); err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudieron leer productos")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h Handler) AdminProductByID(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	id := chi.URLParam(r, "id")
	var item productPayload
	var active bool
	err := h.Pool.QueryRow(r.Context(), `SELECT slug, name, COALESCE(tagline,''), COALESCE(description,''), COALESCE(scent_family,''), COALESCE(gender_tag,''), COALESCE(concentration,''), COALESCE(top_notes,'{}'), COALESCE(heart_notes,'{}'), COALESCE(base_notes,'{}'), featured, display_order, active FROM products WHERE id=$1`, id).
		Scan(&item.Slug, &item.Name, &item.Tagline, &item.Description, &item.ScentFamily, &item.GenderTag, &item.Concentration, &item.TopNotes, &item.HeartNotes, &item.BaseNotes, &item.Featured, &item.DisplayOrder, &active)
	if err != nil {
		httpx.Error(w, r, http.StatusNotFound, "producto no encontrado")
		return
	}
	item.Active = &active
	rows, err := h.Pool.Query(r.Context(), `SELECT size_ml, price_ars_cents, stock, COALESCE(sku,''), weight_grams FROM product_variants WHERE product_id=$1 AND active=true ORDER BY size_ml`, id)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudieron consultar variantes")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var v variantForm
		if err := rows.Scan(&v.SizeML, &v.PriceARSCents, &v.Stock, &v.SKU, &v.WeightGrams); err != nil {
			httpx.Error(w, r, http.StatusInternalServerError, "no se pudieron leer variantes")
			return
		}
		item.Variants = append(item.Variants, v)
	}
	if err := rows.Err(); err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudieron leer variantes")
		return
	}
	imgRows, err := h.Pool.Query(r.Context(), `SELECT url, COALESCE(alt_text,''), is_primary, sort_order FROM product_images WHERE product_id=$1 ORDER BY sort_order`, id)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudieron consultar imágenes")
		return
	}
	defer imgRows.Close()
	for imgRows.Next() {
		var img imageForm
		if err := imgRows.Scan(&img.URL, &img.AltText, &img.IsPrimary, &img.SortOrder); err != nil {
			httpx.Error(w, r, http.StatusInternalServerError, "no se pudieron leer imágenes")
			return
		}
		item.Images = append(item.Images, img)
	}
	if err := imgRows.Err(); err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudieron leer imágenes")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"id": id, "active": active, "product": item})
}

func (h Handler) SaveProduct(w http.ResponseWriter, r *http.Request) {
	var req productPayload
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	if h.Pool == nil {
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	if err := normalizeProductPayload(&req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}
	id := chi.URLParam(r, "id")
	action := "product.update"
	tx, err := h.Pool.Begin(r.Context())
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo iniciar la transacción")
		return
	}
	defer tx.Rollback(r.Context())
	if id == "" {
		action = "product.create"
		err = tx.QueryRow(r.Context(), `INSERT INTO products (slug,name,tagline,description,scent_family,gender_tag,concentration,top_notes,heart_notes,base_notes,featured,active,display_order) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id`,
			req.Slug, req.Name, req.Tagline, req.Description, req.ScentFamily, req.GenderTag, req.Concentration, req.TopNotes, req.HeartNotes, req.BaseNotes, req.Featured, req.activeOrDefault(true), req.DisplayOrder).Scan(&id)
		if err != nil {
			httpx.Error(w, r, http.StatusBadRequest, "no se pudo crear el producto")
			return
		}
	} else {
		var rowsAffected int64
		if req.Active == nil {
			tag, execErr := tx.Exec(r.Context(), `UPDATE products SET slug=$1,name=$2,tagline=$3,description=$4,scent_family=$5,gender_tag=$6,concentration=$7,top_notes=$8,heart_notes=$9,base_notes=$10,featured=$11,display_order=$12,updated_at=now() WHERE id=$13`,
				req.Slug, req.Name, req.Tagline, req.Description, req.ScentFamily, req.GenderTag, req.Concentration, req.TopNotes, req.HeartNotes, req.BaseNotes, req.Featured, req.DisplayOrder, id)
			err = execErr
			rowsAffected = tag.RowsAffected()
		} else {
			tag, execErr := tx.Exec(r.Context(), `UPDATE products SET slug=$1,name=$2,tagline=$3,description=$4,scent_family=$5,gender_tag=$6,concentration=$7,top_notes=$8,heart_notes=$9,base_notes=$10,featured=$11,active=$12,display_order=$13,updated_at=now() WHERE id=$14`,
				req.Slug, req.Name, req.Tagline, req.Description, req.ScentFamily, req.GenderTag, req.Concentration, req.TopNotes, req.HeartNotes, req.BaseNotes, req.Featured, *req.Active, req.DisplayOrder, id)
			err = execErr
			rowsAffected = tag.RowsAffected()
		}
		if err != nil || rowsAffected == 0 {
			httpx.Error(w, r, http.StatusBadRequest, "no se pudo actualizar el producto")
			return
		}
	}
	if err := h.syncVariants(r, tx, id, req.Slug, req.Variants); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudieron guardar las variantes")
		return
	}
	if err := h.syncImages(r, tx, id, req.Images); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudieron guardar las imágenes")
		return
	}
	if err := tx.Commit(r.Context()); err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo guardar")
		return
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: action, ResourceID: id, Metadata: map[string]any{"slug": req.Slug}})
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"id": id})
}

func (h Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	tag, err := h.Pool.Exec(r.Context(), `UPDATE products SET active=false, updated_at=now() WHERE id=$1`, chi.URLParam(r, "id"))
	if err != nil || tag.RowsAffected() == 0 {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo eliminar")
		return
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "product.delete", ResourceID: chi.URLParam(r, "id")})
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h Handler) UpdateProductActive(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	var req struct {
		Active bool `json:"active"`
	}
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	tag, err := h.Pool.Exec(r.Context(), `UPDATE products SET active=$1, updated_at=now() WHERE id=$2`, req.Active, chi.URLParam(r, "id"))
	if err != nil || tag.RowsAffected() == 0 {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo actualizar")
		return
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "product.active", ResourceID: chi.URLParam(r, "id"), Metadata: map[string]any{"active": req.Active}})
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h Handler) syncVariants(r *http.Request, tx pgx.Tx, productID, slug string, variants []variantForm) error {
	seen := make([]string, 0, len(variants))
	for _, v := range variants {
		sku := strings.TrimSpace(v.SKU)
		if sku == "" {
			sku = fmt.Sprintf("%s-%d", slug, v.SizeML)
		}
		weight := v.WeightGrams
		if weight <= 0 {
			weight = 200
		}
		tag, err := tx.Exec(r.Context(), `
	INSERT INTO product_variants (product_id,size_ml,price_ars_cents,stock,sku,active,weight_grams)
	VALUES ($1,$2,$3,$4,$5,true,$6)
	ON CONFLICT (sku) DO UPDATE SET
		size_ml=EXCLUDED.size_ml,
		price_ars_cents=EXCLUDED.price_ars_cents,
		stock=EXCLUDED.stock,
		active=true,
		weight_grams=EXCLUDED.weight_grams
	WHERE product_variants.product_id=EXCLUDED.product_id`,
			productID, v.SizeML, v.PriceARSCents, v.Stock, sku, weight)
		if err != nil {
			return err
		}
		if tag.RowsAffected() != 1 {
			return fmt.Errorf("sku %s pertenece a otro producto", sku)
		}
		seen = append(seen, sku)
	}
	if len(seen) == 0 {
		_, err := tx.Exec(r.Context(), `UPDATE product_variants SET active=false WHERE product_id=$1`, productID)
		return err
	}
	_, err := tx.Exec(r.Context(), `UPDATE product_variants SET active=false WHERE product_id=$1 AND COALESCE(sku,'') <> ALL($2)`, productID, seen)
	return err
}

func (h Handler) syncImages(r *http.Request, tx pgx.Tx, productID string, images []imageForm) error {
	seen := make([]string, 0, len(images))
	for _, img := range images {
		url := strings.TrimSpace(img.URL)
		if url == "" {
			continue
		}
		tag, err := tx.Exec(r.Context(), `UPDATE product_images SET alt_text=$3, is_primary=$4, sort_order=$5 WHERE product_id=$1 AND url=$2`, productID, url, img.AltText, img.IsPrimary, img.SortOrder)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			if _, err := tx.Exec(r.Context(), `INSERT INTO product_images (product_id,url,alt_text,is_primary,sort_order) VALUES ($1,$2,$3,$4,$5)`, productID, url, img.AltText, img.IsPrimary, img.SortOrder); err != nil {
				return err
			}
		}
		seen = append(seen, url)
	}
	if len(seen) == 0 {
		_, err := tx.Exec(r.Context(), `DELETE FROM product_images WHERE product_id=$1`, productID)
		return err
	}
	_, err := tx.Exec(r.Context(), `DELETE FROM product_images WHERE product_id=$1 AND url <> ALL($2)`, productID, seen)
	return err
}

func (h Handler) ImportProducts(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "archivo inválido")
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "archivo requerido")
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil || len(records) < 2 {
		httpx.Error(w, r, http.StatusBadRequest, "csv inválido")
		return
	}
	imported := 0
	errs := []string{}
	for i, row := range records[1:] {
		rowNum := i + 2
		if len(row) < 10 {
			errs = append(errs, fmt.Sprintf("row %d: columnas insuficientes", rowNum))
			continue
		}
		size, err := strconv.Atoi(row[6])
		if err != nil {
			errs = append(errs, fmt.Sprintf("row %d: size_ml inválido", rowNum))
			continue
		}
		priceARS, err := strconv.ParseInt(row[7], 10, 64)
		if err != nil {
			errs = append(errs, fmt.Sprintf("row %d: price_ars inválido", rowNum))
			continue
		}
		stock, err := strconv.Atoi(row[8])
		if err != nil {
			errs = append(errs, fmt.Sprintf("row %d: stock inválido", rowNum))
			continue
		}
		tx, err := h.Pool.Begin(r.Context())
		if err != nil {
			errs = append(errs, fmt.Sprintf("row %d: no se pudo iniciar transacción", rowNum))
			continue
		}
		var productID string
		err = tx.QueryRow(r.Context(), `INSERT INTO products (name, slug, tagline, scent_family, concentration, gender_tag, top_notes, heart_notes, base_notes) VALUES ($1,$2,$3,$4,$5,$6,'{}','{}','{}') ON CONFLICT (slug) DO UPDATE SET name=EXCLUDED.name, tagline=EXCLUDED.tagline, scent_family=EXCLUDED.scent_family, concentration=EXCLUDED.concentration, gender_tag=EXCLUDED.gender_tag, updated_at=now() RETURNING id`,
			row[0], row[1], row[2], row[3], row[4], row[5]).Scan(&productID)
		if err == nil {
			_, err = tx.Exec(r.Context(), `INSERT INTO product_variants (product_id,size_ml,price_ars_cents,stock,sku) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (sku) DO UPDATE SET price_ars_cents=EXCLUDED.price_ars_cents, stock=EXCLUDED.stock`, productID, size, priceARS*100, stock, row[1]+"-"+row[6])
		}
		if err == nil && strings.TrimSpace(row[9]) != "" {
			_, err = tx.Exec(r.Context(), `INSERT INTO product_images (product_id,url,alt_text,is_primary,sort_order) VALUES ($1,$2,$3,true,0)`, productID, row[9], row[0])
		}
		if err != nil {
			_ = tx.Rollback(r.Context())
			errs = append(errs, fmt.Sprintf("row %d: %v", rowNum, err))
			continue
		}
		if err := tx.Commit(r.Context()); err != nil {
			errs = append(errs, fmt.Sprintf("row %d: %v", rowNum, err))
			continue
		}
		imported++
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "product.import", Metadata: map[string]any{"imported": imported, "errors": len(errs)}})
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"imported": imported, "errors": errs})
}

func (h Handler) Orders(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": h.orders(r, 100)})
}

func (h Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	var req orderUpdateRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	req.Status = strings.TrimSpace(req.Status)
	req.TrackingNumber = strings.TrimSpace(req.TrackingNumber)
	req.ShippingCarrier = strings.TrimSpace(req.ShippingCarrier)
	req.InternalNotes = strings.TrimSpace(req.InternalNotes)
	if !validateOrderStatus(req.Status) {
		httpx.Error(w, r, http.StatusBadRequest, "estado de orden inválido")
		return
	}
	if len(req.TrackingNumber) > 120 || len(req.ShippingCarrier) > 80 || len(req.InternalNotes) > 1000 {
		httpx.Error(w, r, http.StatusBadRequest, "revisá tracking, correo o notas internas")
		return
	}
	if req.Status == "shipped" && req.ShippedAt == nil {
		now := time.Now()
		req.ShippedAt = &now
	}
	tag, err := h.Pool.Exec(r.Context(), `
		UPDATE orders
		SET status=$1, tracking_number=$2, shipping_carrier=$3, shipped_at=$4, internal_notes=$5, updated_at=now()
		WHERE id=$6`,
		req.Status, req.TrackingNumber, req.ShippingCarrier, req.ShippedAt, req.InternalNotes, chi.URLParam(r, "id"))
	if err != nil || tag.RowsAffected() == 0 {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo actualizar")
		return
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "order.update", ResourceID: chi.URLParam(r, "id"), Metadata: map[string]any{"status": req.Status, "tracking": req.TrackingNumber, "carrier": req.ShippingCarrier}})
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h Handler) ExportOrders(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="orders-`+time.Now().Format("2006-01-02")+`.csv"`)
	writer := csv.NewWriter(w)
	defer writer.Flush()
	_ = writer.Write([]string{"external_reference", "status", "customer_name", "customer_email", "customer_phone", "subtotal_ars_cents", "shipping_cost_ars_cents", "discount_ars_cents", "total_ars_cents", "payment_status", "tracking_number", "shipping_carrier", "created_at"})
	rows, err := h.Pool.Query(r.Context(), `
		SELECT o.external_reference, o.status, o.customer_name, o.customer_email, COALESCE(o.customer_phone,''),
			o.subtotal_ars_cents, o.shipping_cost_ars_cents, o.discount_ars_cents, o.total_ars_cents,
			COALESCE(pe.mp_status,''), COALESCE(o.tracking_number,''), COALESCE(o.shipping_carrier,''), o.created_at
		FROM orders o
		LEFT JOIN LATERAL (
			SELECT mp_status FROM payment_events pe WHERE pe.order_id=o.id ORDER BY processed_at DESC LIMIT 1
		) pe ON true
		ORDER BY o.created_at DESC`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var ref, status, name, email, phone, paymentStatus, tracking, carrier string
		var subtotal, shipping, discount, total int64
		var created time.Time
		if rows.Scan(&ref, &status, &name, &email, &phone, &subtotal, &shipping, &discount, &total, &paymentStatus, &tracking, &carrier, &created) == nil {
			_ = writer.Write([]string{ref, status, name, email, phone, strconv.FormatInt(subtotal, 10), strconv.FormatInt(shipping, 10), strconv.FormatInt(discount, 10), strconv.FormatInt(total, 10), paymentStatus, tracking, carrier, created.Format(time.RFC3339)})
		}
	}
}

func (h Handler) Discounts(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		if r.Method == http.MethodGet {
			httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": []any{}})
			return
		}
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	if r.Method == http.MethodGet {
		h.listTable(w, r, `SELECT id, code, discount_type, discount_value, min_order_cents, max_uses, uses, active, expires_at FROM discount_codes ORDER BY created_at DESC`)
		return
	}
	var req discountWriteRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	if err := normalizeDiscountPayload(&req, true); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}
	active := true
	if req.Active != nil {
		active = *req.Active
	}
	_, err := h.Pool.Exec(r.Context(), `INSERT INTO discount_codes (code,discount_type,discount_value,min_order_cents,max_uses,expires_at,active) VALUES ($1,$2,$3,$4,$5,$6,$7)`, req.Code, req.DiscountType, req.DiscountValue, req.MinOrderCents, req.MaxUses, req.ExpiresAt, active)
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo crear")
		return
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "discount.create", ResourceID: req.Code})
	httpx.WriteJSON(w, http.StatusCreated, map[string]bool{"ok": true})
}

func (h Handler) DiscountByID(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	if r.Method == http.MethodDelete {
		tag, err := h.Pool.Exec(r.Context(), `DELETE FROM discount_codes WHERE id=$1`, chi.URLParam(r, "id"))
		if err != nil || tag.RowsAffected() == 0 {
			httpx.Error(w, r, http.StatusBadRequest, "no se pudo eliminar")
			return
		}
		audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "discount.delete", ResourceID: chi.URLParam(r, "id")})
		httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
		return
	}
	var req discountWriteRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	if err := normalizeDiscountPayload(&req, false); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}
	active := true
	if req.Active != nil {
		active = *req.Active
	}
	tag, err := h.Pool.Exec(r.Context(), `UPDATE discount_codes SET code=$1, discount_type=$2, discount_value=$3, min_order_cents=$4, max_uses=$5, expires_at=$6, active=$7 WHERE id=$8`, req.Code, req.DiscountType, req.DiscountValue, req.MinOrderCents, req.MaxUses, req.ExpiresAt, active, chi.URLParam(r, "id"))
	if err != nil || tag.RowsAffected() == 0 {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo actualizar")
		return
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "discount.update", ResourceID: chi.URLParam(r, "id"), Metadata: map[string]any{"code": req.Code, "active": active}})
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h Handler) Homepage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.listTable(w, r, `SELECT hero_heading, hero_subheading, hero_image_url, COALESCE(hero_image_mode, 'product_covers') AS hero_image_mode, COALESCE(hero_rotation_interval_ms, 8000) AS hero_rotation_interval_ms, hero_cta_label, hero_cta_url, editorial_heading, editorial_body, editorial_image_url FROM homepage_settings WHERE id=1`)
		return
	}
	var req homepageRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	if h.Pool == nil {
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	if err := normalizeHomepagePayload(&req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}
	_, err := h.Pool.Exec(r.Context(), `INSERT INTO homepage_settings (id, hero_heading, hero_subheading, hero_image_url, hero_image_mode, hero_rotation_interval_ms, hero_cta_label, hero_cta_url, editorial_heading, editorial_body, editorial_image_url) VALUES (1,$1,$2,$3,$4,$5,$6,$7,$8,$9,$10) ON CONFLICT (id) DO UPDATE SET hero_heading=$1, hero_subheading=$2, hero_image_url=$3, hero_image_mode=$4, hero_rotation_interval_ms=$5, hero_cta_label=$6, hero_cta_url=$7, editorial_heading=$8, editorial_body=$9, editorial_image_url=$10, updated_at=now()`,
		req.HeroHeading, req.HeroSubheading, req.HeroImageURL, req.HeroImageMode, req.HeroRotationIntervalMS, req.HeroCTALabel, req.HeroCTAURL, req.EditorialHeading, req.EditorialBody, req.EditorialImageURL)
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo guardar")
		return
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "homepage.update"})
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h Handler) PublicHomepage(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.WriteJSON(w, http.StatusOK, homepageRequest{})
		return
	}
	var item homepageRequest
	err := h.Pool.QueryRow(r.Context(), `SELECT COALESCE(hero_heading,''), COALESCE(hero_subheading,''), COALESCE(hero_image_url,''), COALESCE(hero_image_mode,'product_covers'), COALESCE(hero_rotation_interval_ms,8000), COALESCE(hero_cta_label,''), COALESCE(hero_cta_url,''), COALESCE(editorial_heading,''), COALESCE(editorial_body,''), COALESCE(editorial_image_url,'') FROM homepage_settings WHERE id=1`).
		Scan(&item.HeroHeading, &item.HeroSubheading, &item.HeroImageURL, &item.HeroImageMode, &item.HeroRotationIntervalMS, &item.HeroCTALabel, &item.HeroCTAURL, &item.EditorialHeading, &item.EditorialBody, &item.EditorialImageURL)
	if err != nil {
		httpx.WriteJSON(w, http.StatusOK, homepageRequest{})
		return
	}
	httpx.WriteJSON(w, http.StatusOK, item)
}

func (h Handler) Settings(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		settings, err := h.loadSiteSettings(r)
		if err != nil {
			httpx.Error(w, r, http.StatusInternalServerError, "no se pudo consultar la configuración")
			return
		}
		httpx.WriteJSON(w, http.StatusOK, settings)
		return
	}
	var req siteSettings
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	if h.Pool == nil {
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	if err := normalizeSiteSettingsPayload(&req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}
	faqJSON, _ := json.Marshal(req.FAQItems)
	navJSON, _ := json.Marshal(req.NavbarProductCategories)
	_, err := h.Pool.Exec(r.Context(), `
	INSERT INTO site_settings (
		id, footer_instagram_url, footer_tiktok_url, footer_whatsapp_url,
		announcement_bar_text, announcement_bar_active,
		about_title, about_description, about_location, about_phone,
		faq_items, return_policy_html, navbar_product_categories, low_stock_threshold, updated_at
	) VALUES (1,$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,now())
	ON CONFLICT (id) DO UPDATE SET
		footer_instagram_url=$1,
		footer_tiktok_url=$2,
		footer_whatsapp_url=$3,
		announcement_bar_text=$4,
		announcement_bar_active=$5,
		about_title=$6,
		about_description=$7,
		about_location=$8,
		about_phone=$9,
		faq_items=$10,
		return_policy_html=$11,
		navbar_product_categories=$12,
		low_stock_threshold=$13,
		updated_at=now()`,
		req.FooterInstagramURL, req.FooterTikTokURL, req.FooterWhatsAppURL,
		req.AnnouncementBarText, req.AnnouncementBarActive,
		req.AboutTitle, req.AboutDescription, req.AboutLocation, req.AboutPhone,
		faqJSON, req.ReturnPolicyHTML, navJSON, req.LowStockThreshold)
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo guardar")
		return
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "settings.update"})
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h Handler) PublicSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.loadSiteSettings(r)
	if err != nil {
		settings = defaultSiteSettings()
	}
	w.Header().Set("Cache-Control", "public, max-age=60, stale-while-revalidate=300")
	httpx.WriteJSON(w, http.StatusOK, settings)
}

func (h Handler) loadSiteSettings(r *http.Request) (siteSettings, error) {
	if h.Pool == nil {
		return defaultSiteSettings(), nil
	}
	var s siteSettings
	var faqRaw, navRaw []byte
	err := h.Pool.QueryRow(r.Context(), `
	SELECT
		COALESCE(footer_instagram_url,''),
		COALESCE(footer_tiktok_url,''),
		COALESCE(footer_whatsapp_url,''),
		COALESCE(announcement_bar_text,''),
		COALESCE(announcement_bar_active,true),
		COALESCE(about_title,''),
		COALESCE(about_description,''),
		COALESCE(about_location,''),
		COALESCE(about_phone,''),
		COALESCE(faq_items,'[]'::jsonb),
		COALESCE(return_policy_html,''),
		COALESCE(navbar_product_categories,'[]'::jsonb),
		COALESCE(low_stock_threshold, 5)
	FROM site_settings WHERE id=1`).
		Scan(&s.FooterInstagramURL, &s.FooterTikTokURL, &s.FooterWhatsAppURL, &s.AnnouncementBarText, &s.AnnouncementBarActive, &s.AboutTitle, &s.AboutDescription, &s.AboutLocation, &s.AboutPhone, &faqRaw, &s.ReturnPolicyHTML, &navRaw, &s.LowStockThreshold)
	if err != nil {
		return s, err
	}
	defaults := defaultSiteSettings()
	if strings.TrimSpace(s.FooterInstagramURL) == "" {
		s.FooterInstagramURL = defaults.FooterInstagramURL
	}
	if strings.TrimSpace(s.FooterTikTokURL) == "" {
		s.FooterTikTokURL = defaults.FooterTikTokURL
	}
	if err := json.Unmarshal(faqRaw, &s.FAQItems); err != nil {
		s.FAQItems = defaults.FAQItems
	}
	if err := json.Unmarshal(navRaw, &s.NavbarProductCategories); err != nil {
		s.NavbarProductCategories = defaults.NavbarProductCategories
	}
	if s.LowStockThreshold <= 0 {
		s.LowStockThreshold = defaults.LowStockThreshold
	}
	return s, nil
}

func (h Handler) Contact(w http.ResponseWriter, r *http.Request) {
	h.listTable(w, r, `SELECT id, name, email, COALESCE(subject,''), message, read, created_at FROM contact_messages ORDER BY read ASC, created_at DESC LIMIT 100`)
}

func (h Handler) MarkContactRead(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	read := true
	if r.ContentLength != 0 {
		var req struct {
			Read *bool `json:"read"`
		}
		if err := httpx.DecodeStrict(r, &req); err != nil {
			httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
			return
		}
		if req.Read != nil {
			read = *req.Read
		}
	}
	tag, err := h.Pool.Exec(r.Context(), `UPDATE contact_messages SET read=$1 WHERE id=$2`, read, chi.URLParam(r, "id"))
	if err != nil || tag.RowsAffected() == 0 {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo actualizar el mensaje")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h Handler) LowStock(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": []any{}, "threshold": 5})
		return
	}
	threshold := h.lowStockThreshold(r)
	rows, err := h.Pool.Query(r.Context(), `SELECT v.id, p.id AS product_id, p.name, v.size_ml, v.stock FROM product_variants v JOIN products p ON p.id=v.product_id WHERE v.stock <= $1 AND v.active=true AND p.active=true ORDER BY v.stock ASC, p.name ASC LIMIT 200`, threshold)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo consultar")
		return
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var id, productID, name string
		var size, stock int
		if err := rows.Scan(&id, &productID, &name, &size, &stock); err != nil {
			httpx.Error(w, r, http.StatusInternalServerError, "no se pudo leer stock")
			return
		}
		items = append(items, map[string]any{"id": id, "product_id": productID, "name": name, "size_ml": size, "stock": stock})
	}
	if err := rows.Err(); err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo leer stock")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items, "threshold": threshold})
}

func (h Handler) lowStockThreshold(r *http.Request) int {
	if h.Pool == nil {
		return 5
	}
	var threshold int
	if err := h.Pool.QueryRow(r.Context(), `SELECT COALESCE(low_stock_threshold, 5) FROM site_settings WHERE id=1`).Scan(&threshold); err != nil || threshold <= 0 {
		return 5
	}
	return threshold
}

func (h Handler) listTable(w http.ResponseWriter, r *http.Request, query string) {
	if h.Pool == nil {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": []any{}})
		return
	}
	rows, err := h.Pool.Query(r.Context(), query)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo consultar")
		return
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		values, _ := rows.Values()
		fields := rows.FieldDescriptions()
		item := map[string]any{}
		for i, f := range fields {
			item[string(f.Name)] = values[i]
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo leer la consulta")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h Handler) recentOrders(r *http.Request) []map[string]any {
	items := h.orders(r, 10)
	return items
}

func (h Handler) orders(r *http.Request, limit int) []map[string]any {
	if h.Pool == nil {
		return []map[string]any{}
	}
	rows, err := h.Pool.Query(r.Context(), `
		SELECT o.id, o.external_reference, o.status, o.customer_name, o.customer_email, COALESCE(o.customer_phone,''),
			COALESCE(o.shipping_address,'{}'::jsonb), o.shipping_cost_ars_cents, o.subtotal_ars_cents, o.total_ars_cents,
			COALESCE(o.discount_code,''), o.discount_ars_cents, COALESCE(o.currency,'ARS'),
			COALESCE(o.tracking_number,''), COALESCE(o.shipping_carrier,''), o.shipped_at, COALESCE(o.internal_notes,''), o.created_at,
			pe.mp_payment_id, pe.mp_preference_id, pe.mp_status, pe.mp_status_detail, pe.processed_at
		FROM orders o
		LEFT JOIN LATERAL (
			SELECT mp_payment_id, mp_preference_id, mp_status, mp_status_detail, processed_at
			FROM payment_events pe
			WHERE pe.order_id=o.id
			ORDER BY pe.processed_at DESC
			LIMIT 1
		) pe ON true
		ORDER BY o.created_at DESC
		LIMIT $1`, limit)
	if err != nil {
		return []map[string]any{}
	}
	defer rows.Close()
	items := []map[string]any{}
	orderIDs := []string{}
	for rows.Next() {
		var id, ref, status, name, email, phone, discountCode, currency, tracking, carrier, internalNotes string
		var shipping []byte
		var created time.Time
		var shippedAt sql.NullTime
		var paymentID, preferenceID, paymentStatus, paymentDetail sql.NullString
		var paymentProcessed sql.NullTime
		var shippingCost, subtotal, total, discount int64
		if rows.Scan(&id, &ref, &status, &name, &email, &phone, &shipping, &shippingCost, &subtotal, &total, &discountCode, &discount, &currency, &tracking, &carrier, &shippedAt, &internalNotes, &created, &paymentID, &preferenceID, &paymentStatus, &paymentDetail, &paymentProcessed) == nil {
			var address map[string]any
			_ = json.Unmarshal(shipping, &address)
			orderIDs = append(orderIDs, id)
			item := map[string]any{
				"id": id, "external_reference": ref, "status": status,
				"customer_name": name, "customer_email": email, "customer_phone": phone,
				"shipping_address": address, "shipping_cost_ars_cents": shippingCost,
				"subtotal_ars_cents": subtotal, "total_ars_cents": total, "discount_code": discountCode,
				"discount_ars_cents": discount, "currency": currency, "tracking_number": tracking,
				"shipping_carrier": carrier, "internal_notes": internalNotes, "created_at": created,
				"items": []map[string]any{},
				"payment": map[string]any{
					"mp_payment_id": maybeString(paymentID), "mp_preference_id": maybeString(preferenceID),
					"mp_status": maybeString(paymentStatus), "mp_status_detail": maybeString(paymentDetail),
					"processed_at": maybeTime(paymentProcessed),
				},
			}
			if shippedAt.Valid {
				item["shipped_at"] = shippedAt.Time
			}
			items = append(items, item)
		}
	}
	itemsByOrder := h.orderItemsBatch(r, orderIDs)
	for _, item := range items {
		if id, ok := item["id"].(string); ok {
			item["items"] = itemsByOrder[id]
		}
	}
	return items
}

func (h Handler) orderItemsBatch(r *http.Request, orderIDs []string) map[string][]map[string]any {
	out := make(map[string][]map[string]any, len(orderIDs))
	if len(orderIDs) == 0 {
		return out
	}
	rows, err := h.Pool.Query(r.Context(), `SELECT order_id, product_name, size_ml, quantity, unit_price_ars_cents, subtotal_ars_cents FROM order_items WHERE order_id::text = ANY($1) ORDER BY order_id`, orderIDs)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var orderID, name string
		var size, qty int
		var unit, subtotal int64
		if rows.Scan(&orderID, &name, &size, &qty, &unit, &subtotal) == nil {
			out[orderID] = append(out[orderID], map[string]any{"product_name": name, "size_ml": size, "quantity": qty, "unit_price_ars_cents": unit, "subtotal_ars_cents": subtotal})
		}
	}
	return out
}

func (h Handler) orderItems(r *http.Request, orderID string) []map[string]any {
	rows, err := h.Pool.Query(r.Context(), `SELECT product_name, size_ml, quantity, unit_price_ars_cents, subtotal_ars_cents FROM order_items WHERE order_id=$1`, orderID)
	if err != nil {
		return []map[string]any{}
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var name string
		var size, qty int
		var unit, subtotal int64
		if rows.Scan(&name, &size, &qty, &unit, &subtotal) == nil {
			items = append(items, map[string]any{"product_name": name, "size_ml": size, "quantity": qty, "unit_price_ars_cents": unit, "subtotal_ars_cents": subtotal})
		}
	}
	return items
}

type productPayload struct {
	Slug          string        `json:"slug"`
	Name          string        `json:"name"`
	Tagline       string        `json:"tagline"`
	Description   string        `json:"description"`
	ScentFamily   string        `json:"scent_family"`
	GenderTag     string        `json:"gender_tag"`
	Concentration string        `json:"concentration"`
	TopNotes      []string      `json:"top_notes"`
	HeartNotes    []string      `json:"heart_notes"`
	BaseNotes     []string      `json:"base_notes"`
	Featured      bool          `json:"featured"`
	Active        *bool         `json:"active,omitempty"`
	DisplayOrder  int           `json:"display_order"`
	Variants      []variantForm `json:"variants"`
	Images        []imageForm   `json:"images"`
}

func (p productPayload) activeOrDefault(defaultValue bool) bool {
	if p.Active == nil {
		return defaultValue
	}
	return *p.Active
}

type variantForm struct {
	SizeML        int    `json:"size_ml"`
	PriceARSCents int64  `json:"price_ars_cents"`
	Stock         int    `json:"stock"`
	SKU           string `json:"sku"`
	WeightGrams   int    `json:"weight_grams"`
}

type imageForm struct {
	URL       string `json:"url"`
	AltText   string `json:"alt_text"`
	IsPrimary bool   `json:"is_primary"`
	SortOrder int    `json:"sort_order"`
}

type orderUpdateRequest struct {
	Status          string     `json:"status"`
	TrackingNumber  string     `json:"tracking_number"`
	ShippingCarrier string     `json:"shipping_carrier"`
	ShippedAt       *time.Time `json:"shipped_at"`
	InternalNotes   string     `json:"internal_notes"`
}

type discountWriteRequest struct {
	Code          string     `json:"code"`
	DiscountType  string     `json:"discount_type"`
	DiscountValue int64      `json:"discount_value"`
	MinOrderCents int64      `json:"min_order_cents"`
	MaxUses       *int       `json:"max_uses"`
	ExpiresAt     *time.Time `json:"expires_at"`
	Active        *bool      `json:"active"`
}

type homepageRequest struct {
	HeroHeading            string `json:"hero_heading"`
	HeroSubheading         string `json:"hero_subheading"`
	HeroImageURL           string `json:"hero_image_url"`
	HeroImageMode          string `json:"hero_image_mode"`
	HeroRotationIntervalMS int    `json:"hero_rotation_interval_ms"`
	HeroCTALabel           string `json:"hero_cta_label"`
	HeroCTAURL             string `json:"hero_cta_url"`
	EditorialHeading       string `json:"editorial_heading"`
	EditorialBody          string `json:"editorial_body"`
	EditorialImageURL      string `json:"editorial_image_url"`
}

type siteSettings struct {
	FooterInstagramURL      string    `json:"footer_instagram_url"`
	FooterTikTokURL         string    `json:"footer_tiktok_url"`
	FooterWhatsAppURL       string    `json:"footer_whatsapp_url"`
	AnnouncementBarText     string    `json:"announcement_bar_text"`
	AnnouncementBarActive   bool      `json:"announcement_bar_active"`
	AboutTitle              string    `json:"about_title"`
	AboutDescription        string    `json:"about_description"`
	AboutLocation           string    `json:"about_location"`
	AboutPhone              string    `json:"about_phone"`
	FAQItems                []faqItem `json:"faq_items"`
	ReturnPolicyHTML        string    `json:"return_policy_html"`
	NavbarProductCategories []navItem `json:"navbar_product_categories"`
	LowStockThreshold       int       `json:"low_stock_threshold"`
}

type faqItem struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type navItem struct {
	Label string `json:"label"`
	Href  string `json:"href"`
}

func defaultSiteSettings() siteSettings {
	return siteSettings{
		FooterInstagramURL:    "https://www.instagram.com/",
		FooterTikTokURL:       "https://www.tiktok.com/",
		AnnouncementBarText:   "Envíos a todo el país · Empaque discreto · Seguimiento personalizado",
		AnnouncementBarActive: true,
		AboutTitle:            "ELIXIR Exclusive",
		AboutDescription:      "Perfumería argentina de lujo discreto. Fragancias intensas, envíos nacionales y atención privada.",
		AboutLocation:         "Buenos Aires, Argentina",
		FAQItems: []faqItem{
			{Question: "¿Los perfumes son originales?", Answer: "Sí. ELIXIR Exclusive comercializa fragancias seleccionadas y documentadas."},
			{Question: "¿Qué medios de pago aceptan?", Answer: "El checkout opera en ARS mediante MercadoPago."},
			{Question: "¿Hacen envíos?", Answer: "Sí, a CABA, GBA e Interior con seguimiento."},
			{Question: "¿Puedo consultar por WhatsApp?", Answer: "Sí. Recomendamos WhatsApp para asesoramiento rápido."},
		},
		ReturnPolicyHTML: "<p>Los cambios se revisan caso por caso con el producto cerrado, sin uso y dentro de los plazos informados por atención al cliente.</p>",
		NavbarProductCategories: []navItem{
			{Label: "Fragancias Masculinas", Href: "/fragrances?gender=Masculino"},
			{Label: "Fragancias Femeninas", Href: "/fragrances?gender=Femenino"},
		},
		LowStockThreshold: 5,
	}
}

func (h Handler) actor(r *http.Request) string {
	username, _ := h.Sessions.Username(r)
	if username == "" {
		return "unknown"
	}
	return username
}

func maybeString(v sql.NullString) string {
	if !v.Valid {
		return ""
	}
	return v.String
}

func maybeTime(v sql.NullTime) *time.Time {
	if !v.Valid {
		return nil
	}
	return &v.Time
}

func CleanCode(v string) string { return strings.ToUpper(strings.TrimSpace(v)) }
