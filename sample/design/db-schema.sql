-- SQLite用（課題がMySQL/PostgreSQLなら AUTOINCREMENT → AUTO_INCREMENT / SERIAL に変える）

CREATE TABLE IF NOT EXISTS products (
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    name     TEXT    NOT NULL,
    category TEXT    NOT NULL DEFAULT '',
    price    REAL    NOT NULL,
    stock    INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS customers (
    id    INTEGER PRIMARY KEY AUTOINCREMENT,
    name  TEXT NOT NULL,
    email TEXT,
    phone TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sales (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    customer_id  INTEGER REFERENCES customers(id),
    total_amount REAL    NOT NULL,
    sold_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sale_items (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    sale_id    INTEGER NOT NULL REFERENCES sales(id),
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity   INTEGER NOT NULL,
    unit_price REAL    NOT NULL
);

-- テストデータ（動作確認用）
INSERT INTO products (name, category, price, stock) VALUES
    ('4K テレビ 55型', 'テレビ', 89800, 10),
    ('冷蔵庫 500L', '冷蔵庫', 148000, 5),
    ('ロボット掃除機', '掃除機', 54800, 20),
    ('電子レンジ', '調理家電', 22800, 15),
    ('ノートPC', 'パソコン', 128000, 8);
