package products

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchRouteReturnsOK(t *testing.T) {
	h := Handler{Service: Service{Repo: Repository{}}}
	req := httptest.NewRequest(http.MethodGet, "/search?q=oud", nil)
	rec := httptest.NewRecorder()
	h.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
}
