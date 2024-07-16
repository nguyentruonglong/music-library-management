package routes

import (
	"music-library-management/api/controllers"

	"github.com/gin-gonic/gin"
)

// TrackRoutes sets up the routes for the track-related endpoints
func TrackRoutes(router *gin.Engine, trackController *controllers.TrackController) {
	// Group track routes
	trackRoutes := router.Group("/api/tracks")
	{
		// Add a new music track with a cover image and MP3 file
		trackRoutes.POST("/", trackController.AddTrack)

		// View details of a specific music track
		trackRoutes.GET("/:id", trackController.GetTrack)

		// Update an existing music track, including the cover image
		trackRoutes.PUT("/:id", trackController.UpdateTrack)

		// Delete a music track
		trackRoutes.DELETE("/:id", trackController.DeleteTrack)

		// List all music tracks with pagination
		trackRoutes.GET("/", trackController.ListTracks)

		// Play/Pause an MP3 file of a music track
		trackRoutes.POST("/:id/play", trackController.PlayPauseTrack)
	}
}
