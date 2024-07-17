package services

import (
	"context"
	"music-library-management/api/models"
	"music-library-management/api/utils"
	"music-library-management/config"
	"music-library-management/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SearchService struct {
	trackCollection    *mongo.Collection // MongoDB collection for tracks
	playlistCollection *mongo.Collection // MongoDB collection for playlists
}

// NewSearchService creates a new SearchService
func NewSearchService(client *mongo.Client, cfg *config.Config) *SearchService {
	return &SearchService{
		trackCollection:    utils.GetDBCollection(client, cfg, "tracks"),
		playlistCollection: utils.GetDBCollection(client, cfg, "playlists"),
	}
}

// SearchTracks searches for tracks by title, artist, album, or genre
func (s *SearchService) SearchTracks(query string, page, limit int) ([]*models.Track, int64, error) {
	skip := (page - 1) * limit // Calculate the number of documents to skip
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))                            // Set the number of documents to skip
	findOptions.SetLimit(int64(limit))                          // Set the number of documents to return
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by created_at in descending order

	// Create a filter for case-insensitive substring search and not deleted
	filter := bson.M{
		"$and": []bson.M{
			{"is_deleted": false}, // Filter out deleted tracks
			{
				"$or": []bson.M{
					{"title": bson.M{"$regex": query, "$options": "i"}},  // Search by title
					{"artist": bson.M{"$regex": query, "$options": "i"}}, // Search by artist
					{"album": bson.M{"$regex": query, "$options": "i"}},  // Search by album
					{"genre": bson.M{"$regex": query, "$options": "i"}},  // Search by genre
				},
			},
		},
	}

	// Execute the find query
	cursor, err := s.trackCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, 0, errors.ErrDatabaseOperation // Return error if query fails
	}

	var tracks []*models.Track
	if err := cursor.All(context.Background(), &tracks); err != nil {
		return nil, 0, errors.ErrDatabaseOperation // Return error if decoding fails
	}

	// Count total matching documents
	total, err := s.trackCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, errors.ErrDatabaseOperation // Return error if counting fails
	}

	return tracks, total, nil // Return found tracks and total count
}

// SearchPlaylists searches for playlists by name
func (s *SearchService) SearchPlaylists(query string, page, limit int) ([]*models.Playlist, int64, error) {
	skip := (page - 1) * limit // Calculate the number of documents to skip
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))                            // Set the number of documents to skip
	findOptions.SetLimit(int64(limit))                          // Set the number of documents to return
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by created_at in descending order

	// Create a filter for case-insensitive substring search by name and not deleted
	filter := bson.M{
		"$and": []bson.M{
			{"is_deleted": false}, // Filter out deleted playlists
			{"name": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	// Execute the find query
	cursor, err := s.playlistCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, 0, errors.ErrDatabaseOperation // Return error if query fails
	}

	var playlists []*models.Playlist
	if err := cursor.All(context.Background(), &playlists); err != nil {
		return nil, 0, errors.ErrDatabaseOperation // Return error if decoding fails
	}

	// Count total matching documents
	total, err := s.playlistCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, errors.ErrDatabaseOperation // Return error if counting fails
	}

	return playlists, total, nil // Return found playlists and total count
}
