package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Track represents a music track in the library
type Track struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title" binding:"required"`
	CoverImage  string             `bson:"cover_image" json:"cover_image"`
	Artist      string             `bson:"artist" json:"artist" binding:"required"`
	Album       string             `bson:"album" json:"album"`
	Genre       string             `bson:"genre" json:"genre"`
	ReleaseYear int                `bson:"release_year" json:"release_year"`
	Duration    int                `bson:"duration" json:"duration" binding:"required"` // Duration in seconds
	FilePath    string             `bson:"file_path" json:"file_path"`
	IsDeleted   bool               `bson:"is_deleted" json:"is_deleted"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// BeforeCreate sets the CreatedAt and UpdatedAt fields before creating a new track
func (t *Track) BeforeCreate() {
	t.ID = primitive.NewObjectID()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	t.IsDeleted = false
}

// BeforeUpdate sets the UpdatedAt field before updating an existing track
func (t *Track) BeforeUpdate() {
	t.UpdatedAt = time.Now()
}

// SoftDelete sets the DeletedAt and IsDeleted fields to mark the track as deleted
func (t *Track) SoftDelete() {
	now := time.Now()
	t.DeletedAt = &now
	t.IsDeleted = true
}
