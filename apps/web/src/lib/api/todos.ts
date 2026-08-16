import { apiClient } from "./client";
import type { Todo, TodoStatus } from "./todo";

export async function fetchTodos(status: TodoStatus = "all"): Promise<Todo[]> {
  const { data, error } = await apiClient.GET("/todos", {
    params: { query: { status } },
    cache: "no-store",
  });
  if (error) {
    throw new Error(`failed to list todos: ${error.code} ${error.message}`);
  }

  return data.todos;
}

export async function createTodo(title: string): Promise<Todo> {
  const { data, error } = await apiClient.POST("/todos", { body: { title } });
  if (error) {
    throw new Error(`failed to create todo: ${error.code} ${error.message}`);
  }

  return data;
}

export async function deleteTodo(todoId: string): Promise<void> {
  const { error } = await apiClient.DELETE("/todos/{todoId}", {
    params: { path: { todoId } },
  });
  if (error) {
    throw new Error(`failed to delete todo: ${error.code} ${error.message}`);
  }
}
