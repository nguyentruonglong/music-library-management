package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Genre represents a music genre in the library
type Genre struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name" binding:"required"`
	IsDeleted bool               `bson:"is_deleted" json:"is_deleted"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// BeforeCreate sets the CreatedAt and UpdatedAt fields to the current time before inserting a new genre.
func (g *Genre) BeforeCreate() {
	now := time.Now()
	g.CreatedAt = now
	g.UpdatedAt = now
	g.IsDeleted = false
}

// BeforeUpdate sets the UpdatedAt field to the current time before updating an existing genre.
func (g *Genre) BeforeUpdate() {
	g.UpdatedAt = time.Now()
}

// SoftDelete sets the DeletedAt field to the current time and the IsDeleted field to true.
func (g *Genre) SoftDelete() {
	now := time.Now()
	g.UpdatedAt = now
	g.DeletedAt = &now
	g.IsDeleted = true
}
