package main

type ProductRepository interface {
	List() ([]Product, error)
	Get(id int) (Product, error)
	Create(p Product) (Product, error)
	Update(id int, p Product) (Product, error)
	Delete(id int) error
}

type SaleRepository interface {
	ListSales() ([]Sale, error)
	GetSale(id int) (Sale, error)
	CreateSale(req SaleRequest) (Sale, error)
	GetSummary() (Summary, error)
}
