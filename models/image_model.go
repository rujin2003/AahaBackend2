package model

type GalleryImage struct {
	ID          int    `json:"id"`
	ImageBase64 string `json:"image_base64"`
	Description string `json:"description,omitempty"`
}

func NewGalleryImage(id int, imageBase64 string, description string) *GalleryImage {
	return &GalleryImage{
		ID:          id,
		ImageBase64: imageBase64,
		Description: description,
	}
}
