// Command api runs the HTTP API server.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/kentechn/nextjs-go-sample/apps/api/internal/server"
	"github.com/kentechn/nextjs-go-sample/apps/api/internal/todo"
)

const (
	version         = "0.1.0"
	shutdownTimeout = 10 * time.Second
	readTimeout     = 15 * time.Second
	writeTimeout    = 15 * time.Second
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := run(logger); err != nil {
		logger.Error("server exited with error", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	store := todo.NewStore()
	store.Create("read the OpenAPI spec")
	store.Create("run task dev")

	handler, err := server.NewRouter(server.New(store, version), server.Config{
		AllowedOrigins: splitAndTrim(env("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
	})
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:         ":" + env("PORT", "8080"),
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		logger.Info("api listening", slog.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		logger.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		return httpServer.Shutdown(shutdownCtx)
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func splitAndTrim(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}

	return out
}
