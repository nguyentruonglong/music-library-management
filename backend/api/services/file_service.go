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

// FileService handles file management operations
type FileService struct {
	collection *mongo.Collection
	config     *config.Config
}

// NewFileService creates a new FileService
func NewFileService(client *mongo.Client, cfg *config.Config) *FileService {
	return &FileService{
		collection: utils.GetDBCollection(client, cfg, "files"),
		config:     cfg,
	}
}

// GetUploadPath returns the upload path from the config
func (s *FileService) GetUploadPath() string {
	return s.config.UploadPath
}

// SaveFileMetadata saves metadata for an uploaded file
func (s *FileService) SaveFileMetadata(filename string) (*models.File, error) {
	file := &models.File{
		Filename: filename,
		Filepath: s.GetUploadPath() + "/" + filename,
	}
	file.BeforeCreate() // Set default values before creating the file record

	_, err := s.collection.InsertOne(context.Background(), file)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// ListFiles lists all files metadata with pagination
func (s *FileService) ListFiles(page, limit int) ([]models.File, int64, error) {
	skip := (page - 1) * limit
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))                            // Set the number of documents to skip
	findOptions.SetLimit(int64(limit))                          // Set the number of documents to return
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by created_at in descending order

	cursor, err := s.collection.Find(context.Background(), bson.M{"is_deleted": false}, findOptions) // Find files
	if err != nil {
		return nil, 0, errors.ErrDatabaseOperation
	}

	var files []models.File
	if err := cursor.All(context.Background(), &files); err != nil {
		return nil, 0, errors.ErrDatabaseOperation
	}

	// Count total files excluding soft-deleted files
	total, err := s.collection.CountDocuments(context.Background(), bson.M{"is_deleted": false})
	if err != nil {
		return nil, 0, errors.ErrDatabaseOperation
	}

	return files, total, nil
}
