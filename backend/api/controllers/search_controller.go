package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"music-library-management/api/services"
	"music-library-management/api/utils"
	"music-library-management/errors"
)

// SearchController handles search-related HTTP requests
type SearchController struct {
	searchService *services.SearchService // A reference to the search service
}

// NewSearchController creates a new SearchController
func NewSearchController(searchService *services.SearchService) *SearchController {
	return &SearchController{
		searchService: searchService, // Initialize the search service
	}
}

// SearchTracksInput represents the input data for searching tracks
type SearchTracksInput struct {
	Query string `form:"query" binding:"required"` // The search query string
	Page  int    `form:"page"`                     // The page number for pagination
	Limit int    `form:"limit"`                    // The number of items per page for pagination
}

// SearchTracksOutput represents the output data for searching tracks
type SearchTracksOutput struct {
	Page   int           `json:"page"`   // The current page number
	Limit  int           `json:"limit"`  // The number of items per page
	Total  int64         `json:"total"`  // The total number of matching tracks
	Tracks []TrackOutput `json:"tracks"` // The list of matching tracks
}

// SearchPlaylistsInput represents the input data for searching playlists
type SearchPlaylistsInput struct {
	Query string `form:"query" binding:"required"` // The search query string
	Page  int    `form:"page"`                     // The page number for pagination
	Limit int    `form:"limit"`                    // The number of items per page for pagination
}

// SearchPlaylistsOutput represents the output data for searching playlists
type SearchPlaylistsOutput struct {
	Page      int              `json:"page"`      // The current page number
	Limit     int              `json:"limit"`     // The number of items per page
	Total     int64            `json:"total"`     // The total number of matching playlists
	Playlists []PlaylistOutput `json:"playlists"` // The list of matching playlists
}

// SearchTracks handles searching for tracks based on the query
func (sc *SearchController) SearchTracks(c *gin.Context) {
	// Parse query parameters
	var input SearchTracksInput
	if err := c.ShouldBindQuery(&input); err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput) // Handle binding errors
		return
	}

	// Set default pagination values if not provided
	if input.Page == 0 {
		input.Page = 1
	}
	if input.Limit == 0 {
		input.Limit = 10
	}

	// Call the search service to search tracks
	tracks, total, err := sc.searchService.SearchTracks(input.Query, input.Page, input.Limit)
	if err != nil {
		// Handle any errors that occur during the search
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	// Prepare the response data
	output := SearchTracksOutput{
		Page:   input.Page,
		Limit:  input.Limit,
		Total:  total,
		Tracks: make([]TrackOutput, len(tracks)), // Initialize the tracks slice with the appropriate length
	}

	// Populate the output tracks
	for i, track := range tracks {
		output.Tracks[i] = TrackOutput{
			ID:            track.ID.Hex(),
			Title:         track.Title,
			Artist:        track.Artist,
			Album:         track.Album,
			Genre:         track.Genre,
			ReleaseYear:   track.ReleaseYear,
			Duration:      track.Duration,
			CoverImageUrl: track.CoverImageUrl,
			Mp3FileUrl:    track.Mp3FileUrl,
		}
	}

	// Create a success response
	response := utils.NewSuccessResponse("Tracks retrieved successfully", output)

	// Send the response
	c.JSON(http.StatusOK, response)
}

// SearchPlaylists handles searching for playlists based on the query
func (sc *SearchController) SearchPlaylists(c *gin.Context) {
	// Parse query parameters
	var input SearchPlaylistsInput
	if err := c.ShouldBindQuery(&input); err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput) // Handle binding errors
		return
	}

	// Set default pagination values if not provided
	if input.Page == 0 {
		input.Page = 1
	}
	if input.Limit == 0 {
		input.Limit = 10
	}

	// Call the search service to search playlists
	playlists, total, err := sc.searchService.SearchPlaylists(input.Query, input.Page, input.Limit)
	if err != nil {
		// Handle any errors that occur during the search
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	// Prepare the response data
	output := SearchPlaylistsOutput{
		Page:      input.Page,
		Limit:     input.Limit,
		Total:     total,
		Playlists: make([]PlaylistOutput, len(playlists)), // Initialize the playlists slice with the appropriate length
	}

	// Populate the output playlists
	for i, playlist := range playlists {
		output.Playlists[i] = PlaylistOutput{
			ID:   playlist.ID.Hex(),
			Name: playlist.Name,
		}
	}

	// Create a success response
	response := utils.NewSuccessResponse("Playlists retrieved successfully", output)

	// Send the response
	c.JSON(http.StatusOK, response)
}
