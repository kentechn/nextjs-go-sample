# Todo テスト計画書

## 1. 目的
Todo 機能の作成・一覧・取得・削除と、入力検証および主要画面フローを保証する。
完了操作、永続化、認証・複数ユーザー、ページングは保証しない。

## 2. 対象
- 対象機能: `openapi/openapi.yaml` の todos タグ、`apps/api`、`apps/web`、`e2e`
- 変更範囲: Todo の domain / usecase / infrastructure / presentation、Todo 画面、E2E
- 影響範囲: API の OpenAPI 検証、SSR 初期表示、Server Action による作成・削除
- 対象外: 完了 API が無いため Done への遷移、永続化テスト、認証・複数ユーザー、ページング

## 3. テストレベルとタイプ
| レベル | 目的 | ツール | 担当 | 判断基準 |
| --- | --- | --- | --- | --- |
| ユニット（domain / usecase） | Todo の業務ルールと手順の回帰 | go test + testify | 実装者 | BR-01〜BR-07 の正常系・異常系を下位で確認 |
| インフラ（memory） | 保存、重複、並び順、削除の回帰 | go test | 実装者 | BR-06、BR-07、BR-09 をリポジトリで確認 |
| 結合（presentation） | HTTP 契約、仕様検証、エラー形 | go test + httptest（ルータ全体） | 実装者 | operation ごとに代表ケースを確認。`getTodo` 成功と `status=all/done` は未カバー |
| フロントユニット | Zod の trim、境界値、status | Vitest | 実装者 | title の 0 / 1 / 200 / 201 rune と `archived` を確認 |
| コンポーネント | TodoForm / TodoList の状態と操作 | Storybook（play + addon-a11y） | 実装者 | Default、エラー、空、200 文字タイトルを確認 |
| E2E | SSR、作成、削除、主要な入力エラー、health | Playwright | 実装者 | golden path とユーザーに見える失敗系を確認 |
| 静的検査 | 仕様と生成物の同期 | `task gen:check` / `task docs:lint` | CI | コマンドと CI 対応は共通方針に従う |

完了 API が無いため Done 側の状態遷移は対象外とする。永続化も実装対象外のため、DB との整合性は検証しない。
コマンドと CI ジョブの対応は `docs/testing-policy.md` を参照し、ここでは重複して記載しない。

## 4. 環境
| 環境 | 用途 | 構成 | データ |
| --- | --- | --- | --- |
| ローカル | 開発・単体 | `task dev` | インメモリ |
| CI | 全ジョブ | GitHub Actions | 毎回クリーン |

## 5. テストデータ方針
- E2E のタイトルは `Date.now()` などで一意化し、インメモリの実行間状態と衝突させない。
- domain / usecase の時刻と ID は `WithClock` / `WithIDGenerator` で固定値を注入する。
- seed 2 件を前提にする E2E 以外は、共有状態に依存しない。

## 6. 開始条件 / 完了条件
- 開始: 仕様書と OpenAPI の対象範囲が確認済み、設計書のケースが実装先に割り当て済み。
- 完了: CI 6 ジョブすべて green / 設計書のケースが全て実行済み / 未解決の Critical・High が 0。

## 7. スケジュールと体制
| 作業 | 担当 | 期間 |
| --- | --- | --- |
| 設計書レビュー | 実装者・レビュー担当 | 案件開始時 |
| 下位テスト実行 | 実装者 | 実装中 |
| 結合・E2E・CI 確認 | 実装者 | PR 作成後 |

## 8. リスクと対策
| リスク | 影響 | 対策 |
| --- | --- | --- |
| インメモリで実行間の状態が残る | E2E が既存データに依存する | E2E のタイトルを `Date.now()` で一意化する |
| E2E の不安定化 | CI が信用されなくなる | 下位テストへ網羅を寄せ、共通方針の retry / trace 設定に従う |
| 完了 API 未実装 | done フィルタの実質的な振る舞いを検証できない | usecase / repository のフィルタを確認し、HTTP の未検証を設計書に記録する |

## 9. 不具合の扱い
- 重大度は Critical / High / Medium / Low とし、Critical・High は完了条件までに解消する。
- 「バグを直す前に再現テストを追加する」ルールに従う。

## 10. 成果物
- Todo 仕様書、Todo テスト設計書、CI レポート、テスト結果サマリ
