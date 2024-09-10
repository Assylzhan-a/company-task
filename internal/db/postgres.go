package db

import (
	"context"
	"github.com/assylzhan-a/company-task/pkg/logger"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPostgresConnection(dbURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		logger.Error("Unable to parse database URL", "error", err)
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		logger.Error("Unable to connect to database", "error", err)
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		logger.Error("Unable to ping database", "error", err)
		return nil, err
	}

	logger.Info("Successfully connected to database")
	return pool, nil
}
