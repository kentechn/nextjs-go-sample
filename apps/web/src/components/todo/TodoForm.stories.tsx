import type { Meta, StoryObj } from "@storybook/nextjs-vite";
import { expect, userEvent, within } from "storybook/test";

import { TodoForm } from "./TodoForm";

const meta = {
  title: "todo/TodoForm",
  component: TodoForm,
  parameters: { layout: "centered" },
  args: {
    action: async () => ({}),
  },
} satisfies Meta<typeof TodoForm>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {};

// The error only exists in the action state, so the story has to submit the
// form once to show it.
export const WithError: Story = {
  args: {
    action: async () => ({ error: "タイトルを入力してください" }),
  },
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement);

    await userEvent.click(canvas.getByRole("button", { name: "追加" }));
    await expect(await canvas.findByTestId("todo-form-error")).toHaveTextContent(
      "タイトルを入力してください",
    );
  },
};
