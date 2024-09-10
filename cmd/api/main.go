package main

import (
	"context"
	"errors"
	"flag"
	"github.com/assylzhan-a/company-task/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/assylzhan-a/company-task/internal/auth"
	companyHandler "github.com/assylzhan-a/company-task/internal/company/delivery/handler"
	companyRepository "github.com/assylzhan-a/company-task/internal/company/repository"
	companyUsecase "github.com/assylzhan-a/company-task/internal/company/usecase"
	"github.com/assylzhan-a/company-task/internal/db"
	"github.com/assylzhan-a/company-task/internal/kafka"
	userHandler "github.com/assylzhan-a/company-task/internal/user/delivery/handler"
	userRepository "github.com/assylzhan-a/company-task/internal/user/repository"
	userUsecase "github.com/assylzhan-a/company-task/internal/user/usecase"
	"github.com/assylzhan-a/company-task/internal/worker"
	"github.com/assylzhan-a/company-task/pkg/config"
)

func main() {
	cfg := config.Load()

	log := logger.NewLogger(cfg.LogLevel)

	// Parse command line flags
	migrateFlag := flag.Bool("migrate", false, "Run database migrations")
	flag.Parse()

	// Connect to the database
	dbPool, err := db.NewPostgresConnection(cfg.DatabaseURL, log)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Run migrations if the flag is set
	if *migrateFlag {
		if err := db.RunMigrations(cfg.DatabaseURL); err != nil {
			log.Error("Failed to run migrations", "error", err)
			os.Exit(1)
		}
		log.Info("Migrations completed successfully")
		return
	}

	// Initialize repositories
	companyRepo := companyRepository.NewPostgresRepository(dbPool)
	userRepo := userRepository.NewPostgresUserRepository(dbPool)

	// Initialize Kafka producer
	kafkaProducer := kafka.NewProducer(cfg.KafkaBrokers, log)
	defer kafkaProducer.Close()

	// Initialize use cases
	companyUseCase := companyUsecase.NewCompanyUseCase(companyRepo, log)
	userUseCase := userUsecase.NewUserUseCase(userRepo)

	// Initialize handlers
	companyHandler := companyHandler.NewCompanyHandler(companyUseCase)
	userHandler := userHandler.NewUserHandler(userUseCase)

	// Initialize and start outbox worker
	outboxWorker := worker.NewOutboxWorker(companyRepo, kafkaProducer, log)
	go outboxWorker.Start(context.Background())

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// User routes
	r.Post("/register", userHandler.Register)
	r.Post("/login", userHandler.Login)

	// Public company route
	r.Get("/api/v1/companies/{id}", companyHandler.Get)

	// Protected company routes
	r.Route("/api/v1/companies", func(r chi.Router) {
		r.Use(auth.JWTAuth)
		r.Post("/", companyHandler.Create)
		r.Patch("/{id}", companyHandler.Patch)
		r.Delete("/{id}", companyHandler.Delete)
	})

	// Prometheus metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

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
