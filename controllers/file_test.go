package controllers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownloadFile(t *testing.T) {
	filename := DownloadFile("Waterflow.pdf")
	require.FileExists(t, "/home/stephen/go/src/github.com/E_learning/controllers/Waterflow.pdf")
	require.NotNil(t, filename)
}

func TestDeleteFile(t *testing.T) {
	err := DeleteFile(context.Background(), "61c85fa2cf57cf61ee5f3fda")
	//testing for a file that doesn.t exist
	require.Error(t, err)
}
func TestFindFile(t *testing.T) {
	cursor, err := FindFile(context.Background(), "61c85fa2cf57cf61ee5f3fda")
	require.NoError(t, err)
	require.NotNil(t, cursor)
}
