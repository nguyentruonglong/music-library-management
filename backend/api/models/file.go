package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// File represents the metadata for an uploaded file
type File struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Filename  string             `bson:"filename" json:"filename"`
	Filepath  string             `bson:"filepath" json:"filepath"`
	FileUrl   string             `bson:"file_url" json:"file_url"`
	IsDeleted bool               `bson:"is_deleted" json:"is_deleted"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at" json:"deleted_at"`
}

// BeforeCreate sets default values before creating a new file record
func (f *File) BeforeCreate() {
	now := time.Now()
	f.ID = primitive.NewObjectID()
	f.CreatedAt = now
	f.UpdatedAt = now
	f.DeletedAt = nil
	f.IsDeleted = false
}

// BeforeUpdate sets the updated_at field before updating a file record
func (f *File) BeforeUpdate() {
	f.UpdatedAt = time.Now()
}

// SoftDelete marks the file record as deleted without removing it from the database
func (f *File) SoftDelete() {
	now := time.Now()
	f.UpdatedAt = now
	f.DeletedAt = &now
	f.IsDeleted = true
}
