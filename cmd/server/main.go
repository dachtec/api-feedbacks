package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dev/api-feedbacks/internal/config"
	"github.com/dev/api-feedbacks/internal/handler"
	"github.com/dev/api-feedbacks/internal/repository/postgres"
	"github.com/dev/api-feedbacks/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Configure structured logging
	logLevel := slog.LevelInfo
	switch cfg.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	// Connect to PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("failed to ping database", "error", err)
		os.Exit(1)
	}
	slog.Info("connected to database")

	// Run migrations
	if err := runMigrations(ctx, pool); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("migrations completed")

	// Initialize layers
	repo := postgres.NewFeedbackRepo(pool)
	svc := service.NewFeedbackService(repo)
	feedbackHandler := handler.NewFeedbackHandler(svc)
	router := handler.NewRouter(feedbackHandler, cfg.APIKey)

	// Start HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		slog.Info("server starting", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	slog.Info("shutting down server", "signal", sig.String())

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped gracefully")
}

// runMigrations executes SQL migration files embedded in the binary.
func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migration := `
	CREATE TABLE IF NOT EXISTS feedbacks (
		feedback_id   VARCHAR(10)  PRIMARY KEY,
		user_id       VARCHAR(10)  NOT NULL,
		feedback_type VARCHAR(50)  NOT NULL CHECK (feedback_type IN ('bug','sugerencia','elogio','duda','queja')),
		rating        INTEGER      NOT NULL CHECK (rating >= 1 AND rating <= 5),
		comment       TEXT         NOT NULL,
		created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
		updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_feedbacks_user_id    ON feedbacks(user_id);
	CREATE INDEX IF NOT EXISTS idx_feedbacks_type       ON feedbacks(feedback_type);
	CREATE INDEX IF NOT EXISTS idx_feedbacks_rating     ON feedbacks(rating);
	CREATE INDEX IF NOT EXISTS idx_feedbacks_created_at ON feedbacks(created_at);
	`

	_, err := pool.Exec(ctx, migration)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
