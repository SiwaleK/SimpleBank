package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/SiwaleK/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account { // doesn't have the test prefix so it won't be run as unit test. it should return the created account record so other unit can have enough data to perform their own operation
	arg := CreateAccountParams{
		Owner:      util.RandomOwner(),
		Balance:    util.RandomMoney(),
		CurrencyAt: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)      //check and err must be nil if it's not it return fail
	require.NotEmpty(t, account) // return account should not be empty

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.CurrencyAt, account.CurrencyAt) // to check that the account owner and etc matches with the input
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt) // check if acc id is automatically generate
	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.CurrencyAt, account2.CurrencyAt)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.CurrencyAt, account2.CurrencyAt)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error()) // to find it the database , in this case the call should return an error
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	arg := ListAccountParams{
		Limit:  5,
		Offset: 5,
	}
	account, err := testQueries.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, account, 5)

	for _, account := range account {
		require.NotEmpty(t, account)
	}
}
