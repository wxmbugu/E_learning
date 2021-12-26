package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CourseCategory struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name               string             `json:"name" bson:"name,omitempty"`
	Description        string             `json:"description" bson:"description,omitempty"`
	CourseCategoryName CourseCategoryName `json:"courseCategoryNameList" bson:"courseCategoryNameList,omitempty"`
}

type CourseCategoryName struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
}
