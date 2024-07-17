package controllers

import (
	"music-library-management/api/models"
	"music-library-management/api/services"
	"music-library-management/api/utils"
	"music-library-management/errors"
	"net/http"

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

// AddPlaylistInput represents the input data for adding a new playlist
type AddPlaylistInput struct {
	Name string `json:"name" binding:"required"`
}

// UpdatePlaylistInput represents the input data for updating a playlist
type UpdatePlaylistInput struct {
	Name string `json:"name" binding:"required"`
}

// ListPlaylistsInput represents the input data for listing playlists
type ListPlaylistsInput struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

// PlaylistOutput represents the output data for a playlist
type PlaylistOutput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PaginatedPlaylistsOutput represents the output data for paginated playlists
type PaginatedPlaylistsOutput struct {
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalCount int64            `json:"total_count"`
	Playlists  []PlaylistOutput `json:"playlists"`
}

// AddPlaylist handles adding a new playlist
func (pc *PlaylistController) AddPlaylist(c *gin.Context) {
	var input AddPlaylistInput
	var playlist models.Playlist

	// Bind JSON input to the AddPlaylistInput struct
	if err := c.ShouldBindJSON(&input); err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Copy input data to playlist model
	playlist.Name = input.Name

	// Call service to add the playlist
	createdPlaylist, err := pc.playlistService.AddPlaylist(&playlist)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	// Prepare output data
	output := PlaylistOutput{
		ID:   createdPlaylist.ID.Hex(),
		Name: createdPlaylist.Name,
	}

	// Respond with success message and created playlist
	response := utils.NewSuccessResponse("Playlist added successfully", output)
	c.JSON(http.StatusCreated, response)
}

// GetPlaylist handles retrieving a playlist by ID
func (pc *PlaylistController) GetPlaylist(c *gin.Context) {
	playlistId := c.Param("playlistId")

	// Call service to get the playlist
	playlist, err := pc.playlistService.GetPlaylist(playlistId)
	if err != nil {
		errors.HandleError(c, http.StatusNotFound, err)
		return
	}

	// Prepare output data
	output := PlaylistOutput{
		ID:   playlist.ID.Hex(),
		Name: playlist.Name,
	}

	// Respond with success message and retrieved playlist
	response := utils.NewSuccessResponse("Playlist retrieved successfully", output)
	c.JSON(http.StatusOK, response)
}

// UpdatePlaylist handles updating an existing playlist
func (pc *PlaylistController) UpdatePlaylist(c *gin.Context) {
	playlistId := c.Param("playlistId")

	// Check if the playlist exists
	_, err := pc.playlistService.GetPlaylist(playlistId)
	if err != nil {
		errors.HandleError(c, http.StatusNotFound, err)
		return
	}

	var input UpdatePlaylistInput
	var updatedPlaylist models.Playlist

	// Bind JSON input to the UpdatePlaylistInput struct
	if err := c.ShouldBindJSON(&input); err != nil {
		errors.HandleError(c, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	// Copy input data to updatedPlaylist model
	updatedPlaylist.Name = input.Name

	// Call service to update the playlist
	playlist, err := pc.playlistService.UpdatePlaylist(playlistId, &updatedPlaylist)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	// Prepare output data
	output := PlaylistOutput{
		ID:   playlist.ID.Hex(),
		Name: playlist.Name,
	}

	// Respond with success message and updated playlist
	response := utils.NewSuccessResponse("Playlist updated successfully", output)
	c.JSON(http.StatusOK, response)
}

// DeletePlaylist handles deleting a playlist
func (pc *PlaylistController) DeletePlaylist(c *gin.Context) {
	playlistId := c.Param("playlistId")

	// Call service to delete the playlist
	err := pc.playlistService.DeletePlaylist(playlistId)
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
	var input ListPlaylistsInput

	// Bind query parameters to ListPlaylistsInput struct
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

	// Call service to list playlists
	playlists, totalCount, err := pc.playlistService.ListPlaylists(input.Page, input.Limit)
	if err != nil {
		errors.HandleError(c, http.StatusInternalServerError, err)
		return
	}

	// Prepare output data
	output := PaginatedPlaylistsOutput{
		Page:       input.Page,
		Limit:      input.Limit,
		TotalCount: totalCount,
		Playlists:  make([]PlaylistOutput, len(playlists)),
	}

	for i, playlist := range playlists {
		output.Playlists[i] = PlaylistOutput{
			ID:   playlist.ID.Hex(),
			Name: playlist.Name,
		}
	}

	// Respond with success message and list of playlists
	response := utils.NewSuccessResponse("Playlists retrieved successfully", output)
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
