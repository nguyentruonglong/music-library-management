package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"music-library-management/api/controllers"
	"music-library-management/api/routes"
	"music-library-management/api/services"
	"music-library-management/api/utils"
	"music-library-management/config"
)

func main() {
	// Set Gin mode based on environment
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode) // Set Gin to release mode in production
	}

	// Load configuration
	cfg, err := config.LoadConfig() // Load the application configuration
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err) // Log and exit if there is an error loading the configuration
	}

	// Connect to MongoDB
	client, err := utils.ConnectDB(cfg) // Connect to the MongoDB database
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err) // Log and exit if there is an error connecting to MongoDB
	}

	// Get MongoDB database
	db := utils.GetDatabase(client, cfg) // Get the MongoDB database instance

	// Initialize collections
	err = utils.InitializeCollections(db) // Initialize the MongoDB collections
	if err != nil {
		log.Fatalf("Error initializing collections: %v", err) // Log and exit if there is an error initializing collections
	}

	// Seed genres collection with sample data
	utils.SeedGenres(client, cfg) // Seed the genres collection with sample data

	// Initialize Gin router
	router := gin.Default() // Create a new Gin router with default settings

	// Set trusted proxies
	router.SetTrustedProxies(nil) // Set trusted proxies to nil

	// Middleware
	// router.Use(middleware.CORSMiddleware()) // Uncomment this line to enable CORS middleware

	// Serve static files from the uploads directory
	router.Static("/uploads", "./uploads") // Serve static files from the "uploads" directory

	// Initialize services and controllers
	fileService := services.NewFileService(client, cfg)          // Create a new FileService instance
	fileController := controllers.NewFileController(fileService) // Create a new FileController instance

	trackService := services.NewTrackService(client, cfg)                        // Create a new TrackService instance
	trackController := controllers.NewTrackController(trackService, fileService) // Create a new TrackController instance

	playlistService := services.NewPlaylistService(client, cfg, trackService) // Create a new PlaylistService instance
	playlistController := controllers.NewPlaylistController(playlistService)  // Create a new PlaylistController instance

	genreService := services.NewGenreService(client, cfg)           // Create a new GenreService instance
	genreController := controllers.NewGenreController(genreService) // Create a new GenreController instance

	searchService := services.NewSearchService(client, cfg)            // Create a new SearchService instance
	searchController := controllers.NewSearchController(searchService) // Create a new SearchController instance

	// Initialize routes
	routes.FileRoutes(router, fileController)         // Initialize file routes
	routes.TrackRoutes(router, trackController)       // Initialize track routes
	routes.PlaylistRoutes(router, playlistController) // Initialize playlist routes
	routes.GenreRoutes(router, genreController)       // Initialize genre routes
	routes.SearchRoutes(router, searchController)     // Initialize search routes

	// Start the server
	log.Fatal(router.Run(":" + cfg.Port)) // Start the Gin server on the configured port and log any fatal errors
}
