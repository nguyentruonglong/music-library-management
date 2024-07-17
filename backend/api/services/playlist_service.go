package services

import (
	"context"
	"music-library-management/api/models"
	"music-library-management/api/utils"
	"music-library-management/config"
	"music-library-management/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PlaylistService handles operations related to playlists
type PlaylistService struct {
	collection   *mongo.Collection // MongoDB collection for playlists
	trackService *TrackService     // TrackService to handle track-related operations
}

// NewPlaylistService creates a new instance of PlaylistService
func NewPlaylistService(client *mongo.Client, cfg *config.Config, trackService *TrackService) *PlaylistService {
	return &PlaylistService{
		collection:   utils.GetDBCollection(client, cfg, "playlists"), // Get the playlists collection
		trackService: trackService,                                    // Initialize trackService for track-related operations
	}
}

// AddPlaylist adds a new playlist to the database
func (s *PlaylistService) AddPlaylist(playlist *models.Playlist) (*models.Playlist, error) {
	playlist.BeforeCreate() // Set default values before creating a playlist

	_, err := s.collection.InsertOne(context.Background(), playlist) // Insert playlist into the database
	if err != nil {
		return nil, errors.ErrDatabaseOperation
	}

	return playlist, nil
}

// GetPlaylist retrieves a playlist by its ID
func (s *PlaylistService) GetPlaylist(playlistId string) (*models.Playlist, error) {
	objectID, err := primitive.ObjectIDFromHex(playlistId) // Convert string ID to ObjectID
	if err != nil {
		return nil, errors.ErrInvalidObjectID
	}

	var playlist models.Playlist
	err = s.collection.FindOne(context.Background(), bson.M{"_id": objectID, "is_deleted": false}).Decode(&playlist) // Find playlist by ID
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrPlaylistNotFound
		}
		return nil, errors.ErrDatabaseOperation
	}

	return &playlist, nil
}

// UpdatePlaylist updates an existing playlist
func (s *PlaylistService) UpdatePlaylist(playlistId string, updatedPlaylist *models.Playlist) (*models.Playlist, error) {
	objectID, err := primitive.ObjectIDFromHex(playlistId) // Convert string ID to ObjectID
	if err != nil {
		return nil, errors.ErrInvalidObjectID
	}

	// Retrieve the existing playlist to preserve the old values
	existingPlaylist, err := s.GetPlaylist(playlistId)
	if err != nil {
		return nil, err
	}

	// Preserve the old values for fields that are not updated
	if updatedPlaylist.Name == "" {
		updatedPlaylist.Name = existingPlaylist.Name
	}
	updatedPlaylist.ID = existingPlaylist.ID
	updatedPlaylist.CreatedAt = existingPlaylist.CreatedAt
	updatedPlaylist.BeforeUpdate() // Set updated values before updating the playlist

	filter := bson.M{"_id": objectID, "is_deleted": false}
	update := bson.M{
		"$set": updatedPlaylist, // Set updated playlist values
	}

	result := s.collection.FindOneAndUpdate(context.Background(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)) // Update playlist in the database and return the updated document
	if result.Err() != nil {
		return nil, errors.ErrDatabaseOperation
	}

	var playlist models.Playlist
	err = result.Decode(&playlist) // Decode the updated playlist
	if err != nil {
		return nil, errors.ErrDatabaseOperation
	}

	return &playlist, nil
}

// DeletePlaylist soft deletes a playlist by setting is_deleted to true
func (s *PlaylistService) DeletePlaylist(playlistId string) error {
	objectID, err := primitive.ObjectIDFromHex(playlistId) // Convert string ID to ObjectID
	if err != nil {
		return errors.ErrInvalidObjectID
	}

	// Retrieve the existing playlist
	playlist, err := s.GetPlaylist(playlistId)
	if err != nil {
		return err
	}

	// Apply soft delete
	playlist.SoftDelete()

	update := bson.M{
		"$set": bson.M{
			"is_deleted": playlist.IsDeleted,
			"deleted_at": playlist.DeletedAt,
			"updated_at": playlist.UpdatedAt,
		},
	}

	result := s.collection.FindOneAndUpdate(context.Background(), bson.M{"_id": objectID}, update, nil)
	if result.Err() != nil {
		return errors.ErrDatabaseOperation
	}

	return nil
}

