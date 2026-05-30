<?php
// 起動: php -S localhost:8080 index.php

declare(strict_types=1);

// ── DB Setup ─────────────────────────────────────────────

$db = new PDO('sqlite:store.db');
$db->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
$db->setAttribute(PDO::ATTR_DEFAULT_FETCH_MODE, PDO::FETCH_ASSOC);

$db->exec("CREATE TABLE IF NOT EXISTS products (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT    NOT NULL,
    category   TEXT    NOT NULL DEFAULT '',
    price      REAL    NOT NULL,
    stock      INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
)");
$db->exec("CREATE TABLE IF NOT EXISTS sales (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    customer_id  INTEGER,
    total_amount REAL    NOT NULL,
    sold_at      DATETIME DEFAULT CURRENT_TIMESTAMP
)");
$db->exec("CREATE TABLE IF NOT EXISTS sale_items (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    sale_id    INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity   INTEGER NOT NULL,
    unit_price REAL    NOT NULL
)");

// ── Helpers ───────────────────────────────────────────────

function respond(mixed $data, int $status = 200): void {
    http_response_code($status);
    header('Content-Type: application/json');
    echo json_encode($data);
    exit;
}

function body(): array {
    return json_decode(file_get_contents('php://input'), true) ?? [];
}

// ── Router ────────────────────────────────────────────────

$method = $_SERVER['REQUEST_METHOD'];
$path   = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);
$parts  = array_values(array_filter(explode('/', $path)));

// /products  or  /products/:id
if (($parts[0] ?? '') === 'products') {
    $id = isset($parts[1]) ? (int) $parts[1] : null;

    if ($id === null) {
        // GET /products
        if ($method === 'GET') {
            respond($db->query('SELECT * FROM products ORDER BY id')->fetchAll());
        }
        // POST /products
        if ($method === 'POST') {
            $b = body();
            $stmt = $db->prepare('INSERT INTO products (name, category, price, stock) VALUES (?, ?, ?, ?)');
            $stmt->execute([$b['name'], $b['category'] ?? '', $b['price'], $b['stock'] ?? 0]);
            $b['id'] = (int) $db->lastInsertId();
            respond($b, 201);
        }
    } else {
        // GET /products/:id
        if ($method === 'GET') {
            $stmt = $db->prepare('SELECT * FROM products WHERE id = ?');
            $stmt->execute([$id]);
            $row = $stmt->fetch();
            if (!$row) respond(['error' => 'not found'], 404);
            respond($row);
        }
        // PUT /products/:id
        if ($method === 'PUT') {
            $b = body();
            $db->prepare('UPDATE products SET name=?, category=?, price=?, stock=? WHERE id=?')
               ->execute([$b['name'], $b['category'], $b['price'], $b['stock'], $id]);
            $b['id'] = $id;
            respond($b);
        }
        // DELETE /products/:id
        if ($method === 'DELETE') {
            $db->prepare('DELETE FROM products WHERE id = ?')->execute([$id]);
            respond(null, 204);
        }
    }
}

// /sales  or  /sales/:id
if (($parts[0] ?? '') === 'sales') {
    $id = isset($parts[1]) ? (int) $parts[1] : null;

    // GET /sales
    if ($method === 'GET' && $id === null) {
        respond($db->query('SELECT * FROM sales ORDER BY id DESC')->fetchAll());
    }

    // GET /sales/:id
    if ($method === 'GET' && $id !== null) {
        $stmt = $db->prepare('SELECT * FROM sales WHERE id = ?');
        $stmt->execute([$id]);
        $sale = $stmt->fetch();
        if (!$sale) respond(['error' => 'not found'], 404);

        $stmt = $db->prepare('SELECT * FROM sale_items WHERE sale_id = ?');
        $stmt->execute([$id]);
        $sale['items'] = $stmt->fetchAll();
        respond($sale);
    }

    // POST /sales
    if ($method === 'POST') {
        $b = body();
        $items = $b['items'] ?? [];

        try {
            $db->beginTransaction();
            $total = 0.0;
            $processedItems = [];

            foreach ($items as $item) {
                $stmt = $db->prepare('SELECT * FROM products WHERE id = ?');
                $stmt->execute([$item['product_id']]);
                $product = $stmt->fetch();

                if (!$product) throw new RuntimeException("product not found: {$item['product_id']}", 400);
                if ($product['stock'] < $item['quantity']) throw new RuntimeException("insufficient stock for product {$item['product_id']}", 400);

                $item['unit_price'] = (float) $product['price'];
                $total += $item['unit_price'] * $item['quantity'];
                $processedItems[] = $item;
            }

            $db->prepare('INSERT INTO sales (customer_id, total_amount) VALUES (?, ?)')
               ->execute([$b['customer_id'] ?? null, $total]);
            $saleId = (int) $db->lastInsertId();

            foreach ($processedItems as $item) {
                $db->prepare('INSERT INTO sale_items (sale_id, product_id, quantity, unit_price) VALUES (?, ?, ?, ?)')
                   ->execute([$saleId, $item['product_id'], $item['quantity'], $item['unit_price']]);
                $db->prepare('UPDATE products SET stock = stock - ? WHERE id = ?')
                   ->execute([$item['quantity'], $item['product_id']]);
            }

            $db->commit();
            respond([
                'id'           => $saleId,
                'customer_id'  => $b['customer_id'] ?? null,
                'total_amount' => $total,
                'items'        => $processedItems,
            ], 201);
        } catch (RuntimeException $e) {
            $db->rollBack();
            respond(['error' => $e->getMessage()], $e->getCode() ?: 400);
        }
    }
}

// /reports/summary
if (($parts[0] ?? '') === 'reports' && ($parts[1] ?? '') === 'summary') {
    $row = $db->query('SELECT COUNT(*) as total_sales, COALESCE(SUM(total_amount), 0) as total_revenue FROM sales')->fetch();
    respond($row);
}

// /reports/by-category
if (($parts[0] ?? '') === 'reports' && ($parts[1] ?? '') === 'by-category') {
    $rows = $db->query("
        SELECT p.category, COUNT(si.id) as count, COALESCE(SUM(si.quantity * si.unit_price), 0) as revenue
        FROM sale_items si
        JOIN products p ON p.id = si.product_id
        GROUP BY p.category
        ORDER BY revenue DESC
    ")->fetchAll();
    respond($rows);
}

respond(['error' => 'not found'], 404);
