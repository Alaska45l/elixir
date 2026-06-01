package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"elixir/backend/internal/httpx"
)

type LoginLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
	limit    int
	window   time.Duration
}

func NewLoginLimiter(limit int, window time.Duration) *LoginLimiter {
	return &LoginLimiter{attempts: map[string][]time.Time{}, limit: limit, window: window}
}

func NewAPILimiter(limit int, window time.Duration) *LoginLimiter {
	return NewLoginLimiter(limit, window)
}

func (l *LoginLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := requestIP(r)
		if retry := l.retryAfter(ip); retry > 0 {
			w.Header().Set("Retry-After", strconv.Itoa(int(retry.Seconds())))
			httpx.Error(w, r, http.StatusTooManyRequests, "demasiados intentos")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (l *LoginLimiter) Throttle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := requestIP(r)
		if retry := l.retryAfter(ip); retry > 0 {
			w.Header().Set("Retry-After", strconv.Itoa(int(retry.Seconds())))
			httpx.Error(w, r, http.StatusTooManyRequests, "demasiadas solicitudes")
			return
		}
		l.RegisterFailure(ip)
		next.ServeHTTP(w, r)
	})
}

func (l *LoginLimiter) RegisterFailure(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	l.attempts[ip] = append(l.prune(l.attempts[ip], now), now)
}

func (l *LoginLimiter) Reset(ip string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.attempts, ip)
}

func (l *LoginLimiter) retryAfter(ip string) time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	hits := l.prune(l.attempts[ip], now)
	l.attempts[ip] = hits
	if len(hits) < l.limit {
		return 0
	}
	return l.window - now.Sub(hits[0])
}

func (l *LoginLimiter) prune(hits []time.Time, now time.Time) []time.Time {
	cutoff := now.Add(-l.window)
	kept := hits[:0]
	for _, hit := range hits {
		if hit.After(cutoff) {
			kept = append(kept, hit)
		}
	}
	return kept
}

func requestIP(r *http.Request) string {
	if ip := forwardedIP(r.Header.Get("X-Forwarded-For")); ip != "" {
		return ip
	}
	if ip := forwardedIP(r.Header.Get("X-Real-IP")); ip != "" {
		return ip
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func forwardedIP(raw string) string {
	for _, value := range strings.Split(raw, ",") {
		ip := strings.TrimSpace(value)
		if net.ParseIP(ip) != nil {
			return ip
		}
	}
	return ""
}
