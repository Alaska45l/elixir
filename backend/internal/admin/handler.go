package admin

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"elixir/backend/internal/audit"
	"elixir/backend/internal/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	Pool     *pgxpool.Pool
	Auth     AuthService
	Sessions SessionManager
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
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
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"total_products": 0, "total_orders": 0, "paid_revenue_cents": 0, "pending_orders": 0, "recent_orders": []any{}})
		return
	}
	var totalProducts, totalOrders, pendingOrders int
	var revenue int64
	_ = h.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM products WHERE active=true`).Scan(&totalProducts)
	_ = h.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM orders`).Scan(&totalOrders)
	_ = h.Pool.QueryRow(r.Context(), `SELECT COALESCE(SUM(total_ars_cents),0) FROM orders WHERE status='paid'`).Scan(&revenue)
	_ = h.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM orders WHERE status='pending'`).Scan(&pendingOrders)
	recent := h.recentOrders(r)
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"total_products": totalProducts, "total_orders": totalOrders, "paid_revenue_cents": revenue, "pending_orders": pendingOrders, "recent_orders": recent})
}

func (h Handler) AdminProducts(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": []any{}})
		return
	}
	rows, err := h.Pool.Query(r.Context(), `SELECT id, slug, name, COALESCE(tagline,''), featured, active, display_order FROM products ORDER BY display_order ASC, created_at DESC LIMIT 100`)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudieron listar productos")
		return
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var id, slug, name, tagline string
		var featured, active bool
		var order int
		if rows.Scan(&id, &slug, &name, &tagline, &featured, &active, &order) == nil {
			items = append(items, map[string]any{"id": id, "slug": slug, "name": name, "tagline": tagline, "featured": featured, "active": active, "display_order": order})
		}
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
	rows, err := h.Pool.Query(r.Context(), `SELECT size_ml, price_ars_cents, stock, COALESCE(sku,'') FROM product_variants WHERE product_id=$1 ORDER BY size_ml`, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var v variantForm
			if rows.Scan(&v.SizeML, &v.PriceARSCents, &v.Stock, &v.SKU) == nil {
				item.Variants = append(item.Variants, v)
			}
		}
	}
	imgRows, err := h.Pool.Query(r.Context(), `SELECT url, COALESCE(alt_text,''), is_primary, sort_order FROM product_images WHERE product_id=$1 ORDER BY sort_order`, id)
	if err == nil {
		defer imgRows.Close()
		for imgRows.Next() {
			var img imageForm
			if imgRows.Scan(&img.URL, &img.AltText, &img.IsPrimary, &img.SortOrder) == nil {
				item.Images = append(item.Images, img)
			}
		}
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
	id := chi.URLParam(r, "id")
	action := "product.update"
	if id == "" {
		action = "product.create"
		err := h.Pool.QueryRow(r.Context(), `INSERT INTO products (slug,name,tagline,description,scent_family,gender_tag,concentration,top_notes,heart_notes,base_notes,featured,display_order) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id`,
			req.Slug, req.Name, req.Tagline, req.Description, req.ScentFamily, req.GenderTag, req.Concentration, req.TopNotes, req.HeartNotes, req.BaseNotes, req.Featured, req.DisplayOrder).Scan(&id)
		if err != nil {
			httpx.Error(w, r, http.StatusBadRequest, "no se pudo crear el producto")
			return
		}
	} else {
		_, err := h.Pool.Exec(r.Context(), `UPDATE products SET slug=$1,name=$2,tagline=$3,description=$4,scent_family=$5,gender_tag=$6,concentration=$7,top_notes=$8,heart_notes=$9,base_notes=$10,featured=$11,display_order=$12,updated_at=now() WHERE id=$13`,
			req.Slug, req.Name, req.Tagline, req.Description, req.ScentFamily, req.GenderTag, req.Concentration, req.TopNotes, req.HeartNotes, req.BaseNotes, req.Featured, req.DisplayOrder, id)
		if err != nil {
			httpx.Error(w, r, http.StatusBadRequest, "no se pudo actualizar el producto")
			return
		}
		_, _ = h.Pool.Exec(r.Context(), `DELETE FROM product_variants WHERE product_id=$1`, id)
		_, _ = h.Pool.Exec(r.Context(), `DELETE FROM product_images WHERE product_id=$1`, id)
	}
	for _, v := range req.Variants {
		_, _ = h.Pool.Exec(r.Context(), `INSERT INTO product_variants (product_id,size_ml,price_ars_cents,stock,sku,active) VALUES ($1,$2,$3,$4,$5,true)`, id, v.SizeML, v.PriceARSCents, v.Stock, v.SKU)
	}
	for _, img := range req.Images {
		_, _ = h.Pool.Exec(r.Context(), `INSERT INTO product_images (product_id,url,alt_text,is_primary,sort_order) VALUES ($1,$2,$3,$4,$5)`, id, img.URL, img.AltText, img.IsPrimary, img.SortOrder)
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: action, ResourceID: id, Metadata: map[string]any{"slug": req.Slug}})
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"id": id})
}

func (h Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	_, err := h.Pool.Exec(r.Context(), `UPDATE products SET active=false, updated_at=now() WHERE id=$1`, chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo eliminar")
		return
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "product.delete", ResourceID: chi.URLParam(r, "id")})
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
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
	var req struct {
		Status         string `json:"status"`
		TrackingNumber string `json:"tracking_number"`
	}
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	_, err := h.Pool.Exec(r.Context(), `UPDATE orders SET status=$1, tracking_number=$2, updated_at=now() WHERE id=$3`, req.Status, req.TrackingNumber, chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo actualizar")
		return
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "order.status_change", ResourceID: chi.URLParam(r, "id"), Metadata: map[string]any{"status": req.Status}})
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
	_ = writer.Write([]string{"external_reference", "status", "customer_name", "customer_email", "total_ars_cents", "created_at"})
	rows, err := h.Pool.Query(r.Context(), `SELECT external_reference, status, customer_name, customer_email, total_ars_cents, created_at FROM orders ORDER BY created_at DESC`)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var ref, status, name, email string
		var total int64
		var created time.Time
		if rows.Scan(&ref, &status, &name, &email, &total, &created) == nil {
			_ = writer.Write([]string{ref, status, name, email, strconv.FormatInt(total, 10), created.Format(time.RFC3339)})
		}
	}
}

func (h Handler) Discounts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.listTable(w, r, `SELECT id, code, discount_type, discount_value, min_order_cents, max_uses, uses, active, expires_at FROM discount_codes ORDER BY created_at DESC`)
		return
	}
	var req discountCreateRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	_, err := h.Pool.Exec(r.Context(), `INSERT INTO discount_codes (code,discount_type,discount_value,min_order_cents,max_uses,expires_at,active) VALUES (UPPER($1),$2,$3,$4,$5,$6,true)`, req.Code, req.DiscountType, req.DiscountValue, req.MinOrderCents, req.MaxUses, req.ExpiresAt)
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo crear")
		return
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "discount.create", ResourceID: req.Code})
	httpx.WriteJSON(w, http.StatusCreated, map[string]bool{"ok": true})
}

func (h Handler) DiscountByID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		_, _ = h.Pool.Exec(r.Context(), `DELETE FROM discount_codes WHERE id=$1`, chi.URLParam(r, "id"))
		audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "discount.delete", ResourceID: chi.URLParam(r, "id")})
		httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
		return
	}
	var req discountUpdateRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	_, err := h.Pool.Exec(r.Context(), `UPDATE discount_codes SET active=COALESCE($1, active) WHERE id=$2`, req.Active, chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo actualizar")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h Handler) Homepage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.listTable(w, r, `SELECT hero_heading, hero_subheading, hero_image_url, hero_cta_label, hero_cta_url, editorial_heading, editorial_body, editorial_image_url FROM homepage_settings WHERE id=1`)
		return
	}
	var req homepageRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	_, err := h.Pool.Exec(r.Context(), `INSERT INTO homepage_settings (id, hero_heading, hero_subheading, hero_image_url, hero_cta_label, hero_cta_url, editorial_heading, editorial_body, editorial_image_url) VALUES (1,$1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT (id) DO UPDATE SET hero_heading=$1, hero_subheading=$2, hero_image_url=$3, hero_cta_label=$4, hero_cta_url=$5, editorial_heading=$6, editorial_body=$7, editorial_image_url=$8, updated_at=now()`,
		req.HeroHeading, req.HeroSubheading, req.HeroImageURL, req.HeroCTALabel, req.HeroCTAURL, req.EditorialHeading, req.EditorialBody, req.EditorialImageURL)
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo guardar")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h Handler) PublicHomepage(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.WriteJSON(w, http.StatusOK, homepageRequest{})
		return
	}
	var item homepageRequest
	err := h.Pool.QueryRow(r.Context(), `SELECT COALESCE(hero_heading,''), COALESCE(hero_subheading,''), COALESCE(hero_image_url,''), COALESCE(hero_cta_label,''), COALESCE(hero_cta_url,''), COALESCE(editorial_heading,''), COALESCE(editorial_body,''), COALESCE(editorial_image_url,'') FROM homepage_settings WHERE id=1`).
		Scan(&item.HeroHeading, &item.HeroSubheading, &item.HeroImageURL, &item.HeroCTALabel, &item.HeroCTAURL, &item.EditorialHeading, &item.EditorialBody, &item.EditorialImageURL)
	if err != nil {
		httpx.WriteJSON(w, http.StatusOK, homepageRequest{})
		return
	}
	httpx.WriteJSON(w, http.StatusOK, item)
}

func (h Handler) Contact(w http.ResponseWriter, r *http.Request) {
	h.listTable(w, r, `SELECT id, name, email, COALESCE(subject,''), message, read, created_at FROM contact_messages ORDER BY read ASC, created_at DESC LIMIT 100`)
}

func (h Handler) MarkContactRead(w http.ResponseWriter, r *http.Request) {
	_, _ = h.Pool.Exec(r.Context(), `UPDATE contact_messages SET read=true WHERE id=$1`, chi.URLParam(r, "id"))
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h Handler) LowStock(w http.ResponseWriter, r *http.Request) {
	h.listTable(w, r, `SELECT v.id, p.id AS product_id, p.name, v.size_ml, v.stock FROM product_variants v JOIN products p ON p.id=v.product_id WHERE v.stock <= 5 AND v.active=true AND p.active=true ORDER BY v.stock ASC`)
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
	rows, err := h.Pool.Query(r.Context(), `SELECT id, external_reference, status, customer_name, customer_email, COALESCE(customer_phone,''), COALESCE(shipping_address,'{}'::jsonb), total_ars_cents, COALESCE(tracking_number,''), created_at FROM orders ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return []map[string]any{}
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var id, ref, status, name, email, phone, tracking string
		var shipping []byte
		var created time.Time
		var total int64
		if rows.Scan(&id, &ref, &status, &name, &email, &phone, &shipping, &total, &tracking, &created) == nil {
			var address map[string]any
			_ = json.Unmarshal(shipping, &address)
			items = append(items, map[string]any{"id": id, "external_reference": ref, "status": status, "customer_name": name, "customer_email": email, "customer_phone": phone, "shipping_address": address, "total_ars_cents": total, "tracking_number": tracking, "created_at": created, "items": h.orderItems(r, id)})
		}
	}
	return items
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
	DisplayOrder  int           `json:"display_order"`
	Variants      []variantForm `json:"variants"`
	Images        []imageForm   `json:"images"`
}

