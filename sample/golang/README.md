# Golang 実装

## 起動方法

```bash
go mod tidy
go run main.go
```

## 注意点

- `go-sqlite3` はCGOが必要。GCC/MinGWが必要（WindowsはTDM-GCCなど）
- 課題のDBがMySQL/PostgreSQLなら以下を変える:
  - import: `_ "github.com/go-sql-driver/mysql"` または `_ "github.com/lib/pq"`
  - `sql.Open("sqlite3", "./store.db")` → `sql.Open("mysql", "user:pass@/dbname")`
  - `INTEGER PRIMARY KEY AUTOINCREMENT` → MySQL: `AUTO_INCREMENT`, PG: `SERIAL`
  - `REAL` → `DECIMAL(10,2)`

## ポイント解説

- `sql.NullInt64`: NULL許容のINTEGERカラムをスキャンするときに必要
- トランザクション: `db.Begin()` → 処理 → `tx.Rollback()` or `tx.Commit()`
- ルーティング: 標準の `net/http` だけで `/products` と `/products/:id` を分岐
