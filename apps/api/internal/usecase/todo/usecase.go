// Package todo implements the todo application use cases. It depends on the
// domain only, never on HTTP or on a concrete storage implementation.
package todo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/kentechn/nextjs-go-sample/apps/api/internal/domain/todo"
)

// UseCase groups the todo use cases.
type UseCase struct {
	repo  todo.Repository
	now   func() time.Time
	newID func() uuid.UUID
}

// Option overrides a UseCase dependency; used by tests to freeze time and ids.
type Option func(*UseCase)

// WithClock replaces the clock.
func WithClock(now func() time.Time) Option {
	return func(u *UseCase) { u.now = now }
}

// WithIDGenerator replaces the id generator.
func WithIDGenerator(newID func() uuid.UUID) Option {
	return func(u *UseCase) { u.newID = newID }
}

// New creates a UseCase backed by the given repository.
func New(repo todo.Repository, options ...Option) *UseCase {
	useCase := &UseCase{repo: repo, now: time.Now, newID: uuid.New}
	for _, option := range options {
		option(useCase)
	}

	return useCase
}

// List returns the todos matching the input, oldest first.
func (u *UseCase) List(ctx context.Context, in ListInput) ([]todo.Todo, error) {
	items, err := u.repo.List(ctx, todo.Filter{Done: in.Done})
	if err != nil {
		return nil, fmt.Errorf("list todos: %w", err)
	}

	return items, nil
}

// Get returns a single todo.
func (u *UseCase) Get(ctx context.Context, id uuid.UUID) (todo.Todo, error) {
	item, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return todo.Todo{}, fmt.Errorf("find todo: %w", err)
	}

	return item, nil
}

// Create validates the title and stores a new todo.
func (u *UseCase) Create(ctx context.Context, in CreateInput) (todo.Todo, error) {
	item, err := todo.New(u.newID(), in.Title, u.now().UTC())
	if err != nil {
		return todo.Todo{}, fmt.Errorf("new todo: %w", err)
	}

	if err := u.repo.Create(ctx, item); err != nil {
		return todo.Todo{}, fmt.Errorf("create todo: %w", err)
	}

	return item, nil
}

// Delete removes a todo.
func (u *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}

	return nil
}
