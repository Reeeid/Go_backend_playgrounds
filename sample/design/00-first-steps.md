# 設計の進め方（最初の15分でやること）

## Step 1: エンティティを洗い出す（3分）

家電量販店の販売管理 → 何が「モノ」として存在するか？

```
商品 (products)       ← 何を売るか
販売 (sales)          ← いつ・いくら売ったか
販売明細 (sale_items) ← 販売の中身（何を何個）
顧客 (customers)      ← 誰に売ったか（後回しOK）
```

## Step 2: 各エンティティのカラムを決める（5分）

最低限のカラムだけ書く。あとで追加できる。

```
products: id, name, category, price, stock
sales: id, customer_id(NULL可), total_amount, sold_at
sale_items: id, sale_id, product_id, quantity, unit_price
```

## Step 3: APIエンドポイントを列挙する（5分）

```
GET    /products       一覧
POST   /products       登録
GET    /products/:id   詳細
PUT    /products/:id   更新
DELETE /products/:id   削除

POST   /sales          販売登録（在庫チェック・在庫減算・合計計算）
GET    /sales          販売履歴一覧
GET    /sales/:id      販売詳細

GET    /reports/summary  売上サマリ
```

## Step 4: 販売登録のロジックを整理する（2分）

ここが一番複雑なので事前に頭の中で整理しておく。

```
POST /sales リクエスト例:
{
  "customer_id": 1,  // null可
  "items": [
    { "product_id": 1, "quantity": 2 },
    { "product_id": 3, "quantity": 1 }
  ]
}

処理の流れ:
1. トランザクション開始
2. 各itemについて:
   - productsテーブルから価格・在庫を取得
   - 在庫 < quantity なら rollback → 400
   - unit_price = product.price で計算
3. total_amount = Σ(unit_price × quantity)
4. sales INSERT
5. sale_items INSERT（各item）
6. products UPDATE stock = stock - quantity（各item）
7. コミット
```

## よく出る追加要件（余裕があれば対応）

- 売上の日付絞り込み: WHERE sold_at BETWEEN ? AND ?
- カテゴリ別集計: GROUP BY category
- 商品の売上ランキング: JOIN + GROUP BY + ORDER BY
