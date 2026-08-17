// Package todo holds the todo domain model and its rules. It must not depend
// on any other layer (no HTTP, no storage, no generated code).
package todo

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// TitleMaxLength mirrors the maxLength of Todo.title in openapi/openapi.yaml.
const TitleMaxLength = 200

// Todo is a single todo item.
type Todo struct {
	ID        uuid.UUID
	Title     string
	Done      bool
	CreatedAt time.Time
}

// New creates a Todo, rejecting titles that break the domain rules. The id and
// the creation time are passed in so that the domain stays deterministic.
func New(id uuid.UUID, title string, createdAt time.Time) (Todo, error) {
	trimmed := strings.TrimSpace(title)
	switch {
	case trimmed == "":
		return Todo{}, ErrEmptyTitle
	case utf8.RuneCountInString(trimmed) > TitleMaxLength:
		return Todo{}, ErrTitleTooLong
	}

	return Todo{ID: id, Title: trimmed, CreatedAt: createdAt}, nil
}

// Filter narrows a listing of todos.
type Filter struct {
	// Done selects todos by completion state; nil means "any".
	Done *bool
}

// Matches reports whether the todo satisfies the filter.
func (f Filter) Matches(item Todo) bool {
	return f.Done == nil || item.Done == *f.Done
}
