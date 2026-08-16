import type { Meta, StoryObj } from "@storybook/nextjs-vite";

import { TodoList } from "./TodoList";

const meta = {
  title: "todos/TodoList",
  component: TodoList,
  parameters: { layout: "centered" },
} satisfies Meta<typeof TodoList>;

export default meta;

type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    todos: [
      {
        id: "0f4a1a54-4a1e-4f3a-9f1a-1b2c3d4e5f60",
        title: "OpenAPI 仕様を読む",
        done: false,
        createdAt: "2026-01-01T00:00:00Z",
      },
      {
        id: "1f4a1a54-4a1e-4f3a-9f1a-1b2c3d4e5f61",
        title: "task dev を実行する",
        done: true,
        createdAt: "2026-01-02T00:00:00Z",
      },
    ],
  },
};

export const Empty: Story = {
  args: { todos: [] },
};

export const LongTitle: Story = {
  args: {
    todos: [
      {
        id: "2f4a1a54-4a1e-4f3a-9f1a-1b2c3d4e5f62",
        title: "あ".repeat(200),
        done: false,
        createdAt: "2026-01-03T00:00:00Z",
      },
    ],
    onDelete: () => {},
  },
};
