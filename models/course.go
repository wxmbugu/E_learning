package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `bson:"Name,omitempty"`
	Author           []string           `bson:"Author"`
	Description      string             `bson:"Description,omitempty"`
	CreatedAt        time.Time          `bson:"Created_at"`
	UpdatedAt        time.Time          `bson:"Updated_at"`
	CourseMaterialID primitive.ObjectID `json:"CourseMaterialID" bson:"CourseMaterialID,omitempty"`
}
