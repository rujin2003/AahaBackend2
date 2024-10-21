package storage

import (
	model "AahaFeltBackend2/models"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

// MARK: productimg table
func (s *PostgresStorage) InitProductImage() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS productimg (
            id SERIAL PRIMARY KEY,
            image_base64 TEXT NOT NULL,
            name TEXT NOT NULL,
            type TEXT,
            product_id INT
        );
    `)
	if err != nil {
		fmt.Println("Error creating productimg table:", err)
		return err
	}
	fmt.Println("Product Image table created successfully")
	return nil
}

func (s *PostgresStorage) AddProductImages(r *http.Request) ([]int, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, fmt.Errorf("error parsing form data: %v", err)
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	var imageIDs []int

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("error retrieving the file: %v", err)
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("error reading the file: %v", err)
		}

		imageBase64 := base64.StdEncoding.EncodeToString(fileBytes)
		name := r.FormValue("name")
		imageType := r.FormValue("type")
		productID := r.FormValue("product_id")

		var imageID int

		err = s.db.QueryRow(`
			INSERT INTO productimg (image_base64, name, type, product_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id`, imageBase64, name, imageType, productID).Scan(&imageID)

		if err != nil {
			return nil, fmt.Errorf("error inserting image into database: %v", err)
		}
		imageIDs = append(imageIDs, imageID)
	}

	return imageIDs, nil
}

func (s *PostgresStorage) GetImagesByProductName(productName string) ([]model.ProductImage, error) {
	var images []model.ProductImage
	query := `
        SELECT id,image_base64,name, type,product_id
        FROM productimg
        WHERE name = $1
    `

	rows, err := s.db.Query(query, productName)
	if err != nil {
		fmt.Printf("Error fetching images from database: %v\n", err)
		return nil, fmt.Errorf("error fetching images from database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var image model.ProductImage
		err = rows.Scan(&image.ID, &image.ImageBase64, &image.Name, &image.Type, &image.Product_id)
		if err != nil {
			fmt.Printf("Error scanning image row: %v\n", err)
			return nil, fmt.Errorf("error scanning image row: %v", err)
		}
		fmt.Printf("Fetched image: %+v\n", image)
		images = append(images, image)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Error iterating rows: %v\n", err)
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return images, nil
}

func (s *PostgresStorage) DeleteProductImageByName(productName string) error {
	query := `
		DELETE FROM productimg 
		WHERE name = $1
	`

	result, err := s.db.Exec(query, productName)
	if err != nil {
		return fmt.Errorf("error deleting image from database: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error fetching affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no image found with the given product name")
	}

	return nil
}

func (s *PostgresStorage) GetProductImageByID(id int) (*model.GalleryImage, error) {
	var image model.GalleryImage

	err := s.db.QueryRow(`
		SELECT id, image_base64, description 
		FROM productimg 
		WHERE id = $1`, id).
		Scan(&image.ID, &image.ImageBase64, &image.Description)
	if err != nil {
		return nil, fmt.Errorf("error retrieving image from database: %v", err)
	}

	return &image, nil
}
