package main

import (
	"fit-byte/config"
	"fit-byte/db"
	v1Handlers "fit-byte/handlers/v1"
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

	s3Client := v1Handlers.InitS3Client()
	r := routes.SetupRouter(cfg, db.DB, s3Client)

	fmt.Printf("Starting server on port %s...\n", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
