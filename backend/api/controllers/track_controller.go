package controllers

import (
	"music-library-management/api/models"
	"music-library-management/api/services"
	"music-library-management/api/utils"
	"music-library-management/errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TrackController handles HTTP requests for tracks
type TrackController struct {
	trackService *services.TrackService // A reference to the track service
	fileService  *services.FileService  // A reference to the file service
}

// NewTrackController creates a new TrackController
func NewTrackController(trackService *services.TrackService, fileService *services.FileService) *TrackController {
	return &TrackController{
		trackService: trackService, // Initialize the track service
		fileService:  fileService,  // Initialize the file service
	}
}

// AddTrackInput represents the input data for adding a new track
type AddTrackInput struct {
	Title       string `form:"title" binding:"required"`    // The title of the track, required field
	Artist      string `form:"artist" binding:"required"`   // The artist of the track, required field
	Album       string `form:"album"`                       // The album of the track
	Genre       string `form:"genre"`                       // The genre of the track
	ReleaseYear int    `form:"release_year"`                // The release year of the track
	Duration    int    `form:"duration" binding:"required"` // The duration of the track, required field
}

// UpdateTrackInput represents the input data for updating a track
type UpdateTrackInput struct {
	Title       string `form:"title"`        // The updated title of the track
	Artist      string `form:"artist"`       // The updated artist of the track
	Album       string `form:"album"`        // The updated album of the track
	Genre       string `form:"genre"`        // The updated genre of the track
	ReleaseYear int    `form:"release_year"` // The updated release year of the track
	Duration    int    `form:"duration"`     // The updated duration of the track
}

// ListTracksInput represents the input data for listing tracks
type ListTracksInput struct {
	Page  int `form:"page"`  // The page number for pagination
	Limit int `form:"limit"` // The number of items per page for pagination
}

// PlayPauseTrackInput represents the input data for playing or pausing a track
type PlayPauseTrackInput struct {
	Action string `json:"action" binding:"required,oneof=play pause"` // The action to perform, required field with validation
}

// TrackOutput represents the output data for a track
type TrackOutput struct {
	ID            string `json:"id"`              // The ID of the track
	Title         string `json:"title"`           // The title of the track
	Artist        string `json:"artist"`          // The artist of the track
	Album         string `json:"album"`           // The album of the track
	Genre         string `json:"genre"`           // The genre of the track
	ReleaseYear   int    `json:"release_year"`    // The release year of the track
	Duration      int    `json:"duration"`        // The duration of the track
	CoverImageUrl string `json:"cover_image_url"` // The URL of the cover image
	Mp3FileUrl    string `json:"mp3_file_url"`    // The URL of the MP3 file
}

// PaginatedTracksOutput represents the output data for paginated tracks
type PaginatedTracksOutput struct {
	Page       int           `json:"page"`        // The current page number
	Limit      int           `json:"limit"`       // The number of items per page
	TotalCount int64         `json:"total_count"` // The total number of tracks
	Tracks     []TrackOutput `json:"tracks"`      // The list of tracks
}

// AddTrack handles adding a new track
func (tc *TrackController) AddTrack(c *gin.Context) {
	var input AddTrackInput // Declare a variable to hold the input data
	var track models.Track  // Declare a variable to hold the track model

	// Parse multipart form data
	err := c.Request.ParseMultipartForm(10 << 20) // Limit to 10 MB
	if err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput) // Handle errors if form data is invalid
		return
	}

	// Bind form data to input struct
	if err := c.ShouldBind(&input); err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput) // Handle errors if binding fails
		return
	}

	// Copy input data to track model
	track.Title = input.Title
	track.Artist = input.Artist
	track.Album = input.Album
	track.Genre = input.Genre
	track.ReleaseYear = input.ReleaseYear
	track.Duration = input.Duration

	// Handle cover image upload
	coverImage, err := c.FormFile("cover_image") // Get the cover image file
	if err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput) // Handle errors if the cover image is not provided
		return
	}
	coverImageExt := filepath.Ext(coverImage.Filename)                      // Get the file extension
	coverImageName := uuid.New().String() + coverImageExt                   // Generate a unique file name
	coverImagePath := tc.fileService.GetUploadPath() + "/" + coverImageName // Generate a unique file path
	if err := c.SaveUploadedFile(coverImage, coverImagePath); err != nil {  // Save the cover image file
		errors.HandleError(c, http.StatusInternalServerError, errors.ErrInternalServer) // Handle errors if saving fails
		return
	}
	track.CoverImageUrl = utils.GetScheme(c) + "://" + c.Request.Host + "/" + coverImagePath // Set the cover image URL

	// Save cover image metadata
	_, err = tc.fileService.SaveFileMetadata(c, coverImageName)
	if err != nil {
		os.Remove(coverImagePath) // Remove the uploaded cover image file
		errors.HandleError(c, http.StatusInternalServerError, errors.ErrInternalServer)
		return
	}

	// Handle MP3 file upload
	mp3File, err := c.FormFile("mp3_file") // Get the MP3 file
	if err != nil {
		os.Remove(coverImagePath)                                            // Remove the uploaded cover image file
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput) // Handle errors if the MP3 file is not provided
		return
	}
	mp3FileExt := filepath.Ext(mp3File.Filename)                      // Get the file extension
	mp3FileName := uuid.New().String() + mp3FileExt                   // Generate a unique file name
	mp3FilePath := tc.fileService.GetUploadPath() + "/" + mp3FileName // Generate a unique file path
	if err := c.SaveUploadedFile(mp3File, mp3FilePath); err != nil {  // Save the MP3 file
		os.Remove(coverImagePath)                                                       // Remove the uploaded cover image file
		errors.HandleError(c, http.StatusInternalServerError, errors.ErrInternalServer) // Handle errors if saving fails
		return
	}
	track.Mp3FileUrl = utils.GetScheme(c) + "://" + c.Request.Host + "/" + mp3FilePath // Set the MP3 file URL

	// Save MP3 file metadata
	_, err = tc.fileService.SaveFileMetadata(c, mp3FileName)
	if err != nil {
		os.Remove(coverImagePath) // Remove the uploaded cover image file
		os.Remove(mp3FilePath)    // Remove the uploaded MP3 file
		errors.HandleError(c, http.StatusInternalServerError, errors.ErrInternalServer)
		return
	}

	// Add track to the database
	createdTrack, err := tc.trackService.AddTrack(&track) // Call service to add the track
	if err != nil {
		os.Remove(coverImagePath)                                  // Remove the uploaded cover image file
		os.Remove(mp3FilePath)                                     // Remove the uploaded MP3 file
		errors.HandleError(c, http.StatusInternalServerError, err) // Handle errors from the service
		return
	}

	// Prepare output data
	output := TrackOutput{
		ID:            createdTrack.ID.Hex(),
		Title:         createdTrack.Title,
		Artist:        createdTrack.Artist,
		Album:         createdTrack.Album,
		Genre:         createdTrack.Genre,
		ReleaseYear:   createdTrack.ReleaseYear,
		Duration:      createdTrack.Duration,
		CoverImageUrl: createdTrack.CoverImageUrl,
		Mp3FileUrl:    createdTrack.Mp3FileUrl,
	}

	response := utils.NewSuccessResponse("Track added successfully", output) // Create a success response
	c.JSON(http.StatusCreated, response)                                     // Send the response
}

