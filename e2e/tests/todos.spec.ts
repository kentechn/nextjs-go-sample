import { expect, test } from "@playwright/test";

test.describe("todos", () => {
  test("SSR page renders todos seeded by the API", async ({ page }) => {
    await page.goto("/");

    await expect(page.getByRole("heading", { name: "Todos" })).toBeVisible();
    await expect(page.getByTestId("todo-item")).not.toHaveCount(0);
  });

  test("creates a todo through the server action", async ({ page }) => {
    const title = `e2e todo ${Date.now()}`;

    await page.goto("/");
    await page.getByLabel("やること").fill(title);
    await page.getByRole("button", { name: "追加" }).click();

    await expect(page.getByText(title)).toBeVisible();
  });

  test("rejects an empty title with a validation message", async ({ page }) => {
    await page.goto("/");
    await page.getByRole("button", { name: "追加" }).click();

    await expect(page.getByTestId("todo-form-error")).toHaveText("タイトルを入力してください");
  });

  test("deletes a todo", async ({ page }) => {
    const title = `delete me ${Date.now()}`;

    await page.goto("/");
    await page.getByLabel("やること").fill(title);
    await page.getByRole("button", { name: "追加" }).click();

    const item = page.getByTestId("todo-item").filter({ hasText: title });
    await expect(item).toHaveCount(1);

    await item.getByRole("button", { name: "削除" }).click();
    await expect(page.getByTestId("todo-item").filter({ hasText: title })).toHaveCount(0);
  });
});

test("api health endpoint responds", async ({ request }) => {
  const response = await request.get("http://127.0.0.1:8080/health");

  expect(response.ok()).toBeTruthy();
  expect(await response.json()).toMatchObject({ status: "ok" });
});
