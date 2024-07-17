package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"music-library-management/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ConnectDB initializes a connection to MongoDB
func ConnectDB(cfg *config.Config) (*mongo.Client, error) {
	// Construct MongoDB URI from configuration values
	mongoURI := fmt.Sprintf("mongodb://%s:%s/%s", cfg.MongoHost, cfg.MongoPort, cfg.MongoDB)

	// Set MongoDB client options using the constructed URI
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Create a new MongoDB client
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err // Return an error if the client creation fails
	}

	// Create a context with a timeout of 10 seconds for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context is canceled to free up resources

	// Connect the MongoDB client to the server
	err = client.Connect(ctx)
	if err != nil {
		return nil, err // Return an error if the connection fails
	}

	// Ping the primary server to verify the connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err // Return an error if the ping fails
	}

	// Log a message indicating a successful connection
	log.Println("Connected to MongoDB!")
	return client, nil // Return the MongoDB client
}

// GetDatabase returns a MongoDB database instance
func GetDatabase(client *mongo.Client, cfg *config.Config) *mongo.Database {
	// Return a database instance using the configured database name
	return client.Database(cfg.MongoDB)
}

// InitializeCollections ensures that the required collections exist
func InitializeCollections(db *mongo.Database) error {
	// Define a list of required collections
	collections := []string{"tracks", "playlists", "genres", "files"}

	// Iterate over each collection name
	for _, collection := range collections {
		// Check if the collection exists
		exists, err := collectionExists(db, collection)
		if err != nil {
			return err // Return an error if the existence check fails
		}

		// If the collection does not exist, create it
		if !exists {
			err = db.CreateCollection(context.Background(), collection)
			if err != nil {
				// Return an error if the collection creation fails
				return fmt.Errorf("failed to create collection %s: %v", collection, err)
			}
		}
	}
	return nil // Return nil if all collections are initialized successfully
}

// collectionExists checks if a collection exists in the database
func collectionExists(db *mongo.Database, collectionName string) (bool, error) {
	// List all collection names in the database that match the given collection name
	collections, err := db.ListCollectionNames(context.Background(), bson.D{{Key: "name", Value: collectionName}})
	if err != nil {
		return false, err // Return an error if listing collection names fails
	}

	// Check if the collection name is in the list of existing collections
	for _, name := range collections {
		if name == collectionName {
			return true, nil // Return true if the collection exists
		}
	}
	return false, nil // Return false if the collection does not exist
}

// GetDBCollection returns a MongoDB collection instance
func GetDBCollection(client *mongo.Client, cfg *config.Config, collectionName string) *mongo.Collection {
	// Return a collection instance using the configured database and collection name
	return client.Database(cfg.MongoDB).Collection(collectionName)
}
