package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string // "development" | "production"
	Port        string
	BaseURL     string
	FrontendURL string // where OAuth callbacks redirect back to

	DatabaseURL string

	JWTSecret       []byte
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration

	GoogleClientID     string
	GoogleClientSecret string

	GitHubClientID     string
	GitHubClientSecret string

	// Cloudflare Turnstile
	TurnstileSecret string
	CaptchaEnabled  bool

	AllowedOrigins []string

	RateLimitMax        int
	RateLimitExpiration time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:              getEnv("APP_ENV", "development"),
		Port:                getEnv("PORT", "8080"),
		BaseURL:             getEnv("BASE_URL", "http://localhost:8080"),
		FrontendURL:         getEnv("FRONTEND_URL", "http://localhost:3000"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		AccessTokenTTL:      getDuration("ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTokenTTL:     getDuration("REFRESH_TOKEN_TTL", 720*time.Hour), // 30 days
		GoogleClientID:      os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:  os.Getenv("GOOGLE_CLIENT_SECRET"),
		GitHubClientID:      os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret:  os.Getenv("GITHUB_CLIENT_SECRET"),
		TurnstileSecret:     os.Getenv("TURNSTILE_SECRET"),
		CaptchaEnabled:      getBool("CAPTCHA_ENABLED", true),
		AllowedOrigins:      splitCSV(getEnv("ALLOWED_ORIGINS", "http://localhost:3000")),
		RateLimitMax:        getInt("RATE_LIMIT_MAX", 60),
		RateLimitExpiration: getDuration("RATE_LIMIT_WINDOW", time.Minute),
	}

	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET wajib di-set dan minimal 32 karakter (untuk keamanan signing token)")
	}
	cfg.JWTSecret = []byte(secret)

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL wajib di-set")
	}

	if cfg.CaptchaEnabled && cfg.TurnstileSecret == "" {
		return nil, fmt.Errorf("CAPTCHA_ENABLED=true tetapi TURNSTILE_SECRET kosong")
	}

	return cfg, nil
}

func (c *Config) IsProduction() bool { return c.AppEnv == "production" }

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
