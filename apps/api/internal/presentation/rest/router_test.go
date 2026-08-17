package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kentechn/nextjs-go-sample/apps/api/internal/infrastructure/memory"
	"github.com/kentechn/nextjs-go-sample/apps/api/internal/openapi"
	"github.com/kentechn/nextjs-go-sample/apps/api/internal/presentation/rest"
	todousecase "github.com/kentechn/nextjs-go-sample/apps/api/internal/usecase/todo"
)

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()

	handler, err := rest.NewRouter(
		rest.NewHandler(todousecase.New(memory.NewTodoRepository()), "test"),
		rest.Config{AllowedOrigins: []string{"http://localhost:3000"}},
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

// Path parameters that cannot be parsed are reported in the spec error shape.
func TestGetTodoRejectsMalformedID(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	newTestHandler(t).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/todos/not-a-uuid", nil))

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var body openapi.Error
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Equal(t, "invalid_request", body.Code)
}

func TestGetTodoNotFound(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	path := "/todos/2f1c3f8a-0b1d-4d5e-8a7b-9c0d1e2f3a4b"
	newTestHandler(t).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))

	require.Equal(t, http.StatusNotFound, rec.Code)

	var body openapi.Error
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Equal(t, "not_found", body.Code)
}

func TestCreateThenDeleteTodo(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t)

	createRec := httptest.NewRecorder()
	createReq := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(`{"title":"delete me"}`))
	createReq.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(createRec, createReq)
	require.Equal(t, http.StatusCreated, createRec.Code)

	var created openapi.Todo
	require.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &created))

	deleteRec := httptest.NewRecorder()
	handler.ServeHTTP(deleteRec, httptest.NewRequest(http.MethodDelete, "/todos/"+created.Id.String(), nil))
	require.Equal(t, http.StatusNoContent, deleteRec.Code)

	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, httptest.NewRequest(http.MethodGet, "/todos/"+created.Id.String(), nil))
	require.Equal(t, http.StatusNotFound, getRec.Code)
}
