package controllers

import (
	"context"
	"fmt"
	"testing"
	"time"

	//"github.com/E_learning/models"
	"github.com/E_learning/models"
	"github.com/E_learning/util"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func createInstructorModel() models.User {
	password := util.RandomString(6)
	hashpassword, _ := util.HashPassword(password)
	return models.User{
		ID:        primitive.NewObjectID(),
		FirstName: util.RandomString(6),
		LastName:  util.RandomString(6),
		UserName:  util.RandomString(6),
		Email:     util.RandomEmail(),
		Password:  hashpassword,
		CreatedAt: time.Now(),
	}
}

func TestCreateInstructor(t *testing.T) {
	args := createInstructorModel()
	instructor, err := controllers.User.CreateInstructor(context.Background(), &args)
	require.NoError(t, err)
	require.NotNil(t, instructor)
	require.WithinDuration(t, instructor.CreatedAt, args.CreatedAt, 10)
	instructor2, err := controllers.User.CreateInstructor(context.Background(), &args)
	require.Error(t, err)
	require.NotNil(t, instructor2)
}

func TestFindInstructor(t *testing.T) {
	args := createInstructorModel()
	instructor1, err := controllers.User.CreateInstructor(context.Background(), &args)
	fmt.Println(instructor1.ID.String())
	require.NoError(t, err)
	require.NotEmpty(t, instructor1)
	instructor2, err := controllers.User.FindInstructorbyId(context.Background(), instructor1.ID.Hex())
	require.NoError(t, err)
	require.NotEmpty(t, instructor2)
	require.Equal(t, instructor1.ID, instructor2.ID)
	require.Equal(t, instructor1.FirstName, instructor2.FirstName)
	require.Equal(t, instructor1.LastName, instructor2.LastName)
	require.Equal(t, instructor1.UserName, instructor2.UserName)
	require.Equal(t, instructor1.Password, instructor2.Password)
	require.Equal(t, instructor1.Email, instructor2.Email)
	require.WithinDuration(t, instructor1.CreatedAt, instructor2.CreatedAt, time.Second)
}

func TestUpdateInstructor(t *testing.T) {
	args := createInstructorModel()
	instructor1, err := controllers.User.CreateInstructor(context.Background(), &args)
	require.NoError(t, err)
	require.NotEmpty(t, instructor1)
	updateargs := UpdateInstructorParams{
		ID:        instructor1.ID.Hex(),
		FirstName: util.RandomString(6),
		LastName:  util.RandomString(6),
		UserName:  util.RandomAuthor(),
		Email:     util.RandomEmail(),
		Password:  util.RandomString(10),
	}
	result, err := controllers.User.UpdateInstructor(context.Background(), updateargs)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestDeleteInstructor(t *testing.T) {
	args := createInstructorModel()
	instructor1, err := controllers.User.CreateInstructor(context.Background(), &args)
	require.NoError(t, err)
	require.NotEmpty(t, instructor1)
	err = controllers.User.DeleteInstructor(context.Background(), args.ID.Hex())
	require.NoError(t, err)
	instructor2, err := controllers.User.FindInstructor(context.Background(), "")
	require.EqualError(t, err, mongo.ErrNoDocuments.Error())
	require.Empty(t, instructor2)
}

func TestListInstructors(t *testing.T) {
	args := ListParams{
		Limit: 1,
		Skip:  2,
	}
	results, err := controllers.User.ListInstructors(context.Background(), args)
	require.NoError(t, err)
	require.NotNil(t, results)
	require.NotEmpty(t, results)
}
