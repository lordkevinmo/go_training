package handlers

import (
	"classifieds/data"
	"log"
	"net/http"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle the request for a list of products
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	// Handle the request for adding a new Product
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	// catch all.
	// if no method is specified, return error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST product")
}

func (p *Products) getProducts(rw http.ResponseWriter, req *http.Request) {
	p.l.Println("Handle products get")

	// Fetch the products from datastore
	lp := data.GetProducts()

	// Serialize list data to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshal JSON", http.StatusInternalServerError)
	}
}
