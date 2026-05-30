const express = require('express');
const Database = require('better-sqlite3');

const app = express();
const db = new Database('store.db');

app.use(express.json());

// ── DB Setup ─────────────────────────────────────────────

db.exec(`
  CREATE TABLE IF NOT EXISTS products (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT    NOT NULL,
    category   TEXT    NOT NULL DEFAULT '',
    price      REAL    NOT NULL,
    stock      INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  CREATE TABLE IF NOT EXISTS sales (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    customer_id  INTEGER,
    total_amount REAL    NOT NULL,
    sold_at      DATETIME DEFAULT CURRENT_TIMESTAMP
  );
  CREATE TABLE IF NOT EXISTS sale_items (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    sale_id    INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity   INTEGER NOT NULL,
    unit_price REAL    NOT NULL
  );
`);

// ── Products ──────────────────────────────────────────────

app.get('/products', (req, res) => {
  res.json(db.prepare('SELECT * FROM products ORDER BY id').all());
});

app.post('/products', (req, res) => {
  const { name, category = '', price, stock = 0 } = req.body;
  const result = db.prepare(
    'INSERT INTO products (name, category, price, stock) VALUES (?, ?, ?, ?)'
  ).run(name, category, price, stock);
  res.status(201).json({ id: result.lastInsertRowid, name, category, price, stock });
});

app.get('/products/:id', (req, res) => {
  const product = db.prepare('SELECT * FROM products WHERE id = ?').get(req.params.id);
  if (!product) return res.status(404).json({ error: 'not found' });
  res.json(product);
});

app.put('/products/:id', (req, res) => {
  const { name, category, price, stock } = req.body;
  db.prepare('UPDATE products SET name=?, category=?, price=?, stock=? WHERE id=?')
    .run(name, category, price, stock, req.params.id);
  res.json({ id: Number(req.params.id), name, category, price, stock });
});

app.delete('/products/:id', (req, res) => {
  db.prepare('DELETE FROM products WHERE id = ?').run(req.params.id);
  res.status(204).send();
});

// ── Sales ─────────────────────────────────────────────────

app.get('/sales', (req, res) => {
  res.json(db.prepare('SELECT * FROM sales ORDER BY id DESC').all());
});

app.post('/sales', (req, res) => {
  const { customer_id = null, items } = req.body;

  const createSale = db.transaction((items) => {
    let total = 0;
    const processedItems = items.map(item => {
      const product = db.prepare('SELECT * FROM products WHERE id = ?').get(item.product_id);
      if (!product) throw Object.assign(new Error(`product not found: ${item.product_id}`), { status: 400 });
      if (product.stock < item.quantity) throw Object.assign(new Error(`insufficient stock for product ${item.product_id}`), { status: 400 });
      total += product.price * item.quantity;
      return { ...item, unit_price: product.price };
    });

    const sale = db.prepare('INSERT INTO sales (customer_id, total_amount) VALUES (?, ?)').run(customer_id, total);
    const saleId = sale.lastInsertRowid;

    for (const item of processedItems) {
      db.prepare('INSERT INTO sale_items (sale_id, product_id, quantity, unit_price) VALUES (?, ?, ?, ?)')
        .run(saleId, item.product_id, item.quantity, item.unit_price);
      db.prepare('UPDATE products SET stock = stock - ? WHERE id = ?').run(item.quantity, item.product_id);
    }

    return { id: saleId, customer_id, total_amount: total, items: processedItems };
  });

  try {
    res.status(201).json(createSale(items));
  } catch (err) {
    res.status(err.status || 500).json({ error: err.message });
  }
});

app.get('/sales/:id', (req, res) => {
  const sale = db.prepare('SELECT * FROM sales WHERE id = ?').get(req.params.id);
  if (!sale) return res.status(404).json({ error: 'not found' });
  sale.items = db.prepare('SELECT * FROM sale_items WHERE sale_id = ?').all(req.params.id);
  res.json(sale);
});

// ── Reports ───────────────────────────────────────────────

app.get('/reports/summary', (req, res) => {
  res.json(db.prepare('SELECT COUNT(*) as total_sales, COALESCE(SUM(total_amount), 0) as total_revenue FROM sales').get());
});

app.get('/reports/by-category', (req, res) => {
  res.json(db.prepare(`
    SELECT p.category, COUNT(si.id) as count, COALESCE(SUM(si.quantity * si.unit_price), 0) as revenue
    FROM sale_items si
    JOIN products p ON p.id = si.product_id
    GROUP BY p.category
    ORDER BY revenue DESC
  `).all());
});

// ── Start ─────────────────────────────────────────────────

app.listen(3000, () => console.log('Server running on :3000'));
