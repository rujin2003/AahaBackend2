package storage

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	model "AahaFeltBackend/models"
)

func (s *PostgresStorage) InitGallery() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS gallery (
            id SERIAL PRIMARY KEY,
            image_base64 TEXT NOT NULL,
            description TEXT
        );
    `)
	fmt.Println("Gallery table created")
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) AddImage(r *http.Request) (int, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return 0, fmt.Errorf("error parsing form data: %v", err)
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		return 0, fmt.Errorf("error retrieving the file: %v", err)
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, fmt.Errorf("error reading the file: %v", err)
	}

	imageBase64 := base64.StdEncoding.EncodeToString(fileBytes)
	description := r.FormValue("description")

	var id int
	err = s.db.QueryRow(`
		INSERT INTO gallery (image_base64, description)
		VALUES ($1, $2)
		RETURNING id`, imageBase64, description).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("error inserting image into database: %v", err)
	}

	return id, nil
}

func (s *PostgresStorage) GetImageByID(id int) (*model.GalleryImage, error) {
	var image model.GalleryImage

	err := s.db.QueryRow(`
		SELECT id, image_base64, description 
		FROM gallery 
		WHERE id = $1`, id).
		Scan(&image.ID, &image.ImageBase64, &image.Description)
	if err != nil {
		return nil, fmt.Errorf("error retrieving image from database: %v", err)
	}

	return &image, nil
}

// MARK: Get gallery image by id

func (s *PostgresStorage) GetAllImageIDs() ([]int, error) {
	rows, err := s.db.Query(`SELECT id FROM gallery`)
	if err != nil {
		return nil, fmt.Errorf("error retrieving image IDs from database: %v", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error scanning image ID: %v", err)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through image IDs: %v", err)
	}

	return ids, nil
}

func (s *PostgresStorage) DeleteImageByID(id int) error {
	result, err := s.db.Exec(`DELETE FROM gallery WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("error deleting image from database: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error fetching affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no image found with the given ID")
	}

	return nil
}
