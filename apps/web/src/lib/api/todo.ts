import { z } from "zod";

import type { components } from "./schema.gen";

export type Todo = components["schemas"]["Todo"];
export type TodoStatus = components["schemas"]["TodoStatus"];

/**
 * Input validation for the create-todo form. The constraints mirror
 * CreateTodoRequest in openapi/openapi.yaml.
 */
export const createTodoSchema = z.object({
  title: z
    .string()
    .trim()
    .min(1, "タイトルを入力してください")
    .max(200, "200文字以内で入力してください"),
});

export type CreateTodoInput = z.infer<typeof createTodoSchema>;

export const todoStatusSchema = z.enum(["all", "open", "done"]);
