package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration values.
type Config struct {
	ServerPort   string
	DatabaseURL  string
	APIKey       string
	LogLevel     string
	CORSOrigins  string
	RateLimitRPS int
}

// Load reads configuration from environment variables with sensible defaults.
// It returns an error if required values are missing.
func Load() (*Config, error) {
	cfg := &Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		APIKey:       os.Getenv("API_KEY"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		CORSOrigins:  getEnv("CORS_ORIGINS", "*"),
		RateLimitRPS: getEnvInt("RATE_LIMIT_RPS", 100),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API_KEY environment variable is required")
	}

	return cfg, nil
}

// getEnv reads an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt reads an integer environment variable or returns a default value.
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
