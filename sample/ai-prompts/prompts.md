# AIアシスタント活用ガイド（TRACK対応）

## 基本戦略（推奨フロー）

```
Step 1: 業務要件を渡してDBスキーマを出させる   ← ここを自分でやると時間を食う
Step 2: スキーマを確認・修正する（2分）
Step 3: スキーマを固定してAPIを実装させる
Step 4: 一機能ずつcurlで確認しながら進める
```

---

## Step 1: DBスキーマ設計プロンプト（最初にこれを送る）

```
家電量販店の販売管理システムを設計します。

## 業務要件
- 商品を登録・管理する（名前、カテゴリ、価格、在庫数）
- 商品を販売する（誰が・何を・何個・いくらで）
- 販売時に在庫を自動で減らす
- 売上を集計できる（件数・合計金額・カテゴリ別など）
- 顧客情報は任意（なくても販売できる）

上記を満たすSQLiteのCREATE TABLEを設計してください。
各テーブルのカラム名・型・制約も含めて、そのまま実行できるSQLで出力してください。
```

→ AIが出したスキーマを見て以下を確認する：
- `sale_items` テーブルはあるか（販売明細）
- `unit_price` カラムはあるか（販売時点の価格を記録）
- 在庫カラム（stock）はあるか
- 外部キー制約はあるか（なくてもOK）

---

## Step 2: スキーマ確定後の実装プロンプト

スキーマが確定したらこれを送る：

```
以下のスキーマを使って[Golang / Node.js(Express) / PHP]でREST APIを実装してください。

## 確定スキーマ
[Step 1で出てきたCREATE TABLE をそのまま貼る]

## 言語・環境
- 言語: [Golang / Node.js(Express) / PHP]
- DB: SQLite
- ポート: 8080
- 単一ファイルで実装すること
- DBの初期化（CREATE TABLE IF NOT EXISTS）もコードに含めること

## 実装するAPI
1. GET/POST /products
2. GET/PUT/DELETE /products/:id
3. POST /sales（在庫チェック・在庫減算・合計計算をトランザクションで）
4. GET /sales
5. GET /reports/summary（件数・合計金額）

## POST /sales の仕様
リクエスト: { "customer_id": null, "items": [{"product_id": 1, "quantity": 2}] }
- 販売時点の価格をunit_priceとして記録する
- 在庫不足は {"error": "insufficient stock"} を400で返す
- 成功したら201で販売情報（合計金額・明細）を返す
```

---

## 機能別プロンプト（部分的に追加したいとき）

### 商品一覧・登録を追加
```
上記のコードに以下を追加してください：
GET /products → productsテーブルの全件をJSON配列で返す
POST /products → name, category, price, stockを受け取りINSERT、作成したレコードを201で返す
```

### 販売登録（最重要）
```
POST /sales エンドポイントを実装してください。
リクエスト: { "customer_id": null, "items": [{"product_id": 1, "quantity": 2}] }

処理:
1. トランザクション開始
2. 各itemについてproductsから price, stock を取得
3. stock < quantity なら rollback して {"error": "insufficient stock"} を400で返す
4. total_amount = Σ(price × quantity)
5. sales にINSERT
6. sale_items にINSERT（unit_price = そのときのprice）
7. products の stock を quantity分減算
8. コミット → 201で作成した販売情報を返す
```

### レポート追加
```
GET /reports/summary エンドポイントを追加してください。
salesテーブルから件数(total_sales)と合計金額(total_revenue)を返す。
{"total_sales": 42, "total_revenue": 3850000}
```

### バグ修正依頼の書き方
```
以下のコードでエラーが発生しています。
エラー: [エラーメッセージをそのまま貼る]
コード: [該当箇所を貼る]
修正してください。
```

---

## AIを使うときの注意点

### やってはいけないこと
- ❌ 「全部実装して」→ 長すぎて品質が落ちる。機能ごとに頼む
- ❌ コードを確認せずに次の機能を依頼する → バグが積み重なる
- ❌ エラーを「直してください」だけで送る → エラーメッセージを必ず貼る

### 効率的な進め方
- ✅ まず動く骨格（DBセットアップ + ヘルスチェック）を作ってもらう
- ✅ `curl` でレスポンスを確認してから次へ
- ✅ 「〜の形式でレスポンスを返して」と出力形式を明示する
- ✅ トランザクションが必要な処理は明示的に「トランザクションで」と書く

---

## 自分で書く場合のチェックリスト

```
[ ] DBの初期化（CREATE TABLE IF NOT EXISTS）は最初にやる
[ ] JSONのデコードエラーをちゃんとハンドリング
[ ] 404 Not Found を返す
[ ] 販売登録はトランザクションを使う
[ ] 在庫不足チェック（stock < quantity）
[ ] in庫減算（UPDATE products SET stock = stock - quantity）
[ ] Content-Type: application/json ヘッダを設定
[ ] 正しいHTTPステータスコード（201, 204, 400, 404, 500）
```
