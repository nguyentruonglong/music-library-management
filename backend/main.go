package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"music-library-management/api/middleware"
	"music-library-management/api/utils"
	"music-library-management/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Connect to MongoDB
	client, err := utils.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Get MongoDB database
	db := utils.GetDatabase(client, cfg)

	// Initialize collections
	err = utils.InitializeCollections(db)
	if err != nil {
		log.Fatalf("Error initializing collections: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware())

	// Start the server
	log.Fatal(router.Run(":" + cfg.Port))
}
