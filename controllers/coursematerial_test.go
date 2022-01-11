package controllers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)



func TestCourseMaterial(t *testing.T) {

	arg := createcoursematerial()
	courseMaterial, err := CreateCourseMaterial(context.Background(), &arg)
	require.NoError(t, err)
	require.NotEmpty(t, &courseMaterial)
	require.Equal(t, courseMaterial.Author, arg.Author)
	require.Equal(t, courseMaterial.PdfFileID, arg.PdfFileID)
	require.WithinDuration(t, courseMaterial.CreatedAt, arg.CreatedAt, 10)
	require.WithinDuration(t, courseMaterial.UpdatedAt, arg.UpdatedAt, 10)
}

func TestFindCourseMaterial(t *testing.T) {
	material, err := FindCourseMaterial(context.Background(), "61cca5e671cf508291edbacd")
	require.Error(t, err)
	require.Empty(t, material)
}
func TestUpdateCourseMaterial(t *testing.T) {
	err := UpdateCourseMaterial(context.Background(), "61cca5e671cf508291edbacd", "C#", "Idk you bitch!")
	require.NoError(t, err)
}
func TestDeleteCourseMaterial(t *testing.T) {
	err := DeleteCourseMaterial(context.Background(), "61cca72e1f7151c5052e5fab")
	require.NoError(t, err)
}
