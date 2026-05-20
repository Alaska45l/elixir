package payments

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"elixir/backend/internal/httpx"
)

type Handler struct {
	Service Service
}

func (h Handler) CreatePreference(w http.ResponseWriter, r *http.Request) {
	var req PreferenceRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	out, err := h.Service.CreatePreference(r.Context(), req.ExternalReference)
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo crear la preferencia")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, out)
}

func (h Handler) Webhook(w http.ResponseWriter, r *http.Request) {
	paymentID := r.URL.Query().Get("id")
	if paymentID == "" {
		paymentID = r.URL.Query().Get("data.id")
	}
	if paymentID == "" {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if data, ok := body["data"].(map[string]any); ok {
			paymentID = stringFromAny(data["id"])
		}
		if paymentID == "" {
			paymentID = stringFromAny(body["id"])
		}
	}
	if paymentID != "" {
		if err := h.Service.ProcessWebhook(r.Context(), paymentID); err != nil {
			slog.Error("mercadopago webhook", "error", err, "request_id", httpx.RequestID(r))
		}
	}
	w.WriteHeader(http.StatusOK)
}

func stringFromAny(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return stringValue(t)
	default:
		return ""
	}
}
