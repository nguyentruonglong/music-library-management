package controllers

import (
	"music-library-management/api/models"
	"music-library-management/api/services"
	"music-library-management/api/utils"
	"music-library-management/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PlaylistController handles HTTP requests for playlists
type PlaylistController struct {
	playlistService *services.PlaylistService
}

// NewPlaylistController creates a new PlaylistController
func NewPlaylistController(playlistService *services.PlaylistService) *PlaylistController {
	return &PlaylistController{
		playlistService: playlistService,
	}
}

// AddPlaylist handles adding a new playlist
func (pc *PlaylistController) AddPlaylist(c *gin.Context) {
	var playlist models.Playlist

	// Bind JSON input to the playlist model
	if err := c.ShouldBindJSON(&playlist); err != nil {
		errors.HandleError(c, http.StatusBadRequest, err)
		return
	}

	// Call service to add the playlist
	createdPlaylist, err := pc.playlistService.AddPlaylist(&playlist)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	// Respond with success message and created playlist
	response := utils.NewSuccessResponse("Playlist added successfully", createdPlaylist)
	c.JSON(http.StatusCreated, response)
}

// GetPlaylist handles retrieving a playlist by ID
func (pc *PlaylistController) GetPlaylist(c *gin.Context) {
	id := c.Param("playlistId")

	// Call service to get the playlist
	playlist, err := pc.playlistService.GetPlaylist(id)
	if err != nil {
		errors.HandleError(c, http.StatusNotFound, err)
		return
	}

	// Respond with success message and retrieved playlist
	response := utils.NewSuccessResponse("Playlist retrieved successfully", playlist)
	c.JSON(http.StatusOK, response)
}

// UpdatePlaylist handles updating an existing playlist
func (pc *PlaylistController) UpdatePlaylist(c *gin.Context) {
	id := c.Param("playlistId")

	// Check if the playlist exists
	_, err := pc.playlistService.GetPlaylist(id)
	if err != nil {
		errors.HandleError(c, http.StatusNotFound, err)
		return
	}

	var updatedPlaylist models.Playlist
	// Bind JSON input to the updated playlist model
	if err := c.ShouldBindJSON(&updatedPlaylist); err != nil {
		errors.HandleError(c, http.StatusBadRequest, err)
		return
	}

	// Call service to update the playlist
	playlist, err := pc.playlistService.UpdatePlaylist(id, &updatedPlaylist)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	// Respond with success message and updated playlist
	response := utils.NewSuccessResponse("Playlist updated successfully", playlist)
	c.JSON(http.StatusOK, response)
}

// DeletePlaylist handles deleting a playlist
func (pc *PlaylistController) DeletePlaylist(c *gin.Context) {
	id := c.Param("playlistId")

	// Call service to delete the playlist
	err := pc.playlistService.DeletePlaylist(id)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	// Respond with success message
	response := utils.NewSuccessResponse("Playlist deleted successfully", nil)
	c.JSON(http.StatusOK, response)
}

// ListPlaylists handles listing all playlists with pagination
func (pc *PlaylistController) ListPlaylists(c *gin.Context) {
	// Parse pagination parameters from query
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Call service to list playlists
	playlists, err := pc.playlistService.ListPlaylists(page, limit)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	// Respond with success message and list of playlists
	response := utils.NewSuccessResponse("Playlists retrieved successfully", playlists)
	c.JSON(http.StatusOK, response)
}

// AddTrackToPlaylist handles adding a track to a playlist
func (pc *PlaylistController) AddTrackToPlaylist(c *gin.Context) {
	playlistId := c.Param("playlistId")
	trackId := c.Param("trackId")

	// Call service to add the track to the playlist
	err := pc.playlistService.AddTrackToPlaylist(playlistId, trackId)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	// Respond with success message
	response := utils.NewSuccessResponse("Track added to playlist successfully", nil)
	c.JSON(http.StatusOK, response)
}

// RemoveTrackFromPlaylist handles removing a track from a playlist
func (pc *PlaylistController) RemoveTrackFromPlaylist(c *gin.Context) {
	playlistId := c.Param("playlistId")
	trackId := c.Param("trackId")

	// Call service to remove the track from the playlist
	err := pc.playlistService.RemoveTrackFromPlaylist(playlistId, trackId)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	// Respond with success message
	response := utils.NewSuccessResponse("Track removed from playlist successfully", nil)
	c.JSON(http.StatusOK, response)
}
