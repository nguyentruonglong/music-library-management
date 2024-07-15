package errors

import "errors"

// Common error messages
var (
	ErrNotFound       = errors.New("resource not found")
	ErrInternalServer = errors.New("internal server error")
	ErrBadRequest     = errors.New("bad request")
	ErrUnauthorized   = errors.New("unauthorized")
)

// CustomError represents a custom error type
type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

// NewError creates a new error with a given message
func NewError(message string) error {
	return &CustomError{message}
}
