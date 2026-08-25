package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Address            string
	DatabaseURL        string
	FrontendURL        string
	SessionSecret      string
	SessionTTL         time.Duration
	CookieSecure       bool
	GoogleClientID     string
	GoogleClientSecret string
}

func Load() Config {
	loadDotEnv(".env", "../.env")

	return Config{
		Address:            env("ADDRESS", "localhost:8080"),
		DatabaseURL:        env("DATABASE_URL", "postgres://form_builder:form_builder@localhost:5432/form_builder?sslmode=disable"),
		FrontendURL:        env("FRONTEND_URL", "http://localhost:5173"),
		SessionSecret:      env("SESSION_SECRET", "dev-session-secret-change-me-before-production"),
		SessionTTL:         envHours("SESSION_TTL_HOURS", 168),
		CookieSecure:       envBool("COOKIE_SECURE", false),
		GoogleClientID:     env("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: env("GOOGLE_CLIENT_SECRET", ""),
	}
}

func (cfg Config) Validate() error {
	if strings.TrimSpace(cfg.Address) == "" {
		return fmt.Errorf("ADDRESS is required")
	}

	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if len(cfg.SessionSecret) < 32 {
		return fmt.Errorf("SESSION_SECRET must have at least 32 characters")
	}

	if cfg.SessionTTL <= 0 {
		return fmt.Errorf("SESSION_TTL_HOURS must be greater than zero")
	}

	return nil
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func envBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func envHours(key string, fallback int) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return time.Duration(fallback) * time.Hour
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return time.Duration(fallback) * time.Hour
	}

	return time.Duration(parsed) * time.Hour
}

func loadDotEnv(paths ...string) {
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			continue
		}

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
			if key == "" || os.Getenv(key) != "" {
				continue
			}

			_ = os.Setenv(key, value)
		}

		_ = file.Close()
	}
}
