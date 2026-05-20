package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv          string
	Port            string
	FrontendURL     string
	BackendURL      string
	DatabaseURL     string
	SessionSecret   string
	SessionDuration time.Duration
	MPAccessToken   string
	MPWebhookSecret string
	AllowedOrigins  []string
	WhatsAppNumber  string
	LogLevel        string
}

func Load() Config {
	loadDotEnv(".env")
	loadDotEnv("backend/.env")
	hours, _ := strconv.Atoi(env("SESSION_DURATION_HOURS", "8"))
	databaseURL := env("DATABASE_URL", "")
	if strings.Contains(databaseURL, "user:pass@host/dbname") {
		databaseURL = ""
	}
	return Config{
		AppEnv:          env("APP_ENV", "development"),
		Port:            env("PORT", "8080"),
		FrontendURL:     env("FRONTEND_URL", "http://localhost:5173"),
		BackendURL:      env("BACKEND_URL", "http://localhost:8080"),
		DatabaseURL:     databaseURL,
		SessionSecret:   env("SESSION_SECRET", "dev-secret-change-me-dev-secret-change-me-dev-secret-change-me"),
		SessionDuration: time.Duration(hours) * time.Hour,
		MPAccessToken:   os.Getenv("MP_ACCESS_TOKEN"),
		MPWebhookSecret: os.Getenv("MP_WEBHOOK_SECRET"),
		AllowedOrigins:  split(env("ALLOWED_ORIGINS", "http://localhost:5173")),
		WhatsAppNumber:  env("WHATSAPP_NUMBER", "5491100000000"),
		LogLevel:        env("LOG_LEVEL", "info"),
	}
}

func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key != "" && os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
}

func env(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func split(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if s := strings.TrimSpace(part); s != "" {
			out = append(out, s)
		}
	}
	return out
}
