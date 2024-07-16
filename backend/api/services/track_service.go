package services

import (
	"context"
	"errors"
	"music-library-management/api/models"
	"music-library-management/api/utils"
	"music-library-management/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TrackService struct {
	collection *mongo.Collection
}

func NewTrackService(client *mongo.Client, cfg *config.Config) *TrackService {
	return &TrackService{
		collection: utils.GetDBCollection(client, cfg, "tracks"),
	}
}

func (s *TrackService) AddTrack(track *models.Track) (*models.Track, error) {
	track.BeforeCreate()

	_, err := s.collection.InsertOne(context.Background(), track)
	if err != nil {
		return nil, err
	}

	return track, nil
}

func (s *TrackService) GetTrack(id string) (*models.Track, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var track models.Track
	err = s.collection.FindOne(context.Background(), bson.M{"_id": objectID, "is_deleted": false}).Decode(&track)
	if err != nil {
		return nil, err
	}

	return &track, nil
}

func (s *TrackService) UpdateTrack(id string, updatedTrack *models.Track) (*models.Track, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Retrieve the existing track to preserve the old values
	existingTrack, err := s.GetTrack(id)
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
	updatedTrack.BeforeUpdate()

	filter := bson.M{"_id": objectID, "is_deleted": false}
	update := bson.M{
		"$set": updatedTrack,
	}

	result := s.collection.FindOneAndUpdate(context.Background(), filter, update, nil)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var track models.Track
	err = result.Decode(&track)
	if err != nil {
		return nil, err
	}

	return &track, nil
}

func (s *TrackService) DeleteTrack(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Soft delete the track
	var track models.Track
	err = s.collection.FindOne(context.Background(), bson.M{"_id": objectID, "is_deleted": false}).Decode(&track)
	if err != nil {
		return err
	}
	track.SoftDelete()

	update := bson.M{
		"$set": track,
	}

	result := s.collection.FindOneAndUpdate(context.Background(), bson.M{"_id": objectID}, update, nil)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (s *TrackService) ListTracks(page, limit int) ([]*models.Track, error) {
	skip := (page - 1) * limit
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	cursor, err := s.collection.Find(context.Background(), bson.M{"is_deleted": false}, findOptions)
	if err != nil {
		return nil, err
	}

	var tracks []*models.Track
	err = cursor.All(context.Background(), &tracks)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (s *TrackService) PlayPauseTrack(id string, action string) error {
	if action != "play" && action != "pause" {
		return errors.New("invalid action")
	}

	return nil
}
