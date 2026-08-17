// Package rest is the HTTP presentation layer: it implements the generated
// OpenAPI contract by translating requests into use case calls and results back
// into spec-shaped responses. It is the only layer that knows about HTTP and
// about the generated code.
package rest

import (
	"context"
	"errors"

	domain "github.com/kentechn/nextjs-go-sample/apps/api/internal/domain/todo"
	"github.com/kentechn/nextjs-go-sample/apps/api/internal/openapi"
	usecase "github.com/kentechn/nextjs-go-sample/apps/api/internal/usecase/todo"
)

// Handler implements openapi.StrictServerInterface. Because the interface is
// generated from the spec, a spec change that is not implemented fails the build.
type Handler struct {
	todos   *usecase.UseCase
	version string
}

var _ openapi.StrictServerInterface = (*Handler)(nil)

// NewHandler creates a Handler backed by the todo use cases.
func NewHandler(todos *usecase.UseCase, version string) *Handler {
	return &Handler{todos: todos, version: version}
}

// GetHealth implements the liveness probe.
func (h *Handler) GetHealth(
	_ context.Context,
	_ openapi.GetHealthRequestObject,
) (openapi.GetHealthResponseObject, error) {
	return openapi.GetHealth200JSONResponse{Status: "ok", Version: h.version}, nil
}

// ListTodos returns todos, optionally filtered by status.
func (h *Handler) ListTodos(
	ctx context.Context,
	request openapi.ListTodosRequestObject,
) (openapi.ListTodosResponseObject, error) {
	items, err := h.todos.List(ctx, usecase.ListInput{Done: doneFilter(request.Params.Status)})
	if err != nil {
		return nil, err
	}

	return openapi.ListTodos200JSONResponse(toAPITodoList(items)), nil
}

// CreateTodo creates a todo.
func (h *Handler) CreateTodo(
	ctx context.Context,
	request openapi.CreateTodoRequestObject,
) (openapi.CreateTodoResponseObject, error) {
	if request.Body == nil {
		return openapi.CreateTodo400JSONResponse{
			ErrorJSONResponse: invalidArgument(domain.ErrEmptyTitle),
		}, nil
	}

	item, err := h.todos.Create(ctx, usecase.CreateInput{Title: request.Body.Title})
	if err != nil {
		if isInvalidArgument(err) {
			return openapi.CreateTodo400JSONResponse{ErrorJSONResponse: invalidArgument(err)}, nil
		}

		return nil, err
	}

	return openapi.CreateTodo201JSONResponse(toAPITodo(item)), nil
}

// GetTodo returns a single todo.
func (h *Handler) GetTodo(
	ctx context.Context,
	request openapi.GetTodoRequestObject,
) (openapi.GetTodoResponseObject, error) {
	item, err := h.todos.Get(ctx, request.TodoId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return openapi.GetTodo404JSONResponse{ErrorJSONResponse: notFound()}, nil
		}

		return nil, err
	}

	return openapi.GetTodo200JSONResponse(toAPITodo(item)), nil
}

// DeleteTodo deletes a todo.
func (h *Handler) DeleteTodo(
	ctx context.Context,
	request openapi.DeleteTodoRequestObject,
) (openapi.DeleteTodoResponseObject, error) {
	if err := h.todos.Delete(ctx, request.TodoId); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return openapi.DeleteTodo404JSONResponse{ErrorJSONResponse: notFound()}, nil
		}

		return nil, err
	}

	return openapi.DeleteTodo204Response{}, nil
}
