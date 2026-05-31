package main

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

type Summary struct {
	TotalSales   int     `json:"total_sales"`
	TotalRevenue float64 `json:"total_revenue"`
}
