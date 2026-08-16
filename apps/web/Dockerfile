# syntax=docker/dockerfile:1.19
# ビルドコンテキストはリポジトリルート（pnpm workspace 全体が必要）。
ARG NODE_VERSION=22.23.2

FROM node:${NODE_VERSION}-alpine AS base
ENV PNPM_HOME=/pnpm PATH=/pnpm:$PATH NEXT_TELEMETRY_DISABLED=1
RUN corepack enable
WORKDIR /app

# マニフェストだけを先に置くことで、ソース変更では依存インストール層を再実行しない。
FROM base AS deps
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY apps/web/package.json apps/web/
COPY e2e/package.json e2e/
RUN --mount=type=cache,id=pnpm-store,target=/pnpm/store \
    pnpm install --frozen-lockfile --filter web...

FROM deps AS build
COPY openapi openapi
COPY apps/web apps/web
RUN --mount=type=cache,id=next-cache,target=/app/apps/web/.next/cache \
    pnpm --filter web build

# ローカル開発用: next dev（Fast Refresh）。ソースは compose でバインドマウントする。
FROM deps AS dev
ENV NODE_ENV=development
EXPOSE 3000
CMD ["pnpm", "--filter", "web", "dev"]

FROM node:${NODE_VERSION}-alpine AS runtime
WORKDIR /app
ENV NODE_ENV=production NEXT_TELEMETRY_DISABLED=1 PORT=3000 HOSTNAME=0.0.0.0
COPY --from=build /app/apps/web/.next/standalone ./
COPY --from=build /app/apps/web/.next/static ./apps/web/.next/static
COPY --from=build /app/apps/web/public ./apps/web/public
EXPOSE 3000
USER node
CMD ["node", "apps/web/server.js"]
