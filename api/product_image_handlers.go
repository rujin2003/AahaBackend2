package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *ApiServer) GetProductImagesByNameHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	productName := vars["product_name"]

	// Call the storage function to get images by product name
	images, err := s.store.GetImagesByProductName(productName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve images: %v", err), http.StatusInternalServerError)
		return nil
	}

	// Create an array to hold the image links
	var imageLinks []string
	for _, image := range images {
		// Generate the URL for each image
		link := fmt.Sprintf("http://localhost:3000/productimg/%d", image.ID)
		imageLinks = append(imageLinks, link)
	}

	// Set response header to JSON and encode the image links as a response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string][]string{"image_links": imageLinks}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return nil
	}

	return nil
}

// AddProductImagesHandler handles requests to add new product images
func (s *ApiServer) AddProductImagesHandler(w http.ResponseWriter, r *http.Request) error {

	imageIDs, err := s.store.AddProductImages(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add images: %v", err), http.StatusInternalServerError)
		return nil
	}

	// Generate URLs for the uploaded images
	var imageLinks []map[string]string
	for _, imageID := range imageIDs {
		imageLink := map[string]string{
			"image_id": fmt.Sprintf("%d", imageID),
			"url":      fmt.Sprintf("http://localhost:3000/productimg/%d", imageID),
		}
		imageLinks = append(imageLinks, imageLink)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(imageLinks)
	return nil
}

func (s *ApiServer) getProductImageHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return nil
	}

	image, err := s.store.GetProductImageByID(id)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return nil
	}

	imageBytes, err := base64.StdEncoding.DecodeString(image.ImageBase64)
	if err != nil {
		http.Error(w, "Failed to decode image", http.StatusInternalServerError)
		return nil
	}

	w.Header().Set("Content-Disposition", "inline; filename=image.jpg")
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(imageBytes)
	return nil

}

func (s *ApiServer) DeleteProductImagesByNameHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	productName := vars["product_name"]

	err := s.store.DeleteProductImageByName(productName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete images: %v", err), http.StatusInternalServerError)
		return nil
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
