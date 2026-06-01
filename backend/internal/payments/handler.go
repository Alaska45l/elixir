package payments

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"elixir/backend/internal/httpx"
)

type Handler struct {
	Service       Service
	WebhookSecret string
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
	if err := h.validateWebhookSignature(r); err != nil {
		slog.Warn("mercadopago webhook rejected", "error", err, "request_id", httpx.RequestID(r))
		httpx.Error(w, r, http.StatusUnauthorized, "firma de webhook inválida")
		return
	}

	paymentID := r.URL.Query().Get("id")
	if paymentID == "" {
		paymentID = r.URL.Query().Get("data.id")
	}
	var body map[string]any
	if paymentID == "" {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil && !errors.Is(err, io.EOF) {
			httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
			return
		}
		if data, ok := body["data"].(map[string]any); ok {
			paymentID = stringValue(data["id"])
		}
		if paymentID == "" {
			paymentID = stringValue(body["id"])
		}
	}
	if paymentID == "" {
		httpx.Error(w, r, http.StatusBadRequest, "id de pago requerido")
		return
	}
	if err := h.Service.ProcessWebhook(r.Context(), paymentID); err != nil {
		slog.Error("mercadopago webhook", "error", err, "request_id", httpx.RequestID(r))
		httpx.Error(w, r, http.StatusBadGateway, "no se pudo procesar la notificación")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h Handler) validateWebhookSignature(r *http.Request) error {
	secret := strings.TrimSpace(h.WebhookSecret)
	if secret == "" {
		return errors.New("webhook secret not configured")
	}
	ts, signature, err := parseMercadoPagoSignature(r.Header.Get("x-signature"))
	if err != nil {
		return err
	}
	if err := validateMercadoPagoTimestamp(ts, 15*time.Minute); err != nil {
		return err
	}
	manifest := mercadoPagoSignatureManifest(r.URL.Query().Get("data.id"), r.Header.Get("x-request-id"), ts)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(manifest))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(strings.ToLower(signature)), []byte(expected)) {
		return errors.New("signature mismatch")
	}
	return nil
}

func parseMercadoPagoSignature(header string) (string, string, error) {
	var ts, v1 string
	for _, part := range strings.Split(header, ",") {
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		switch strings.TrimSpace(key) {
		case "ts":
			ts = strings.TrimSpace(value)
		case "v1":
			v1 = strings.TrimSpace(value)
		}
	}
	if ts == "" || v1 == "" {
		return "", "", errors.New("missing signature fields")
	}
	if _, err := hex.DecodeString(v1); err != nil {
		return "", "", fmt.Errorf("invalid signature encoding: %w", err)
	}
	return ts, v1, nil
}

func mercadoPagoSignatureManifest(dataID, requestID, ts string) string {
	var b strings.Builder
	if dataID != "" {
		b.WriteString("id:")
		b.WriteString(dataID)
		b.WriteByte(';')
	}
	if requestID != "" {
		b.WriteString("request-id:")
		b.WriteString(requestID)
		b.WriteByte(';')
	}
	if ts != "" {
		b.WriteString("ts:")
		b.WriteString(ts)
		b.WriteByte(';')
	}
	return b.String()
}

func validateMercadoPagoTimestamp(ts string, tolerance time.Duration) error {
	ms, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return errors.New("invalid signature timestamp")
	}
	signedAt := time.UnixMilli(ms)
	now := time.Now()
	if signedAt.Before(now.Add(-tolerance)) || signedAt.After(now.Add(tolerance)) {
		return errors.New("signature timestamp outside tolerance")
	}
	return nil
}
