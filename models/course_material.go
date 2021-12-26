package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CourseMaterial struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Author      []string             `json:"Author" bson:"Author,omitempty"`
	PdfFileID   []primitive.ObjectID `json:"PdfFileID" bson:"PdfFileID,omitempty"`
	VideoFileID []primitive.ObjectID `json:"VideoFileID" bson:"VideoFileID,omitempty"`
	CreatedAt   time.Time            `json:"Created_at" bson:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" bson:"updated_at"`
}
