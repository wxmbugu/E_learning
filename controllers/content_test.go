package controllers

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/E_learning/models"
	"github.com/E_learning/util"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func createContent() models.Content {
	content := models.Content{
		ID:    primitive.NewObjectID(),
		Title: util.RandomString(10),
		Video: util.RandomString(18),
	}
	content, err := controllers.Content.AddContent(context.Background(), content)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

func TestFindContent(t *testing.T) {
	content := createContent()
	foundcontent, err := controllers.Content.FindContent(context.Background(), content.ID.Hex())
	require.NoError(t, err)
	require.Equal(t, foundcontent, content)
	_, err = controllers.Content.FindContent(context.Background(), "")
	require.EqualError(t, err, mongo.ErrNoDocuments.Error())
}

func TestDeleteContent(t *testing.T) {
	content := createContent()
	err := controllers.Content.DeleteContent(context.Background(), content.ID.Hex())
	require.NoError(t, err)
	_, err = controllers.Content.FindContent(context.Background(), content.ID.Hex())
	require.EqualError(t, err, mongo.ErrNoDocuments.Error())
}

func TestAddContent(t *testing.T) {
	content := models.Content{
		ID:    primitive.NewObjectID(),
		Title: util.RandomString(10),
		Video: util.RandomString(18),
	}
	newcontent, err := controllers.Content.AddContent(context.Background(), content)
	require.Equal(t, content.Title, newcontent.Title)
	require.NoError(t, err)

}

func TestUpdateSectionTitle(t *testing.T) {
	content := createContent()
	result, err := controllers.Content.UpdateContentTitle(context.Background(), content.ID.Hex(), "newtitle")
	fmt.Println(result)
	require.Equal(t, 1, result.ModifiedCount)
	require.NoError(t, err)
}
func TestUpdateContentVideo(t *testing.T) {
	content := createContent()
	result, err := controllers.Content.UpdateContentVideo(context.Background(), content.ID.Hex(), "newvideourl")
	fmt.Println(result)
	require.Equal(t, 1, result.ModifiedCount)
	require.NoError(t, err)
}
