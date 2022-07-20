package controllers

import (
	"context"
	"testing"

	"github.com/E_learning/models"
	"github.com/E_learning/util"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewSection() models.Section {
	section := models.Section{
		ID:    primitive.NewObjectID(),
		Title: util.RandomString(10),
	}
	results, _ := controllers.Section.AddSection(context.Background(), section)
	return results
}
func TestFindSection(t *testing.T) {
	section := NewSection()
	foundsection, err := controllers.Section.FindSection(context.Background(), section.ID.Hex())
	require.NoError(t, err)
	require.Equal(t, section, foundsection)
	_, err = controllers.Section.FindSection(context.Background(), "")
	require.EqualError(t, err, mongo.ErrNoDocuments.Error())
}

func TestDeleteSection(t *testing.T) {
	section := NewSection()
	err := controllers.Section.DeleteSection(context.Background(), section.ID.Hex())
	require.NoError(t, err)
	_, err = controllers.Section.FindSection(context.Background(), section.ID.Hex())
	require.EqualError(t, err, mongo.ErrNoDocuments.Error())
}

func TestAddSection(t *testing.T) {
	args := NewSection()
	require.NotEmpty(t, args)

}

func TestUpdateSection(t *testing.T) {
	args := NewSection()
	result, err := controllers.Section.UpdateSection(context.Background(), args.ID.Hex(), "NewSectionTitle")
	require.NoError(t, err)
	require.Equal(t, 1, result.ModifiedCount)
}
