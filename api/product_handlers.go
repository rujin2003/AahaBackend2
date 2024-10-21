package api

import (
	model "AahaFeltBackend/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func verifyToken(tokenString string) error {
	password := strings.TrimSpace(os.Getenv("PASSWORD"))

	if tokenString != password {
		return fmt.Errorf("invalid token")
	}
	return nil
}

func makeHandler(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Missing Authorization header"}`, http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error": "Invalid Authorization header format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(authHeader[len("Bearer "):])
		if err := verifyToken(tokenString); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusUnauthorized)
			return
		}

		if err := fn(w, r); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func (s *ApiServer) handleGetProducts(w http.ResponseWriter, r *http.Request) error {
	products, err := s.store.GetProducts()
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, products)
}
func (s *ApiServer) handlePostProducts(w http.ResponseWriter, r *http.Request) error {
	product := model.Product{}
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		return fmt.Errorf("invalid request payload: %v", err)
	}

	pro := model.NewProduct(
		product.Name,
		product.Category,
		product.Description,
		product.Image,
		product.Images,
		product.Stock,
		product.Price,
		product.Listed,
		product.Offer,
		product.Sizes,
		product.Highlights,
		product.Color,
		product.Discount,
	)

	if err := s.store.AddProducts(*pro); err != nil {
		return fmt.Errorf("failed to add product: %v", err)
	}

	return writeJSON(w, http.StatusOK, pro)
}

func (s *ApiServer) handleGetProductsById(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid product ID: %v", err)
	}

	product, err := s.store.GetProductsById(id)
	if err != nil {
		return fmt.Errorf("failed to retrieve product: %v", err)
	}

	return writeJSON(w, http.StatusOK, product)
}

func (s *ApiServer) UpdateProductHandler(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid product ID: %v", err)
	}

	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		return fmt.Errorf("invalid request payload: %v", err)
	}

	if err := s.store.UpdateProductById(id, product); err != nil {
		return fmt.Errorf("failed to update product: %v", err)
	}

	return writeJSON(w, http.StatusOK, "Product updated successfully")
}

func (s *ApiServer) handleDeleteProduct(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid product ID: %v", err)
	}

	if err := s.store.DeleteProductById(id); err != nil {
		return fmt.Errorf("failed to delete product: %v", err)
	}

	return writeJSON(w, http.StatusOK, "Product deleted successfully")
}
func (s *ApiServer) handleGetCategories(w http.ResponseWriter, r *http.Request) error {
	categories := s.store.GetCategories()

	return writeJSON(w, http.StatusOK, categories)
}
