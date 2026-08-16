package rest

import (
	"errors"
	"net/http"

	domain "github.com/kentechn/nextjs-go-sample/apps/api/internal/domain/todo"
	"github.com/kentechn/nextjs-go-sample/apps/api/internal/openapi"
)

// Error codes returned in the spec's Error schema.
const (
	codeInvalidRequest  = "invalid_request"
	codeInvalidArgument = "invalid_argument"
	codeNotFound        = "not_found"
	codeInternal        = "internal"
)

// isInvalidArgument reports whether the error is a domain rule violation, which
// maps to 400 rather than 500.
func isInvalidArgument(err error) bool {
	return errors.Is(err, domain.ErrEmptyTitle) || errors.Is(err, domain.ErrTitleTooLong)
}

func invalidArgument(err error) openapi.ErrorJSONResponse {
	return openapi.ErrorJSONResponse{Code: codeInvalidArgument, Message: err.Error()}
}

func notFound() openapi.ErrorJSONResponse {
	return openapi.ErrorJSONResponse{Code: codeNotFound, Message: domain.ErrNotFound.Error()}
}

// ErrorResponse writes a spec-shaped error body. It is used by the middleware
// and handler error hooks, which run outside the generated handlers.
func ErrorResponse(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, openapi.Error{Code: code, Message: message})
}
