package controllers

import (
	"context"
	"fmt"
	"testing"

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
		result, err := FindSection(context.Background(), args.Name, section.ID.Hex())
		fmt.Println(result)
		fmt.Println(args.Name)
		require.NoError(t, err)
		require.NotNil(t, result)
		_, err = FindSection(context.Background(), "", "ok")
		require.Error(t, err)
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
	for _, section := range argsec.Section {
		argdel := DelSection{
			Name: course.Name,
			Id:   section.ID.Hex(),
		}
		res, err := DeleteSection(context.Background(), argdel)
		require.NoError(t, err)
		require.NotEmpty(t, res)
		_, err = FindSection(context.Background(), course.Name, section.ID.Hex())
		require.Error(t, err)

	}
}
