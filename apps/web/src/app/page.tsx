import { createTodoAction, deleteTodoAction } from "@/app/actions";
import { TodoForm } from "@/components/todo/TodoForm";
import { TodoList } from "@/components/todo/TodoList";
import { fetchTodos } from "@/lib/api/todos";

// Rendered on every request (SSR) so the list always reflects the Go API.
export const dynamic = "force-dynamic";

export default async function Home() {
  const todos = await fetchTodos();

  return (
    <main className="mx-auto flex w-full max-w-xl flex-col gap-6 px-6 py-16">
      <header className="flex flex-col gap-1">
        <h1 className="text-2xl font-bold">Todos</h1>
        <p className="text-sm text-gray-500">Next.js (SSR) + Go / OpenAPI 仕様ファースト</p>
      </header>
      <TodoForm action={createTodoAction} />
      <TodoList todos={todos} onDelete={deleteTodoAction} />
    </main>
  );
}
