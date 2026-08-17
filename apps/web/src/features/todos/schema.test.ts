import { describe, expect, it } from "vitest";

import { createTodoSchema, todoStatusSchema } from "./schema";

describe("createTodoSchema", () => {
  it("trims the title", () => {
    const parsed = createTodoSchema.parse({ title: "  ship it  " });

    expect(parsed.title).toBe("ship it");
  });

  it.each([
    ["empty", ""],
    ["whitespace only", "   "],
  ])("rejects a %s title", (_name, title) => {
    const result = createTodoSchema.safeParse({ title });

    expect(result.success).toBe(false);
    expect(result.error?.issues[0].message).toBe("タイトルを入力してください");
  });

  // The boundary mirrors maxLength: 200 in openapi/openapi.yaml.
  it("accepts 200 characters and rejects 201", () => {
    expect(createTodoSchema.safeParse({ title: "あ".repeat(200) }).success).toBe(true);

    const tooLong = createTodoSchema.safeParse({ title: "あ".repeat(201) });
    expect(tooLong.success).toBe(false);
    expect(tooLong.error?.issues[0].message).toBe("200文字以内で入力してください");
  });

  it("rejects a missing title", () => {
    expect(createTodoSchema.safeParse({}).success).toBe(false);
  });
});

describe("todoStatusSchema", () => {
  it("accepts the statuses defined in the spec", () => {
    expect(todoStatusSchema.parse("open")).toBe("open");
    expect(todoStatusSchema.safeParse("archived").success).toBe(false);
  });
});
