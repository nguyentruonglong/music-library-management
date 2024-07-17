package controllers

import (
	"music-library-management/api/models"
	"music-library-management/api/services"
	"music-library-management/api/utils"
	"music-library-management/errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TrackController struct {
	trackService *services.TrackService
}

func NewTrackController(trackService *services.TrackService) *TrackController {
	return &TrackController{
		trackService: trackService,
	}
}

func (tc *TrackController) AddTrack(c *gin.Context) {
	var track models.Track

	// Parse multipart form data
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		errors.HandleError(c, http.StatusBadRequest, err)
		return
	}

	// Parse form fields
	track.Title = c.PostForm("title")
	track.Artist = c.PostForm("artist")
	track.Album = c.PostForm("album")
	track.Genre = c.PostForm("genre")
	track.ReleaseYear, _ = strconv.Atoi(c.PostForm("release_year"))
	track.Duration, _ = strconv.Atoi(c.PostForm("duration"))

	// Handle cover image upload
	coverImage, err := c.FormFile("cover_image")
	if err != nil {
		errors.HandleError(c, http.StatusBadRequest, err)
		return
	}
	coverImageExt := filepath.Ext(coverImage.Filename)
	coverImagePath := "uploads/" + uuid.New().String() + coverImageExt
	if err := c.SaveUploadedFile(coverImage, coverImagePath); err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}
	track.CoverImageUrl = c.Request.Host + "/" + coverImagePath

	// Handle MP3 file upload
	mp3File, err := c.FormFile("mp3_file")
	if err != nil {
		os.Remove(coverImagePath)
		errors.HandleError(c, http.StatusBadRequest, err)
		return
	}
	mp3FileExt := filepath.Ext(mp3File.Filename)
	mp3FilePath := "uploads/" + uuid.New().String() + mp3FileExt
	if err := c.SaveUploadedFile(mp3File, mp3FilePath); err != nil {
		os.Remove(coverImagePath)
		errors.HandleError(c, http.StatusInternalServerError, err)
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

	response := utils.NewSuccessResponse("Track added successfully", createdTrack)
	c.JSON(http.StatusCreated, response)
}

func (tc *TrackController) GetTrack(c *gin.Context) {
	id := c.Param("trackId")

	track, err := tc.trackService.GetTrack(id)
	if err != nil {
		errors.HandleError(c, http.StatusNotFound, err)
		return
	}

	response := utils.NewSuccessResponse("Track retrieved successfully", track)
	c.JSON(http.StatusOK, response)
}

func (tc *TrackController) UpdateTrack(c *gin.Context) {
	id := c.Param("trackId")

	// Check if the track exists
	_, err := tc.trackService.GetTrack(id)
	if err != nil {
		errors.HandleError(c, http.StatusNotFound, err)
		return
	}

	// Parse multipart form data
	err = c.Request.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		errors.HandleError(c, http.StatusBadRequest, err)
		return
	}

	var updatedTrack models.Track
	var coverImagePath, mp3FilePath string

	// Parse form fields
	updatedTrack.Title = c.PostForm("title")
	updatedTrack.Artist = c.PostForm("artist")
	updatedTrack.Album = c.PostForm("album")
	updatedTrack.Genre = c.PostForm("genre")
	updatedTrack.ReleaseYear, _ = strconv.Atoi(c.PostForm("release_year"))
	updatedTrack.Duration, _ = strconv.Atoi(c.PostForm("duration"))

	// Handle cover image upload
	coverImage, err := c.FormFile("cover_image")
	if err == nil {
		coverImageExt := filepath.Ext(coverImage.Filename)
		coverImagePath = "uploads/" + uuid.New().String() + coverImageExt
		if err := c.SaveUploadedFile(coverImage, coverImagePath); err != nil {
			errors.HandleError(c, http.StatusInternalServerError, err)
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
			errors.HandleError(c, http.StatusInternalServerError, err)
			return
		}
		updatedTrack.Mp3FileUrl = c.Request.Host + "/" + mp3FilePath
	}

	// Update track in the database
	track, err := tc.trackService.UpdateTrack(id, &updatedTrack)
	if err != nil {
		os.Remove(coverImagePath)
		os.Remove(mp3FilePath)
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	response := utils.NewSuccessResponse("Track updated successfully", track)
	c.JSON(http.StatusOK, response)
}

func (tc *TrackController) DeleteTrack(c *gin.Context) {
	id := c.Param("trackId")

	err := tc.trackService.DeleteTrack(id)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	response := utils.NewSuccessResponse("Track deleted successfully", nil)
	c.JSON(http.StatusOK, response)
}

func (tc *TrackController) ListTracks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	tracks, err := tc.trackService.ListTracks(page, limit)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	response := utils.NewSuccessResponse("Tracks retrieved successfully", tracks)
	c.JSON(http.StatusOK, response)
}

func (tc *TrackController) PlayPauseTrack(c *gin.Context) {
	id := c.Param("trackId")
	var requestBody struct {
		Action string `json:"action"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		errors.HandleError(c, http.StatusBadRequest, err)
		return
	}

	err := tc.trackService.PlayPauseTrack(id, requestBody.Action)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	response := utils.NewSuccessResponse("Track action performed successfully", nil)
	c.JSON(http.StatusOK, response)
}
