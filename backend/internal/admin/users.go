package admin

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"elixir/backend/internal/audit"
	"elixir/backend/internal/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type adminUserResponse struct {
	Username    string     `json:"username"`
	CreatedAt   time.Time  `json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

type adminUserWriteRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type adminPasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type adminPasswordResetRequest struct {
	NewPassword string `json:"new_password"`
}

func (h Handler) AdminUsers(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		if r.Method == http.MethodGet {
			httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": []any{}})
			return
		}
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	if r.Method == http.MethodGet {
		rows, err := h.Pool.Query(r.Context(), `SELECT username, created_at, last_login_at FROM admin_users ORDER BY username ASC`)
		if err != nil {
			httpx.Error(w, r, http.StatusInternalServerError, "no se pudieron listar usuarios")
			return
		}
		defer rows.Close()
		items := []adminUserResponse{}
		for rows.Next() {
			var item adminUserResponse
			var lastLogin sql.NullTime
			if err := rows.Scan(&item.Username, &item.CreatedAt, &lastLogin); err != nil {
				httpx.Error(w, r, http.StatusInternalServerError, "no se pudieron leer usuarios")
				return
			}
			if lastLogin.Valid {
				item.LastLoginAt = &lastLogin.Time
			}
			items = append(items, item)
		}
		httpx.WriteJSON(w, http.StatusOK, map[string]any{"items": items})
		return
	}

	var req adminUserWriteRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	if err := validateUsername(req.Username); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}
	if err := validatePassword(req.Password, req.Username); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo proteger la contraseña")
		return
	}
	_, err = h.Pool.Exec(r.Context(), `INSERT INTO admin_users (username, password_hash) VALUES ($1,$2)`, req.Username, string(hash))
	if err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "no se pudo crear el usuario; revisá que no exista")
		return
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "admin_user.create", ResourceID: req.Username})
	httpx.WriteJSON(w, http.StatusCreated, map[string]bool{"ok": true})
}

func (h Handler) AdminUserByUsername(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	username := strings.TrimSpace(chi.URLParam(r, "username"))
	if err := validateUsername(username); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}
	switch r.Method {
	case http.MethodDelete:
		current := h.actor(r)
		if username == current {
			httpx.Error(w, r, http.StatusBadRequest, "no podés eliminar el usuario con la sesión abierta")
			return
		}
		count, err := h.adminUserCount(r)
		if err != nil {
			httpx.Error(w, r, http.StatusInternalServerError, "no se pudo validar usuarios")
			return
		}
		if count <= 1 {
			httpx.Error(w, r, http.StatusBadRequest, "debe quedar al menos un usuario administrador")
			return
		}
		tag, err := h.Pool.Exec(r.Context(), `DELETE FROM admin_users WHERE username=$1`, username)
		if err != nil || tag.RowsAffected() == 0 {
			httpx.Error(w, r, http.StatusBadRequest, "no se pudo eliminar el usuario")
			return
		}
		audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: current, Action: "admin_user.delete", ResourceID: username})
		httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
	case http.MethodPost:
		var req adminPasswordResetRequest
		if err := httpx.DecodeStrict(r, &req); err != nil {
			httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
			return
		}
		if err := validatePassword(req.NewPassword, username); err != nil {
			httpx.Error(w, r, http.StatusBadRequest, err.Error())
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			httpx.Error(w, r, http.StatusInternalServerError, "no se pudo proteger la contraseña")
			return
		}
		tag, err := h.Pool.Exec(r.Context(), `UPDATE admin_users SET password_hash=$1 WHERE username=$2`, string(hash), username)
		if err != nil || tag.RowsAffected() == 0 {
			httpx.Error(w, r, http.StatusBadRequest, "no se pudo actualizar la contraseña")
			return
		}
		audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: h.actor(r), Action: "admin_user.password_reset", ResourceID: username})
		httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
	default:
		httpx.Error(w, r, http.StatusMethodNotAllowed, "método no permitido")
	}
}

func (h Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if h.Pool == nil {
		httpx.Error(w, r, http.StatusServiceUnavailable, "base de datos no configurada")
		return
	}
	username := h.actor(r)
	if username == "" || username == "unknown" {
		httpx.Error(w, r, http.StatusUnauthorized, "sesión requerida")
		return
	}
	var req adminPasswordRequest
	if err := httpx.DecodeStrict(r, &req); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, "cuerpo inválido")
		return
	}
	if err := validatePassword(req.NewPassword, username); err != nil {
		httpx.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}
	var hash string
	err := h.Pool.QueryRow(r.Context(), `SELECT password_hash FROM admin_users WHERE username=$1`, username).Scan(&hash)
	if errors.Is(err, pgx.ErrNoRows) {
		httpx.Error(w, r, http.StatusUnauthorized, "usuario no encontrado")
		return
	}
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.CurrentPassword)) != nil {
		httpx.Error(w, r, http.StatusUnauthorized, "la contraseña actual no coincide")
		return
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo proteger la contraseña")
		return
	}
	if _, err := h.Pool.Exec(r.Context(), `UPDATE admin_users SET password_hash=$1 WHERE username=$2`, string(newHash), username); err != nil {
		httpx.Error(w, r, http.StatusInternalServerError, "no se pudo actualizar la contraseña")
		return
	}
	audit.Log(r.Context(), h.Pool, audit.Event{ActorUsername: username, Action: "admin_user.password_change", ResourceID: username})
	httpx.WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h Handler) adminUserCount(r *http.Request) (int, error) {
	var count int
	err := h.Pool.QueryRow(r.Context(), `SELECT COUNT(*) FROM admin_users`).Scan(&count)
	return count, err
}
