import type { Meta, StoryObj } from "@storybook/nextjs-vite";

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

export const WithError: Story = {
  args: {
    action: async () => ({ error: "タイトルを入力してください" }),
  },
};
