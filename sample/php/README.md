# PHP 実装

## 起動方法

```bash
php -S localhost:8080 index.php
```

## 課題がMySQL/PostgreSQLの場合

```php
// MySQL
$db = new PDO('mysql:host=localhost;dbname=store;charset=utf8mb4', 'user', 'pass');

// PostgreSQL
$db = new PDO('pgsql:host=localhost;dbname=store', 'user', 'pass');

// AUTO_INCREMENT → MySQL は同じ、PostgreSQL は SERIAL または GENERATED ALWAYS AS IDENTITY
```

## ポイント解説

- `PDO::ATTR_DEFAULT_FETCH_MODE = PDO::FETCH_ASSOC`: 全クエリでキー付き配列が返る
- `PDO::ATTR_ERRMODE = PDO::ERRMODE_EXCEPTION`: エラーを例外として受け取る
- Prepared statement: `$db->prepare(...)->execute([...])` → SQLインジェクション対策済み
- `$db->lastInsertId()`: 直前のINSERTで生成されたIDを取得
- `respond(null, 204)`: 204はボディなし。json_encode(null) は "null" になるが問題ない
