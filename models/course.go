package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `json:"name" binding:"required" bson:"Name,omitempty"`
	Author           string             `json:"author" binding:"required" bson:"Author"`
	Description      string             `json:"description" binding:"required" bson:"Description,omitempty"`
	CreatedAt        time.Time          `json:"created_at" bson:"Created_at"`
	UpdatedAt        time.Time          `json:"updated_at,omitempty" bson:"Updated_at,omitempty"`
	CourseMaterialID primitive.ObjectID `json:"CourseMaterialID,omitempty" bson:"CourseMaterialID,omitempty"`
	Section          []Section          `json:"Section,omitempty" bson:"Section,omitempty"`
	StudentsEnrolled []string           `json:"StudentsEnrolled" bson:"StudentsEnrolled,omitempty"`
}

type Section struct {
	ID      primitive.ObjectID `bson:"SectionId,omitempty"`
	Title   string             `json:"Title" binding:"required" bson:"Title,omitempty"`
	Content string             `json:"Content" binding:"required" bson:"Content,omitempty"`
}
