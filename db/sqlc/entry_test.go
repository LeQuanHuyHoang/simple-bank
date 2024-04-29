package db

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"simple-bank/utils"
	"testing"
	"time"
)

var testuuid uuid.UUID

func createRandomEntry(t *testing.T) Entry {
	testAccount := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: testAccount.ID,
		Amount:    utils.RandomMoney(),
	}

	testEntry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, testEntry)

	require.NotZero(t, testEntry.ID)
	require.Equal(t, arg.AccountID, testEntry.AccountID)
	require.Equal(t, arg.Amount, testEntry.Amount)
	require.NotZero(t, testEntry.CreatedAt)

	return testEntry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	testEntry1 := createRandomEntry(t)
	testEntry2, err := testQueries.GetEntry(context.Background(), testEntry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, testEntry2)

	require.Equal(t, testEntry1.ID, testEntry2.ID)
	require.Equal(t, testEntry1.AccountID, testEntry2.AccountID)
	require.Equal(t, testEntry1.Amount, testEntry2.Amount)
	require.WithinDuration(t, testEntry1.CreatedAt, testEntry2.CreatedAt, time.Second)
}

func TestDeleteEntry(t *testing.T) {
	testEntry1 := createRandomEntry(t)
	err := testQueries.DeleteAccount(context.Background(), testEntry1.ID)
	require.NoError(t, err)

	testEntry2, err := testQueries.GetAccount(context.Background(), testEntry1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, testEntry2)
}
