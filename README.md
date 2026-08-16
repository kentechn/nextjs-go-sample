# nextjs-go-sample

Next.js (SSR) + Go のモノレポ。HTTP の契約は `openapi/openapi.yaml` が唯一の真実で、
フロント/バックの型は両方そこから自動生成する（生成物は手で編集しない）。

## スタック

| レイヤ | 技術 |
| --- | --- |
| フロント | Next.js 16.3.1 (App Router / SSR), TypeScript, Tailwind CSS 4 |
| バリデーション | Zod 4 |
| Lint / Format | Biome 2 (TS), golangci-lint 2 (Go) |
| UI カタログ | Storybook 10 (`@storybook/nextjs-vite`) |
| バックエンド | Go 1.26, chi, oapi-codegen v2 (strict server), OpenAPI リクエスト検証ミドルウェア |
| 型生成 | Go: oapi-codegen / TS: openapi-typescript + openapi-fetch |
| E2E | Playwright 1.62 |
| API ドキュメント | Redocly CLI 2 (lint / build-docs / bundle) |
| タスクランナー | Task (Taskfile) |
| コンテナ | Docker マルチステージ + Buildx (bake) / Compose でホットリロード |
| CI | GitHub Actions (lint / typecheck / build / codegen 差分 / docs / docker / e2e) |

## ディレクトリ

```
openapi/openapi.yaml         API 仕様（単一の真実）
apps/web                     Next.js アプリ
  src/lib/api/schema.gen.ts  openapi-typescript の生成物
  src/lib/api/client.ts      openapi-fetch の型付きクライアント
apps/api                     Go API
  internal/openapi/api.gen.go  oapi-codegen の生成物（型 + サーバIF + 埋め込みspec）
  internal/server            生成 IF の実装（StrictServerInterface）
  internal/todo              ドメイン / ストア（インメモリ）
e2e                          Playwright テスト
Taskfile.yml                 task コマンドの入口
taskfiles/app.yml            アプリ（web / api / e2e / docker）のタスク
taskfiles/docs.yml           OpenAPI（Redocly）のタスク
docker-bake.hcl              buildx bake の定義（本番イメージ）
compose.yaml                 ローカル開発（ホットリロード）
```

## セットアップ

前提: [Task](https://taskfile.dev) と、ローカル実行する場合は Node 22 系 + pnpm 11、Go 1.26。
Docker で動かす場合は Docker + Buildx だけあればよい。

```bash
task setup          # pnpm install + go mod download
cp apps/web/.env.example apps/web/.env.local
task --list-all     # 全タスク一覧
```

## 開発

ローカル実行:

```bash
task dev              # API (:8080, air) と Web (:3000, next dev) を同時起動
task app:api:dev
task app:web:dev
task app:storybook    # http://localhost:6006
```

Docker（ホットリロード付き）:

```bash
task app:docker:up    # api は air、web は next dev。ソースはバインドマウント
task app:docker:logs
task app:docker:down
```

本番イメージ（マルチステージ + buildx キャッシュ）:

```bash
task app:docker:build       # docker buildx bake（api + web）
task app:docker:build:api
task app:docker:build:web
```

- `apps/api/Dockerfile`: `deps`（go.mod/go.sum のみ先にコピー）→ `build`（BuildKit の module/build キャッシュマウント）→ `runtime`（distroless static, nonroot, 約 12MB）。`dev` ステージは air。
- `apps/web/Dockerfile`: `deps`（package.json / lockfile のみ先にコピー、pnpm store をキャッシュマウント）→ `build`（Next.js standalone）→ `runtime`（standalone のみをコピー）。`dev` ステージは next dev。
- `docker-bake.hcl` はレイヤキャッシュを `.buildx-cache/` に import/export する。キャッシュのエクスポートには docker-container ドライバが必要なので、`task app:docker:build` が専用ビルダーを自動作成する。

## API ドキュメント（Redocly）

```bash
task docs:lint      # redocly lint（ルールは redocly.yaml）
task docs:build     # redocly build-docs → docs/api.html
task docs:preview   # 生成した HTML を :4000 で配信
task docs:bundle    # redocly bundle → docs/openapi.bundled.yaml
task docs:stats     # redocly stats
```

生成物（`docs/api.html`, `docs/openapi.bundled.yaml`）はコミットせず、CI の `docs` ジョブが artifact として出力する。

## 仕様変更の流れ（spec-first）

1. `openapi/openapi.yaml` を編集する
2. `task gen` で Go/TS の型を再生成する
3. Go 側は生成された `StrictServerInterface` を満たすまでコンパイルエラーになる → 実装する
4. フロントは `apiClient` の型が変わるので、呼び出し側を直す

CI の `codegen` ジョブが「生成物が仕様と一致しているか」を検証するため、
生成し忘れたままのコミットは落ちる。

## よく使うコマンド

```bash
task lint           # biome + golangci-lint + redocly lint
task fmt            # 自動修正
task test           # tsc --noEmit + go test
task e2e            # Playwright（api/web を自動起動）
task build
task gen:check
task clean
```
