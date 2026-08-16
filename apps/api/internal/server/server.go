// Package server wires the generated OpenAPI contract to the domain logic.
package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/kentechn/nextjs-go-sample/apps/api/internal/openapi"
	"github.com/kentechn/nextjs-go-sample/apps/api/internal/todo"
)

// Server implements openapi.StrictServerInterface. Because the interface is
// generated from the spec, a spec change that is not implemented fails the build.
type Server struct {
	todos   *todo.Store
	version string
}

// New creates a Server backed by the given store.
func New(todos *todo.Store, version string) *Server {
	return &Server{todos: todos, version: version}
}

var _ openapi.StrictServerInterface = (*Server)(nil)

// GetHealth implements the liveness probe.
func (s *Server) GetHealth(
	_ context.Context,
	_ openapi.GetHealthRequestObject,
) (openapi.GetHealthResponseObject, error) {
	return openapi.GetHealth200JSONResponse{Status: "ok", Version: s.version}, nil
}

// ListTodos returns todos, optionally filtered by status.
func (s *Server) ListTodos(
	_ context.Context,
	request openapi.ListTodosRequestObject,
) (openapi.ListTodosResponseObject, error) {
	var done *bool
	if request.Params.Status != nil {
		switch *request.Params.Status {
		case "open":
			value := false
			done = &value
		case "done":
			value := true
			done = &value
		case "all":
		}
	}

	items := s.todos.List(done)
	response := openapi.TodoList{Todos: make([]openapi.Todo, 0, len(items))}
	for _, item := range items {
		response.Todos = append(response.Todos, toAPITodo(item))
	}

	return openapi.ListTodos200JSONResponse(response), nil
}

// CreateTodo creates a todo.
func (s *Server) CreateTodo(
	_ context.Context,
	request openapi.CreateTodoRequestObject,
) (openapi.CreateTodoResponseObject, error) {
	if request.Body == nil || request.Body.Title == "" {
		return openapi.CreateTodo400JSONResponse{ErrorJSONResponse: openapi.ErrorJSONResponse{
			Code:    "invalid_argument",
			Message: "title must not be empty",
		}}, nil
	}

	return openapi.CreateTodo201JSONResponse(toAPITodo(s.todos.Create(request.Body.Title))), nil
}

// GetTodo returns a single todo.
func (s *Server) GetTodo(
	_ context.Context,
	request openapi.GetTodoRequestObject,
) (openapi.GetTodoResponseObject, error) {
	item, err := s.todos.Get(request.TodoId)
	if err != nil {
		if errors.Is(err, todo.ErrNotFound) {
			return openapi.GetTodo404JSONResponse{ErrorJSONResponse: notFound()}, nil
		}

		return nil, err
	}

	return openapi.GetTodo200JSONResponse(toAPITodo(item)), nil
}

// DeleteTodo deletes a todo.
func (s *Server) DeleteTodo(
	_ context.Context,
	request openapi.DeleteTodoRequestObject,
) (openapi.DeleteTodoResponseObject, error) {
	if err := s.todos.Delete(request.TodoId); err != nil {
		if errors.Is(err, todo.ErrNotFound) {
			return openapi.DeleteTodo404JSONResponse{ErrorJSONResponse: notFound()}, nil
		}

		return nil, err
	}

	return openapi.DeleteTodo204Response{}, nil
}

func notFound() openapi.ErrorJSONResponse {
	return openapi.ErrorJSONResponse{Code: "not_found", Message: "todo not found"}
}

func toAPITodo(item todo.Todo) openapi.Todo {
	return openapi.Todo{
		Id:        item.ID,
		Title:     item.Title,
		Done:      item.Done,
		CreatedAt: item.CreatedAt,
	}
}

// ErrorResponse writes a spec-shaped error body.
func ErrorResponse(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, openapi.Error{Code: code, Message: message})
}
