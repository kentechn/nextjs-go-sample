package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"

	"github.com/kentechn/nextjs-go-sample/apps/api/internal/openapi"
)

// Config configures the HTTP router.
type Config struct {
	// AllowedOrigins are the origins allowed to call the API from a browser.
	AllowedOrigins []string
}

// NewRouter builds the HTTP handler: request validation against the embedded
// spec runs before any handler, so invalid requests never reach the domain.
func NewRouter(srv *Server, cfg Config) (http.Handler, error) {
	spec, err := openapi.GetSpec()
	if err != nil {
		return nil, fmt.Errorf("load openapi spec: %w", err)
	}
	spec.Servers = nil

	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.Logger, middleware.Recoverer)
	router.Use(corsMiddleware(cfg.AllowedOrigins))
	router.Use(nethttpmiddleware.OapiRequestValidatorWithOptions(spec, &nethttpmiddleware.Options{
		SilenceServersWarning: true,
		ErrorHandlerWithOpts: func(
			_ context.Context,
			err error,
			w http.ResponseWriter,
			_ *http.Request,
			opts nethttpmiddleware.ErrorHandlerOpts,
		) {
			ErrorResponse(w, opts.StatusCode, "invalid_request", err.Error())
		},
	}))

	handler := openapi.NewStrictHandlerWithOptions(srv, nil, openapi.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
			ErrorResponse(w, http.StatusBadRequest, "invalid_request", err.Error())
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
			ErrorResponse(w, http.StatusInternalServerError, "internal", err.Error())
		},
	})

	return openapi.HandlerWithOptions(handler, openapi.ChiServerOptions{
		BaseRouter: router,
		ErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
			ErrorResponse(w, http.StatusBadRequest, "invalid_request", err.Error())
		},
	}), nil
}

func corsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowed[origin] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if _, ok := allowed[origin]; ok && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.Header().Set("Vary", "Origin")
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)

				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
