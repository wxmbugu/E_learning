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
	InstructorCollection = "Instructors"
)

//creates the course collection
func CollectionInstructor() *mongo.Collection {
	db, err := db.DBInstance()
	if err != nil {
		log.Fatal(err)
	}
	collection := db.OpenCollection(context.Background(), InstructorCollection)
	return collection
}

func CreateInstructor(ctx context.Context, instructor *models.User) (*models.User, error) {
	collection := CollectionInstructor()
	_, err := collection.InsertOne(ctx, instructor)
	return instructor, err
}

// find one course
func FindInstructor(ctx context.Context, username string) (*models.User, error) {
	collection := CollectionInstructor()
	var results models.User
	err := collection.FindOne(ctx, bson.M{"Username": username}).Decode(&results)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			log.Print("No such document")
		}
	}
	return &results, err
}

type UpdateInstructorParams struct {
	ID        string    `bson:"_id,omitempty"`
	FirstName string    `json:"Firstname" binding:"required" bson:"Firstname,omitempty"`
	LastName  string    `json:"Lastname" binding:"required" bson:"Lastname"`
	UserName  string    `bson:"Username" json:"Username" binding:"required"`
	Email     string    `bson:"Email" json:"Email" binding:"required"`
	Password  string    `bson:"Password" json:"Password" binding:"required"`
	UpdatedAt time.Time `json:"Updated_at,omitempty" bson:"Updated_at,omitempty"`
}

func UpdateInstructor(ctx context.Context, arg UpdateInstructorParams) (*mongo.UpdateResult, error) {
	collection := CollectionInstructor()
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "Firstname", Value: arg.FirstName}, {Key: "Lastname", Value: arg.LastName}, {Key: "Username", Value: arg.UserName}, {Key: "Password", Value: arg.Password}, {Key: "Updated_at", Value: time.Now()}}},
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

func DeleteInstructor(ctx context.Context, id string) error {
	collection := CollectionInstructor()
	iuud, _ := primitive.ObjectIDFromHex(id)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": iuud})
	if err != nil {
		log.Fatal(err)
	}
	return err
}

type ListParams struct {
	//Owner  string `json:"owner"`
	Limit int64
	Skip  int64
}

//Find multiple documents
func ListInstructors(ctx context.Context, arg ListParams) ([]models.User, error) {
	collection := CollectionInstructor()
	//check the connection

	//find records
	//pass these options to the Find method
	findOptions := options.Find()
	//Set the limit of the number of record to find
	findOptions.SetLimit(arg.Limit)
	findOptions.SetSkip(arg.Skip)
	//Define an array in which you can store the decoded documents
	var results []models.User

	//Passing the bson.D{{}} as the filter matches  documents in the collection
	cur, err := collection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	//Finding multiple documents returns a cursor
	//Iterate through the cursor allows us to decode documents one at a time

	for cur.Next(ctx) {
		//Create a value into which the single document can be decoded
		var elem models.User
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
