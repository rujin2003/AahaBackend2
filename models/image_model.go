package models

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

type ProductImage struct {
	ID          int    `json:"id"`
	ImageBase64 string `json:"image_base64"`
	Name        string `json:"image_name"`
	Type        string `json:"image_type"`
	Product_id  int    `json:"product_id"`
}

func NewProductImage(id int, imageBase64 string, name string, imageType string, productId int) *ProductImage {
	return &ProductImage{
		ID:          id,
		ImageBase64: imageBase64,
		Name:        name,
		Type:        imageType,
		Product_id:  productId,
	}
}
