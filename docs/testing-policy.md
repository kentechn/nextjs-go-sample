# テスト方針

## 1. 目的

- 仕様（`openapi/openapi.yaml`）と実装のズレを機械的に検出する。
- 業務ルールの回帰を、速く実行できる下位のテストで捕まえる。
- E2E は「主要フローが繋がっていること」の確認に限定し、網羅はしない。

## 2. 何をどこでテストするか

| 対象 | 種類 | ツール | 置き場所 |
| --- | --- | --- | --- |
| domain（エンティティ / 不変条件 / ドメインエラー） | ユニット | `go test` + testify | `internal/domain/<name>/*_test.go` |
| usecase（手順 / 異常系） | ユニット（リポジトリはフェイク） | `go test` | `internal/usecase/<name>/*_test.go` |
| infrastructure（リポジトリ実装） | ユニット / 結合 | `go test` | `internal/infrastructure/<driver>/*_test.go` |
| presentation（HTTP・仕様検証・エラー形） | 結合 | `go test` + `httptest`（ルータ全体） | `internal/presentation/http/*_test.go` |
| Zod スキーマ / feature のロジック | ユニット | Vitest | `apps/web/src/**/*.test.ts` |
| UI コンポーネント | インタラクション / a11y | Storybook（`play` + addon-a11y） | `*.stories.tsx` |
| 主要ユーザーフロー | E2E | Playwright | `e2e/tests/*.spec.ts` |
| 仕様と生成物の同期 | 静的検査 | `task gen:check` / `task docs:lint` | CI |

比率の目安は「domain / usecase のユニット > presentation の結合 > E2E」。

## 3. 方針

### 共通

- **テストするのは振る舞い**。内部実装（呼び出し回数、privateな構造）に依存させない。
- **生成コードはテストしない**（`api.gen.go` / `schema.gen.ts`）。生成の正しさは `task gen:check` で担保する。
- **仕様で表現できる検証を二重に書かない。** 必須項目・型・`minLength` / `maxLength` / `format: uuid` の
  拒否は presentation の結合テストで代表 1〜2 ケースだけ確認し、値の組み合わせは domain のテストで持つ。
- テストは並列で動かす（Go は `t.Parallel()` + `go test -race`）。共有状態を持ち込まない。

### Go

- **usecase のテストで DB / HTTP を触らない。** domain のリポジトリインターフェースに対する
  フェイク（インメモリ実装かテスト内の struct）を注入する。
  依存性逆転はこのために入れている（[バックエンドアーキテクチャ](architecture/backend.md)）。
- **時刻と ID を固定する。** usecase は `now func() time.Time` / `newID func() uuid.UUID` を
  受け取るので、テストでは固定値を渡し、`time.Now()` / `uuid.New()` をテスト内で呼ばない。
- **presentation はルータ全体（`NewRouter`）を通して叩く。** ハンドラを直接呼ぶと
  OpenAPI 検証ミドルウェアやエラーハンドラを通らず、実際の応答と乖離するため。
- エラーの検証は `require.ErrorIs`。メッセージ文字列に依存しない。
- 表明は testify（`require` = 失敗したら中断、`assert` = 続行）。

### フロント

- **Vitest**: Zod スキーマの境界値（空文字 / 前後空白 / 200 文字 / 201 文字）や、
  純粋な変換ロジック。React のレンダリングを伴うものは Storybook 側に寄せる。
- **Storybook**: コンポーネントの状態を story として列挙し（Default / Empty / エラー / 極端な入力）、
  ユーザー操作を伴うものは `play` で検証する。`task app:storybook:build` が CI で全 story を実行する。
- **`api.ts` のテストは書かない**（薄いラッパであり、型と E2E で足りる）。
  ロジックが増えたら Vitest で `apiClient` をモックする。
- Client Component は action を props で受け取る設計にして、story から差し替えられる状態を保つ。

### E2E

- 対象は **golden path + ユーザーに見える主要な失敗系**。現状は
  SSR の初期表示 / 作成 / 空タイトルの検証エラー / 削除 / `/health`。
- 網羅・境界値は下位のテストに任せ、ここには足さない（遅く壊れやすいため）。
- 要素の指定は role / label を優先し、必要な場合のみ `data-testid`。CSS クラスに依存しない。
- テストデータは `Date.now()` などで一意にし、他のテストと衝突させない
  （API はインメモリで、実行間の状態がリセットされない前提で書く）。
- Playwright が api / web（本番ビルド）を自動起動する。CI では retry 2 回、trace は初回リトライ時のみ。

## 4. コマンド

```bash
task test              # フロントの型チェック + Vitest + Go テスト
task app:api:test      # go test -race ./...
task app:web:test      # Vitest
task app:web:typecheck # tsc --noEmit
task app:storybook:build  # 全 story のビルド + play の実行
task e2e               # Playwright（api / web を自動起動）
task lint              # biome + golangci-lint + redocly lint
task gen:check         # 生成物が仕様と同期しているか
```

## 5. CI との対応

| ジョブ | 実行内容 |
| --- | --- |
| frontend | `app:web:lint` / `app:web:typecheck` / `app:web:test` / `app:web:build` / `app:storybook:build` |
| backend | `app:api:build` / `app:api:test` / golangci-lint |
| docs | `docs:lint` / `docs:build` / `docs:bundle`（生成物は artifact） |
| codegen | `gen:check` |
| docker | api / web イメージのビルド |
| e2e | `app:e2e:install` / `app:web:build` / `app:e2e`（レポートは artifact） |

新しい種類のテストを追加したら、必ずどれかのジョブに載せる（ローカル専用のテストは腐る）。

## 6. カバレッジ

数値目標は置かない。代わりに次を必須とする。

- domain / usecase に追加したロジックには、必ず正常系 1 + 異常系 1 以上。
- OpenAPI に operation を追加したら、presentation の結合テストを 1 本以上。
- コンポーネントを追加したら story を 1 本以上。
- バグを直したら、そのバグを再現するテストを先に追加する。

## 7. 失敗したときの調べ方

| 失敗 | 見るところ |
| --- | --- |
| `codegen` ジョブ | `task gen` を実行してコミットし忘れていないか |
| presentation の 400 が想定外 | OpenAPI の制約と、検証ミドルウェアが返す `invalid_request` の内容 |
| E2E | `e2e/playwright-report`（CI では artifact）の trace / スクリーンショット |
| Storybook の `play` | `task app:storybook` でブラウザから該当 story を開く |
| Docker ビルド | `task app:docker:build` をローカルで再実行（buildx キャッシュは `.buildx-cache/`） |
