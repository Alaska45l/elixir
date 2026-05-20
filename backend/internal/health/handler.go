package health

import (
	"net/http"
	"time"

	"elixir/backend/internal/db"
	"elixir/backend/internal/httpx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	Pool    *pgxpool.Pool
	Started time.Time
	Version string
}

func (h Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", h.Health)
	return mux
}

func (h Handler) Health(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	dbStatus := "ok"
	dbPool := map[string]int32{"total": 0, "idle": 0, "acquired": 0}
	if err := db.Ping(r.Context(), h.Pool); err != nil {
		dbStatus = "unavailable"
		if h.Pool != nil {
			status = "degraded"
		}
	} else if h.Pool != nil {
		stat := h.Pool.Stat()
		dbPool = map[string]int32{"total": stat.TotalConns(), "idle": stat.IdleConns(), "acquired": stat.AcquiredConns()}
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"status":         status,
		"db":             dbStatus,
		"db_pool":        dbPool,
		"version":        h.Version,
		"uptime_seconds": int(time.Since(h.Started).Seconds()),
	})
}
