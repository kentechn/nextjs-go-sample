import type { Todo } from "@/lib/api/todo";

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
          className="flex items-center justify-between rounded-md border border-gray-200 px-4 py-3"
        >
          <div className="flex flex-col">
            <span className={todo.done ? "line-through text-gray-400" : undefined}>
              {todo.title}
            </span>
            <time className="text-xs text-gray-400" dateTime={todo.createdAt}>
              {todo.createdAt}
            </time>
          </div>
          {onDelete ? (
            <form action={onDelete}>
              <input type="hidden" name="todoId" value={todo.id} />
              <button type="submit" className="text-sm text-red-600 hover:underline">
                削除
              </button>
            </form>
          ) : null}
        </li>
      ))}
    </ul>
  );
}
