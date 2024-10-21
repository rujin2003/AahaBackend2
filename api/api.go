package api

import (
	"AahaFeltBackend/storage"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type ApiServer struct {
	address string
	store   storage.Storage
}

func NewApiServer(address string, store storage.Storage) *ApiServer {
	return &ApiServer{
		address: address,
		store:   store,
	}
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func (s *ApiServer) Start() {
	router := mux.NewRouter()

	//MARK: Products
	router.HandleFunc("/products", makeHandler(s.handleGetProducts)).Methods("GET")
	router.HandleFunc("/products", makeHandler(s.handlePostProducts)).Methods("POST")
	router.HandleFunc("/products/{id}", makeHandler(s.handleGetProductsById)).Methods("GET")
	router.HandleFunc("/products/{id}", makeHandler(s.UpdateProductHandler)).Methods("POST")
	router.HandleFunc("/products/{id}", makeHandler(s.handleDeleteProduct)).Methods("DELETE")

	// MARK: Gallery
	router.HandleFunc("/gallery-images", makeHandler(s.addImageHandler)).Methods("POST")
	router.HandleFunc("/gallery-images/{id}", makeHandler(s.getImageHandler)).Methods("GET")
	router.HandleFunc("/gallery-images", makeHandler(s.getAllImageLinksHandler)).Methods("GET")
	router.HandleFunc("/gallery-images/{id}", makeHandler(s.deleteImageHandler)).Methods("DELETE")

	// MARK: Product Images
	router.HandleFunc("/productimage", makeHandler(s.AddProductImagesHandler)).Methods("POST")
	router.HandleFunc("/productimage/{product_name}", makeHandler(s.GetProductImagesByNameHandler)).Methods("GET")
	router.HandleFunc("/productimg/{id}", makeHandler(s.getImageHandler)).Methods("GET")
	router.HandleFunc("/productimage/{product_name}", makeHandler(s.DeleteProductImagesByNameHandler)).Methods("DELETE")

	// MARK: Status
	router.HandleFunc("/status", makeHandler(s.handleGetStatus)).Methods("GET")
	router.HandleFunc("/status", makeHandler(s.handleUpdateStatus)).Methods("POST")

	// CORS configuration
	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedOrigins([]string{"http://localhost:5173"}),
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	fmt.Printf("Server is starting on %s...\n", s.address)
	if err := http.ListenAndServe(s.address, corsMiddleware(router)); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
