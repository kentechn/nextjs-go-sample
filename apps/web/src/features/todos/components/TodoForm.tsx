"use client";

import { useActionState } from "react";

import type { FormState } from "../types";

type TodoFormProps = {
  action: (state: FormState, formData: FormData) => Promise<FormState>;
};

export function TodoForm({ action }: TodoFormProps) {
  const [state, formAction, isPending] = useActionState<FormState, FormData>(action, {});

  return (
    <form action={formAction} className="flex flex-col gap-2">
      <div className="flex gap-2">
        <input
          type="text"
          name="title"
          placeholder="やることを入力"
          aria-label="やること"
          className="flex-1 rounded-md border border-gray-300 px-3 py-2"
        />
        <button
          type="submit"
          disabled={isPending}
          className="rounded-md bg-black px-4 py-2 text-white disabled:opacity-50"
        >
          追加
        </button>
      </div>
      {state.error ? (
        <p role="alert" data-testid="todo-form-error" className="text-sm text-red-600">
          {state.error}
        </p>
      ) : null}
    </form>
  );
}
