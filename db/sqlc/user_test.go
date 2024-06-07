package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"simple-bank/utils"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	hasdedPassword, err := utils.HasdPassword(utils.RandomString(6))
	args := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hasdedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := utils.RandomOwner()
	updateUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, updateUser.FullName, oldUser.FullName)
	require.Equal(t, updateUser.Username, oldUser.Username)
	require.Equal(t, updateUser.Email, oldUser.Email)
	require.Equal(t, updateUser.HashedPassword, oldUser.HashedPassword)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := utils.RandomEmail()
	updateUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, updateUser.Email, oldUser.Email)
	require.Equal(t, updateUser.Username, oldUser.Username)
	require.Equal(t, updateUser.FullName, oldUser.FullName)
	require.Equal(t, updateUser.HashedPassword, oldUser.HashedPassword)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := utils.RandomString(6)
	newHasdPassword, err := utils.HasdPassword(newPassword)
	require.NoError(t, err)
	updateUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text{
			String: newHasdPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, updateUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updateUser.HashedPassword, newHasdPassword)
	require.Equal(t, updateUser.Username, oldUser.Username)
	require.Equal(t, updateUser.FullName, oldUser.FullName)
	require.Equal(t, updateUser.Email, oldUser.Email)
}

func TestUpdateUserAllField(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := utils.RandomOwner()
	newEmail := utils.RandomEmail()
	newPassword := utils.RandomString(6)
	newHasdPassword, err := utils.HasdPassword(newPassword)
	require.NoError(t, err)
	updateUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text{
			String: newHasdPassword,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, updateUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updateUser.HashedPassword, newHasdPassword)
	require.Equal(t, updateUser.FullName, newFullName)
	require.Equal(t, updateUser.Email, newEmail)
	require.Equal(t, updateUser.Username, oldUser.Username)
	require.NotEqual(t, updateUser.FullName, oldUser.FullName)
	require.NotEqual(t, updateUser.Email, oldUser.Email)
}
