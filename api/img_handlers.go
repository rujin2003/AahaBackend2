package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *ApiServer) addImageHandler(w http.ResponseWriter, r *http.Request) error {

	id, err := s.store.AddImage(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add image: %v", err), http.StatusInternalServerError)
		return nil
	}

	url := fmt.Sprintf("http://localhost:8080/images/%d", id)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"image_url": "%s"}`, url)))
	return nil
}

func (s *ApiServer) getImageHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return nil
	}

	image, err := s.store.GetImageByID(id)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return nil
	}

	// Ensure image is not nil
	if image == nil {
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

func (s *ApiServer) getAllImageLinksHandler(w http.ResponseWriter, r *http.Request) error {
	ids, err := s.store.GetAllImageIDs()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve image links: %v", err), http.StatusInternalServerError)
		return nil
	}

	var links []string
	for _, id := range ids {
		link := fmt.Sprintf("http://localhost:3000/gallery-images/%d", id)
		links = append(links, link)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string][]string{"image_links": links}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return nil
	}
	return nil
}
func (s *ApiServer) deleteImageHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return nil
	}

	err = s.store.DeleteImageByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete image: %v", err), http.StatusInternalServerError)
		return nil
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
