package main

import (
	"database/sql"
	"fmt"
)

type SQLiteStore struct {
	db DB
}

func NewSQLiteStore(db DB) *SQLiteStore {
	return &SQLiteStore{db: db}
}

// ── Products ──────────────────────────────────────────────

func (s *SQLiteStore) List() ([]Product, error) {
	rows, err := s.db.Query(`SELECT id, name, category, price, stock FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Stock)
		products = append(products, p)
	}
	return products, nil
}

func (s *SQLiteStore) Get(id int) (Product, error) {
	var p Product
	err := s.db.QueryRow(`SELECT id, name, category, price, stock FROM products WHERE id = ?`, id).
		Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Stock)
	return p, err
}

func (s *SQLiteStore) Create(p Product) (Product, error) {
	res, err := s.db.Exec(`INSERT INTO products (name, category, price, stock) VALUES (?, ?, ?, ?)`,
		p.Name, p.Category, p.Price, p.Stock)
	if err != nil {
		return Product{}, err
	}
	id, _ := res.LastInsertId()
	p.ID = int(id)
	return p, nil
}

func (s *SQLiteStore) Update(id int, p Product) (Product, error) {
	_, err := s.db.Exec(`UPDATE products SET name=?, category=?, price=?, stock=? WHERE id=?`,
		p.Name, p.Category, p.Price, p.Stock, id)
	if err != nil {
		return Product{}, err
	}
	p.ID = id
	return p, nil
}

func (s *SQLiteStore) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM products WHERE id = ?`, id)
	return err
}

// ── Sales ─────────────────────────────────────────────────

func (s *SQLiteStore) ListSales() ([]Sale, error) {
	rows, err := s.db.Query(`SELECT id, customer_id, total_amount, sold_at FROM sales ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sales := []Sale{}
	for rows.Next() {
		var sale Sale
		var cid sql.NullInt64
		rows.Scan(&sale.ID, &cid, &sale.TotalAmount, &sale.SoldAt)
		if cid.Valid {
			v := int(cid.Int64)
			sale.CustomerID = &v
		}
		sales = append(sales, sale)
	}
	return sales, nil
}

func (s *SQLiteStore) GetSale(id int) (Sale, error) {
	var sale Sale
	var cid sql.NullInt64
	err := s.db.QueryRow(`SELECT id, customer_id, total_amount, sold_at FROM sales WHERE id = ?`, id).
		Scan(&sale.ID, &cid, &sale.TotalAmount, &sale.SoldAt)
	if err != nil {
		return Sale{}, err
	}
	if cid.Valid {
		v := int(cid.Int64)
		sale.CustomerID = &v
	}

	rows, _ := s.db.Query(`SELECT product_id, quantity, unit_price FROM sale_items WHERE sale_id = ?`, id)
	defer rows.Close()
	for rows.Next() {
		var item SaleItem
		rows.Scan(&item.ProductID, &item.Quantity, &item.UnitPrice)
		sale.Items = append(sale.Items, item)
	}
	return sale, nil
}

func (s *SQLiteStore) CreateSale(req SaleRequest) (Sale, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return Sale{}, err
	}

	var total float64
	for i, item := range req.Items {
		var price float64
		var stock int
		err := tx.QueryRow(`SELECT price, stock FROM products WHERE id = ?`, item.ProductID).
			Scan(&price, &stock)
		if err == sql.ErrNoRows {
			tx.Rollback()
			return Sale{}, fmt.Errorf("product not found: %d", item.ProductID)
		}
		if stock < item.Quantity {
			tx.Rollback()
			return Sale{}, fmt.Errorf("insufficient stock for product %d", item.ProductID)
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
		return Sale{}, err
	}

	return Sale{
		ID:          int(saleID),
		CustomerID:  req.CustomerID,
		TotalAmount: total,
		Items:       req.Items,
	}, nil
}

func (s *SQLiteStore) GetSummary() (Summary, error) {
	var sum Summary
	err := s.db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(total_amount), 0) FROM sales`).
		Scan(&sum.TotalSales, &sum.TotalRevenue)
	return sum, err
}
