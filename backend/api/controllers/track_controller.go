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

type TrackController struct {
	trackService *services.TrackService
}

// NewTrackController creates a new TrackController
func NewTrackController(trackService *services.TrackService) *TrackController {
	return &TrackController{
		trackService: trackService,
	}
}

// AddTrackInput represents the input data for adding a new track
type AddTrackInput struct {
	Title       string `form:"title" binding:"required"`
	Artist      string `form:"artist" binding:"required"`
	Album       string `form:"album"`
	Genre       string `form:"genre"`
	ReleaseYear int    `form:"release_year"`
	Duration    int    `form:"duration" binding:"required"`
}

// UpdateTrackInput represents the input data for updating a track
type UpdateTrackInput struct {
	Title       string `form:"title"`
	Artist      string `form:"artist"`
	Album       string `form:"album"`
	Genre       string `form:"genre"`
	ReleaseYear int    `form:"release_year"`
	Duration    int    `form:"duration"`
}

// ListTracksInput represents the input data for listing tracks
type ListTracksInput struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

// PlayPauseTrackInput represents the input data for playing or pausing a track
type PlayPauseTrackInput struct {
	Action string `json:"action" binding:"required,oneof=play pause"`
}

// TrackOutput represents the output data for a track
type TrackOutput struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Artist        string `json:"artist"`
	Album         string `json:"album"`
	Genre         string `json:"genre"`
	ReleaseYear   int    `json:"release_year"`
	Duration      int    `json:"duration"`
	CoverImageUrl string `json:"cover_image_url"`
	Mp3FileUrl    string `json:"mp3_file_url"`
}

// PaginatedTracksOutput represents the output data for paginated tracks
type PaginatedTracksOutput struct {
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalCount int64         `json:"total_count"`
	Tracks     []TrackOutput `json:"tracks"`
}

// AddTrack handles adding a new track
func (tc *TrackController) AddTrack(c *gin.Context) {
	var input AddTrackInput
	var track models.Track

	// Parse multipart form data
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Bind form data to input struct
	if err := c.ShouldBind(&input); err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput)
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
	coverImage, err := c.FormFile("cover_image")
	if err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}
	coverImageExt := filepath.Ext(coverImage.Filename)
	coverImagePath := "uploads/" + uuid.New().String() + coverImageExt
	if err := c.SaveUploadedFile(coverImage, coverImagePath); err != nil {
		errors.HandleError(c, http.StatusInternalServerError, errors.ErrInternalServer)
		return
	}
	track.CoverImageUrl = c.Request.Host + "/" + coverImagePath

	// Handle MP3 file upload
	mp3File, err := c.FormFile("mp3_file")
	if err != nil {
		os.Remove(coverImagePath)
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}
	mp3FileExt := filepath.Ext(mp3File.Filename)
	mp3FilePath := "uploads/" + uuid.New().String() + mp3FileExt
	if err := c.SaveUploadedFile(mp3File, mp3FilePath); err != nil {
		os.Remove(coverImagePath)
		errors.HandleError(c, http.StatusInternalServerError, errors.ErrInternalServer)
		return
	}
	track.Mp3FileUrl = c.Request.Host + "/" + mp3FilePath

	// Add track to the database
	createdTrack, err := tc.trackService.AddTrack(&track)
	if err != nil {
		os.Remove(coverImagePath)
		os.Remove(mp3FilePath)
		errors.HandleError(c, http.StatusInternalServerError, err)
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

	response := utils.NewSuccessResponse("Track added successfully", output)
	c.JSON(http.StatusCreated, response)
}

// GetTrack handles retrieving a track by its ID
func (tc *TrackController) GetTrack(c *gin.Context) {
	trackId := c.Param("trackId")

	track, err := tc.trackService.GetTrack(trackId)
	if err != nil {
		errors.HandleError(c, http.StatusNotFound, err)
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

	response := utils.NewSuccessResponse("Track retrieved successfully", output)
	c.JSON(http.StatusOK, response)
}

// UpdateTrack handles updating an existing track
func (tc *TrackController) UpdateTrack(c *gin.Context) {
	trackId := c.Param("trackId")

	// Check if the track exists
	_, err := tc.trackService.GetTrack(trackId)
	if err != nil {
		errors.HandleError(c, http.StatusNotFound, err)
		return
	}

	// Parse multipart form data
	err = c.Request.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	var input UpdateTrackInput
	var updatedTrack models.Track
	var coverImagePath, mp3FilePath string

	// Bind form data to input struct
	if err := c.ShouldBind(&input); err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput)
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
	coverImage, err := c.FormFile("cover_image")
	if err == nil {
		coverImageExt := filepath.Ext(coverImage.Filename)
		coverImagePath = "uploads/" + uuid.New().String() + coverImageExt
		if err := c.SaveUploadedFile(coverImage, coverImagePath); err != nil {
			errors.HandleError(c, http.StatusInternalServerError, errors.ErrInternalServer)
			return
		}
		updatedTrack.CoverImageUrl = c.Request.Host + "/" + coverImagePath
	}

	// Handle MP3 file upload
	mp3File, err := c.FormFile("mp3_file")
	if err == nil {
		mp3FileExt := filepath.Ext(mp3File.Filename)
		mp3FilePath = "uploads/" + uuid.New().String() + mp3FileExt
		if err := c.SaveUploadedFile(mp3File, mp3FilePath); err != nil {
			os.Remove(coverImagePath)
			errors.HandleError(c, http.StatusInternalServerError, errors.ErrInternalServer)
			return
		}
		updatedTrack.Mp3FileUrl = c.Request.Host + "/" + mp3FilePath
	}

	// Update track in the database
	track, err := tc.trackService.UpdateTrack(trackId, &updatedTrack)
	if err != nil {
		os.Remove(coverImagePath)
		os.Remove(mp3FilePath)
		errors.HandleError(c, http.StatusInternalServerError, err)
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

	response := utils.NewSuccessResponse("Track updated successfully", output)
	c.JSON(http.StatusOK, response)
}

// DeleteTrack handles deleting a track
func (tc *TrackController) DeleteTrack(c *gin.Context) {
	trackId := c.Param("trackId")

	err := tc.trackService.DeleteTrack(trackId)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	response := utils.NewSuccessResponse("Track deleted successfully", nil)
	c.JSON(http.StatusOK, response)
}

// ListTracks handles listing all tracks with pagination
func (tc *TrackController) ListTracks(c *gin.Context) {
	var input ListTracksInput

	// Bind query parameters to input struct
	if err := c.ShouldBindQuery(&input); err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	if input.Page == 0 {
		input.Page = 1
	}
	if input.Limit == 0 {
		input.Limit = 10
	}

	// Call service to list tracks
	tracks, totalCount, err := tc.trackService.ListTracks(input.Page, input.Limit)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	// Prepare output data
	output := PaginatedTracksOutput{
		Page:       input.Page,
		Limit:      input.Limit,
		TotalCount: totalCount,
		Tracks:     make([]TrackOutput, len(tracks)),
	}

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

	response := utils.NewSuccessResponse("Tracks retrieved successfully", output)
	c.JSON(http.StatusOK, response)
}

// PlayPauseTrack handles playing or pausing a track
func (tc *TrackController) PlayPauseTrack(c *gin.Context) {
	trackId := c.Param("trackId")
	var input PlayPauseTrackInput

	// Bind JSON input to struct
	if err := c.ShouldBindJSON(&input); err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	err := tc.trackService.PlayPauseTrack(trackId, input.Action)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	response := utils.NewSuccessResponse("Track action performed successfully", nil)
	c.JSON(http.StatusOK, response)
}
