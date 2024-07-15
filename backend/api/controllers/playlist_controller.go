package controllers

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type PlaylistController struct {
	collection *mongo.Collection
}
