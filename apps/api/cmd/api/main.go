// Command api runs the HTTP API server. This file is the composition root: it
// is the only place that picks concrete implementations and wires the layers.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/kentechn/nextjs-go-sample/apps/api/internal/infrastructure/memory"
	"github.com/kentechn/nextjs-go-sample/apps/api/internal/presentation/rest"
	todousecase "github.com/kentechn/nextjs-go-sample/apps/api/internal/usecase/todo"
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
	todos := todousecase.New(memory.NewTodoRepository())
	if err := seed(todos); err != nil {
		return err
	}

	handler, err := rest.NewRouter(rest.NewHandler(todos, version), rest.Config{
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

// seed inserts sample todos so that the SSR page has something to render.
func seed(todos *todousecase.UseCase) error {
	for _, title := range []string{"read the OpenAPI spec", "run task dev"} {
		if _, err := todos.Create(context.Background(), todousecase.CreateInput{Title: title}); err != nil {
			return fmt.Errorf("seed todos: %w", err)
		}
	}

	return nil
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
