package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCSRFAllowsConfiguredCrossSiteOrigin(t *testing.T) {
	handler := CSRF([]string{"https://elixir-demo-front.onrender.com"}, "https://elixir-demo-api.onrender.com")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/admin/login", nil)
	req.Header.Set("Origin", "https://elixir-demo-front.onrender.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected allowed origin to pass, got status %d and body %q", rec.Code, rec.Body.String())
	}
}

func TestCSRFRejectsUnconfiguredCrossSiteOrigin(t *testing.T) {
	handler := CSRF([]string{"https://elixir-demo-front.onrender.com"}, "https://elixir-demo-api.onrender.com")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/admin/login", nil)
	req.Header.Set("Origin", "https://evil.example")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected unconfigured origin to be rejected, got status %d", rec.Code)
	}
}

func TestCSRFRejectsCrossSiteWithoutOrigin(t *testing.T) {
	handler := CSRF([]string{"https://elixir-demo-front.onrender.com"}, "https://elixir-demo-api.onrender.com")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/admin/login", nil)
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected cross-site request without origin to be rejected, got status %d", rec.Code)
	}
}
