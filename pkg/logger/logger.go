// pkg/logger/logger.go

package logger

import (
	"log/slog"
	"os"
)

var log *slog.Logger

func InitLogger(logLevel string) {
	level := getLogLevel(logLevel)
	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	log = slog.New(handler)
}

func getLogLevel(logLevel string) slog.Level {
	switch logLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func Info(msg string, args ...any) {
	log.Info(msg, args...)
}

func Error(msg string, args ...any) {
	log.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	log.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	log.Warn(msg, args...)
}
