# Todo テスト設計書

## 1. 対象と参照
- 仕様書: `docs/specs/todo.md`
- OpenAPI: `getHealth` / `listTodos` / `createTodo` / `getTodo` / `deleteTodo`
- テスト方針: `docs/testing-policy.md`

## 2. テスト観点一覧
| 観点 ID | 観点 | 技法 | 実施レベル |
| --- | --- | --- | --- |
| V-01 | title の trim と空値 | 同値分割 | domain / usecase / Vitest |
| V-02 | title の rune 数境界 | 境界値分析 | domain / Vitest |
| V-03 | Todo の ID・時刻・初期状態 | 同値分割 | usecase |
| V-04 | status の同値分割 | 同値分割 | usecase / infrastructure / Vitest |
| V-05 | 一覧の作成順と削除 | 順序テスト | infrastructure / presentation / E2E |
| V-06 | UUID とリクエスト仕様の検証 | 同値分割 | presentation / Vitest |
| V-07 | 存在しない Todo とリポジトリエラー | エラー推測 | usecase / infrastructure / presentation |
| V-08 | 画面の状態と Server Action | 状態網羅 | Storybook / E2E |
| V-09 | API の主要フロー | シナリオ | presentation / E2E |

## 3. 境界値・同値分割
| 項目 | 有効同値 | 無効同値 | 境界値 |
| --- | --- | --- | --- |
| title の rune 数 | 1〜200 rune、マルチバイト 200 rune | 0 rune、201 rune | 0, 1, 200, 201 rune |
| title の空白 | 前後空白を含む値は trim 後に非空 | 空文字、空白のみ | `""`, `"   "`, `"  title  "` |

| status | 有効同値 | 無効同値 |
| --- | --- | --- |
| `status` | 未指定 / `all`（絞り込みなし）、`open`（`done=false`）、`done`（`done=true`） | `archived` → 400 `invalid_request` |

## 4. 状態遷移表（必要な機能のみ）
| 現状態 \ 操作 | 作成 | 削除 | 完了 |
| --- | --- | --- | --- |
| 未作成 | Open | 不可(404) | 不可（API なし） |
| Open | 不可（同一 id は通常サーバ生成） | 終了 | 不可（API なし） |

Done への遷移手段は未実装であり、Done 状態を作成する API は存在しない。

