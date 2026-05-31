package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv            string
	Port              string
	FrontendURL       string
	BackendURL        string
	DatabaseURL       string
	SessionSecret     string
	SessionDuration   time.Duration
	SessionSameSite   string
	MPAccessToken     string
	MPWebhookSecret   string
	AllowedOrigins    []string
	WhatsAppNumber    string
	LogLevel          string
	CorreoArgAPIKey   string
	CorreoArgClientID string
	AndreaniUser      string
	AndreaniPassword  string
	AndreaniClientID  string
	OriginPostalCode  string
	R2AccountID       string
	R2AccessKey       string
	R2SecretKey       string
	R2BucketName      string
	R2PublicURL       string
}

func Load() Config {
	loadDotEnv(".env")
	loadDotEnv("backend/.env")
	hours, _ := strconv.Atoi(env("SESSION_DURATION_HOURS", "8"))
	databaseURL := env("DATABASE_URL", "")
	if strings.Contains(databaseURL, "user:pass@host/dbname") {
		databaseURL = ""
	}
	frontendURL := env("FRONTEND_URL", "http://localhost:5173")
	return Config{
		AppEnv:            env("APP_ENV", "development"),
		Port:              env("PORT", "8080"),
		FrontendURL:       frontendURL,
		BackendURL:        env("BACKEND_URL", "http://localhost:8080"),
		DatabaseURL:       databaseURL,
		SessionSecret:     env("SESSION_SECRET", "dev-secret-change-me-dev-secret-change-me-dev-secret-change-me"),
		SessionDuration:   time.Duration(hours) * time.Hour,
		SessionSameSite:   env("SESSION_SAME_SITE", "strict"),
		MPAccessToken:     os.Getenv("MP_ACCESS_TOKEN"),
		MPWebhookSecret:   os.Getenv("MP_WEBHOOK_SECRET"),
		AllowedOrigins:    allowedOrigins(env("ALLOWED_ORIGINS", ""), frontendURL),
		WhatsAppNumber:    env("WHATSAPP_NUMBER", "5491100000000"),
		LogLevel:          env("LOG_LEVEL", "info"),
		CorreoArgAPIKey:   os.Getenv("CORREO_ARG_API_KEY"),
		CorreoArgClientID: os.Getenv("CORREO_ARG_CLIENT_ID"),
		AndreaniUser:      os.Getenv("ANDREANI_USER"),
		AndreaniPassword:  os.Getenv("ANDREANI_PASSWORD"),
		AndreaniClientID:  os.Getenv("ANDREANI_CLIENT_ID"),
		OriginPostalCode:  env("ORIGIN_POSTAL_CODE", "1000"),
		R2AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		R2AccessKey:       os.Getenv("R2_ACCESS_KEY"),
		R2SecretKey:       os.Getenv("R2_SECRET_KEY"),
		R2BucketName:      os.Getenv("R2_BUCKET_NAME"),
		R2PublicURL:       os.Getenv("R2_PUBLIC_URL"),
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

func allowedOrigins(rawAllowedOrigins, frontendURL string) []string {
	return appendUniqueOrigins(split(rawAllowedOrigins), frontendURL)
}

func appendUniqueOrigins(origins []string, candidates ...string) []string {
	seen := make(map[string]bool, len(origins)+len(candidates))
	out := make([]string, 0, len(origins)+len(candidates))
	for _, origin := range append(origins, candidates...) {
		origin = strings.TrimRight(strings.TrimSpace(origin), "/")
		if origin == "" || seen[origin] {
			continue
		}
		seen[origin] = true
		out = append(out, origin)
	}
	return out
}
