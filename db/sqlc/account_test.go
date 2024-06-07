package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simple-bank/utils"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testStore.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	testAccount1 := createRandomAccount(t)
	testAccount2, err := testStore.GetAccount(context.Background(), testAccount1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, testAccount2)

	require.Equal(t, testAccount1.ID, testAccount2.ID)
	require.Equal(t, testAccount1.Owner, testAccount2.Owner)
	require.Equal(t, testAccount1.Balance, testAccount2.Balance)
	require.Equal(t, testAccount1.Currency, testAccount2.Currency)
	//require.Equal(t, testAccount1.CreatedAt, testAccount2.Currency)
	require.WithinDuration(t, testAccount1.CreatedAt, testAccount2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	testAccount1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      testAccount1.ID,
		Balance: utils.RandomMoney(),
	}

	testAccount2, err := testStore.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, testAccount2)

	require.Equal(t, testAccount1.ID, testAccount2.ID)
	require.Equal(t, testAccount1.Owner, testAccount2.Owner)
	require.Equal(t, arg.Balance, testAccount2.Balance)
	require.Equal(t, testAccount1.Currency, testAccount2.Currency)
	//require.Equal(t, testAccount1.CreatedAt, testAccount2.Currency)
	require.WithinDuration(t, testAccount1.CreatedAt, testAccount2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	testAccount1 := createRandomAccount(t)
	err := testStore.DeleteAccount(context.Background(), testAccount1.ID)
	require.NoError(t, err)

	testAccount2, err := testStore.GetAccount(context.Background(), testAccount1.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, testAccount2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testStore.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, account.Owner, lastAccount.Owner)
	}
}
