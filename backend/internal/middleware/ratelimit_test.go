package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoginRateLimiterSixthAttemptReturns429(t *testing.T) {
	limiter := NewLoginLimiter(5, 10*time.Minute)
	for i := 0; i < 5; i++ {
		limiter.RegisterFailure("127.0.0.1")
	}
	req := httptest.NewRequest(http.MethodPost, "/api/admin/login", nil)
	req.RemoteAddr = "127.0.0.1:5000"
	rec := httptest.NewRecorder()
	limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rec.Code)
	}
}
