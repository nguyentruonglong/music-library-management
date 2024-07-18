package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Track represents a music track in the library
type Track struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title         string             `bson:"title" json:"title" binding:"required"`
	CoverImageUrl string             `bson:"cover_image_url" json:"cover_image_url"`
	Artist        string             `bson:"artist" json:"artist" binding:"required"`
	Album         string             `bson:"album" json:"album"`
	Genre         string             `bson:"genre" json:"genre"`
	ReleaseYear   int                `bson:"release_year" json:"release_year"`
	Duration      int                `bson:"duration" json:"duration" binding:"required"` // Duration in seconds
	Mp3FileUrl    string             `bson:"mp3_file_url" json:"mp3_file_url"`
	IsDeleted     bool               `bson:"is_deleted" json:"is_deleted"` // Soft delete flag
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"` // Creation timestamp
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"` // Last update timestamp
	DeletedAt     *time.Time         `bson:"deleted_at" json:"deleted_at"` // Deletion timestamp
}

// BeforeCreate sets the CreatedAt and UpdatedAt fields before creating a new track
func (t *Track) BeforeCreate() {
	now := time.Now()
	t.ID = primitive.NewObjectID()
	t.CreatedAt = now
	t.UpdatedAt = now
	t.DeletedAt = nil
	t.IsDeleted = false
}

// BeforeUpdate sets the UpdatedAt field before updating an existing track
func (t *Track) BeforeUpdate() {
	t.UpdatedAt = time.Now()
}

// SoftDelete sets the DeletedAt and IsDeleted fields to mark the track as deleted
func (t *Track) SoftDelete() {
	now := time.Now()
	t.UpdatedAt = now
	t.DeletedAt = &now
	t.IsDeleted = true
}
