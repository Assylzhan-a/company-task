// pkg/config/config.go

package config

import (
	"os"
)

type Config struct {
	Environment   string
	DatabaseURL   string
	ServerAddress string
	JWTSecret     string
	LogLevel      string
}

func Load() *Config {
	return &Config{
		Environment:   getEnv("ENVIRONMENT", "development"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/company_db?sslmode=disable"),
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		JWTSecret:     getEnv("JWT_SECRET", "your-default-secret-key"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
