package controllers

import (
	"context"
	"fmt"
	"log"

	"github.com/E_learning/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func UpdateSection(ctx context.Context, name string, arg *models.Section) (*mongo.UpdateResult, error) {
	collection := CourseCollection()
	filter := bson.D{primitive.E{Key: "Name", Value: name}}
	arrayFilters := options.ArrayFilters{Filters: bson.A{bson.M{"x._id": arg.ID}}}
	upsert := true
	opts := options.UpdateOptions{
		ArrayFilters: &arrayFilters,
		Upsert:       &upsert,
	}
	update := bson.M{
		"$set": bson.M{
			"Section.$[x].Title":   arg.Title,
			"Section.$[x].Content": arg.Content,
		},
	}
	result, err := collection.UpdateOne(ctx, filter, update, &opts)
	if err != nil {
		fmt.Printf("error updating db: %+v\n", err)
	}
	return result, err
}

type DelSection struct {
	Name string
	Id   string
}

func DeleteSection(ctx context.Context, arg DelSection) (*mongo.UpdateResult, error) {
	collection := CourseCollection()
	filter := bson.D{primitive.E{Key: "Name", Value: arg.Name}}
	iuud, _ := primitive.ObjectIDFromHex(arg.Id)
	update := bson.M{
		"$pull": bson.M{
			"Section": bson.D{primitive.E{Key: "_id", Value: iuud}},
		},
	}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return result, err
}

func FindSection(ctx context.Context, name string, id string) (*models.Section, error) {
	var section models.Section

	collection := CourseCollection()
	iuud, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"Name": name, "Section._id": iuud}
	err := collection.FindOne(ctx, filter).Decode(&section)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			log.Print("No such document")
		}
	}
	fmt.Println(section)
	return &section, err
}