// GetTrack handles retrieving a track by its ID
func (tc *TrackController) GetTrack(c *gin.Context) {
	trackId := c.Param("trackId") // Get the track ID from the URL parameter

	track, err := tc.trackService.GetTrack(trackId) // Call service to get the track
	if err != nil {
		errors.HandleError(c, http.StatusNotFound, err) // Handle errors if the track is not found
		return
	}

	// Prepare output data
	output := TrackOutput{
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

	response := utils.NewSuccessResponse("Track retrieved successfully", output) // Create a success response
	c.JSON(http.StatusOK, response)                                              // Send the response
}

// UpdateTrack handles updating an existing track
func (tc *TrackController) UpdateTrack(c *gin.Context) {
	trackId := c.Param("trackId") // Get the track ID from the URL parameter

	// Check if the track exists
	_, err := tc.trackService.GetTrack(trackId) // Call service to check if the track exists
	if err != nil {
		errors.HandleError(c, http.StatusNotFound, err) // Handle errors if the track is not found
		return
	}

	// Parse multipart form data
	err = c.Request.ParseMultipartForm(10 << 20) // Limit to 10 MB
	if err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput) // Handle errors if form data is invalid
		return
	}

	var input UpdateTrackInput             // Declare a variable to hold the input data
	var updatedTrack models.Track          // Declare a variable to hold the updated track model
	var coverImagePath, mp3FilePath string // Variables to hold the file paths

	// Bind form data to input struct
	if err := c.ShouldBind(&input); err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput) // Handle errors if binding fails
		return
	}

	// Copy input data to updatedTrack model
	updatedTrack.Title = input.Title
	updatedTrack.Artist = input.Artist
	updatedTrack.Album = input.Album
	updatedTrack.Genre = input.Genre
	updatedTrack.ReleaseYear = input.ReleaseYear
	updatedTrack.Duration = input.Duration

	// Handle cover image upload
	coverImage, err := c.FormFile("cover_image") // Get the cover image file
	if err == nil {
		coverImageExt := filepath.Ext(coverImage.Filename)                     // Get the file extension
		coverImageName := uuid.New().String() + coverImageExt                  // Generate a unique file name
		coverImagePath = tc.fileService.GetUploadPath() + "/" + coverImageName // Generate a unique file path
		if err := c.SaveUploadedFile(coverImage, coverImagePath); err != nil { // Save the cover image file
			errors.HandleError(c, http.StatusInternalServerError, errors.ErrInternalServer) // Handle errors if saving fails
			return
		}
		updatedTrack.CoverImageUrl = utils.GetScheme(c) + "://" + c.Request.Host + "/" + coverImagePath // Set the cover image URL

		// Save cover image metadata
		_, err = tc.fileService.SaveFileMetadata(c, coverImageName)
		if err != nil {
			os.Remove(coverImagePath) // Remove the uploaded cover image file
			errors.HandleError(c, http.StatusInternalServerError, errors.ErrInternalServer)
			return
		}
	}

	// Handle MP3 file upload
	mp3File, err := c.FormFile("mp3_file") // Get the MP3 file
	if err == nil {
		mp3FileExt := filepath.Ext(mp3File.Filename)                     // Get the file extension
		mp3FileName := uuid.New().String() + mp3FileExt                  // Generate a unique file name
		mp3FilePath = tc.fileService.GetUploadPath() + "/" + mp3FileName // Generate a unique file path
		if err := c.SaveUploadedFile(mp3File, mp3FilePath); err != nil { // Save the MP3 file
			os.Remove(coverImagePath)                                                       // Remove the uploaded cover image file
			errors.HandleError(c, http.StatusInternalServerError, errors.ErrInternalServer) // Handle errors if saving fails
			return
		}
		updatedTrack.Mp3FileUrl = utils.GetScheme(c) + "://" + c.Request.Host + "/" + mp3FilePath // Set the MP3 file URL

		// Save MP3 file metadata
		_, err = tc.fileService.SaveFileMetadata(c, mp3FileName)
		if err != nil {
			os.Remove(coverImagePath) // Remove the uploaded cover image file
			os.Remove(mp3FilePath)    // Remove the uploaded MP3 file
			errors.HandleError(c, http.StatusInternalServerError, errors.ErrInternalServer)
			return
		}
	}

	// Update track in the database
	track, err := tc.trackService.UpdateTrack(trackId, &updatedTrack) // Call service to update the track
	if err != nil {
		os.Remove(coverImagePath)                                  // Remove the uploaded cover image file
		os.Remove(mp3FilePath)                                     // Remove the uploaded MP3 file
		errors.HandleError(c, http.StatusInternalServerError, err) // Handle errors from the service
		return
	}

	// Prepare output data
	output := TrackOutput{
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

	response := utils.NewSuccessResponse("Track updated successfully", output) // Create a success response
	c.JSON(http.StatusOK, response)                                            // Send the response
}

