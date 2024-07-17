package routes

import (
	"music-library-management/api/controllers"

	"github.com/gin-gonic/gin"
)

// FileRoutes sets up the routes for the file-related endpoints
func FileRoutes(router *gin.Engine, fileController *controllers.FileController) {
	// Group file routes
	files := router.Group("/api/files")
	{
		// List all files
		files.GET("/", fileController.ListFiles)
	}
}
