"use server";

import { revalidatePath } from "next/cache";

import { createTodo, deleteTodo } from "./api";
import { createTodoSchema } from "./schema";
import type { FormState } from "./types";

export async function createTodoAction(_state: FormState, formData: FormData): Promise<FormState> {
  const parsed = createTodoSchema.safeParse({ title: formData.get("title") });
  if (!parsed.success) {
    return { error: parsed.error.issues[0].message };
  }

  try {
    await createTodo(parsed.data.title);
  } catch (error) {
    return { error: error instanceof Error ? error.message : "unknown error" };
  }

  revalidatePath("/");

  return {};
}

export async function deleteTodoAction(formData: FormData): Promise<void> {
  const todoId = formData.get("todoId");
  if (typeof todoId !== "string") {
    return;
  }

  await deleteTodo(todoId);
  revalidatePath("/");
}
