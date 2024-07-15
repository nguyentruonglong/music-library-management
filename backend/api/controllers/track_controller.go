package controllers

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type TrackController struct {
	collection *mongo.Collection
}
