package config

import "os"

type Config struct {
	Address            string
	DatabaseURL        string
	FrontendURL        string
	SessionSecret      string
	GoogleClientID     string
	GoogleClientSecret string
}

func Load() Config {
	return Config{
		Address:            env("ADDRESS", "localhost:8080"),
		DatabaseURL:        env("DATABASE_URL", ""),
		FrontendURL:        env("FRONTEND_URL", "http://localhost:5173"),
		SessionSecret:      env("SESSION_SECRET", ""),
		GoogleClientID:     env("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: env("GOOGLE_CLIENT_SECRET", ""),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
