package main

import (
	"context"
	"errors"
	"github.com/assylzhan-a/company-task/config"
	"github.com/assylzhan-a/company-task/internal/db/repository"
	handler "github.com/assylzhan-a/company-task/internal/delivery/http"
	uc "github.com/assylzhan-a/company-task/internal/domain/usecase"
	"github.com/assylzhan-a/company-task/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/assylzhan-a/company-task/internal/db"
	"github.com/assylzhan-a/company-task/internal/kafka"
	"github.com/assylzhan-a/company-task/internal/worker"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()
	log := logger.NewLogger(cfg.LogLevel)

	// Connect to the database
	dbPool, err := db.NewPostgresConnection(cfg.DatabaseURL, log)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Run migrations
	if err := db.RunMigrations("migrations", cfg.DatabaseURL); err != nil {
		log.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	log.Info("Migrations completed successfully")

	r := chi.NewRouter()

	//middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// Prometheus' metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	// repositories
	userRepo := repository.NewUserRepository(dbPool)
	companyRepo := repository.NewCompanyRepository(dbPool)

	// use cases
	userUseCase := uc.NewUserUseCase(userRepo)
	companyUseCase := uc.NewCompanyUseCase(companyRepo, log)

	// Initialize handlers
	handler.NewUserHandler(r, userUseCase)
	handler.NewCompanyHandler(r, companyUseCase)

	// Initialize Kafka producer
	kafkaProducer := kafka.NewProducer(cfg.KafkaBrokers, log)
	defer kafkaProducer.Close()

	// Initialize and start outbox worker
	outboxWorker := worker.NewOutboxWorker(companyRepo, kafkaProducer, log)
	go outboxWorker.Start(context.Background())

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	// Start server
	go func() {
		log.Info("Starting server", "address", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	log.Info("Server exiting")
}
