package controllers

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindSetion(t *testing.T) {
<<<<<<< HEAD
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
=======
	_, err := FindSection(context.Background(), "btuemk", "61d5e29b1f61c9a1dd1f9e7d")
	require.Error(t, err)
	
>>>>>>> 323c2ab0bb4ea66a9424d89bd0119f51116147b7
}
