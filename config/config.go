package config

import (
	"strings"

	"github.com/spf13/viper"
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

func Load() Config {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return Config{
		Environment:      viper.GetString("ENVIRONMENT"),
		DatabaseURL:      viper.GetString("DATABASE_URL"),
		ServerAddress:    viper.GetString("SERVER_ADDRESS"),
		JWTSecret:        viper.GetString("JWT_SECRET"),
		LogLevel:         viper.GetString("LOG_LEVEL"),
		KafkaBrokers:     strings.Split(viper.GetString("KAFKA_BROKERS"), ","),
		KafkaClientID:    viper.GetString("KAFKA_CLIENT_ID"),
		OutboxWorkerTick: viper.GetString("OUTBOX_WORKER_TICK"),
	}
}
