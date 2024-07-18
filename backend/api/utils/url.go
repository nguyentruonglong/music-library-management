package utils

import "github.com/gin-gonic/gin"

// GetScheme determines the scheme (http or https) based on the request
func GetScheme(c *gin.Context) string {
	if c.Request.TLS != nil {
		return "https"
	}
	return "http"
}
