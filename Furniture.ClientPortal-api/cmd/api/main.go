package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/yourusername/furniture-api/internal/handlers"
)

func main() {
	// Initialize logger
	zerolog.TimeFieldFormat = time.RFC3339
	logOutput := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(logOutput).With().Timestamp().Logger()

	// Load configuration (in production, use environment variables or config file)
	cfg := struct {
		DB struct {
			Host     string
			Port     string
			User     string
			Password string
			DBName   string
			SSLMode  string
		}
		Server struct {
			Port string
		}
	}{
		DB: struct {
			Host     string
			Port     string
			User     string
			Password string
			DBName   string
			SSLMode  string
		}{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "1881"),
			DBName:   getEnv("DB_NAME", "furniture_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: struct {
			Port string
		}{
			Port: getEnv("PORT", "8080"),
		},
	}

	// Initialize database connection (degraded mode if connection fails)
	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode,
	))
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to database; starting in degraded mode (health only)")
		db = nil
	} else {
		defer db.Close()

		// Test database connection
		if err := db.Ping(); err != nil {
			logger.Error().Err(err).Msg("Failed to ping database; starting in degraded mode (health only)")
			db = nil
		} else {
			// Set connection pool settings
			db.SetMaxOpenConns(25)
			db.SetMaxIdleConns(5)
			db.SetConnMaxLifetime(5 * time.Minute)
		}
	}

	// Initialize router
	r := mux.NewRouter()

	// Add CORS middleware
	cors := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"http://localhost:5000", "http://localhost:5001", "https://localhost:61341", "http://localhost:61342"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// API routes (only register if DB is available)
	if db != nil {
		api := r.PathPrefix("/api").Subrouter()

		itemHandler := handlers.NewItemHandler(db)
		itemHandler.RegisterRoutes(api)

		categoryHandler := handlers.NewCategoryHandler(db)
		categoryHandler.RegisterRoutes(api)

		inventoryHandler := handlers.NewInventoryHandler(db)
		inventoryHandler.RegisterRoutes(api)

		userHandler := handlers.NewUserHandler(db.DB)
		userHandler.RegisterRoutes(api)

		orderHandler := handlers.NewOrderHandler(db)
		orderHandler.RegisterRoutes(api)

		customerHandler := handlers.NewCustomerHandler(db)
		customerHandler.RegisterRoutes(api)

		// New handlers for payments, ratings, and profiles
		paymentHandler := handlers.NewPaymentHandler(db)
		paymentHandler.RegisterRoutes(api)

		ratingHandler := handlers.NewRatingHandler(db)
		ratingHandler.RegisterRoutes(api)

		profileHandler := handlers.NewProfileHandler(db)
		profileHandler.RegisterRoutes(api)
	} else {
		logger.Warn().Msg("Database not connected; API routes disabled. Only /health endpoint is available.")
	}

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      cors(r),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info().Str("port", cfg.Server.Port).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server error")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server stopped")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
