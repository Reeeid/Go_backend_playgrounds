package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// ── Models ──────────────────────────────────────────────

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
}

type SaleItem struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price,omitempty"`
}

type SaleRequest struct {
	CustomerID *int       `json:"customer_id"`
	Items      []SaleItem `json:"items"`
}

type Sale struct {
	ID          int        `json:"id"`
	CustomerID  *int       `json:"customer_id"`
	TotalAmount float64    `json:"total_amount"`
	SoldAt      string     `json:"sold_at,omitempty"`
	Items       []SaleItem `json:"items,omitempty"`
}

// ── DB Setup ─────────────────────────────────────────────

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./store.db")
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
		}
	}
}

// ── Helpers ───────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func errJSON(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// pathID extracts the trailing integer from a URL path.
// e.g. "/products/42" with prefix "/products/" → 42
func pathID(r *http.Request, prefix string) (int, bool) {
	seg := strings.TrimPrefix(r.URL.Path, prefix)
	seg = strings.TrimSuffix(seg, "/")
	id, err := strconv.Atoi(seg)
	return id, err == nil
}

// ── Handlers: Products ────────────────────────────────────

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		rows, err := db.Query(`SELECT id, name, category, price, stock FROM products ORDER BY id`)
		if err != nil {
			errJSON(w, 500, err.Error()); return
		}
		defer rows.Close()

		products := []Product{}
		for rows.Next() {
			var p Product
			rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Stock)
			products = append(products, p)
		}
		writeJSON(w, 200, products)

	case http.MethodPost:
		var p Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			errJSON(w, 400, "invalid json"); return
		}
		res, err := db.Exec(`INSERT INTO products (name, category, price, stock) VALUES (?, ?, ?, ?)`,
			p.Name, p.Category, p.Price, p.Stock)
		if err != nil {
			errJSON(w, 500, err.Error()); return
		}
		id, _ := res.LastInsertId()
		p.ID = int(id)
		writeJSON(w, 201, p)

	default:
		errJSON(w, 405, "method not allowed")
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "/products/")
	if !ok {
		errJSON(w, 400, "invalid id"); return
	}

	switch r.Method {

	case http.MethodGet:
		var p Product
		err := db.QueryRow(`SELECT id, name, category, price, stock FROM products WHERE id = ?`, id).
			Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Stock)
		if err == sql.ErrNoRows {
			errJSON(w, 404, "not found"); return
		}
		writeJSON(w, 200, p)

	case http.MethodPut:
		var p Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			errJSON(w, 400, "invalid json"); return
		}
		_, err := db.Exec(`UPDATE products SET name=?, category=?, price=?, stock=? WHERE id=?`,
			p.Name, p.Category, p.Price, p.Stock, id)
		if err != nil {
			errJSON(w, 500, err.Error()); return
		}
		p.ID = id
		writeJSON(w, 200, p)

	case http.MethodDelete:
		db.Exec(`DELETE FROM products WHERE id = ?`, id)
		w.WriteHeader(204)

	default:
		errJSON(w, 405, "method not allowed")
	}
}

// ── Handlers: Sales ───────────────────────────────────────

func salesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		rows, _ := db.Query(`SELECT id, customer_id, total_amount, sold_at FROM sales ORDER BY id DESC`)
		defer rows.Close()

		sales := []Sale{}
		for rows.Next() {
			var s Sale
			var cid sql.NullInt64
			rows.Scan(&s.ID, &cid, &s.TotalAmount, &s.SoldAt)
			if cid.Valid {
				v := int(cid.Int64)
				s.CustomerID = &v
			}
			sales = append(sales, s)
		}
		writeJSON(w, 200, sales)

	case http.MethodPost:
		var req SaleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errJSON(w, 400, "invalid json"); return
		}

		tx, err := db.Begin()
		if err != nil {
			errJSON(w, 500, err.Error()); return
		}

		var total float64
		for i, item := range req.Items {
			var price float64
			var stock int
			err := tx.QueryRow(`SELECT price, stock FROM products WHERE id = ?`, item.ProductID).
				Scan(&price, &stock)
			if err == sql.ErrNoRows {
				tx.Rollback()
				errJSON(w, 400, "product not found: "+strconv.Itoa(item.ProductID)); return
			}
			if stock < item.Quantity {
				tx.Rollback()
				errJSON(w, 400, "insufficient stock for product "+strconv.Itoa(item.ProductID)); return
			}
			req.Items[i].UnitPrice = price
			total += price * float64(item.Quantity)
		}

		var customerID any = nil
		if req.CustomerID != nil {
			customerID = *req.CustomerID
		}
		res, _ := tx.Exec(`INSERT INTO sales (customer_id, total_amount) VALUES (?, ?)`, customerID, total)
		saleID, _ := res.LastInsertId()

		for _, item := range req.Items {
			tx.Exec(`INSERT INTO sale_items (sale_id, product_id, quantity, unit_price) VALUES (?, ?, ?, ?)`,
				saleID, item.ProductID, item.Quantity, item.UnitPrice)
			tx.Exec(`UPDATE products SET stock = stock - ? WHERE id = ?`, item.Quantity, item.ProductID)
		}

		if err := tx.Commit(); err != nil {
			errJSON(w, 500, err.Error()); return
		}

		writeJSON(w, 201, Sale{
			ID:          int(saleID),
			CustomerID:  req.CustomerID,
			TotalAmount: total,
			Items:       req.Items,
		})

	default:
		errJSON(w, 405, "method not allowed")
	}
}

func saleDetailHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "/sales/")
	if !ok {
		errJSON(w, 400, "invalid id"); return
	}

	var s Sale
	var cid sql.NullInt64
	err := db.QueryRow(`SELECT id, customer_id, total_amount, sold_at FROM sales WHERE id = ?`, id).
		Scan(&s.ID, &cid, &s.TotalAmount, &s.SoldAt)
	if err == sql.ErrNoRows {
		errJSON(w, 404, "not found"); return
	}
	if cid.Valid {
		v := int(cid.Int64)
		s.CustomerID = &v
	}

	rows, _ := db.Query(`SELECT product_id, quantity, unit_price FROM sale_items WHERE sale_id = ?`, id)
	defer rows.Close()
	for rows.Next() {
		var item SaleItem
		rows.Scan(&item.ProductID, &item.Quantity, &item.UnitPrice)
		s.Items = append(s.Items, item)
	}

	writeJSON(w, 200, s)
}

// ── Handler: Reports ─────────────────────────────────────

func reportSummaryHandler(w http.ResponseWriter, r *http.Request) {
	var count int
	var total float64
	db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(total_amount), 0) FROM sales`).Scan(&count, &total)
	writeJSON(w, 200, map[string]any{
		"total_sales":   count,
		"total_revenue": total,
	})
}

// ── Router ────────────────────────────────────────────────

func router(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	switch {
	case p == "/products" || p == "/products/":
		productsHandler(w, r)
	case strings.HasPrefix(p, "/products/"):
		productHandler(w, r)
	case p == "/sales" || p == "/sales/":
		salesHandler(w, r)
	case strings.HasPrefix(p, "/sales/"):
		saleDetailHandler(w, r)
	case p == "/reports/summary":
		reportSummaryHandler(w, r)
	default:
		errJSON(w, 404, "not found")
	}
}

// ── Main ──────────────────────────────────────────────────

func main() {
	initDB()
	http.HandleFunc("/", router)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
