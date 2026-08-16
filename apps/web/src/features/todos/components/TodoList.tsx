import type { Todo } from "../types";

type TodoListProps = {
  todos: Todo[];
  onDelete?: (formData: FormData) => void | Promise<void>;
};

export function TodoList({ todos, onDelete }: TodoListProps) {
  if (todos.length === 0) {
    return (
      <p data-testid="todo-empty" className="text-sm text-gray-500">
        Todo はまだありません。
      </p>
    );
  }

  return (
    <ul data-testid="todo-list" className="flex flex-col gap-2">
      {todos.map((todo) => (
        <li
          key={todo.id}
          data-testid="todo-item"
          className="flex items-center justify-between gap-4 rounded-md border border-gray-200 px-4 py-3"
        >
          <div className="flex min-w-0 flex-col">
            <span className={`break-words ${todo.done ? "line-through text-gray-400" : ""}`}>
              {todo.title}
            </span>
            <time className="text-xs text-gray-400" dateTime={todo.createdAt}>
              {todo.createdAt}
            </time>
          </div>
          {onDelete ? (
            <form action={onDelete}>
              <input type="hidden" name="todoId" value={todo.id} />
              <button type="submit" className="shrink-0 text-sm text-red-600 hover:underline">
                削除
              </button>
            </form>
          ) : null}
        </li>
      ))}
    </ul>
  );
}
