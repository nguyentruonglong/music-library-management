package controllers

import (
	"music-library-management/api/models"
	"music-library-management/api/services"
	"music-library-management/api/utils"
	"music-library-management/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GenreController handles HTTP requests for genres
type GenreController struct {
	genreService *services.GenreService // A reference to the genre service
}

// NewGenreController creates a new GenreController
func NewGenreController(genreService *services.GenreService) *GenreController {
	return &GenreController{
		genreService: genreService, // Initialize the genre service
	}
}

// AddGenreInput represents the input data for adding a new genre
type AddGenreInput struct {
	Name string `json:"name" binding:"required"` // The name of the genre, required field
}

// UpdateGenreInput represents the input data for updating a genre
type UpdateGenreInput struct {
	Name string `json:"name" binding:"required"` // The updated name of the genre, required field
}

// ListGenresInput represents the input data for listing genres
type ListGenresInput struct {
	Page  int `form:"page"`  // The page number for pagination
	Limit int `form:"limit"` // The number of items per page for pagination
}

// GenreOutput represents the output data for a genre
type GenreOutput struct {
	ID   string `json:"id"`   // The ID of the genre
	Name string `json:"name"` // The name of the genre
}

// PaginatedGenresOutput represents the output data for paginated genres
type PaginatedGenresOutput struct {
	Page       int           `json:"page"`        // The current page number
	Limit      int           `json:"limit"`       // The number of items per page
	TotalCount int64         `json:"total_count"` // The total number of genres
	Genres     []GenreOutput `json:"genres"`      // The list of genres
}

// AddGenre handles adding a new genre
func (gc *GenreController) AddGenre(c *gin.Context) {
	var input AddGenreInput
	var genre models.Genre

	// Bind JSON input to the AddGenreInput struct
	if err := c.ShouldBindJSON(&input); err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput) // Handle binding errors
		return
	}

	// Copy input data to genre model
	genre.Name = input.Name

	// Call service to add the genre
	createdGenre, err := gc.genreService.AddGenre(&genre)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err) // Handle errors from the service
		return
	}

	// Prepare output data
	output := GenreOutput{
		ID:   createdGenre.ID.Hex(),
		Name: createdGenre.Name,
	}

	// Respond with success message and created genre
	response := utils.NewSuccessResponse("Genre added successfully", output)
	c.JSON(http.StatusCreated, response)
}

// GetGenre handles retrieving a genre by ID
func (gc *GenreController) GetGenre(c *gin.Context) {
	genreId := c.Param("genreId") // Get the genre ID from the URL parameter

	// Call service to get the genre
	genre, err := gc.genreService.GetGenre(genreId)
	if err != nil {
		errors.HandleError(c, http.StatusNotFound, err) // Handle errors if the genre is not found
		return
	}

	// Prepare output data
	output := GenreOutput{
		ID:   genre.ID.Hex(),
		Name: genre.Name,
	}

	// Respond with success message and retrieved genre
	response := utils.NewSuccessResponse("Genre retrieved successfully", output)
	c.JSON(http.StatusOK, response)
}

// UpdateGenre handles updating an existing genre
func (gc *GenreController) UpdateGenre(c *gin.Context) {
	genreId := c.Param("genreId") // Get the genre ID from the URL parameter

	// Check if the genre exists
	_, err := gc.genreService.GetGenre(genreId)
	if err != nil {
		errors.HandleError(c, http.StatusNotFound, err) // Handle errors if the genre is not found
		return
	}

	var input UpdateGenreInput
	var updatedGenre models.Genre

	// Bind JSON input to the UpdateGenreInput struct
	if err := c.ShouldBindJSON(&input); err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput) // Handle binding errors
		return
	}

	// Copy input data to updatedGenre model
	updatedGenre.Name = input.Name

	// Call service to update the genre
	genre, err := gc.genreService.UpdateGenre(genreId, &updatedGenre)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err) // Handle errors from the service
		return
	}

	// Prepare output data
	output := GenreOutput{
		ID:   genre.ID.Hex(),
		Name: genre.Name,
	}

	// Respond with success message and updated genre
	response := utils.NewSuccessResponse("Genre updated successfully", output)
	c.JSON(http.StatusOK, response)
}

// DeleteGenre handles deleting a genre
func (gc *GenreController) DeleteGenre(c *gin.Context) {
	genreId := c.Param("genreId") // Get the genre ID from the URL parameter

	// Call service to delete the genre
	err := gc.genreService.DeleteGenre(genreId)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err) // Handle errors from the service
		return
	}

	// Respond with success message
	response := utils.NewSuccessResponse("Genre deleted successfully", nil)
	c.JSON(http.StatusOK, response)
}

// ListGenres handles listing all genres with pagination
func (gc *GenreController) ListGenres(c *gin.Context) {
	var input ListGenresInput

	// Bind query parameters to ListGenresInput struct
	if err := c.ShouldBindQuery(&input); err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput) // Handle binding errors
		return
	}

	// Set default pagination values if not provided
	if input.Page == 0 {
		input.Page = 1
	}
	if input.Limit == 0 {
		input.Limit = 100
	}

	// Call service to list genres
	genres, totalCount, err := gc.genreService.ListGenres(input.Page, input.Limit)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err) // Handle errors from the service
		return
	}

	// Prepare output data
	output := PaginatedGenresOutput{
		Page:       input.Page,
		Limit:      input.Limit,
		TotalCount: totalCount,
		Genres:     make([]GenreOutput, len(genres)), // Initialize the genres slice with the appropriate length
	}

	// Populate the output genres
	for i, genre := range genres {
		output.Genres[i] = GenreOutput{
			ID:   genre.ID.Hex(),
			Name: genre.Name,
		}
	}

	// Respond with success message and list of genres
	response := utils.NewSuccessResponse("Genres retrieved successfully", output)
	c.JSON(http.StatusOK, response)
}
