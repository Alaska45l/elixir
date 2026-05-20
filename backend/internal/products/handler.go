package products

import (
	"errors"
	"net/http"
	"strconv"

	"elixir/backend/internal/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	Service Service
}

func (h Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Get("/search", h.Search)
	r.Get("/{slug}", h.Detail)
	return r
}

func (h Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var featured *bool
	if q.Has("featured") {
		v := q.Get("featured") == "true"
		featured = &v
	}
	filters := ListFilters{
		Featured:      featured,
		Family:        q.Get("family"),
		Gender:        q.Get("gender"),
		Concentration: q.Get("concentration"),
		InStock:       q.Get("in_stock") == "true",
		Search:        q.Get("search"),
		Limit:         intParam(q.Get("limit"), 24),
		Offset:        intParam(q.Get("offset"), 0),
		MinPrice:      int64Param(q.Get("min_price")),
		MaxPrice:      int64Param(q.Get("max_price")),
	}
	items, err := h.Service.List(r.Context(), filters)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo obtener el catálogo")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items, "limit": filters.Limit, "offset": filters.Offset})
}

func (h Handler) Detail(w http.ResponseWriter, r *http.Request) {
	item, err := h.Service.Detail(r.Context(), chi.URLParam(r, "slug"))
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.Error(w, r, http.StatusNotFound, "fragancia no encontrada")
		return
	}
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo obtener la fragancia")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, item)
}

func (h Handler) Search(w http.ResponseWriter, r *http.Request) {
	items, err := h.Service.Search(r.Context(), r.URL.Query().Get("q"))
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo buscar")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
}

func intParam(v string, fallback int) int {
	if n, err := strconv.Atoi(v); err == nil {
		return n
	}
	return fallback
}

func int64Param(v string) int64 {
	if n, err := strconv.ParseInt(v, 10, 64); err == nil {
		return n
	}
	return 0
}
