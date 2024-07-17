package routes

import (
	"music-library-management/api/controllers"

	"github.com/gin-gonic/gin"
)

// SearchRoutes sets up the routes for the search-related endpoints
func SearchRoutes(router *gin.Engine, searchController *controllers.SearchController) {
	// Group search routes
	search := router.Group("/api/search")
	{
		// Search all tracks with pagination
		search.GET("/tracks", searchController.SearchTracks)

		// Search all playlists with pagination
		search.GET("/playlists", searchController.SearchPlaylists)
	}
}
