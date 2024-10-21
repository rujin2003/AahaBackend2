package storage

import (
	model "AahaFeltBackend2/models"
	"database/sql"
	"net/http"

	_ "github.com/lib/pq"
)

type Storage interface {
	AddProducts(product model.Product) error
	GetProducts() ([]model.Product, error)
	GetProductsById(id int) (*model.Product, error)
	UpdateProductById(id int, product model.Product) error
	DeleteProductById(id int) error
	GetCategories() error
	// gallery
	AddImage(r *http.Request) (int, error)
	GetImageByID(id int) (*model.GalleryImage, error)
	GetAllImageIDs() ([]int, error)
	DeleteImageByID(id int) error

	// product image
	AddProductImages(r *http.Request) ([]int, error)

	// GetImagesByProductName(productName string) ([]model.ProductImage, error)
	GetImagesByProductName(string) ([]model.ProductImage, error)
	DeleteProductImageByName(productName string) error

	GetProductImageByID(id int) (*model.GalleryImage, error)

	// AddSellStatus(status string) error

	AddSellStatus(status string) error
	UpdateSellStatus(status string) error

	Close()
	Init() error
}

type PostgresStorage struct {
	db *sql.DB
}
