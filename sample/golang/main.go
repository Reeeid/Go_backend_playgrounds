package main

import (
	"log"
	"net/http"
	"os"
)

func bearerAuth(next http.Handler) http.Handler {
	token := os.Getenv("API_TOKEN")
	if token == "" {
		token = "secret"
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+token {
			errJSON(w, 401, "unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	db, err := NewDB("./store.db")
	if err != nil {
		log.Fatal(err)
	}

	store := NewSQLiteStore(db)
	h := &Handler{
		products: store,
		sales:    store,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products",         h.listProducts)
	mux.HandleFunc("POST /products",        h.createProduct)
	mux.HandleFunc("GET /products/{id}",    h.getProduct)
	mux.HandleFunc("PUT /products/{id}",    h.updateProduct)
	mux.HandleFunc("DELETE /products/{id}", h.deleteProduct)
	mux.HandleFunc("GET /sales",            h.listSales)
	mux.HandleFunc("POST /sales",           h.createSale)
	mux.HandleFunc("GET /sales/{id}",       h.getSale)
	mux.HandleFunc("GET /reports/summary",  h.reportSummary)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", bearerAuth(mux)))
}
