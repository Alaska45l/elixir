package shipping

import (
	"net/http"
	"strings"

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

func (h Handler) Quote(w http.ResponseWriter, r *http.Request) {
	var req QuoteRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	if strings.TrimSpace(req.DestinationPostalCode) == "" {
		httpx.Error(w, r, http.StatusBadRequest, "código postal requerido")
		return
	}
	options, err := h.Service.Quote(r.Context(), req)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo cotizar el envío")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": options})
}
