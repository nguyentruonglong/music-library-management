package utils

import (
	"context"
	"log"
	"music-library-management/api/models"
	"music-library-management/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SeedGenres seeds the genres collection with sample data
func SeedGenres(client *mongo.Client, cfg *config.Config) {
	collection := GetDBCollection(client, cfg, "genres")

	// Sample genres to be added
	sampleGenres := []models.Genre{
		{Name: "Rock"},
		{Name: "Pop"},
		{Name: "Jazz"},
		{Name: "Classical"},
		{Name: "Hip Hop"},
	}

	var genresToInsert []interface{}
	for _, genre := range sampleGenres {
		// Check if the genre already exists
		var existingGenre models.Genre
		err := collection.FindOne(context.Background(), bson.M{"name": genre.Name}).Decode(&existingGenre)
		if err == mongo.ErrNoDocuments {
			// If the genre does not exist, prepare it for insertion
			genre.BeforeCreate()
			genresToInsert = append(genresToInsert, genre)
		}
	}

	// Perform bulk insert if there are genres to insert
	if len(genresToInsert) > 0 {
		opts := options.InsertMany().SetOrdered(false)
		_, err := collection.InsertMany(context.Background(), genresToInsert, opts)
		if err != nil {
			log.Printf("Error seeding genres: %v", err)
		} else {
			log.Printf("Seeded %d genres", len(genresToInsert))
		}
	}
}
