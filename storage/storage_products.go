package storage

import (
	model "AahaFeltBackend/models"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "user=postgres password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Check if the database exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = 'products')").Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		_, err = db.Exec("CREATE DATABASE products")
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
	}

	db, err = sql.Open("postgres", connStr+" dbname=products")
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Close() {
	s.db.Close()
}

func (s *PostgresStorage) Init() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS products (
            id SERIAL PRIMARY KEY,
            name TEXT,
            category TEXT,
            description TEXT,
            image TEXT,
            images TEXT[],
            stock TEXT,
            price TEXT,
            listed TEXT,
            offer TEXT,
            sizes TEXT[],
            highlights TEXT,
            color TEXT,
            discount TEXT
        )
    `)
	return err
}

func (s *PostgresStorage) AddProducts(product model.Product) error {
	_, err := s.db.Exec(`
        INSERT INTO products (name, category, description, image, images, stock, price, listed, offer, sizes, highlights, color, discount)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		product.Name, product.Category, product.Description, product.Image, pq.Array(strings.Split(product.Images, ",")), product.Stock, product.Price, product.Listed, product.Offer, pq.Array(strings.Split(product.Sizes, ",")), product.Highlights, product.Color, product.Discount)
	return err
}

func (s *PostgresStorage) GetProducts() ([]model.Product, error) {
	rows, err := s.db.Query(`
        SELECT id, name, category, description, image, images, stock, price, listed, offer, sizes, highlights, color, discount
        FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		var images, sizes []string
		err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.Description, &product.Image, pq.Array(&images), &product.Stock, &product.Price, &product.Listed, &product.Offer, pq.Array(&sizes), &product.Highlights, &product.Color, &product.Discount)
		if err != nil {
			return nil, err
		}
		product.Images = strings.Join(images, ",")
		product.Sizes = strings.Join(sizes, ",")
		products = append(products, product)
	}
	return products, nil
}

func (s *PostgresStorage) GetProductsById(id int) (*model.Product, error) {
	row := s.db.QueryRow(`
        SELECT id, name, category, description, image, images, stock, price, listed, offer, sizes, highlights, color, discount
        FROM products WHERE id = $1`, id)
	product := &model.Product{}
	var images, sizes []string
	err := row.Scan(&product.ID, &product.Name, &product.Category, &product.Description, &product.Image, pq.Array(&images), &product.Stock, &product.Price, &product.Listed, &product.Offer, pq.Array(&sizes), &product.Highlights, &product.Color, &product.Discount)
	if err != nil {
		return nil, err
	}
	product.Images = strings.Join(images, ",")
	product.Sizes = strings.Join(sizes, ",")
	return product, nil
}

func (s *PostgresStorage) UpdateProductById(id int, product model.Product) error {
	_, err := s.db.Exec(`
        UPDATE products
        SET name = $1, category = $2, description = $3, image = $4, images = $5, stock = $6, price = $7, listed = $8, offer = $9, sizes = $10, highlights = $11, color = $12, discount = $13
        WHERE id = $14`,
		product.Name, product.Category, product.Description, product.Image, pq.Array(strings.Split(product.Images, ",")), product.Stock, product.Price, product.Listed, product.Offer, pq.Array(strings.Split(product.Sizes, ",")), product.Highlights, product.Color, product.Discount, id)
	return err
}

func (s *PostgresStorage) DeleteProductById(id int) error {
	_, err := s.db.Exec(`DELETE FROM products WHERE id = $1`, id)
	return err
}

func (s *PostgresStorage) GetCategories() ([]string, error) {
	rows, err := s.db.Query(`SELECT DISTINCT category FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
