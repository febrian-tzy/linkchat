package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	APIPort         string
	SessionSecret   string
}

// Load reads environment variables and returns a populated Config.
// In production (Railway), env vars are set directly — no .env file needed.
func Load() (*Config, error) {
	// Best-effort: load .env if it exists (local dev only)
	for _, path := range []string{".env", "../../.env", "../../../.env"} {
		if err := godotenv.Load(path); err == nil {
			break
		}
	}

	cfg := &Config{
		DBHost:        getEnv("DB_HOST", "127.0.0.1"),
		DBPort:        getEnv("DB_PORT", "3306"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "linkchat"),
		APIPort:       getEnv("PORT", getEnv("API_PORT", "8080")), // Railway uses PORT
		SessionSecret: getEnv("SESSION_SECRET", ""),
	}

	if cfg.SessionSecret == "" {
		return nil, fmt.Errorf("SESSION_SECRET must be set")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
