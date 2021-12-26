package controllers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPdfUpload(t *testing.T) {
	file := "/home/stephen/Documents/Waterflow.pdf"
	_, err := Pdf(file)
	require.NoError(t, err)

}

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
	material, err := FindCourseMaterial(context.Background(), "61c857530b1847aa31fe4d1b")
	require.NoError(t, err)
	require.NotNil(t, material)
}
func TestUpdateCourseMaterial(t *testing.T) {
	err := UpdateCourseMaterial(context.Background(), "61c8571a156059a708b2f8f2", "C#", "Idk you bitch!")
	require.NoError(t, err)
}
func TestDeleteCourseMaterial(t *testing.T) {
	err := DeleteCourseMaterial(context.Background(), "61c8571a156059a708b2f8f2")
	require.NoError(t, err)
}
