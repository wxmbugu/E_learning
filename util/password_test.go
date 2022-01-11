package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(10)

	hashpassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotNil(t, hashpassword)

	err = CheckPassword(password, hashpassword)
	require.NoError(t, err)

	password2 := RandomString(10)
	err = CheckPassword(password2, hashpassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashpassword1, err := HashPassword(password2)
	require.NoError(t, err)
	require.NotEqual(t, hashpassword, hashpassword1)
}
