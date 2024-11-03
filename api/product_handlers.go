package api

import (
	model "AahaFeltBackend/models"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (s *ApiServer) handlePostProducts(w http.ResponseWriter, r *http.Request) error {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return fmt.Errorf("error parsing multipart form: %v", err)
	}

	name := r.FormValue("name")
	category := r.FormValue("category")
	description := r.FormValue("description")
	stock := r.FormValue("stock")
	price := r.FormValue("price")
	listed := r.FormValue("listed")
	offer := r.FormValue("offer")
	sizes := r.FormValue("sizes")
	highlights := r.FormValue("highlights")
	color := r.FormValue("color")
	discount := r.FormValue("discount")

	imageFile, _, err := r.FormFile("image")
	if err != nil {
		return fmt.Errorf("error reading image file: %v", err)
	}
	defer imageFile.Close()

	imageBytes, err := ioutil.ReadAll(imageFile)
	if err != nil {
		return fmt.Errorf("error reading image bytes: %v", err)
	}
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	var imagesBase64 []string
	for _, fileHeader := range r.MultipartForm.File["images"] {
		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("error reading additional image file: %v", err)
		}
		defer file.Close()

		imageBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return fmt.Errorf("error reading additional image bytes: %v", err)
		}
		imagesBase64 = append(imagesBase64, base64.StdEncoding.EncodeToString(imageBytes))
	}

	imagesBase64Str := strings.Join(imagesBase64, ",")

	// Create a new product
	product := model.NewProduct(
		name,
		category,
		description,
		imageBase64,
		imagesBase64Str,
		stock,
		price,
		listed,
		offer,
		sizes,
		highlights,
		color,
		discount,
	)

	if err := s.store.AddProducts(*product); err != nil {
		return fmt.Errorf("failed to add product: %v", err)
	}

	return writeJSON(w, http.StatusOK, product)
}

func (s *ApiServer) handleGetProducts(w http.ResponseWriter, r *http.Request) error {
	products, err := s.store.GetProducts()
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, products)
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
	categories, err := s.store.GetCategories()
	if err != nil {

	}

	return writeJSON(w, http.StatusOK, categories)
}