// DeleteTrack handles deleting a track
func (tc *TrackController) DeleteTrack(c *gin.Context) {
	trackId := c.Param("trackId") // Get the track ID from the URL parameter

	err := tc.trackService.DeleteTrack(trackId) // Call service to delete the track
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err) // Handle errors from the service
		return
	}

	response := utils.NewSuccessResponse("Track deleted successfully", nil) // Create a success response
	c.JSON(http.StatusOK, response)                                         // Send the response
}

// ListTracks handles listing all tracks with pagination
func (tc *TrackController) ListTracks(c *gin.Context) {
	var input ListTracksInput // Declare a variable to hold the input data

	// Bind query parameters to input struct
	if err := c.ShouldBindQuery(&input); err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput) // Handle errors if binding fails
		return
	}

	if input.Page == 0 {
		input.Page = 1
	}
	if input.Limit == 0 {
		input.Limit = 10
	}

	// Call service to list tracks
	tracks, totalCount, err := tc.trackService.ListTracks(input.Page, input.Limit) // Call service to list tracks
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err) // Handle errors from the service
		return
	}

	// Prepare output data
	output := PaginatedTracksOutput{
		Page:       input.Page,
		Limit:      input.Limit,
		TotalCount: totalCount,
		Tracks:     make([]TrackOutput, len(tracks)), // Initialize the tracks slice with the appropriate length
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

	response := utils.NewSuccessResponse("Tracks retrieved successfully", output) // Create a success response
	c.JSON(http.StatusOK, response)                                               // Send the response
}

// PlayPauseTrack handles playing or pausing a track
func (tc *TrackController) PlayPauseTrack(c *gin.Context) {
	trackId := c.Param("trackId") // Get the track ID from the URL parameter
	var input PlayPauseTrackInput // Declare a variable to hold the input data

	// Bind JSON input to struct
	if err := c.ShouldBindJSON(&input); err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput) // Handle errors if binding fails
		return
	}

	err := tc.trackService.PlayPauseTrack(trackId, input.Action) // Call service to play or pause the track
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err) // Handle errors from the service
		return
	}

	response := utils.NewSuccessResponse("Track action performed successfully", nil) // Create a success response
	c.JSON(http.StatusOK, response)                                                  // Send the response
}
