package handlers

import (
	"classifieds/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

// Products is a type defined to represent product
type Products struct {
	l *log.Logger
}

// NewProducts defined an DI to log events occured in product api
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

	// Handle the request for editing an existing product
	if r.Method == http.MethodPut {
		p.l.Println("Handle PUT request", r.URL.Path)

		// expect the ID in the URI
		reg := regexp.MustCompile(`([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		id := g[0][1]
		idString, err := strconv.Atoi(id)

		if err != nil {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.updateProduct(idString, rw, r)
		return
	}

	// catch all.
	// if no method is specified, return error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
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

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
