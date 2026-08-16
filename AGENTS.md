# AGENTS.md

このリポジトリで作業するエージェント / 開発者向けの約束事。

## 原則

- `openapi/openapi.yaml` が HTTP 契約の唯一の真実。API を変えるときは必ず仕様から。
- 生成物 (`apps/api/internal/openapi/api.gen.go`, `apps/web/src/shared/api/schema.gen.ts`) は手で編集しない。
  仕様を直して `task gen` を実行する。
- 型は自前定義せず生成型（`components["schemas"][...]` / `openapi.Todo`）を使う。

## コマンド

タスクランナーは [Task](https://taskfile.dev)（`Taskfile.yml` + `taskfiles/app.yml` / `taskfiles/docs.yml`）。

```bash
task setup   # 依存インストール
task gen     # OpenAPI から Go/TS を再生成
task lint    # biome + golangci-lint + redocly lint
task test    # tsc --noEmit + go test
task e2e     # Playwright
task --list-all  # 全タスク
```

## 規約

- フロント: Biome の設定に従う（整形は `task fmt`）。ESLint/Prettier は使わない。
- 入力値検証は Zod、スキーマは feature 配下 (`apps/web/src/features/<feature>/schema.ts`) に置く。
- Go: `internal/{domain,usecase,infrastructure,presentation}` の 4 層構成。詳細は `docs/architecture/backend.md`。
- フロント: 機能単位で `apps/web/src/features/<feature>/`、共通は `apps/web/src/shared/`。詳細は `docs/architecture/frontend.md`。
- コンポーネントを追加したら Storybook のストーリーも追加する。
- PR 前に `task lint test` を通す。E2E に影響する変更なら `task e2e` も実行する。
