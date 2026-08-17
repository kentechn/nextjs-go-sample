import type { components } from "@/shared/api/schema.gen";

/** Types come from the generated schema; never redefine the API shape by hand. */
export type Todo = components["schemas"]["Todo"];
export type TodoStatus = components["schemas"]["TodoStatus"];

/**
 * State returned by the create-todo Server Action. It lives here rather than in
 * actions.ts so Client Components can type their props without importing a
 * "use server" module.
 */
export type FormState = { error?: string };
