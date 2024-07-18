package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Playlist represents a playlist in the library
type Playlist struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string               `bson:"name" json:"name" binding:"required"`
	Tracks    []primitive.ObjectID `bson:"tracks" json:"tracks"`         // Tracks in the playlist
	IsDeleted bool                 `bson:"is_deleted" json:"is_deleted"` // Soft delete flag
	CreatedAt time.Time            `bson:"created_at" json:"created_at"` // Creation timestamp
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"` // Last update timestamp
	DeletedAt *time.Time           `bson:"deleted_at" json:"deleted_at"` // Deletion timestamp
}

// BeforeCreate sets the CreatedAt, UpdatedAt fields and initializes Tracks before creating a new playlist
func (p *Playlist) BeforeCreate() {
	now := time.Now()
	p.ID = primitive.NewObjectID()
	p.CreatedAt = now
	p.UpdatedAt = now
	p.DeletedAt = nil
	p.IsDeleted = false
	p.Tracks = []primitive.ObjectID{} // Initialize Tracks as an empty array
}

// BeforeUpdate sets the UpdatedAt field before updating an existing playlist
func (p *Playlist) BeforeUpdate() {
	p.UpdatedAt = time.Now()
}

// SoftDelete sets the DeletedAt and IsDeleted fields to mark the playlist as deleted
func (p *Playlist) SoftDelete() {
	now := time.Now()
	p.UpdatedAt = now
	p.DeletedAt = &now
	p.IsDeleted = true
}
