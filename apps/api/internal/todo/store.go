// Package todo holds the domain logic and storage for todos.
package todo

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ErrNotFound is returned when a todo does not exist.
var ErrNotFound = errors.New("todo not found")

// Todo is a single todo item.
type Todo struct {
	ID        uuid.UUID
	Title     string
	Done      bool
	CreatedAt time.Time
}

// Store is an in-memory todo repository. Replace it with a real database
// implementation by satisfying the same method set.
type Store struct {
	mu    sync.RWMutex
	items map[uuid.UUID]Todo
	now   func() time.Time
}

// NewStore creates an empty Store.
func NewStore() *Store {
	return &Store{items: make(map[uuid.UUID]Todo), now: time.Now}
}

// List returns todos sorted by creation time, optionally filtered by done state.
func (s *Store) List(done *bool) []Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Todo, 0, len(s.items))
	for _, item := range s.items {
		if done != nil && item.Done != *done {
			continue
		}
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })

	return out
}

// Get returns the todo with the given id.
func (s *Store) Get(id uuid.UUID) (Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	if !ok {
		return Todo{}, ErrNotFound
	}

	return item, nil
}

// Create stores a new todo with the given title.
func (s *Store) Create(title string) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := Todo{ID: uuid.New(), Title: title, CreatedAt: s.now().UTC()}
	s.items[item.ID] = item

	return item
}

// Delete removes the todo with the given id.
func (s *Store) Delete(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)

	return nil
}
