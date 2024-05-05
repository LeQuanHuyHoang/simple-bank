package utils

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswword(t *testing.T) {
	password := RandomString(6)

	hasdedPassword1, err := HasdPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hasdedPassword1)

	err = CheckPassword(password, hasdedPassword1)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hasdedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hasdedPassword2, err := HasdPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hasdedPassword2)
	require.NotEqual(t, hasdedPassword2, hasdedPassword1)
}
