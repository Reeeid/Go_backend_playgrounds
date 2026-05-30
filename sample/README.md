# 家電量販店 販売管理システム - 180分攻略ガイド

## 時間配分（推奨）

| フェーズ | 時間 | やること |
|---|---|---|
| 設計 | 15分 | DB設計・API設計を紙に書く（迷わない基盤） |
| DB + サーバ起動 | 15分 | テーブル作成・動作確認 |
| 商品CRUD | 30分 | GET/POST/PUT/DELETE /products |
| 在庫更新付き販売登録 | 35分 | POST /sales + トランザクション |
| 販売履歴 | 15分 | GET /sales, GET /sales/:id |
| 売上レポート | 15分 | GET /reports/summary |
| バッファ・テスト | 15分 | curl/Postmanで動作確認 |

## 最初に決めること（15分で完了）

### 1. 課題文で「使えるDB」を確認してからドライバを決める

外部ライブラリのインストールは失敗すると時間を大きく溶かす。
最初にDB環境を確認してから必要最小限だけ入れる。

| DB | Go | Node.js | PHP |
|---|---|---|---|
| SQLite | `go-sqlite3`（CGO必要） | `better-sqlite3` | PDO標準組み込み |
| MySQL | `go-sql-driver/mysql` | `mysql2` | PDO標準組み込み |
| PostgreSQL | `lib/pq` | `pg` | PDO標準組み込み |

**PHPはPDOが標準組み込みなのでインストール作業ゼロ。** 環境が不安なときはPHPが一番リスクが低い。

### 2. ライブラリは最小限にする

ORM（GORM, Prisma, Eloquent）は開発速度が上がるが：
- セットアップ・スキーマファイル生成などに時間がかかる
- 知らないライブラリは確認コストが高い
- Prepared Statement（`?` プレースホルダ）を使えばSQLiは防げる

**知っているORMがあれば使う。知らないなら素のSQLで十分。**

### 3. 実装スコープを決める

商品管理 → 販売登録 → レポート の順で確実に進める。

## 実装優先度

```
Must（必ず実装）
  ├── 商品一覧 GET /products
  ├── 商品登録 POST /products
  ├── 在庫確認付き販売登録 POST /sales
  └── 売上合計 GET /reports/summary

Should（時間があれば）
  ├── 商品更新 PUT /products/:id
  ├── 商品削除 DELETE /products/:id
  └── 販売詳細 GET /sales/:id

Nice to have（余裕があれば）
  ├── 顧客管理
  ├── カテゴリ別集計
  └── 期間指定レポート
```

## AIに書かせる場合のコツ

→ `ai-prompts/prompts.md` を参照

## 言語別最小実装

- Golang: `golang/main.go`（単一ファイル、sqlite3使用）
- Node.js: `nodejs/app.js`（Express + better-sqlite3）
- PHP: `php/index.php`（PDO sqlite、ビルトインサーバ）

## よくある罠

- トランザクションを忘れる → 在庫が不整合になる
- NULL許容カラムのスキャンでパニック → `sql.NullInt64` を使う
- 在庫不足チェックを忘れる → マイナス在庫になる
- エラー時にHTTPステータスを200で返す → NG、4xx/5xxを使う
