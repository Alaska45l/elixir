package httpx

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("json response encode failed", "error", err)
	}
}

func Error(w http.ResponseWriter, r *http.Request, status int, message string) {
	WriteJSON(w, status, ErrorResponse{Error: message, RequestID: RequestID(r)})
}

func DecodeStrict(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if dec.More() {
		return errors.New("invalid json body")
	}
	return nil
}

func RequestID(r *http.Request) string {
	if v, ok := r.Context().Value(requestIDKey{}).(string); ok {
		return v
	}
	return r.Header.Get("X-Request-ID")
}

type requestIDKey struct{}

func WithRequestIDValue() any {
	return requestIDKey{}
}
