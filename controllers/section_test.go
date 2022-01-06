package controllers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindSetion(t *testing.T) {
	v, err := FindSection(context.Background(), "btuemk", "61d5e29b1f61c9a1dd1f9e7d")
	require.NoError(t, err)
	require.NotNil(t, v)
}
