package errors

import "errors"

// Common error messages
var (
	ErrNotFound               = errors.New("resource not found")                   // Error when a requested resource is not found
	ErrInternalServer         = errors.New("internal server error")                // Error for general server issues
	ErrBadRequest             = errors.New("bad request")                          // Error for invalid client requests
	ErrUnauthorized           = errors.New("unauthorized")                         // Error for unauthorized access
	ErrInvalidObjectID        = errors.New("invalid ID format")                    // Error for invalid ID format
	ErrDatabaseOperation      = errors.New("database operation failed")            // Error for database operation failures
	ErrTrackNotFound          = errors.New("track not found")                      // Error when a track is not found
	ErrPlaylistNotFound       = errors.New("playlist not found")                   // Error when a playlist is not found
	ErrTrackAlreadyInPlaylist = errors.New("track already exists in the playlist") // Error when a track is already in a playlist
	ErrTrackNotInPlaylist     = errors.New("track does not exist in the playlist") // Error when a track is not in a playlist
	ErrInvalidInput           = errors.New("invalid input")                        // Error for invalid input
)

// CustomError represents a custom error type
type CustomError struct {
	Message string // Message holds the custom error message
}

func (e *CustomError) Error() string {
	return e.Message // Returns the custom error message
}

// NewError creates a new error with a given message
func NewError(message string) error {
	return &CustomError{message} // Returns a new custom error
}
