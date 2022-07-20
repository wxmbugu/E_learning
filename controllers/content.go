package controllers

import (
	"context"
	//	"errors"
	//	"fmt"
	"log"

	"github.com/E_learning/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	//	"go.mongodb.org/mongo-driver/mongo/options"
)

type Content struct {
	client *mongo.Client
}

const (
	collectionContent = "Content"
)

func (c *Content) contentCollection(ctx context.Context) *mongo.Collection {
	collection := c.client.Database(dbname).Collection((collectionContent))
	return collection
}

//changes id string to primitive.ObjectID
func IDFromHex(id string) primitive.ObjectID {
	iuud, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	return iuud
}

//find content of a section in a course
func (c *Content) FindContent(ctx context.Context, id string) (models.Content, error) {
	collection := c.contentCollection(ctx)
	var results models.Content
	iuud := IDFromHex(id)
	err := collection.FindOne(ctx, bson.M{"_id": iuud}).Decode(&results)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			log.Print("No such document")
		}
	}
	return results, err
}

type DelContent struct {
	CourseName   string
	SubsectionId string
}

//delete content of a section in a course
func (c *Content) DeleteContent(ctx context.Context, id string) error {
	collection := c.contentCollection(ctx)
	iuud := IDFromHex(id)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": iuud})
	if err != nil {
		log.Fatal(err)
	}
	return err
}

type CourseSubSection struct {
	Content []*models.Content `json:"Content"`
}

//add content of a section in a course
func (c *Content) AddContent(ctx context.Context, args models.Content) (models.Content, error) {
	collection := c.contentCollection(ctx)
	_, err := collection.InsertOne(ctx, args)
	return args, err
}

//update content of a section in a course
func (c *Content) UpdateContentVideo(ctx context.Context, id, video string) (*mongo.UpdateResult, error) {
	collection := c.contentCollection(ctx)
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "Video", Value: video}}},
	}
	iuud := IDFromHex(id)
	updateResult, err := collection.UpdateByID(ctx, iuud, update)
	if err != nil {
		if we, ok := err.(mongo.WriteException); ok {
			for _, e := range we.WriteErrors {
				if e.Index == 0 {
					log.Print(err)
				}
			}
		}
	}
	//fmt.Print(updateResult)
	return updateResult, err
}

//update content of a section in a course
func (c *Content) UpdateContentTitle(ctx context.Context, id, title string) (*mongo.UpdateResult, error) {
	collection := c.contentCollection(ctx)
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "Title", Value: title}}},
	}
	iuud := IDFromHex(id)
	updateResult, err := collection.UpdateByID(ctx, iuud, update)
	if err != nil {
		if we, ok := err.(mongo.WriteException); ok {
			for _, e := range we.WriteErrors {
				if e.Index == 0 {
					log.Print(err)
				}
			}
		}
	}
	//fmt.Print(updateResult)
	return updateResult, err
}
