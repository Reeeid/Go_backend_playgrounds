# Node.js 実装

## 起動方法

```bash
npm install
npm start
# または開発用（ファイル変更で自動再起動）
npm run dev
```

## better-sqlite3 vs pg/mysql2

| | better-sqlite3 | pg / mysql2 |
|---|---|---|
| 同期/非同期 | 同期（async不要） | 非同期（async/await必要） |
| セットアップ | ゼロ | DBサーバ必要 |
| トランザクション | `db.transaction()` | `BEGIN` / `COMMIT` SQL |

## 課題がMySQL/PostgreSQLの場合

```js
// pg の場合
const { Pool } = require('pg');
const pool = new Pool({ connectionString: process.env.DATABASE_URL });

// クエリ（$1 プレースホルダ）
const result = await pool.query('SELECT * FROM products WHERE id = $1', [id]);
const row = result.rows[0];

// トランザクション
const client = await pool.connect();
try {
  await client.query('BEGIN');
  // ... 処理 ...
  await client.query('COMMIT');
} catch (e) {
  await client.query('ROLLBACK');
  throw e;
} finally {
  client.release();
}
```

## ポイント解説

- `db.transaction()`: better-sqlite3のトランザクション。例外が出ると自動ロールバック
- `Object.assign(new Error(...), { status: 400 })`: エラーにHTTPステータスを添付するパターン
- `COALESCE(SUM(...), 0)`: salesが0件のとき SUM は NULL になるので0に変換
