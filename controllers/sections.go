package controllers

import (
	"context"
	"fmt"
	"log"

	"github.com/E_learning/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseSec struct {
	Name    string            `uri:"name"  binding:"required"`
	Section []*models.Section `json:"Section,omitempty" bson:"Section,omitempty"`
}

func AddSection(ctx context.Context, arg CourseSec) (*mongo.UpdateResult, error) {
	collection := CourseCollection()
	match := bson.M{"Name": arg.Name}
	change := bson.M{"$push": bson.M{"Section": bson.M{"$each": arg.Section}}}
	result, err := collection.UpdateOne(ctx, match, change)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("success")
	}
	return result, err
}

type SectionUpd struct {
	Name    string `uri:"name"  binding:"required"`
	Id      string `uri:"id"   binding:"required"`
	Title   string `json:"Title"`
	Content string `json:"Content"`
}

func UpdateSection(ctx context.Context, arg *SectionUpd) (*mongo.UpdateResult, error) {
	collection := CourseCollection()
	var result *mongo.UpdateResult
	var err error
	iuud, _ := primitive.ObjectIDFromHex(arg.Id)
	//match := bson.M{"Section.$.id": iuud}
	change := bson.M{
		"$set": bson.M{
			"Section.$.Title":   arg.Title,
			"Section.$.Content": arg.Content,
		},
	}
	result, err = collection.UpdateByID(ctx, iuud, change)
	if err != nil {
		log.Fatal(err)
	}

	return result, err
}
