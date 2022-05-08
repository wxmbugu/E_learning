package controllers

import (
	"context"
	"testing"

	"github.com/E_learning/models"
	"github.com/E_learning/util"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestFindContent(t *testing.T) {
	args := createcourse()
	section := NewSection()
	argsec := CourseSec{
		Name:    args.Name,
		Section: section,
	}
	args.Section = section
	course, _ := CreateCourse(context.Background(), &args)
	for _, section := range argsec.Section {
		result, err := FindSection(context.Background(), argsec.Name, course.Author, section.ID.Hex())
		require.NoError(t, err)
		require.NotNil(t, result)
		//var id string
		for _, v := range section.Content {
			content, err := FindContent(context.Background(), argsec.Name, v.ID.Hex())
			require.NoError(t, err)
			require.NotEmpty(t, content)
			require.Equal(t, content.SubTitle, v.SubTitle)
			require.Equal(t, v.ID.Hex(), content.ID.Hex())
			content2, err2 := FindContent(context.Background(), argsec.Name, "1")
			require.NoError(t, err2)
			require.Empty(t, content2)

		}
		//require.Equal(t, section.Title, result.Title)
		_, err = FindSection(context.Background(), "", "", "")
		require.EqualError(t, err, mongo.ErrNoDocuments.Error())
	}

}

func TestDeleteContent(t *testing.T) {
	args := createcourse()
	section := NewSection()
	argsec := CourseSec{
		Name:    args.Name,
		Section: section,
	}
	args.Section = section
	course, _ := CreateCourse(context.Background(), &args)
	require.NotEmpty(t, course)
	for _, section := range argsec.Section {
		result, err := FindSection(context.Background(), argsec.Name, course.Author, section.ID.Hex())
		require.NoError(t, err)
		require.NotNil(t, result)
		for _, v := range section.Content {
			args := DelContent{
				CourseName:   course.Name,
				SubsectionId: v.ID.Hex(),
			}

			_, err := DeleteContent(context.Background(), args)
			require.NoError(t, err)
			content2, err := FindContent(context.Background(), argsec.Name, v.ID.Hex())
			require.NoError(t, err)
			require.Empty(t, content2)
		}
		//require.Equal(t, section.Title, result.Title)
		_, err = FindSection(context.Background(), "", "", "")
		require.EqualError(t, err, mongo.ErrNoDocuments.Error())
	}
}

func TestAddContent(t *testing.T) {
	args := createcourse()
	section := NewSection()
	argsec := CourseSec{
		Name:    args.Name,
		Section: section,
	}
	args.Section = section
	course, _ := CreateCourse(context.Background(), &args)
	require.NotEmpty(t, course)
	for _, section := range argsec.Section {
		args := CourseSubSection{
			Id:         section.ID.Hex(),
			CourseName: course.Name,
			Content:    section.Content,
		}
		result, err := AddContent(context.Background(), args, course.Author)
		require.NoError(t, err)
		require.NotNil(t, result)
	}

}

func TestUpdateCourseSubSection(t *testing.T) {
	args := createcourse()
	section := NewSection()
	argsec := CourseSec{
		Name:    args.Name,
		Section: section,
	}
	args.Section = section
	course, _ := CreateCourse(context.Background(), &args)
	require.NotEmpty(t, course)
	for _, section := range argsec.Section {
		args := models.Content{
			SubTitle:   util.RandomString(4),
			SubContent: util.RandomString(100),
		}
		for _, cont := range section.Content {
			result, err := UpdateSectionContent(context.Background(), course.Name, cont.ID.Hex(), section.Title, &args)
			require.NoError(t, err)
			require.NotNil(t, result)
		}
	}
}
