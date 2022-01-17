package controllers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
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
		require.Equal(t, section.Title, result.Title)
		_, err = FindSection(context.Background(), "", "", "")
		require.EqualError(t, err, mongo.ErrNoDocuments.Error())
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
		sec, _ := FindSection(context.Background(), argdel.Name, "", section.ID.Hex())
		require.Empty(t, sec)

	}

}
