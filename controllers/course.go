package controllers

import (
	"context"
	"log"
	"time"

	"github.com/E_learning/db"
	"github.com/E_learning/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionCourse = "Course"
)

//creates the course collection
func CourseCollection(ctx context.Context) *mongo.Collection {
	db, err := db.DBInstance()
	if err != nil {
		log.Fatal(err)
	}
	collection := db.OpenCollection(ctx, collectionCourse)
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"Name": 1},
		Options: options.Index().SetUnique(true),
	},
	)
	if err != nil {
		log.Fatal(err)
	}
	return collection
}

//create course
func CreateCourse(ctx context.Context, course *models.Course) (*models.Course, error) {
	collection := CourseCollection(ctx)
	_, err := collection.InsertOne(ctx, course)
	return course, err
}

// find one course by id
func FindCourse(ctx context.Context, id string) (models.Course, error) {
	collection := CourseCollection(ctx)
	var results models.Course
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

//find course by name
func FindCoursebyName(ctx context.Context, name string) (models.Course, error) {
	collection := CourseCollection(ctx)
	var results models.Course
	err := collection.FindOne(ctx, bson.M{"Name": name}).Decode(&results)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			log.Print("No such document")
		}
	}
	return results, err
}

type UpdateCourseParams struct {
	ID          string `bson:"_id,omitempty"`
	Name        string `json:"name" binding:"required" bson:"Name,omitempty"`
	Description string `json:"description" binding:"required" bson:"Description,omitempty"`
}

//update course
func UpdateCourse(ctx context.Context, arg UpdateCourseParams) (*mongo.UpdateResult, error) {
	collection := CourseCollection(ctx)
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "Name", Value: arg.Name}, {Key: "Description", Value: arg.Description}, {Key: "Updated_at", Value: time.Now()}}},
	}
	iuud, _ := primitive.ObjectIDFromHex(arg.ID)

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

//delete course
func DeleteCourse(ctx context.Context, id string) error {
	collection := CourseCollection(ctx)
	iuud, _ := primitive.ObjectIDFromHex(id)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": iuud})
	if err != nil {
		log.Fatal(err)
	}
	return err
}

type ListCourseParams struct {
	Owner string `json:"owner"`
	Limit int64
	Skip  int64
}

//Find multiple documents of courses
func ListCourses(ctx context.Context, arg ListCourseParams) ([]models.Course, error) {
	collection := CourseCollection(ctx)
	//check the connection

	//find records
	//pass these options to the Find method
	findOptions := options.Find()
	//Set the limit of the number of record to find
	findOptions.SetLimit(arg.Limit)
	findOptions.SetSkip(arg.Skip)
	//Define an array in which you can store the decoded documents
	var results []models.Course

	//Passing the bson.D{{}} as the filter matches  documents in the collection
	cur, err := collection.Find(ctx, bson.D{{Key: "Author", Value: arg.Owner}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	//Finding multiple documents returns a cursor
	//Iterate through the cursor allows us to decode documents one at a time

	for cur.Next(ctx) {
		//Create a value into which the single document can be decoded
		var elem models.Course
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)

	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	//Close the cursor once finished
	cur.Close(ctx)
	return results, err
}
