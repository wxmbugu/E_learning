package controllers

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/E_learning/models"
	"github.com/E_learning/util"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createcourse() models.Course {
	ar := createcoursematerial()

	material, err := CreateCourseMaterial(context.Background(), &ar)
	if err != nil {
		log.Fatal(err)
	}
	//section := NewSection()
	arg := models.Course{
		ID:               primitive.NewObjectID(),
		Name:             util.RandomAuthor(),
		Author:           util.RandomAuthor(),
		Description:      util.RandomString(40),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		CourseMaterialID: material.ID,
	}
	return arg
}
func NewSection() []*models.Section {
	result := []*models.Section{}
	section := models.Section{
		ID:      primitive.NewObjectID(),
		Title:   util.RandomString(10),
		Content: util.RandomString(1000),
	}
	result = append(result, &section)
	return result
}

func createcoursematerial() models.CourseMaterial {
	//ide, _ := primitive.ObjectIDFromHex(id)
	material := models.CourseMaterial{
		ID:        primitive.NewObjectID(),
		Author:    []string{util.RandomAuthor()},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return material
}
func TestCourseCollection(t *testing.T) {
	collection := CourseCollection()
	require.NotEmpty(t, collection)
}

func TestCreateCourse(t *testing.T) {
	arg := createcourse()
	course, err := CreateCourse(context.Background(), &arg)
	require.NoError(t, err)
	require.NotEmpty(t, &course)
	require.Equal(t, course.Author, arg.Author)
	require.Equal(t, course.Description, arg.Description)

}

func TestFindCourse(t *testing.T) {
	course := createcourse()
	arg, err := CreateCourse(context.Background(), &course)
	require.NoError(t, err)
	course, err = FindCourse(context.Background(), arg.ID.Hex())
	require.NoError(t, err)
	require.NotNil(t, course)
	require.Equal(t, course.Author, arg.Author)
	require.Equal(t, course.Description, arg.Description)
}
func TestUpdateCourse(t *testing.T) {
	course := createcourse()
	arg, err := CreateCourse(context.Background(), &course)
	require.NoError(t, err)
	updateargs := UpdateCourseParams{
		ID:          arg.ID.Hex(),
		Name:        util.RandomString(6),
		Description: util.RandomString(100),
	}
	results, err := UpdateCourse(context.Background(), updateargs)
	require.NoError(t, err)
	require.NotNil(t, results)
}

func TestDeleteCourse(t *testing.T) {
	course := createcourse()
	arg, err := CreateCourse(context.Background(), &course)
	require.NoError(t, err)
	err = DeleteCourse(context.Background(), arg.ID.Hex())
	require.NoError(t, err)
	result, err := FindCourse(context.Background(), arg.ID.Hex())
	require.Error(t, err)
	require.Empty(t, result)
}

func TestListCourse(t *testing.T) {
	arg := ListParams{
		Limit: 10,
		Skip:  1,
	}
	results, err := ListCourses(context.Background(), arg)
	require.NoError(t, err)
	require.NotNil(t, results)
}