// ListPlaylists lists all playlists with pagination
func (s *PlaylistService) ListPlaylists(page, limit int) ([]*models.Playlist, int64, error) {
	skip := (page - 1) * limit
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))   // Set the number of documents to skip
	findOptions.SetLimit(int64(limit)) // Set the number of documents to return

	cursor, err := s.collection.Find(context.Background(), bson.M{"is_deleted": false}, findOptions) // Find playlists
	if err != nil {
		return nil, 0, errors.ErrDatabaseOperation
	}

	var playlists []*models.Playlist
	err = cursor.All(context.Background(), &playlists) // Decode all playlists
	if err != nil {
		return nil, 0, errors.ErrDatabaseOperation
	}

	total, err := s.collection.CountDocuments(context.Background(), bson.M{"is_deleted": false}) // Get the total number of playlists
	if err != nil {
		return nil, 0, errors.ErrDatabaseOperation
	}

	return playlists, total, nil
}

// AddTrackToPlaylist adds a track to a playlist
func (s *PlaylistService) AddTrackToPlaylist(playlistId, trackId string) error {
	playlistObjectID, err := primitive.ObjectIDFromHex(playlistId) // Convert playlist ID to ObjectID
	if err != nil {
		return errors.ErrInvalidObjectID
	}

	// Check if the track exists
	track, err := s.trackService.GetTrack(trackId) // Get the track by ID
	if err != nil {
		return errors.ErrTrackNotFound
	}

	// Retrieve the existing playlist
	playlist, err := s.GetPlaylist(playlistId)
	if err != nil {
		return err
	}

	// Check if the track already exists in the playlist
	for _, t := range playlist.Tracks {
		if t == track.ID {
			return errors.ErrTrackAlreadyInPlaylist
		}
	}

	filter := bson.M{"_id": playlistObjectID, "is_deleted": false}
	update := bson.M{
		"$addToSet": bson.M{"tracks": track.ID}, // Add track ID to the tracks array in the playlist
	}

	result := s.collection.FindOneAndUpdate(context.Background(), filter, update, nil) // Update the playlist with the new track
	if result.Err() != nil {
		return errors.ErrDatabaseOperation
	}

	return nil
}

// RemoveTrackFromPlaylist removes a track from a playlist
func (s *PlaylistService) RemoveTrackFromPlaylist(playlistId, trackId string) error {
	playlistObjectID, err := primitive.ObjectIDFromHex(playlistId) // Convert playlist ID to ObjectID
	if err != nil {
		return errors.ErrInvalidObjectID
	}

	// Check if the track exists
	track, err := s.trackService.GetTrack(trackId) // Get the track by ID
	if err != nil {
		return errors.ErrTrackNotFound
	}

	// Retrieve the existing playlist
	playlist, err := s.GetPlaylist(playlistId)
	if err != nil {
		return err
	}

	// Check if the track does not exist in the playlist
	found := false
	for _, t := range playlist.Tracks {
		if t == track.ID {
			found = true
			break
		}
	}
	if !found {
		return errors.ErrTrackNotInPlaylist
	}

	filter := bson.M{"_id": playlistObjectID, "is_deleted": false}
	update := bson.M{
		"$pull": bson.M{"tracks": track.ID}, // Remove track ID from the tracks array in the playlist
	}

	result := s.collection.FindOneAndUpdate(context.Background(), filter, update, nil) // Update the playlist by removing the track
	if result.Err() != nil {
		return errors.ErrDatabaseOperation
	}

	return nil
}
