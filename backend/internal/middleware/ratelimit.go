package middleware

import (
	"net"
	"net/http"
	"strconv"
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
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		if ip == "" {
			ip = r.RemoteAddr
		}
		if retry := l.retryAfter(ip); retry > 0 {
			w.Header().Set("Retry-After", strconv.Itoa(int(retry.Seconds())))
			httpx.Error(w, r, http.StatusTooManyRequests, "demasiados intentos")
			return
		}
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
