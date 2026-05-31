package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	products ProductRepository
	sales    SaleRepository
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

func pathIntID(r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	return id, err == nil
}

// ── Products ──────────────────────────────────────────────

func (h *Handler) listProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.products.List()
	if err != nil {
		errJSON(w, 500, err.Error()); return
	}
	writeJSON(w, 200, products)
}

func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		errJSON(w, 400, "invalid json"); return
	}
	created, err := h.products.Create(p)
	if err != nil {
		errJSON(w, 500, err.Error()); return
	}
	writeJSON(w, 201, created)
}

func (h *Handler) getProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := pathIntID(r)
	if !ok {
		errJSON(w, 400, "invalid id"); return
	}
	p, err := h.products.Get(id)
	if err == sql.ErrNoRows {
		errJSON(w, 404, "not found"); return
	}
	writeJSON(w, 200, p)
}

func (h *Handler) updateProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := pathIntID(r)
	if !ok {
		errJSON(w, 400, "invalid id"); return
	}
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		errJSON(w, 400, "invalid json"); return
	}
	updated, err := h.products.Update(id, p)
	if err != nil {
		errJSON(w, 500, err.Error()); return
	}
	writeJSON(w, 200, updated)
}

func (h *Handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := pathIntID(r)
	if !ok {
		errJSON(w, 400, "invalid id"); return
	}
	h.products.Delete(id)
	w.WriteHeader(204)
}

// ── Sales ─────────────────────────────────────────────────

func (h *Handler) listSales(w http.ResponseWriter, r *http.Request) {
	sales, err := h.sales.ListSales()
	if err != nil {
		errJSON(w, 500, err.Error()); return
	}
	writeJSON(w, 200, sales)
}

func (h *Handler) createSale(w http.ResponseWriter, r *http.Request) {
	var req SaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errJSON(w, 400, "invalid json"); return
	}
	sale, err := h.sales.CreateSale(req)
	if err != nil {
		errJSON(w, 400, err.Error()); return
	}
	writeJSON(w, 201, sale)
}

func (h *Handler) getSale(w http.ResponseWriter, r *http.Request) {
	id, ok := pathIntID(r)
	if !ok {
		errJSON(w, 400, "invalid id"); return
	}
	sale, err := h.sales.GetSale(id)
	if err == sql.ErrNoRows {
		errJSON(w, 404, "not found"); return
	}
	writeJSON(w, 200, sale)
}

// ── Reports ───────────────────────────────────────────────

func (h *Handler) reportSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.sales.GetSummary()
	if err != nil {
		errJSON(w, 500, err.Error()); return
	}
	writeJSON(w, 200, summary)
}
