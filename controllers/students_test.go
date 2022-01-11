package controllers

import (
	"context"
	"testing"
	"time"

	"github.com/E_learning/models"
	"github.com/E_learning/util"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func createStudentModel() models.Student {
	password := util.RandomString(6)
	hashpassword, _ := util.HashPassword(password)
	return models.Student{
		ID:        primitive.NewObjectID(),
		FirstName: util.RandomString(6),
		LastName:  util.RandomString(6),
		UserName:  util.RandomString(6),
		Email:     util.RandomEmail(),
		Password:  hashpassword,
		CreatedAt: time.Now(),
	}
}

func TestCreateStudent(t *testing.T) {
	args := createStudentModel()
	student, err := CreateStudent(context.Background(), &args)
	require.NoError(t, err)
	require.NotNil(t, student)
	require.WithinDuration(t, student.CreatedAt, args.CreatedAt, 10)
	student1, err := CreateStudent(context.Background(), &args)
	require.Error(t, err)
	require.NotNil(t, student1)
}

func TestFindStudent(t *testing.T) {
	args := createStudentModel()
	student1, err := CreateStudent(context.Background(), &args)
	require.NoError(t, err)
	require.NotEmpty(t, student1)
	student2, err := FindStudent(context.Background(), student1.ID.Hex())
	require.NoError(t, err)
	require.NotEmpty(t, student2)
	require.Equal(t, student1.ID, student2.ID)
	require.Equal(t, student1.FirstName, student2.FirstName)
	require.Equal(t, student1.LastName, student2.LastName)
	require.Equal(t, student1.UserName, student2.UserName)
	require.Equal(t, student1.Password, student2.Password)
	require.Equal(t, student1.Email, student2.Email)
	require.WithinDuration(t, student1.CreatedAt, student2.CreatedAt, time.Second)
}

func TestUpdateStudent(t *testing.T) {
	args := createStudentModel()
	student, err := CreateStudent(context.Background(), &args)
	require.NoError(t, err)
	require.NotEmpty(t, student)
	updateargs := UpdateInstructorParams{
		ID:        student.ID.Hex(),
		FirstName: util.RandomString(6),
		LastName:  util.RandomString(6),
		UserName:  util.RandomAuthor(),
		Email:     util.RandomEmail(),
		Password:  util.RandomString(10),
	}
	result, err := UpdateInstructor(context.Background(), updateargs)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestDeleteStudent(t *testing.T) {
	args := createStudentModel()
	student, err := CreateStudent(context.Background(), &args)
	require.NoError(t, err)
	require.NotEmpty(t, student)
	err = DeleteStudent(context.Background(), args.ID.Hex())
	require.NoError(t, err)
	instructor2, err := FindStudent(context.Background(), args.ID.Hex())
	require.EqualError(t, err, mongo.ErrNoDocuments.Error())
	require.Empty(t, instructor2)
}

func TestListStudent(t *testing.T) {
	args := ListParams{
		Limit: 1,
		Skip:  10,
	}
	results, err := ListStudents(context.Background(), args)
	require.NoError(t, err)
	require.NotNil(t, results)
	require.NotEmpty(t, results)
}
