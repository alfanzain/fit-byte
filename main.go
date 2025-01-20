package main

import (
	"fit-byte/config"
	"fit-byte/db"
	v1Handlers "fit-byte/handlers/v1"
	"fit-byte/routes"
	"fit-byte/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	// Initialize and register Prometheus metrics
	utils.InitPrometheusMetrics()

	// Simulate metric updates in the background
	utils.SimulateMetrics()

	// Add Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	fmt.Printf("Starting server on port %s...\n", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