## 5. テストケース
| ケース ID | 観点 | 前提 | 操作 | 期待結果 | レベル | 実装先 |
| --- | --- | --- | --- | --- | --- | --- |
| TC-001 | V-01 | Todo を生成できる | 前後に空白を含む title で生成 | 保存値は trim 後の値 | unit | `apps/api/internal/domain/todo/todo_test.go::TestNewTrimsTitle` |
| TC-002 | V-01 | Todo を生成できる | title を空文字で生成 | `ErrEmptyTitle` | unit | `apps/api/internal/domain/todo/todo_test.go::TestNewRejectsInvalidTitle`（empty） |
| TC-003 | V-01 | Todo を生成できる | title を空白のみで生成 | `ErrEmptyTitle` | unit | `apps/api/internal/domain/todo/todo_test.go::TestNewRejectsInvalidTitle`（whitespace only） |
| TC-004 | V-02 | Todo を生成できる | title を 1 rune で生成 | 生成に成功する | unit | 未カバー |
| TC-005 | V-02 | Todo を生成できる | title を 200 rune のマルチバイト文字で生成 | 生成に成功する | unit | `apps/api/internal/domain/todo/todo_test.go::TestNewAcceptsMaxLengthTitle` |
| TC-006 | V-02 | Todo を生成できる | title を 201 rune で生成 | `ErrTitleTooLong` | unit | `apps/api/internal/domain/todo/todo_test.go::TestNewRejectsInvalidTitle`（one rune too long） |
| TC-007 | V-03 | usecase に生成関数を注入する | 固定 clock / ID で Todo を作成 | `createdAt` と `id` が固定値、`done=false` | unit | `apps/api/internal/usecase/todo/usecase_test.go::TestCreateUsesInjectedClockAndID` |
| TC-008 | V-01,V-02 | 不正 title を usecase に渡す | 空 title / 長すぎる title を作成 | ドメインエラーが伝播する | unit | `apps/api/internal/usecase/todo/usecase_test.go::TestCreateRejectsInvalidTitle` |
| TC-009 | V-04 | 複数の Todo が保存済み | `all` / `open` / `done` で一覧する | done 条件に一致する Todo だけ返る | unit | `apps/api/internal/usecase/todo/usecase_test.go::TestListFiltersByDone` |
| TC-010 | V-07 | 対象 id が未登録 | get / delete を実行する | `not_found` 相当のエラー | unit | `apps/api/internal/usecase/todo/usecase_test.go::TestGetAndDeleteMissingTodo` |
| TC-011 | V-07 | repository がエラーを返す | usecase の操作を実行する | repository エラーが伝播する | unit | `apps/api/internal/usecase/todo/usecase_test.go::TestRepositoryErrorIsPropagated` |
| TC-012 | V-05 | memory repository が空 | Todo を作成して id で取得する | 作成した Todo が取得できる | unit | `apps/api/internal/infrastructure/memory/todo_repository_test.go::TestCreateAndFindByID` |
| TC-013 | V-07 | 同一 id の Todo が保存済み | 同一 id で再作成する | `ErrAlreadyExists` | unit | `apps/api/internal/infrastructure/memory/todo_repository_test.go::TestCreateRejectsDuplicateID` |
| TC-014 | V-04,V-05 | 複数の Todo が保存済み | status で絞り込み、一覧する | 作成が古い順で返る | unit | `apps/api/internal/infrastructure/memory/todo_repository_test.go::TestListIsSortedAndFiltered` |
| TC-015 | V-05 | Todo が保存済み | id を指定して削除する | Todo が削除される | unit | `apps/api/internal/infrastructure/memory/todo_repository_test.go::TestDelete` |
| TC-016 | V-09 | router が起動している | health endpoint を GET する | 正常応答を返す | integration | `apps/api/internal/presentation/rest/router_test.go::TestGetHealth` |
| TC-017 | V-05,V-09 | router が起動している | create → list を実行する | 作成した Todo が一覧に現れる | integration | `apps/api/internal/presentation/rest/router_test.go::TestCreateThenListTodo` |
| TC-018 | V-06 | router が起動している | 不正な body で create する | HTTP 400（検証ミドルウェアが拒否。テストは status のみ検証し `invalid_request` の code は未検証） | integration | `apps/api/internal/presentation/rest/router_test.go::TestCreateTodoRejectsInvalidBody` |
| TC-019 | V-06 | router が起動している | UUID 形式でない id を get する | ハンドラに到達せず HTTP 400 `invalid_request` | integration | `apps/api/internal/presentation/rest/router_test.go::TestGetTodoRejectsMalformedID` |
| TC-020 | V-07 | 対象 id が未登録 | id を get する | HTTP 404 `not_found` | integration | `apps/api/internal/presentation/rest/router_test.go::TestGetTodoNotFound` |
| TC-021 | V-05,V-07 | Todo を作成済み | create → delete を実行する | 削除でき、対象は取得できない | integration | `apps/api/internal/presentation/rest/router_test.go::TestCreateThenDeleteTodo` |
| TC-022 | V-01,V-02 | createTodoSchema を使用する | trim、空 / 空白のみ、200 / 201、title 欠落を検証する | Zod の検証結果と指定メッセージが一致する | unit | `apps/web/src/features/todos/schema.test.ts` |
| TC-023 | V-04,V-06 | todoStatusSchema を使用する | `open` と `archived` を検証する | `open` は有効、`archived` は無効 | unit | `apps/web/src/features/todos/schema.test.ts::accepts the statuses defined in the spec` |
| TC-024 | V-08 | TodoForm を表示する | Default と WithError で空送信する | 入力とエラー表示が確認できる | component | `apps/web/src/features/todos/components/TodoForm.stories.tsx`（Default / WithError） |
| TC-025 | V-08 | TodoList を表示する | Default / Empty / LongTitle を表示する | 一覧、空状態、200 文字タイトルを確認できる | component | `apps/web/src/features/todos/components/TodoList.stories.tsx`（Default / Empty / LongTitle） |
| TC-026 | V-08,V-09 | API と web が起動している | SSR の初期表示を開く | seed Todo と `Todos` が表示される | e2e | `e2e/tests/todos.spec.ts::SSR page renders todos seeded by the API` |
| TC-027 | V-08,V-09 | `/` を開いている | title を入力して追加する | Todo が画面に追加される | e2e | `e2e/tests/todos.spec.ts::creates a todo through the server action` |
| TC-028 | V-01,V-08 | `/` を開いている | title 空欄で追加する | `タイトルを入力してください` が表示される | e2e | `e2e/tests/todos.spec.ts::rejects an empty title with a validation message` |
| TC-029 | V-08,V-09 | Todo を追加済み | 「削除」を押す | 対象の `todo-item` が消える | e2e | `e2e/tests/todos.spec.ts::deletes a todo` |
| TC-030 | V-09 | API が起動している | `/health` を GET する | health endpoint が正常応答を返す | e2e | `e2e/tests/todos.spec.ts::api health endpoint responds` |

期待結果は振る舞いで書く（呼び出し回数や内部構造に触れない）。

## 6. トレーサビリティ
| 業務ルール | 観点 | ケース | 未カバー / 備考 |
| --- | --- | --- | --- |
| BR-01 | V-01 | TC-002, TC-003, TC-008, TC-018, TC-022, TC-028 | - |
| BR-02 | V-01 | TC-001, TC-022 | - |
| BR-03 | V-02 | TC-005, TC-006, TC-008, TC-018, TC-022 | - |
| BR-04 | V-03 | TC-007, TC-026 | - |
| BR-05 | V-04 | TC-009, TC-023 | 完了 API は未実装。HTTP の `status=done` は未カバー |
| BR-06 | V-04,V-05 | TC-009, TC-014 | HTTP の `status=all` / `status=done` は未カバー |
| BR-07 | V-07 | TC-010, TC-020, TC-021 | `getTodo` の 200 成功応答は未カバー |
| BR-08 | V-06 | TC-018, TC-019, TC-023 | - |
| BR-09 | V-07 | TC-013 | HTTP 500 `internal` の形は未カバー |

未カバーとして次を記録する。

- `getTodo` 200 の結合テストが無い（作成→一覧、作成→削除→404 はある）。
- title の 1 rune 下限境界が未検証。
- `status=done` / `status=all` の HTTP レベル検証が無い。
- BR-09 は infra のユニットのみで、HTTP 500 の形は未検証。
- `deleteTodoAction` の `todoId` 非文字列時の no-op が未検証。

## 7. 実施結果（実行時に埋める）
| ケース ID | 結果 | 実行日 | 備考 / Issue |
| --- | --- | --- | --- |
