# AGENTS.md

このリポジトリで作業するエージェント / 開発者向けの約束事。

## 原則

- `openapi/openapi.yaml` が HTTP 契約の唯一の真実。API を変えるときは必ず仕様から。
- 生成物 (`apps/api/internal/openapi/api.gen.go`, `apps/web/src/lib/api/schema.gen.ts`) は手で編集しない。
  仕様を直して `make gen` を実行する。
- 型は自前定義せず生成型（`components["schemas"][...]` / `openapi.Todo`）を使う。

## コマンド

```bash
make setup   # 依存インストール
make gen     # OpenAPI から Go/TS を再生成
make lint    # biome + golangci-lint
make test    # tsc --noEmit + go test
make e2e     # Playwright
```

## 規約

- フロント: Biome の設定に従う（整形は `make fmt`）。ESLint/Prettier は使わない。
- 入力値検証は Zod、スキーマは `apps/web/src/lib/api/` 配下に置く。
- Go: `internal/server` は生成 IF の実装のみ。ドメインロジックは `internal/<domain>` に置く。
- コンポーネントを追加したら Storybook のストーリーも追加する。
- PR 前に `make lint test` を通す。E2E に影響する変更なら `make e2e` も実行する。