type variantForm struct {
	SizeML        int    `json:"size_ml"`
	PriceARSCents int64  `json:"price_ars_cents"`
	Stock         int    `json:"stock"`
	SKU           string `json:"sku"`
}

type imageForm struct {
	URL       string `json:"url"`
	AltText   string `json:"alt_text"`
	IsPrimary bool   `json:"is_primary"`
	SortOrder int    `json:"sort_order"`
}

type discountCreateRequest struct {
	Code          string     `json:"code"`
	DiscountType  string     `json:"discount_type"`
	DiscountValue int64      `json:"discount_value"`
	MinOrderCents int64      `json:"min_order_cents"`
	MaxUses       *int       `json:"max_uses"`
	ExpiresAt     *time.Time `json:"expires_at"`
}

type discountUpdateRequest struct {
	Active *bool `json:"active"`
}

type homepageRequest struct {
	HeroHeading       string `json:"hero_heading"`
	HeroSubheading    string `json:"hero_subheading"`
	HeroImageURL      string `json:"hero_image_url"`
	HeroCTALabel      string `json:"hero_cta_label"`
	HeroCTAURL        string `json:"hero_cta_url"`
	EditorialHeading  string `json:"editorial_heading"`
	EditorialBody     string `json:"editorial_body"`
	EditorialImageURL string `json:"editorial_image_url"`
}

func (h Handler) actor(r *http.Request) string {
	username, _ := h.Sessions.Username(r)
	if username == "" {
		return "unknown"
	}
	return username
}

func CleanCode(v string) string { return strings.ToUpper(strings.TrimSpace(v)) }
