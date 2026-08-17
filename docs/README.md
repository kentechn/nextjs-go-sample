# ドキュメント

| ドキュメント | 内容 |
| --- | --- |
| [システム概要書](system-overview.md) | 何のためのリポジトリか、全体構成、開発フロー、現状の制約 |
| [テスト方針](testing-policy.md) | どの層を何でテストするか、CI との対応、テストデータの方針 |
| [バックエンドアーキテクチャ](architecture/backend.md) | domain / usecase / infrastructure / presentation の 4 層と依存性逆転 |
| [フロントエンドアーキテクチャ](architecture/frontend.md) | App Router + feature ベースの構成ルール |

セットアップ手順とコマンド一覧はリポジトリルートの [README](../README.md)、
エージェント / 開発者向けの約束事は [AGENTS.md](../AGENTS.md) を参照。

## 移行状況

アーキテクチャの 2 ドキュメントは**目標構成**を記述しており、サンプル実装（Todo）の
コードは次の PR で追随させる。それまでは以下が実際の状態。

| ドキュメント上 | 現在のコード |
| --- | --- |
| `internal/{domain,usecase,infrastructure,presentation}` | `internal/server`（HTTP + 変換）と `internal/todo`（エンティティ + 永続化） |
| `apps/web/src/features/todos/` | `apps/web/src/components/todo/` + `src/app/actions.ts` + `src/lib/api/todo(s).ts` |
| `apps/web/src/shared/api/` | `apps/web/src/lib/api/` |
| Vitest（`task app:web:test`） | 未導入 |

追随が完了したらこのセクションを削除する。

## 生成物

`docs/api/` は Redocly CLI の出力先（`task docs:build` / `task docs:bundle`）で、
コミットしない。CI の `docs` ジョブが artifact として出力する。
