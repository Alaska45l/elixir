package admin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"
)

const SessionCookie = "elixir_admin_session"

type SessionManager struct {
	Secret   string
	Duration time.Duration
	Secure   bool
}

func (s SessionManager) Create(username string, now time.Time) string {
	exp := now.Add(s.Duration).Unix()
	payload := username + "|" + time.Unix(exp, 0).UTC().Format(time.RFC3339)
	sig := s.sign(payload)
	return base64.RawURLEncoding.EncodeToString([]byte(payload + "|" + sig))
}

func (s SessionManager) Validate(token string, now time.Time) (string, bool) {
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return "", false
	}
	parts := strings.Split(string(raw), "|")
	if len(parts) != 3 {
		return "", false
	}
	payload := parts[0] + "|" + parts[1]
	if !hmac.Equal([]byte(parts[2]), []byte(s.sign(payload))) {
		return "", false
	}
	exp, err := time.Parse(time.RFC3339, parts[1])
	if err != nil || now.After(exp) {
		return "", false
	}
	return parts[0], true
}

func (s SessionManager) SetCookie(w http.ResponseWriter, username string) {
	http.SetCookie(w, &http.Cookie{
		Name: SessionCookie, Value: s.Create(username, time.Now()), Path: "/",
		HttpOnly: true, Secure: s.Secure, SameSite: http.SameSiteStrictMode, MaxAge: int(s.Duration.Seconds()),
	})
}

func (s SessionManager) ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: SessionCookie, Value: "", Path: "/", HttpOnly: true, Secure: s.Secure, SameSite: http.SameSiteStrictMode, MaxAge: -1})
}

func (s SessionManager) Username(r *http.Request) (string, error) {
	cookie, err := r.Cookie(SessionCookie)
	if err != nil {
		return "", err
	}
	if user, ok := s.Validate(cookie.Value, time.Now()); ok {
		return user, nil
	}
	return "", errors.New("invalid session")
}

func (s SessionManager) sign(payload string) string {
	mac := hmac.New(sha256.New, []byte(s.Secret))
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
