# Todo 仕様書

## 1. 概要
- Todo の作成・一覧・取得・削除を提供し、やることを管理する。
- API は OpenAPI、画面は Next.js の SSR と Server Action で提供する。

## 2. スコープ
- やること: Todo の作成、一覧、取得、削除、タイトル入力の検証。
- やらないこと: 完了操作、永続化、認証・複数ユーザー、一覧のページング。

## 3. 用語・ドメインモデル
| 用語 | 定義 | 対応する実装 |
| --- | --- | --- |
| Todo | `Todo{ID uuid.UUID, Title string, Done bool, CreatedAt time.Time}` | `apps/api/internal/domain/todo` |
| TitleMaxLength | タイトルの最大長。rune 数で 200 | `apps/api/internal/domain/todo` |

- `Title` は前後の空白をトリムし、空文字・空白のみを許可しない。
- `Done` は作成時に `false` とする。`ID` はサーバ生成 UUID、`CreatedAt` はサーバ時刻の UTC とする。
- usecase は `now func() time.Time` / `newID func() uuid.UUID` を受け取り、`WithClock` / `WithIDGenerator` でテスト時に固定できる。

## 4. 状態遷移
```mermaid
stateDiagram-v2
    [*] --> Open: 作成
    Open --> [*]: 削除
```

Done への遷移手段は未実装であり、完了に変更する API は存在しない。

## 5. 業務ルール
| ID | ルール | 違反時の振る舞い |
| --- | --- | --- |
| BR-01 | `title` は必須。trim 後に空（空文字 / 空白のみ）を許可しない | `ErrEmptyTitle` → HTTP 400 `invalid_argument` |
| BR-02 | 保存される `title` は前後の空白をトリムした値とする | トリム後の値を保存する |
| BR-03 | `title` の長さは rune 数で 200 まで。201 rune は許可しない | `ErrTitleTooLong` → HTTP 400 `invalid_argument` |
| BR-04 | `id` はサーバ生成 UUID、`createdAt` はサーバ時刻の UTC とする | usecase の注入点で生成する |
| BR-05 | `done` は作成時 `false`。完了に変更する API は存在しない | `GET /todos?status=done` は常に空配列 |
| BR-06 | 一覧は作成が古い順。`status` は `all`（default）、`open`（`done=false`）、`done`（`done=true`） | presentation の `doneFilter` で変換する |
| BR-07 | 存在しない id への `getTodo` / `deleteTodo` は対象なしとする | HTTP 404 `not_found` |
| BR-08 | UUID 形式でない path パラメータや仕様違反リクエストを受け付けない | OpenAPI 検証ミドルウェアが HTTP 400 `invalid_request` を返す |
| BR-09 | 同一 id の重複作成を許可しない | `ErrAlreadyExists` → HTTP 500 `internal` |

## 6. インターフェース
### 6.1 API
OpenAPI の operationId を列挙するだけにする（定義は複製しない）。
| operationId | 概要 | 関連ルール |
| --- | --- | --- |
| `getHealth` | ヘルスチェック | - |
| `listTodos` | Todo 一覧取得（`status` query） | BR-05, BR-06 |
| `createTodo` | Todo 作成 | BR-01, BR-02, BR-03, BR-04, BR-09 |
| `getTodo` | Todo 取得 | BR-07, BR-08 |
| `deleteTodo` | Todo 削除 | BR-07, BR-08 |

### 6.2 画面 / 操作
- `/` のみを提供し、`export const dynamic = "force-dynamic"` で SSR する。
- 見出しは `Todos` とする。
- `TodoForm` は入力の aria-label を「やること」、送信ボタンを「追加」とする。エラーは `role="alert"` と `data-testid="todo-form-error"` で表示する。
- `TodoList` は各行を `data-testid="todo-item"` で識別し、各行に「削除」ボタンを表示する。Todo が無い場合は空状態を表示する。
- 更新は Server Action `createTodoAction` / `deleteTodoAction` で行い、成功後に `revalidatePath("/")` を実行する。
- フロント側の入力検証は Zod の `createTodoSchema` で行う。メッセージは「タイトルを入力してください」/「200文字以内で入力してください」とする。
- `deleteTodoAction` は `todoId` が文字列でなければ何もしない。

## 7. エラー仕様
| コード | HTTP | 発生条件 | ユーザーへの表示 |
| --- | --- | --- | --- |
| `invalid_request` | 400 | 検証ミドルウェアが UUID 形式やリクエスト仕様の違反を検出 | 入力またはリクエストのエラー |
| `invalid_argument` | 400 | ドメイン規則（空タイトル、タイトル長超過）に違反 | 入力欄下にメッセージ |
| `not_found` | 404 | 指定した Todo が存在しない | 対象が存在しないことを表示 |
| `internal` | 500 | 同一 id の重複作成など、内部エラーが発生 | 共通エラー |

## 8. 非機能要件
- データは `internal/infrastructure/memory` のインメモリ実装で保持し、プロセス再起動で消失する。
- 起動時に `"read the OpenAPI spec"` / `"run task dev"` の 2 件を seed する。
- ブラウザは Go API を直接呼ばず、必ずサーバ経由とする。`API_BASE_URL` はサーバ専用で、`NEXT_PUBLIC_` を付けない。
- CORS はブラウザから直接叩いてデバッグする場合の保険とし、既定値は `http://localhost:3000` とする。
- API の read/write timeout は 15s、graceful shutdown は 10s とする。
- ログは slog の JSON とする。

## 9. 前提と未決事項
| # | 論点 | 暫定 | 決定期限 |
| --- | --- | --- | --- |
| 1 | 完了操作（`PATCH /todos/{todoId}` 等） | 未実装 | 未定 |
| 2 | 永続化（インメモリ → DB） | 未実装 | 未定 |
| 3 | 認証・複数ユーザー | 未対応 | 未定 |
| 4 | 一覧のページング | 未対応 | 未定 |

## 10. 変更履歴
| 日付 | 変更 | 対応 PR |
| --- | --- | --- |
| 2026-08-19 | 初版作成（既存実装からの追記） | PR #4 |
