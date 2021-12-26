package controllers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/E_learning/db"
	"github.com/E_learning/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionCourse = "Course"
)

//creates the course collection
func CourseCollection() *mongo.Collection {
	db, err := db.DBInstance()
	if err != nil {
		log.Fatal(err)
	}
	collection := db.OpenCollection(context.Background(), collectionCourse)
	return collection
}

func CreateCourse(ctx context.Context, course *models.Course) (*models.Course, error) {
	collection := CourseCollection()
	_, err := collection.InsertOne(ctx, course)
	return course, err
}

// find one course
func FindCourse(ctx context.Context, name string) (models.Course, error) {
	collection := CourseCollection()
	var results models.Course
	err := collection.FindOne(ctx, bson.M{"Name": name}).Decode(&results)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		log.Fatal(err)
	}
	//log.Fatal(err)
	fmt.Print(results)
	//log.Fatal(err)
	return results, err
}

func UpdateCourse(ctx context.Context, id string, name, description string) error {
	collection := CourseCollection()
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "Name", Value: name}, {Key: "Description", Value: description}, {Key: "Updated_at", Value: time.Now()}}},
	}
	iuud, _ := primitive.ObjectIDFromHex(id)
	updateResult, err := collection.UpdateByID(context.TODO(), iuud, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(updateResult)
	return err
}
func DeleteCourse(ctx context.Context, id string) error {
	collection := CourseCollection()
	iuud, _ := primitive.ObjectIDFromHex(id)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": iuud})
	if err != nil {
		log.Fatal(err)
	}
	return err
}
