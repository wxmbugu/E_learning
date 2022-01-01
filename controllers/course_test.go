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

func createcoursematerial() models.CourseMaterial {
	file := "/home/stephen/Documents/Waterflow.pdf"
	id, _ := Pdf(file)
	//ide, _ := primitive.ObjectIDFromHex(id)
	material := models.CourseMaterial{
		ID:        primitive.NewObjectID(),
		Author:    []string{util.RandomAuthor()},
		PdfFileID: []primitive.ObjectID{id},
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
	arg := createcourse()
	_, err := CreateCourse(context.Background(), &arg)
	if err != nil {
		log.Fatal(err)
	}
	id := arg.ID

	course, err := FindCourse(context.Background(), id.String())
	require.Error(t, err)
	require.NotNil(t, course)
	//require.Equal(t, course.Author, arg.Author)
	//require.Equal(t, course.Description, arg.Description)
}
func TestUpdateCourse(t *testing.T) {
	args := UpdateCourseParams{
		ID:          "32",
		Name:        "Idk",
		Description: "Do better!",
	}
	results, err := UpdateCourse(context.Background(), args)
	require.NoError(t, err)
	require.NotNil(t, results)
}

func TestDeleteCourse(t *testing.T) {
	err := DeleteCourse(context.Background(), "61c6279e2febff924341004c")
	require.NoError(t, err)
}

func TestListCourse(t *testing.T) {
	arg := ListCoursesParams{
		Limit: 10,
		Skip:  1,
	}
	results, err := ListCourses(context.Background(), arg)
	require.NoError(t, err)
	require.NotNil(t, results)
}
