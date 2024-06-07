package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simple-bank/utils"
	"testing"
	"time"
)

func createRandomTransfer(t *testing.T) Transfer {
	testAccount1 := createRandomAccount(t)
	testAccount2 := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: testAccount1.ID,
		ToAccountID:   testAccount2.ID,
		Amount:        utils.RandomMoney(),
	}

	testTransfer, err := testStore.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, testTransfer)

	require.NotZero(t, testTransfer.ID)
	require.NotZero(t, testTransfer.CreatedAt)
	require.Equal(t, arg.FromAccountID, testTransfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, testTransfer.ToAccountID)
	require.Equal(t, arg.Amount, testTransfer.Amount)

	return testTransfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	testTransfer1 := createRandomTransfer(t)
	testTransfer2, err := testStore.GetTransfer(context.Background(), testTransfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, testTransfer2)

	require.Equal(t, testTransfer1.ID, testTransfer2.ID)
	require.Equal(t, testTransfer1.FromAccountID, testTransfer2.FromAccountID)
	require.Equal(t, testTransfer1.ToAccountID, testTransfer2.ToAccountID)
	require.Equal(t, testTransfer1.Amount, testTransfer2.Amount)
	require.WithinDuration(t, testTransfer1.CreatedAt, testTransfer2.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	testTransfer1 := createRandomTransfer(t)
	err := testStore.DeleteAccount(context.Background(), testTransfer1.ID)
	require.NoError(t, err)

	testTransfer2, err := testStore.GetAccount(context.Background(), testTransfer1.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, testTransfer2)
}

//func TestListTransfer(t *testing.T) {
//	for i := 0; i < 10; i++ {
//		createRandomTransfer(t)
//	}
//
//	arg := ListTransferParams{
//		Limit:  5,
//		Offset: 5,
//	}
//
//	transfers, err := testStore.ListTransfer(context.Background(), arg)
//	require.NoError(t, err)
//	require.Len(t, transfers, 5)
//
//	for _, transfer := range transfers {
//		require.NotEmpty(t, transfer)
//	}
//}
