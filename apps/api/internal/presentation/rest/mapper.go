package rest

import (
	domain "github.com/kentechn/nextjs-go-sample/apps/api/internal/domain/todo"
	"github.com/kentechn/nextjs-go-sample/apps/api/internal/openapi"
)

// doneFilter translates the spec's status query parameter into a domain filter.
func doneFilter(status *openapi.TodoStatus) *bool {
	if status == nil {
		return nil
	}

	switch *status {
	case openapi.Open:
		value := false

		return &value
	case openapi.Done:
		value := true

		return &value
	case openapi.All:
		return nil
	}

	return nil
}

func toAPITodoList(items []domain.Todo) openapi.TodoList {
	list := openapi.TodoList{Todos: make([]openapi.Todo, 0, len(items))}
	for _, item := range items {
		list.Todos = append(list.Todos, toAPITodo(item))
	}

	return list
}

func toAPITodo(item domain.Todo) openapi.Todo {
	return openapi.Todo{
		Id:        item.ID,
		Title:     item.Title,
		Done:      item.Done,
		CreatedAt: item.CreatedAt,
	}
}
