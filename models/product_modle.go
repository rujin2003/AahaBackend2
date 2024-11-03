package model

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Images      string `json:"images"`
	Stock       string `json:"stock"`
	Price       string `json:"price"`
	Listed      string `json:"listed"`
	Offer       string `json:"offer"`
	Sizes       string `json:"sizes"`
	Highlights  string `json:"highlights"`
	Color       string `json:"color"`
	Discount    string `json:"discount"`
}

func NewProduct(name, category, description, image, images string, stock, price string, listed string, offer, sizes, highlights, color string, discount string) *Product {
	return &Product{
		Name:        name,
		Category:    category,
		Description: description,
		Image:       image,
		Images:      images,
		Stock:       stock,
		Price:       price,
		Listed:      listed,
		Offer:       offer,
		Sizes:       sizes,
		Highlights:  highlights,
		Color:       color,
		Discount:    discount,
	}
}
