// cmd/api/main.go

package main

import (
	"context"
	"errors"
	"flag"
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
	userHandler "github.com/assylzhan-a/company-task/internal/user/delivery/handler"
	userRepository "github.com/assylzhan-a/company-task/internal/user/repository"
	userUsecase "github.com/assylzhan-a/company-task/internal/user/usecase"
	"github.com/assylzhan-a/company-task/pkg/config"
	"github.com/assylzhan-a/company-task/pkg/logger"
)

func main() {
	cfg := config.Load()

	logger.InitLogger(cfg.LogLevel)

	// Parse command line flags
	migrateFlag := flag.Bool("migrate", false, "Run database migrations")
	flag.Parse()

	// Connect to the database
	dbPool, err := db.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Run migrations if the flag is set
	if *migrateFlag {
		if err := db.RunMigrations(cfg.DatabaseURL); err != nil {
			logger.Error("Failed to run migrations", "error", err)
			os.Exit(1)
		}
		logger.Info("Migrations completed successfully")
		return
	}

	companyRepo := companyRepository.NewPostgresRepository(dbPool)
	companyUseCase := companyUsecase.NewCompanyUseCase(companyRepo)
	companyHandler := companyHandler.NewCompanyHandler(companyUseCase)

	userRepo := userRepository.NewPostgresUserRepository(dbPool)
	userUseCase := userUsecase.NewUserUseCase(userRepo)
	userHandler := userHandler.NewUserHandler(userUseCase)

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
		logger.Info("Starting server", "address", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Server exiting")
}
