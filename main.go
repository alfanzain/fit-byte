package main

import (
	"fit-byte/config"
	"fit-byte/db"
	"fit-byte/routes"
	"fmt"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	db.InitDB(cfg)
	defer func() {
		if err := db.DB.Close(); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
		log.Println("Database connection closed.")
	}()

	r := routes.SetupRouter(cfg, db.DB)

	fmt.Printf("Starting server on port %s...\n", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
