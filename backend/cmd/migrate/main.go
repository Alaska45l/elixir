package main

import (
	"context"
	"log"
	"strings"

	"elixir/backend/internal/config"
	"elixir/backend/internal/db"
)

func main() {
	cfg := config.Load()
	if cfg.DatabaseURL == "" || strings.Contains(cfg.DatabaseURL, "user:pass@host/dbname") {
		log.Fatal("DATABASE_URL is not configured; edit backend/.env with your Neon/PostgreSQL connection string")
	}
	pool, err := db.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	if pool != nil {
		defer pool.Close()
	}
	if err := db.RunMigrations(context.Background(), pool, "migrations"); err != nil {
		log.Fatal(err)
	}
	log.Println("migrations applied")
}
