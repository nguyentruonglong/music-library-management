package routes

import (
	"music-library-management/api/controllers"

	"github.com/gin-gonic/gin"
)

// GenreRoutes sets up the genre-related routes
func GenreRoutes(router *gin.Engine, genreController *controllers.GenreController) {
	// Group genre routes
	genres := router.Group("/api/genres")
	{
		// Add a new genre
		genres.POST("/", genreController.AddGenre)
		// List all genres
		genres.GET("/", genreController.ListGenres)
		// Retrieve a genre by ID
		genres.GET("/:genreId", genreController.GetGenre)
		// Update a genre by ID
		genres.PUT("/:genreId", genreController.UpdateGenre)
		// Delete a genre by ID
		genres.DELETE("/:genreId", genreController.DeleteGenre)
	}
}
