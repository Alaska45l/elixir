package admin

import (
	"context"
	"net"
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

func (a AuthService) Login(ctx context.Context, remoteAddr, username, password string) (bool, error) {
	if a.Pool == nil {
		return false, nil
	}
	username = strings.TrimSpace(username)
	var hash string
	err := a.Pool.QueryRow(ctx, `SELECT password_hash FROM admin_users WHERE username=$1`, username).Scan(&hash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		if a.Limiter != nil {
			a.Limiter.RegisterFailure(clientIP(remoteAddr))
		}
		return false, nil
	}
	if _, err := a.Pool.Exec(ctx, `UPDATE admin_users SET last_login_at=now() WHERE username=$1`, username); err != nil {
		return false, err
	}
	if a.Limiter != nil {
		a.Limiter.Reset(clientIP(remoteAddr))
	}
	return true, nil
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
