package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kentechn/nextjs-go-sample/apps/api/internal/openapi"
	"github.com/kentechn/nextjs-go-sample/apps/api/internal/server"
	"github.com/kentechn/nextjs-go-sample/apps/api/internal/todo"
)

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()

	handler, err := server.NewRouter(
		server.New(todo.NewStore(), "test"),
		server.Config{AllowedOrigins: []string{"http://localhost:3000"}},
	)
	require.NoError(t, err)

	return handler
}

func TestGetHealth(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	newTestHandler(t).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	require.Equal(t, http.StatusOK, rec.Code)

	var body openapi.Health
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Equal(t, openapi.HealthStatus("ok"), body.Status)
}

func TestCreateThenListTodo(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t)

	createRec := httptest.NewRecorder()
	createReq := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(`{"title":"ship it"}`))
	createReq.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(createRec, createReq)
	require.Equal(t, http.StatusCreated, createRec.Code)

	listRec := httptest.NewRecorder()
	handler.ServeHTTP(listRec, httptest.NewRequest(http.MethodGet, "/todos?status=open", nil))
	require.Equal(t, http.StatusOK, listRec.Code)

	var list openapi.TodoList
	require.NoError(t, json.Unmarshal(listRec.Body.Bytes(), &list))
	require.Len(t, list.Todos, 1)
	require.Equal(t, "ship it", list.Todos[0].Title)
}

// The request validator middleware rejects bodies that violate the spec before
// they reach a handler.
func TestCreateTodoRejectsInvalidBody(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(`{"title":""}`))
	req.Header.Set("Content-Type", "application/json")
	newTestHandler(t).ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetTodoNotFound(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	path := "/todos/2f1c3f8a-0b1d-4d5e-8a7b-9c0d1e2f3a4b"
	newTestHandler(t).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))

	require.Equal(t, http.StatusNotFound, rec.Code)
}
