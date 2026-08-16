// Package memory provides in-memory implementations of the domain repositories.
// Replace them with a database implementation by satisfying the same interfaces.
package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/google/uuid"

	"github.com/kentechn/nextjs-go-sample/apps/api/internal/domain/todo"
)

// TodoRepository stores todos in a map guarded by a mutex.
type TodoRepository struct {
	mu    sync.RWMutex
	items map[uuid.UUID]todo.Todo
}

var _ todo.Repository = (*TodoRepository)(nil)

// NewTodoRepository creates an empty repository.
func NewTodoRepository() *TodoRepository {
	return &TodoRepository{items: make(map[uuid.UUID]todo.Todo)}
}

// List returns the matching todos sorted by creation time.
func (r *TodoRepository) List(_ context.Context, filter todo.Filter) ([]todo.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]todo.Todo, 0, len(r.items))
	for _, item := range r.items {
		if filter.Matches(item) {
			out = append(out, item)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })

	return out, nil
}

// FindByID returns the todo with the given id.
func (r *TodoRepository) FindByID(_ context.Context, id uuid.UUID) (todo.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[id]
	if !ok {
		return todo.Todo{}, todo.ErrNotFound
	}

	return item, nil
}

// Create stores a new todo.
func (r *TodoRepository) Create(_ context.Context, item todo.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[item.ID]; ok {
		return todo.ErrAlreadyExists
	}
	r.items[item.ID] = item

	return nil
}

// Delete removes the todo with the given id.
func (r *TodoRepository) Delete(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return todo.ErrNotFound
	}
	delete(r.items, id)

	return nil
}
