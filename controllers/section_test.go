package controllers

import (
	"context"
	"testing"

	"github.com/E_learning/models"
	"github.com/E_learning/util"
	"github.com/stretchr/testify/require"
)

func TestFindSection(t *testing.T) {
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
		result1, err := FindSectionbyTitle(context.Background(), argsec.Name, course.Author, result.Title)
		require.NoError(t, err)
		require.NotNil(t, result1)
		require.Equal(t, result.Title, result1.Title)

	}

}

func TestDeleteSection(t *testing.T) {
	args := createcourse()
	section := NewSection()
	argsec := CourseSec{
		Name:    args.Name,
		Section: section,
	}
	args.Section = section
	course, _ := CreateCourse(context.Background(), &args)
	require.NotEmpty(t, course)
	var argdel DelSection
	for _, section := range argsec.Section {
		argdel = DelSection{
			Name: course.Name,
			Id:   section.ID.Hex(),
		}
		res, err := DeleteSection(context.Background(), argdel)
		require.NoError(t, err)
		require.NotEmpty(t, res)
		sec, err := FindSection(context.Background(), argdel.Name, course.Author, "")
		require.NoError(t, err)
		require.Empty(t, sec)

	}

}

func TestAddSection(t *testing.T) {
	args := createcourse()
	section := NewSection()
	args.Section = section
	course, _ := CreateCourse(context.Background(), &args)
	require.NotEmpty(t, course)
	argsec := CourseSec{
		Name:    args.Name,
		Section: section,
	}
	result, err := AddSection(context.Background(), argsec, course.Author)
	require.NoError(t, err)
	require.NotNil(t, result)

}

func TestUpdateSection(t *testing.T) {
	args := createcourse()
	section := NewSection()
	args.Section = section
	course, _ := CreateCourse(context.Background(), &args)
	require.NotEmpty(t, course)
	argsec := CourseSec{
		Name:    args.Name,
		Section: section,
	}
	for _, section := range argsec.Section {
		args := models.Section{
			Title: util.RandomString(4),
		}
		result, err := UpdateSection(context.Background(), course.Name, section.ID.Hex(), &args)
		require.NoError(t, err)
		require.NotNil(t, result)
	}
}
