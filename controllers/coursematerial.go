package controllers

import (
	"context"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/E_learning/db"
	"github.com/E_learning/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionCourseMaterial = "CourseMaterial"
)

func CourseMaterialCollection() *mongo.Collection {
	db, err := db.DBInstance()
	if err != nil {
		log.Fatal(err)
	}
	collection := db.OpenCollection(context.Background(), collectionCourseMaterial)
	return collection
}



func CreateCourseMaterial(ctx context.Context, material *models.CourseMaterial) (*models.CourseMaterial, error) {
	collection := CourseMaterialCollection()
	_, err := collection.InsertOne(ctx, material)
	if err != nil {
		log.Fatal(err)
	}
	return material, err
}

// find one course
func FindCourseMaterial(ctx context.Context, id string) (models.CourseMaterial, error) {
	collection := CourseMaterialCollection()
	var results models.CourseMaterial
	iuud, _ := primitive.ObjectIDFromHex(id)
	err := collection.FindOne(ctx, bson.M{"_id": iuud}).Decode(&results)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		log.Fatal(err)
	}
	//log.Fatal(err)
	fmt.Print(results)
	//log.Fatal(err)
	return results, err
}

func UpdateCourseMaterial(ctx context.Context, id string, name, description string) error {
	collection := CourseMaterialCollection()
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
func DeleteCourseMaterial(ctx context.Context, id string) error {
	collection := CourseMaterialCollection()
	iuud, _ := primitive.ObjectIDFromHex(id)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": iuud})
	if err != nil {
		log.Fatal(err)
	}
	return err
}
