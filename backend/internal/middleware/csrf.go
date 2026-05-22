package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"elixir/backend/internal/httpx"
)

func CSRF(allowedOrigins []string, backendURL string) func(http.Handler) http.Handler {
	allowed := map[string]bool{}
	for _, origin := range allowedOrigins {
		if origin = strings.TrimRight(strings.TrimSpace(origin), "/"); origin != "" {
			allowed[origin] = true
		}
	}
	if backendURL = strings.TrimRight(strings.TrimSpace(backendURL), "/"); backendURL != "" {
		allowed[backendURL] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}
			if site := r.Header.Get("Sec-Fetch-Site"); site == "cross-site" {
				httpx.Error(w, r, http.StatusForbidden, "origen no permitido")
				return
			}
			origin := strings.TrimRight(r.Header.Get("Origin"), "/")
			if origin == "" {
				origin = refererOrigin(r.Header.Get("Referer"))
			}
			if origin != "" && !allowed[origin] {
				httpx.Error(w, r, http.StatusForbidden, "origen no permitido")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func refererOrigin(raw string) string {
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}
	return parsed.Scheme + "://" + parsed.Host
}
