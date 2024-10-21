package model

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Images      string `json:"images"`
	Stock       int    `json:"stock"`
	Price       int    `json:"price"`
	Listed      bool   `json:"listed"`
	Offer       string `json:"offer"`
	Sizes       string `json:"sizes"`
	Highlights  string `json:"highlights"`
	Color       string `json:"color"`
	Discount    int    `json:"discount"`
}

func NewProduct(name, category, description, image, images string, stock, price int, listed bool, offer, sizes, highlights, color string, discount int) *Product {
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
