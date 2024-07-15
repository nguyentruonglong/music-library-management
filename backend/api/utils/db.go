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
	// Construct MongoDB URI
	mongoURI := fmt.Sprintf("mongodb://%s:%s/%s", cfg.MongoHost, cfg.MongoPort, cfg.MongoDB)

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}

// GetDatabase returns a MongoDB database instance
func GetDatabase(client *mongo.Client, cfg *config.Config) *mongo.Database {
	return client.Database(cfg.MongoDB)
}

// InitializeCollections ensures that the required collections exist
func InitializeCollections(db *mongo.Database) error {
	collections := []string{}
	for _, collection := range collections {
		exists, err := collectionExists(db, collection)
		if err != nil {
			return err
		}
		if !exists {
			err = db.CreateCollection(context.Background(), collection)
			if err != nil {
				return fmt.Errorf("failed to create collection %s: %v", collection, err)
			}
		}
	}
	return nil
}

// collectionExists checks if a collection exists in the database
func collectionExists(db *mongo.Database, collectionName string) (bool, error) {
	collections, err := db.ListCollectionNames(context.Background(), bson.D{{"name", collectionName}})
	if err != nil {
		return false, err
	}
	for _, name := range collections {
		if name == collectionName {
			return true, nil
		}
	}
	return false, nil
}
