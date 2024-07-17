package controllers

import (
	"music-library-management/api/services"
	"music-library-management/api/utils"
	"music-library-management/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FileController handles file-related HTTP requests
type FileController struct {
	fileService *services.FileService // A reference to the file service
}

// NewFileController creates a new FileController
func NewFileController(fileService *services.FileService) *FileController {
	return &FileController{
		fileService: fileService, // Initialize the file service
	}
}

// ListFilesInput represents the input data for listing files
type ListFilesInput struct {
	Page  int `form:"page"`  // The page number for pagination
	Limit int `form:"limit"` // The number of items per page for pagination
}

// FileOutput represents the output data for a file
type FileOutput struct {
	ID       string `json:"id"`       // The ID of the file
	Filename string `json:"filename"` // The filename of the file
	Filepath string `json:"filepath"` // The filepath of the file
}

// PaginatedFilesOutput represents the output data for paginated files
type PaginatedFilesOutput struct {
	Page  int          `json:"page"`  // The current page number
	Limit int          `json:"limit"` // The number of items per page
	Total int64        `json:"total"` // The total number of files
	Files []FileOutput `json:"files"` // The list of files
}

// ListFiles handles listing uploaded files
func (fc *FileController) ListFiles(c *gin.Context) {
	var input ListFilesInput // Declare a variable to hold the input data

	// Bind query parameters to input struct
	if err := c.ShouldBindQuery(&input); err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput) // Handle errors if binding fails
		return
	}

	// Set default values for pagination if not provided
	if input.Page == 0 {
		input.Page = 1
	}
	if input.Limit == 0 {
		input.Limit = 10
	}

	// Retrieve files and total count from the file service
	files, total, err := fc.fileService.ListFiles(input.Page, input.Limit)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err) // Handle errors from the service
		return
	}

	// Prepare the output data
	fileOutputs := make([]FileOutput, len(files))
	for i, file := range files {
		fileOutputs[i] = FileOutput{
			ID:       file.ID.Hex(),
			Filename: file.Filename,
			Filepath: file.Filepath,
		}
	}

	// Create a success response with the paginated files
	response := utils.NewSuccessResponse("Files retrieved successfully", PaginatedFilesOutput{
		Page:  input.Page,
		Limit: input.Limit,
		Total: total,
		Files: fileOutputs,
	})

	// Send the response as JSON
	c.JSON(http.StatusOK, response)
}
