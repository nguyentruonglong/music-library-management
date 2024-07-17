package utils

// SuccessResponse defines the structure for a successful API response
type SuccessResponse struct {
	Status  string      `json:"status"`  // Status of the response, typically "success"
	Message string      `json:"message"` // Message providing additional information about the response
	Data    interface{} `json:"data"`    // Data payload of the response, can be of any type
}

// NewSuccessResponse creates a new SuccessResponse with the given message and data
func NewSuccessResponse(message string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Status:  "success", // Set the status to "success"
		Message: message,   // Set the message to the provided message
		Data:    data,      // Set the data to the provided data
	}
}
