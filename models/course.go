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
	Section          []*Section         `json:"Section,omitempty" bson:"Section,omitempty"`
	StudentsEnrolled []string           `json:"StudentsEnrolled" bson:"StudentsEnrolled,omitempty"`
}

type Section struct {
	ID      primitive.ObjectID `bson:"_id" `
	Title   string             `json:"Title" bson:"Title,omitempty"`
	Content []*Content         `json:"Content" binding:"required" bson:"Content,omitempty"`
}

type Content struct {
	ID         primitive.ObjectID `bson:"subsectionid" `
	SubTitle   string             `json:"Subsection_Title" binding:"required" bson:"Subsection_Title,omitempty"`
	SubContent string             `json:"SubContent" bson:"SubContent,omitempty"`
	Thumbnail  string             `json:"Thumbnail" bson:"Thumbnail,omitempty"`
}
