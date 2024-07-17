package routes

import (
	"music-library-management/api/controllers"

	"github.com/gin-gonic/gin"
)

// PlaylistRoutes sets up the routes for the playlist-related endpoints
func PlaylistRoutes(router *gin.Engine, playlistController *controllers.PlaylistController) {
	// Group playlist routes
	playlistRoutes := router.Group("/api/playlists")
	{
		// Add a new playlist
		playlistRoutes.POST("/", playlistController.AddPlaylist)

		// View details of a specific playlist
		playlistRoutes.GET("/:playlistId", playlistController.GetPlaylist)

		// Update an existing playlist
		playlistRoutes.PUT("/:playlistId", playlistController.UpdatePlaylist)

		// Delete a playlist
		playlistRoutes.DELETE("/:playlistId", playlistController.DeletePlaylist)

		// List all playlists with pagination
		playlistRoutes.GET("/", playlistController.ListPlaylists)

		// Add a track to a playlist
		playlistRoutes.POST("/:playlistId/tracks/:trackId", playlistController.AddTrackToPlaylist)

		// Remove a track from a playlist
		playlistRoutes.DELETE("/:playlistId/tracks/:trackId", playlistController.RemoveTrackFromPlaylist)
	}
}
