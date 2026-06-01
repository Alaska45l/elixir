package config

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

const devSessionSecret = "dev-secret-change-me-dev-secret-change-me-dev-secret-change-me"

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
	hours, err := strconv.Atoi(env("SESSION_DURATION_HOURS", "8"))
	if err != nil {
		hours = 0
	}
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
		SessionSecret:     env("SESSION_SECRET", devSessionSecret),
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

func (c Config) Validate() error {
	appEnv := strings.ToLower(strings.TrimSpace(c.AppEnv))
	if c.SessionDuration <= 0 {
		return errors.New("SESSION_DURATION_HOURS must be an integer greater than zero")
	}
	if len(strings.TrimSpace(c.SessionSecret)) < 32 {
		return errors.New("SESSION_SECRET must be at least 32 characters")
	}
	if appEnv != "development" {
		if strings.TrimSpace(c.DatabaseURL) == "" {
			return errors.New("DATABASE_URL is required outside development")
		}
		if c.SessionSecret == devSessionSecret || strings.Contains(c.SessionSecret, "change-this") {
			return errors.New("SESSION_SECRET must be a production random secret")
		}
		if strings.TrimSpace(c.MPAccessToken) != "" && strings.TrimSpace(c.MPWebhookSecret) == "" {
			return errors.New("MP_WEBHOOK_SECRET is required when MercadoPago is configured")
		}
	}
	return nil
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
