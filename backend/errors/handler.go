package errors

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standard API error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// HandleError is a utility function to handle errors in a standardized way
func HandleError(c *gin.Context, code int, err error) {
	c.JSON(code, ErrorResponse{
		Code:    code,
		Message: err.Error(),
	})
}
