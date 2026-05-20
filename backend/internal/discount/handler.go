package discount

import (
	"net/http"

	"elixir/backend/internal/httpx"
)

type Handler struct {
	Service Service
}

type validateRequest struct {
	Code          string `json:"code"`
	SubtotalCents int64  `json:"subtotal_ars_cents"`
}

func (h Handler) Validate(w http.ResponseWriter, r *http.Request) {
	var req validateRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, h.Service.Validate(r.Context(), req.Code, req.SubtotalCents))
}
