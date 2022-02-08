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

func FindContent(ctx context.Context, name string, subsectionid string) (*models.Content, error) {
	var content models.Content
	collection := CourseCollection()
	iuud, _ := primitive.ObjectIDFromHex(subsectionid)
	pipeline := []bson.M{
		{"$match": bson.M{"Name": name}},
		{"$unwind": "$Section"},
		{"$unwind": "$Section.Content"},
		{"$match": bson.M{"Section.Content.subsectionid": iuud}},
		{"$project": bson.M{
			"Subsection_Title": "$Section.Content.Subsection_Title",
			"SubContent":       "$Section.Content.SubContent",
			"_id":              "$Section.Content.subsectionid",
		}},
	}
	iter, err := collection.Aggregate(ctx, pipeline)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			log.Println("No such document")
			return nil, err
		}
		log.Fatal(err)
	}
	var results []bson.M
	if err = iter.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		fmt.Printf("SubContent %s Subsection_Title %s  _id %s times\n", result["SubContent"], result["Subsection_Title"], result["_id"])
		content.Content = result["SubContent"].(string)
		content.SubTitle = result["Subsection_Title"].(string)
		content.ID = result["_id"].(primitive.ObjectID)
	}
	fmt.Println("Content", content)
	return &content, nil
}

type DelContent struct {
	CourseName   string
	SubsectionId string
}

func DeleteContent(ctx context.Context, arg DelContent) (*mongo.UpdateResult, error) {
	collection := CourseCollection()
	iuud, _ := primitive.ObjectIDFromHex(arg.SubsectionId)
	filter := bson.D{primitive.E{Key: "Name", Value: arg.CourseName}}

	update := bson.M{
		"$pull": bson.M{
			"Section.$[].Content": bson.D{primitive.E{Key: "subsectionid", Value: iuud}},
		},
	}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return result, err
}

type CourseSubSection struct {
	Id         string            `uri:"sectionid"`
	CourseName string            `uri:"name"`
	Content    []*models.Content `json:"Content"`
}

func AddContent(ctx context.Context, arg CourseSubSection, author string) (*mongo.UpdateResult, error) {
	collection := CourseCollection()
	course, err := FindCoursebyName(ctx, arg.CourseName)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoSuchDocument
		}
	}
	if course.Author != author {
		return nil, ErrInvalidUser
	}
	iuud, _ := primitive.ObjectIDFromHex(arg.Id)
	fmt.Println(iuud)
	match := bson.M{"Name": arg.CourseName, "Section._id": iuud}
	for _, v := range arg.Content {
		v.ID = primitive.NewObjectID()
	}
	change := bson.M{"$push": bson.M{"Section.$.SubContent": bson.M{"$each": arg.Content}}}
	result, err := collection.UpdateOne(ctx, match, change)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("success", result)
	}
	return result, err
}

func UpdateSectionContent(ctx context.Context, name string, subsectionid string, sectiontitle string, arg *models.Content) (*mongo.UpdateResult, error) {
	collection := CourseCollection()
	filter := bson.D{primitive.E{Key: "Name", Value: name}}
	iuud, _ := primitive.ObjectIDFromHex(subsectionid)
	arrayFilters := options.ArrayFilters{Filters: bson.A{bson.M{"x.Title": sectiontitle}, bson.M{"y.subsectionid": iuud}}}
	upsert := true
	opts := options.UpdateOptions{
		ArrayFilters: &arrayFilters,
		Upsert:       &upsert,
	}
	update := bson.M{
		"$set": bson.M{
			"Section.$[x].Content.$[y].Subsection_Title": arg.SubTitle,
			"Section.$[x].Content.$[y].SubContent":       arg.Content,
		},
	}
	result, err := collection.UpdateOne(ctx, filter, update, &opts)
	if err != nil {
		fmt.Printf("error updating db: %+v\n", err)
	}
	return result, err
}
