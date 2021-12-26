package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Name           string             `bson:"Name,omitempty"`
	Author         []string           `bson:"Author"`
	Description    string             `bson:"Description,omitempty"`
	CreatedAt      time.Time          `bson:"Created_at"`
	UpdatedAt      time.Time          `bson:"Updated_at"`
	CourseMaterial CourseMaterial     `json:"CourseMaterial" bson:"CourseMaterial,omitempty"`
}
