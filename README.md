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
| CI | GitHub Actions (lint / typecheck / build / codegen 差分 / e2e) |

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
```

## セットアップ

前提: Node 22 系 + pnpm 11、Go 1.26。

```bash
make setup          # pnpm install + go mod download
cp apps/web/.env.example apps/web/.env.local
```

## 開発

```bash
make dev            # API (:8080) と Web (:3000) を同時起動
make dev-api
make dev-web
make storybook      # http://localhost:6006
```

## 仕様変更の流れ（spec-first）

1. `openapi/openapi.yaml` を編集する
2. `make gen` で Go/TS の型を再生成する
3. Go 側は生成された `StrictServerInterface` を満たすまでコンパイルエラーになる → 実装する
4. フロントは `apiClient` の型が変わるので、呼び出し側を直す

CI の `codegen` ジョブが「生成物が仕様と一致しているか」を検証するため、
生成し忘れたままのコミットは落ちる。

## よく使うコマンド

```bash
make lint           # biome + golangci-lint
make fmt            # 自動修正
make test           # tsc --noEmit + go test
make e2e            # Playwright（api/web を自動起動）
make build
make gen-check
```
