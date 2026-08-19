# ドキュメント

| ドキュメント | 内容 |
| --- | --- |
| [システム概要書](system-overview.md) | 何のためのリポジトリか、全体構成、開発フロー、現状の制約 |
| [テスト方針](testing-policy.md) | どの層を何でテストするか、CI との対応、テストデータの方針 |
| [仕様書テンプレート](specs/_template.md) | 機能 / ドメイン単位の仕様書のひな形 |
| [テスト計画書テンプレート](test/plan-_template.md) | 案件・リリース単位のテスト計画書のひな形 |
| [テスト設計書テンプレート](test/design-_template.md) | 機能 / API 単位のテスト設計書のひな形 |
| [テストドキュメントの使い分け](test/README.md) | 仕様書・テスト計画書・テスト設計書の役割分担 |
| [バックエンドアーキテクチャ](architecture/backend.md) | domain / usecase / infrastructure / presentation の 4 層と依存性逆転 |
| [フロントエンドアーキテクチャ](architecture/frontend.md) | App Router + feature ベースの構成ルール |

セットアップ手順とコマンド一覧はリポジトリルートの [README](../README.md)、
エージェント / 開発者向けの約束事は [AGENTS.md](../AGENTS.md) を参照。

## 生成物

`docs/api/` は Redocly CLI の出力先（`task docs:build` / `task docs:bundle`）で、
コミットしない。CI の `docs` ジョブが artifact として出力する。
