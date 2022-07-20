package controllers

import (
	"context"
	"fmt"
	"log"

	"github.com/E_learning/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionSection = "Section"
)

type Section struct {
	client *mongo.Client
}

var (
	ErrInvalidUser    = errors.New("Account doesn't belong to authenticated user")
	ErrNoSuchDocument = errors.New("No such document")
)

func (s *Section) sectionCollection(ctx context.Context) *mongo.Collection {
	collection := s.client.Database(dbname).Collection((collectionSection))
	return collection
}

//function to add section to a course
func (s *Section) AddSection(ctx context.Context, section models.Section) (models.Section, error) {
	collection := s.sectionCollection(ctx)
	_, err := collection.InsertOne(ctx, section)
	return section, err
}

//update course section
func (s *Section) UpdateSection(ctx context.Context, id, title string) (*mongo.UpdateResult, error) {
	collection := s.sectionCollection(ctx)
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "Title", Value: title}}},
	}
	iuud, _ := primitive.ObjectIDFromHex(id)
	updateResult, err := collection.UpdateByID(context.TODO(), iuud, update)
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

func (s *Section) DeleteSection(ctx context.Context, id string) error {
	collection := s.sectionCollection(ctx)
	iuud, _ := primitive.ObjectIDFromHex(id)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": iuud})
	if err != nil {
		log.Fatal(err)
	}
	return err
}

//find course section by id
func (s *Section) FindSection(ctx context.Context, id string) (models.Section, error) {
	collection := s.sectionCollection(ctx)
	var results models.Section
	iuud, _ := primitive.ObjectIDFromHex(id)
	err := collection.FindOne(ctx, bson.M{"_id": iuud}).Decode(&results)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			log.Print("No such document")
		}
	}
	return results, err
}

//find course section by title
func (s *Section) FindSectionbyTitle(ctx context.Context, title string) (models.Section, error) {
	collection := s.sectionCollection(ctx)
	var results models.Section
	err := collection.FindOne(ctx, bson.M{"Title": title}).Decode(&results)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			log.Print("No such document")
		}
	}
	return results, err
}

func (s *Section) Content(ctx context.Context, sectiontitle, id string) (*mongo.UpdateResult, error) {
	collection := s.sectionCollection(ctx)
	section, err := s.FindSectionbyTitle(ctx, sectiontitle)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
	}

	match := bson.M{"Title": sectiontitle}
	change := bson.M{"$push": bson.M{"Content": id}}

	result, err := collection.UpdateOne(ctx, match, change)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("New section added to:", section.Title)
	}
	return result, err
}
