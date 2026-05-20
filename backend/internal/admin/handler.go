package admin

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

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
	if id == "" {
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
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"id": id})
}

func (h Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	_, err := h.Pool.Exec(r.Context(), `UPDATE products SET active=false, updated_at=now() WHERE id=$1`, chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo eliminar")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h Handler) ImportProducts(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"imported": 0, "message": "Importación CSV disponible desde el panel con procesamiento por filas."})
}

func (h Handler) Orders(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": h.recentOrders(r)})
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
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h Handler) Discounts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.listTable(w, r, `SELECT id, code, discount_type, discount_value, min_order_cents, max_uses, uses, active, expires_at FROM discount_codes ORDER BY created_at DESC`)
		return
	}
	var req map[string]any
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	_, err := h.Pool.Exec(r.Context(), `INSERT INTO discount_codes (code,discount_type,discount_value,min_order_cents,max_uses,active) VALUES (UPPER($1),$2,$3,$4,$5,true)`, req["code"], req["discount_type"], req["discount_value"], req["min_order_cents"], req["max_uses"])
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo crear")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, map[string]bool{"ok": true})
}

func (h Handler) DiscountByID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		_, _ = h.Pool.Exec(r.Context(), `DELETE FROM discount_codes WHERE id=$1`, chi.URLParam(r, "id"))
		httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
		return
	}
	var req map[string]any
	_ = json.NewDecoder(r.Body).Decode(&req)
	_, err := h.Pool.Exec(r.Context(), `UPDATE discount_codes SET active=COALESCE($1, active) WHERE id=$2`, req["active"], chi.URLParam(r, "id"))
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
	var req map[string]any
	_ = json.NewDecoder(r.Body).Decode(&req)
	_, err := h.Pool.Exec(r.Context(), `INSERT INTO homepage_settings (id, hero_heading, hero_subheading, hero_image_url, hero_cta_label, hero_cta_url, editorial_heading, editorial_body, editorial_image_url) VALUES (1,$1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT (id) DO UPDATE SET hero_heading=$1, hero_subheading=$2, hero_image_url=$3, hero_cta_label=$4, hero_cta_url=$5, editorial_heading=$6, editorial_body=$7, editorial_image_url=$8, updated_at=now()`,
		req["hero_heading"], req["hero_subheading"], req["hero_image_url"], req["hero_cta_label"], req["hero_cta_url"], req["editorial_heading"], req["editorial_body"], req["editorial_image_url"])
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo guardar")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
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
	if h.Pool == nil {
		return []map[string]any{}
	}
	rows, err := h.Pool.Query(r.Context(), `SELECT id, external_reference, status, customer_name, customer_email, total_ars_cents, created_at FROM orders ORDER BY created_at DESC LIMIT 10`)
	if err != nil {
		return []map[string]any{}
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var id, ref, status, name, email string
		var created time.Time
		var total int64
		if rows.Scan(&id, &ref, &status, &name, &email, &total, &created) == nil {
			items = append(items, map[string]any{"id": id, "external_reference": ref, "status": status, "customer_name": name, "customer_email": email, "total_ars_cents": total, "created_at": created})
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

func CleanCode(v string) string { return strings.ToUpper(strings.TrimSpace(v)) }
