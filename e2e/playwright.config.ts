import { defineConfig, devices } from "@playwright/test";

const webPort = Number(process.env.WEB_PORT ?? 3000);
const apiPort = Number(process.env.API_PORT ?? 8080);
const baseURL = process.env.E2E_BASE_URL ?? `http://127.0.0.1:${webPort}`;

export default defineConfig({
  testDir: "./tests",
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1,
  reporter: process.env.CI ? [["html"], ["github"]] : [["html"]],
  use: {
    baseURL,
    trace: "on-first-retry",
    screenshot: "only-on-failure",
  },
  projects: [{ name: "chromium", use: { ...devices["Desktop Chrome"] } }],
  // Boots the Go API and the Next.js server (production build) before the run.
  webServer: process.env.E2E_BASE_URL
    ? undefined
    : [
        {
          command: "go run ./cmd/api",
          cwd: "../apps/api",
          port: apiPort,
          env: { PORT: String(apiPort) },
          reuseExistingServer: !process.env.CI,
          stdout: "pipe",
          stderr: "pipe",
        },
        {
          command: `pnpm --filter web start --port ${webPort}`,
          cwd: "..",
          port: webPort,
          env: { API_BASE_URL: `http://127.0.0.1:${apiPort}` },
          reuseExistingServer: !process.env.CI,
          stdout: "pipe",
          stderr: "pipe",
        },
      ],
});
