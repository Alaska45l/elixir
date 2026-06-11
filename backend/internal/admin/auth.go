package admin

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Pool     *pgxpool.Pool
	Sessions SessionManager
	Limiter  LoginLimiter
}

type LoginLimiter interface {
	RegisterFailure(ip string)
	Reset(ip string)
}

func (a AuthService) Login(ctx context.Context, clientAddr, username, password string) (bool, error) {
	if a.Pool == nil {
		return false, nil
	}
	username = strings.TrimSpace(username)
	var hash string
	err := a.Pool.QueryRow(ctx, `SELECT password_hash FROM admin_users WHERE username=$1`, username).Scan(&hash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		if a.Limiter != nil {
			a.Limiter.RegisterFailure(clientIP(clientAddr))
		}
		return false, nil
	}
	if _, err := a.Pool.Exec(ctx, `UPDATE admin_users SET last_login_at=now() WHERE username=$1`, username); err != nil {
		return false, err
	}
	if a.Limiter != nil {
		a.Limiter.Reset(clientIP(clientAddr))
	}
	return true, nil
}

func requestClientIP(r *http.Request) string {
	if ip := firstForwardedIP(r.Header.Get("X-Forwarded-For")); ip != "" {
		return ip
	}
	if ip := firstForwardedIP(r.Header.Get("X-Real-IP")); ip != "" {
		return ip
	}
	return clientIP(r.RemoteAddr)
}

func firstForwardedIP(raw string) string {
	for _, value := range strings.Split(raw, ",") {
		ip := strings.TrimSpace(value)
		if net.ParseIP(ip) != nil {
			return ip
		}
	}
	return ""
}

func clientIP(remote string) string {
	ip, _, err := net.SplitHostPort(remote)
	if err == nil && ip != "" {
		return ip
	}
	return remote
}

func IsSecureEnv(appEnv string) bool {
	return appEnv != "development"
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type meResponse struct {
	Username string `json:"username"`
}
