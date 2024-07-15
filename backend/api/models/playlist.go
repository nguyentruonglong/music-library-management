package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Playlist represents a playlist in the library
type Playlist struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string               `bson:"name" json:"name" binding:"required"`
	Tracks    []primitive.ObjectID `bson:"tracks" json:"tracks"`
	IsDeleted bool                 `bson:"is_deleted" json:"is_deleted"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time           `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// BeforeCreate sets the CreatedAt and UpdatedAt fields before creating a new playlist
func (p *Playlist) BeforeCreate() {
	p.ID = primitive.NewObjectID()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	p.IsDeleted = false
}

// BeforeUpdate sets the UpdatedAt field before updating an existing playlist
func (p *Playlist) BeforeUpdate() {
	p.UpdatedAt = time.Now()
}

// SoftDelete sets the DeletedAt and IsDeleted fields to mark the playlist as deleted
func (p *Playlist) SoftDelete() {
	now := time.Now()
	p.DeletedAt = &now
	p.IsDeleted = true
}
