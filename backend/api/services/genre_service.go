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

// GenreService handles operations related to genres
type GenreService struct {
	collection *mongo.Collection
}

// NewGenreService creates a new instance of GenreService
func NewGenreService(client *mongo.Client, cfg *config.Config) *GenreService {
	return &GenreService{
		collection: utils.GetDBCollection(client, cfg, "genres"),
	}
}

// AddGenre adds a new genre to the database
func (s *GenreService) AddGenre(genre *models.Genre) (*models.Genre, error) {
	genre.BeforeCreate() // Set default values before creating a genre

	_, err := s.collection.InsertOne(context.Background(), genre) // Insert genre into the database
	if err != nil {
		return nil, errors.ErrDatabaseOperation
	}

	return genre, nil
}

// GetGenre retrieves a genre by its ID
func (s *GenreService) GetGenre(genreId string) (*models.Genre, error) {
	objectID, err := primitive.ObjectIDFromHex(genreId) // Convert string ID to ObjectID
	if err != nil {
		return nil, errors.ErrInvalidObjectID
	}

	var genre models.Genre
	err = s.collection.FindOne(context.Background(), bson.M{"_id": objectID, "is_deleted": false}).Decode(&genre) // Find genre by ID
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrGenreNotFound
		}
		return nil, errors.ErrDatabaseOperation
	}

	return &genre, nil
}

// UpdateGenre updates an existing genre
func (s *GenreService) UpdateGenre(genreId string, updatedGenre *models.Genre) (*models.Genre, error) {
	objectID, err := primitive.ObjectIDFromHex(genreId) // Convert string ID to ObjectID
	if err != nil {
		return nil, errors.ErrInvalidObjectID
	}

	// Retrieve the existing genre to preserve the old values
	existingGenre, err := s.GetGenre(genreId)
	if err != nil {
		return nil, err
	}

	// Preserve the old values for fields that are not updated
	if updatedGenre.Name == "" {
		updatedGenre.Name = existingGenre.Name
	}
	updatedGenre.ID = existingGenre.ID
	updatedGenre.CreatedAt = existingGenre.CreatedAt
	updatedGenre.BeforeUpdate() // Set updated values before updating the genre

	filter := bson.M{"_id": objectID, "is_deleted": false}
	update := bson.M{
		"$set": updatedGenre, // Set updated genre values
	}

	result := s.collection.FindOneAndUpdate(context.Background(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)) // Update genre in the database and return the updated document
	if result.Err() != nil {
		return nil, errors.ErrDatabaseOperation
	}

	var genre models.Genre
	err = result.Decode(&genre) // Decode the updated genre
	if err != nil {
		return nil, errors.ErrDatabaseOperation
	}

	return &genre, nil
}

// DeleteGenre soft deletes a genre by setting is_deleted to true
func (s *GenreService) DeleteGenre(genreId string) error {
	objectID, err := primitive.ObjectIDFromHex(genreId) // Convert string ID to ObjectID
	if err != nil {
		return errors.ErrInvalidObjectID
	}

	var genre models.Genre
	err = s.collection.FindOne(context.Background(), bson.M{"_id": objectID, "is_deleted": false}).Decode(&genre) // Find genre by ID
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.ErrGenreNotFound
		}
		return errors.ErrDatabaseOperation
	}

	genre.SoftDelete() // Apply soft delete to the genre

	update := bson.M{
		"$set": bson.M{
			"is_deleted": genre.IsDeleted,
			"deleted_at": genre.DeletedAt,
			"updated_at": genre.UpdatedAt,
		},
	}

	result := s.collection.FindOneAndUpdate(context.Background(), bson.M{"_id": objectID}, update, nil) // Update the genre to soft delete it
	if result.Err() != nil {
		return errors.ErrDatabaseOperation
	}

	return nil
}

// ListGenres lists all genres with pagination
func (s *GenreService) ListGenres(page, limit int) ([]*models.Genre, int64, error) {
	skip := (page - 1) * limit
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))                            // Set the number of documents to skip
	findOptions.SetLimit(int64(limit))                          // Set the number of documents to return
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by created_at in descending order

	cursor, err := s.collection.Find(context.Background(), bson.M{"is_deleted": false}, findOptions) // Find genres
	if err != nil {
		return nil, 0, errors.ErrDatabaseOperation
	}

	var genres []*models.Genre
	err = cursor.All(context.Background(), &genres) // Decode all genres
	if err != nil {
		return nil, 0, errors.ErrDatabaseOperation
	}

	totalCount, err := s.collection.CountDocuments(context.Background(), bson.M{"is_deleted": false}) // Get the total number of genres
	if err != nil {
		return nil, 0, errors.ErrDatabaseOperation
	}

	return genres, totalCount, nil
}
