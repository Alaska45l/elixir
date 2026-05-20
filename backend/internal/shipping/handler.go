package shipping

import (
	"net/http"

	"elixir/backend/internal/httpx"
)

type Handler struct {
	Service Service
}

func (h Handler) Zones(w http.ResponseWriter, r *http.Request) {
	zones, err := h.Service.Zones(r.Context())
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudieron obtener los envíos")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": zones})
}
