package storage

import (
	"fmt"
)

func (s *PostgresStorage) InitStatusStorage() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS SellStatus(
			status TEXT
		);
	`)
	fmt.Println("sellStatus table created")
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) AddSellStatus(status string) error {

	_, err := s.db.Exec(`
		INSERT INTO SellStatus(status)
		VALUES($1)
		ON CONFLICT (status)
		DO UPDATE SET status = EXCLUDED.status`, status)
	if err != nil {
		return fmt.Errorf("error inserting status into database: %v", err)
	}
	fmt.Println("Sell status added or updated")
	return nil
}

func (s *PostgresStorage) UpdateSellStatus(status string) error {
	_, err := s.db.Exec(`
		UPDATE SellStatus
		SET status=$1`, status)
	if err != nil {
		return fmt.Errorf("error updating status: %v", err)
	}
	fmt.Println("Sell status updated")
	return nil
}
