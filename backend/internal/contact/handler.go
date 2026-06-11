package contact

import (
	"errors"
	"net/http"
	"strings"

	"elixir/backend/internal/httpx"
)

type Handler struct {
	Service Service
}

func (h Handler) Message(w http.ResponseWriter, r *http.Request) {
	var req MessageRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	if strings.TrimSpace(req.Website) != "" {
		httpx.WriteJSON(w, http.StatusCreated, map[string]bool{"ok": true})
		return
	}
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Message) == "" {
		httpx.Error(w, r, http.StatusBadRequest, "complete los campos obligatorios")
		return
	}
	if err := h.Service.SaveMessage(r.Context(), req); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, errInvalidContactInput) {
			status = http.StatusBadRequest
		}
		httpx.Error(w, r, status, "no se pudo enviar el mensaje")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, map[string]bool{"ok": true})
}

func (h Handler) AbandonedCart(w http.ResponseWriter, r *http.Request) {
	var req AbandonedCartRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	if err := h.Service.SaveAbandoned(r.Context(), req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo guardar el carrito")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, map[string]bool{"ok": true})
}
