package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// DB は store が依存する抽象。*sql.DB と *sql.Tx が両方満たす。
type DB interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Begin() (*sql.Tx, error)
}

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	statements := []string{
		`CREATE TABLE IF NOT EXISTS products (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			name       TEXT    NOT NULL,
			category   TEXT    NOT NULL DEFAULT '',
			price      REAL    NOT NULL,
			stock      INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sales (
			id           INTEGER PRIMARY KEY AUTOINCREMENT,
			customer_id  INTEGER,
			total_amount REAL    NOT NULL,
			sold_at      DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sale_items (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			sale_id    INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity   INTEGER NOT NULL,
			unit_price REAL    NOT NULL
		)`,
	}
	for _, s := range statements {
		if _, err := db.Exec(s); err != nil {
			return nil, err
		}
	}

	return db, nil
}
