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
	// Get the genres collection from the database
	collection := GetDBCollection(client, cfg, "genres")

	// Define sample genres to be added
	sampleGenres := []models.Genre{
		{Name: "Rock"},
		{Name: "Pop"},
		{Name: "Jazz"},
		{Name: "Classical"},
		{Name: "Hip Hop"},
	}

	var genresToInsert []interface{} // Slice to hold genres that need to be inserted
	for _, genre := range sampleGenres {
		// Check if the genre already exists in the collection
		var existingGenre models.Genre
		err := collection.FindOne(context.Background(), bson.M{"name": genre.Name}).Decode(&existingGenre)
		if err == mongo.ErrNoDocuments {
			// If the genre does not exist, prepare it for insertion
			genre.BeforeCreate()                           // Set default values before creating the genre
			genresToInsert = append(genresToInsert, genre) // Add the genre to the list of genres to insert
		}
	}

	// Perform bulk insert if there are genres to insert
	if len(genresToInsert) > 0 {
		opts := options.InsertMany().SetOrdered(false)                              // Set the insert options to unordered
		_, err := collection.InsertMany(context.Background(), genresToInsert, opts) // Insert the genres into the collection
		if err != nil {
			// Log an error if the bulk insert fails
			log.Printf("Error seeding genres: %v", err)
		} else {
			// Log the number of genres seeded
			log.Printf("Seeded %d genres", len(genresToInsert))
		}
	}
}
