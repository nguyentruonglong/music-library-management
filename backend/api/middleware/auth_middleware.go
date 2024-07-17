package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles CORS-related headers and ensures the client is allowed to access the resource
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set the Access-Control-Allow-Origin header to allow all origins (*)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// Set the Access-Control-Allow-Credentials header to true, allowing credentials to be included in the requests
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// Set the Access-Control-Allow-Headers header to specify which headers can be used during the actual request
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// Set the Access-Control-Allow-Methods header to specify the methods allowed when accessing the resource
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// If the request method is OPTIONS, abort the request with a 204 No Content status
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Proceed to the next middleware or handler
		c.Next()
	}
}
