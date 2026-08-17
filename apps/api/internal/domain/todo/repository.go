package todo

import (
	"context"

	"github.com/google/uuid"
)

// Repository persists todos. The interface lives in the domain layer and is
// implemented by the infrastructure layer (dependency inversion), so swapping
// storage does not touch the domain or the use cases.
type Repository interface {
	// List returns the todos matching the filter, oldest first.
	List(ctx context.Context, filter Filter) ([]Todo, error)
	// FindByID returns ErrNotFound when the todo does not exist.
	FindByID(ctx context.Context, id uuid.UUID) (Todo, error)
	// Create stores a new todo and returns ErrAlreadyExists on a duplicate id.
	Create(ctx context.Context, item Todo) error
	// Delete returns ErrNotFound when the todo does not exist.
	Delete(ctx context.Context, id uuid.UUID) error
}
