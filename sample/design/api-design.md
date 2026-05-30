# API設計書

## 商品管理

### GET /products
```json
// Response 200
[
  { "id": 1, "name": "4K テレビ 55型", "category": "テレビ", "price": 89800, "stock": 10 }
]
```

### POST /products
```json
// Request
{ "name": "エアコン", "category": "空調", "price": 68000, "stock": 5 }

// Response 201
{ "id": 6, "name": "エアコン", "category": "空調", "price": 68000, "stock": 5 }
```

### GET /products/:id
```json
// Response 200
{ "id": 1, "name": "4K テレビ 55型", "category": "テレビ", "price": 89800, "stock": 10 }
// Response 404
{ "error": "not found" }
```

### PUT /products/:id
```json
// Request（変更したいフィールドを送る）
{ "name": "4K テレビ 55型", "category": "テレビ", "price": 85000, "stock": 8 }
// Response 200
{ "id": 1, "name": "4K テレビ 55型", "category": "テレビ", "price": 85000, "stock": 8 }
```

### DELETE /products/:id
```
// Response 204 No Content
```

---

## 販売管理

### POST /sales（最重要エンドポイント）
```json
// Request
{
  "customer_id": null,
  "items": [
    { "product_id": 1, "quantity": 1 },
    { "product_id": 3, "quantity": 2 }
  ]
}

// Response 201
{
  "id": 1,
  "customer_id": null,
  "total_amount": 199400,
  "items": [
    { "product_id": 1, "quantity": 1, "unit_price": 89800 },
    { "product_id": 3, "quantity": 2, "unit_price": 54800 }
  ]
}

// Response 400（在庫不足）
{ "error": "insufficient stock for product 3" }
```

### GET /sales
```json
// Response 200
[
  { "id": 1, "customer_id": null, "total_amount": 199400, "sold_at": "2024-01-15T10:30:00Z" }
]
```

### GET /sales/:id
```json
// Response 200
{
  "id": 1,
  "total_amount": 199400,
  "sold_at": "2024-01-15T10:30:00Z",
  "items": [
    { "product_id": 1, "product_name": "4K テレビ 55型", "quantity": 1, "unit_price": 89800 }
  ]
}
```

---

## レポート

### GET /reports/summary
```json
// Response 200
{
  "total_sales": 42,
  "total_revenue": 3850000
}
```

### GET /reports/by-category（余裕があれば）
```json
// Response 200
[
  { "category": "テレビ", "count": 15, "revenue": 1350000 },
  { "category": "冷蔵庫", "count": 8, "revenue": 1184000 }
]
```

---

## curlテスト例

```bash
# 商品一覧
curl http://localhost:8080/products

# 商品登録
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"エアコン","category":"空調","price":68000,"stock":5}'

# 販売登録
curl -X POST http://localhost:8080/sales \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":1,"quantity":1},{"product_id":3,"quantity":2}]}'

# 売上サマリ
curl http://localhost:8080/reports/summary
```
