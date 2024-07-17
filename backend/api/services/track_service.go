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

// TrackService handles operations related to tracks
type TrackService struct {
	collection *mongo.Collection // MongoDB collection for tracks
}

// NewTrackService creates a new TrackService
func NewTrackService(client *mongo.Client, cfg *config.Config) *TrackService {
	return &TrackService{
		collection: utils.GetDBCollection(client, cfg, "tracks"),
	}
}

// AddTrack adds a new track to the database
func (s *TrackService) AddTrack(track *models.Track) (*models.Track, error) {
	track.BeforeCreate() // Set default values before creating a new track

	_, err := s.collection.InsertOne(context.Background(), track) // Insert the track into the database
	if err != nil {
		return nil, errors.ErrDatabaseOperation
	}

	return track, nil
}

// GetTrack retrieves a track by its ID
func (s *TrackService) GetTrack(trackId string) (*models.Track, error) {
	objectID, err := primitive.ObjectIDFromHex(trackId) // Convert string ID to ObjectID
	if err != nil {
		return nil, errors.ErrInvalidObjectID
	}

	var track models.Track
	err = s.collection.FindOne(context.Background(), bson.M{"_id": objectID, "is_deleted": false}).Decode(&track) // Find track by ID and ensure it's not deleted
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrTrackNotFound
		}
		return nil, errors.ErrDatabaseOperation
	}

	return &track, nil
}

// UpdateTrack updates an existing track
func (s *TrackService) UpdateTrack(trackId string, updatedTrack *models.Track) (*models.Track, error) {
	objectID, err := primitive.ObjectIDFromHex(trackId) // Convert string ID to ObjectID
	if err != nil {
		return nil, errors.ErrInvalidObjectID
	}

	existingTrack, err := s.GetTrack(trackId) // Retrieve the existing track
	if err != nil {
		return nil, err
	}

	// Preserve the old values for fields that are not updated
	if updatedTrack.Title == "" {
		updatedTrack.Title = existingTrack.Title
	}
	if updatedTrack.CoverImageUrl == "" {
		updatedTrack.CoverImageUrl = existingTrack.CoverImageUrl
	}
	if updatedTrack.Artist == "" {
		updatedTrack.Artist = existingTrack.Artist
	}
	if updatedTrack.Album == "" {
		updatedTrack.Album = existingTrack.Album
	}
	if updatedTrack.Genre == "" {
		updatedTrack.Genre = existingTrack.Genre
	}
	if updatedTrack.ReleaseYear == 0 {
		updatedTrack.ReleaseYear = existingTrack.ReleaseYear
	}
	if updatedTrack.Duration == 0 {
		updatedTrack.Duration = existingTrack.Duration
	}
	if updatedTrack.Mp3FileUrl == "" {
		updatedTrack.Mp3FileUrl = existingTrack.Mp3FileUrl
	}
	updatedTrack.ID = existingTrack.ID
	updatedTrack.CreatedAt = existingTrack.CreatedAt
	updatedTrack.BeforeUpdate() // Set updated values before updating the track

	filter := bson.M{"_id": objectID, "is_deleted": false} // Filter to find the track by ID and ensure it's not deleted
	update := bson.M{
		"$set": updatedTrack, // Update the track with new values
	}

	result := s.collection.FindOneAndUpdate(context.Background(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)) // Update the track in the database
	if result.Err() != nil {
		return nil, errors.ErrDatabaseOperation
	}

	var track models.Track
	err = result.Decode(&track) // Decode the updated track
	if err != nil {
		return nil, errors.ErrDatabaseOperation
	}

	return &track, nil
}

// DeleteTrack soft deletes a track by setting is_deleted to true
func (s *TrackService) DeleteTrack(trackId string) error {
	objectID, err := primitive.ObjectIDFromHex(trackId) // Convert string ID to ObjectID
	if err != nil {
		return errors.ErrInvalidObjectID
	}

	var track models.Track
	err = s.collection.FindOne(context.Background(), bson.M{"_id": objectID, "is_deleted": false}).Decode(&track) // Find track by ID and ensure it's not deleted
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.ErrTrackNotFound
		}
		return errors.ErrDatabaseOperation
	}

	track.SoftDelete() // Apply soft delete to the track

	update := bson.M{
		"$set": bson.M{
			"is_deleted": track.IsDeleted,
			"deleted_at": track.DeletedAt,
			"updated_at": track.UpdatedAt,
		},
	}

	result := s.collection.FindOneAndUpdate(context.Background(), bson.M{"_id": objectID}, update, nil) // Update the track to soft delete it
	if result.Err() != nil {
		return errors.ErrDatabaseOperation
	}

	return nil
}

// ListTracks lists all tracks with pagination
func (s *TrackService) ListTracks(page, limit int) ([]*models.Track, int64, error) {
	skip := (page - 1) * limit // Calculate the number of documents to skip
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))                            // Set the number of documents to skip
	findOptions.SetLimit(int64(limit))                          // Set the number of documents to return
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by created_at in descending order

	// Execute the find query to get tracks that are not deleted
	cursor, err := s.collection.Find(context.Background(), bson.M{"is_deleted": false}, findOptions)
	if err != nil {
		return nil, 0, errors.ErrDatabaseOperation
	}

	var tracks []*models.Track
	err = cursor.All(context.Background(), &tracks) // Decode all tracks
	if err != nil {
		return nil, 0, errors.ErrDatabaseOperation
	}

	// Count total matching documents to get the total number of tracks that are not deleted
	total, err := s.collection.CountDocuments(context.Background(), bson.M{"is_deleted": false})
	if err != nil {
		return nil, 0, errors.ErrDatabaseOperation
	}

	return tracks, total, nil // Return found tracks and total count
}

// PlayPauseTrack plays or pauses a track based on the action provided
func (s *TrackService) PlayPauseTrack(trackId string, action string) error {
	if action != "play" && action != "pause" { // Validate action
		return errors.ErrBadRequest
	}

	return nil
}
