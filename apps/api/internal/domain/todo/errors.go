package todo

import "errors"

var (
	// ErrNotFound is returned when a todo does not exist.
	ErrNotFound = errors.New("todo not found")
	// ErrEmptyTitle is returned when a title contains no visible characters.
	ErrEmptyTitle = errors.New("title must not be empty")
	// ErrTitleTooLong is returned when a title exceeds TitleMaxLength.
	ErrTitleTooLong = errors.New("title is too long")
	// ErrAlreadyExists is returned when a todo with the same id is stored twice.
	ErrAlreadyExists = errors.New("todo already exists")
)
