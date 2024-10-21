package main

import (
	"AahaFeltBackend/api"
	"AahaFeltBackend/storage"
	"fmt"
)

func main() {
	db, err := storage.NewPostgresStorage()
	if err != nil {
		fmt.Println("Failed to initialize storage:", err)
		return
	}
	defer db.Close()

	// Initialize the database (create tables)
	if err := db.Init(); err != nil {
		fmt.Println("Failed to initialize database:", err)
		return
	}
	db.InitGallery()
	db.InitProductImage()
	server := api.NewApiServer(":3000", db)

	server.Start()
}
