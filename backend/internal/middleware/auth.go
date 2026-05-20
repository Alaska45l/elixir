package middleware

import (
	"net/http"

	"elixir/backend/internal/admin"
	"elixir/backend/internal/httpx"
)

func AdminSession(sessions admin.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := sessions.Username(r); err != nil {
				httpx.Error(w, r, http.StatusUnauthorized, "sesión requerida")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
