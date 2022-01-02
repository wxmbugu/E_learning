package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Instructor struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `json:"Firstname" binding:"required" bson:"Firstname,omitempty"`
	LastName  string             `json:"Lastname" binding:"required" bson:"Lastname"`
	UserName  string             `bson:"Username" json:"Username" binding:"required"`
	Email     string             `bson:"Email" json:"Email" binding:"required"`
	Password  string             `bson:"Password" json:"Password" binding:"required"`
	Courses   []string           `json:"Courses,omitempty" bson:"Courses,omitempty"`
	Roles     []string           `bson:"Roles,omitempty"`
	CreatedAt time.Time          `json:"Created_at" bson:"Created_at"`
	UpdatedAt time.Time          `json:"Updated_at,omitempty" bson:"Updated_at,omitempty"`
}
