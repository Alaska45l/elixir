package audit

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Event struct {
	ActorUsername string
	Action        string
	ResourceID    string
	Metadata      map[string]any
}

func Log(ctx context.Context, pool *pgxpool.Pool, e Event) {
	if pool == nil || e.ActorUsername == "" || e.Action == "" {
		return
	}
	metadata, err := json.Marshal(e.Metadata)
	if err != nil {
		slog.Error("audit metadata marshal failed", "error", err)
		metadata = []byte("{}")
	}
	go func() {
		logCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if _, err := pool.Exec(logCtx, `INSERT INTO audit_log (actor, action, resource_id, metadata) VALUES ($1,$2,$3,$4)`, e.ActorUsername, e.Action, e.ResourceID, metadata); err != nil {
			slog.Error("audit log failed", "error", err)
		}
	}()
}
