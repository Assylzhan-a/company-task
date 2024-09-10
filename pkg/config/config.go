package config

import (
	"os"
	"strings"
)

type Config struct {
	Environment      string
	DatabaseURL      string
	ServerAddress    string
	JWTSecret        string
	LogLevel         string
	KafkaBrokers     []string
	KafkaClientID    string
	OutboxWorkerTick string
}

func Load() *Config {
	return &Config{
		Environment:      getEnv("ENVIRONMENT", "development"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/company_db?sslmode=disable"),
		ServerAddress:    getEnv("SERVER_ADDRESS", ":8080"),
		JWTSecret:        getEnv("JWT_SECRET", "your-default-secret-key"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		KafkaBrokers:     strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaClientID:    getEnv("KAFKA_CLIENT_ID", "company-service"),
		OutboxWorkerTick: getEnv("OUTBOX_WORKER_TICK", "5s"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
