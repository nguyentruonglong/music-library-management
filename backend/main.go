package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"music-library-management/api/controllers"
	// "music-library-management/api/middleware"
	"music-library-management/api/routes"
	"music-library-management/api/services"
	"music-library-management/api/utils"
	"music-library-management/config"
)

func main() {
	// Set Gin mode based on environment
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

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

	// Set trusted proxies
	router.SetTrustedProxies(nil)

	// Middleware
	// router.Use(middleware.CORSMiddleware())

	// Initialize services and controllers
	trackService := services.NewTrackService(client, cfg)
	trackController := controllers.NewTrackController(trackService)

	playlistService := services.NewPlaylistService(client, cfg, trackService)
	playlistController := controllers.NewPlaylistController(playlistService)

	// Initialize routes
	routes.TrackRoutes(router, trackController)
	routes.PlaylistRoutes(router, playlistController)

	// Start the server
	log.Fatal(router.Run(":" + cfg.Port))
}
