package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/db/util"
	"testing"
	"time"
)

//Basic setup util
func CreateRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(), //Randomly generate this?
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	//Require test package methods for handling errors
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	//Checking Account exists
	require.NotZero(t, account.ID)

	//Checking TimeStamp
	require.NotZero(t, account.CreatedAt)

	return account
}

//Testing Create Account
func TestQueries_CreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

//Testing Read Account
func TestQueries_GetAccount(t *testing.T) {
	//Creating test accounts
	account1 := CreateRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	//Making sure the error isn't nil in the creation/duplication thing
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	//Testing all the things...
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

//Testing Update Account
func TestQueries_UpdateAccount(t *testing.T) {
	//Creating test account
	account1 := CreateRandomAccount(t)

	//Setting up update arguments
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	//Creating account2 or an error
	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	//Testing to see if there was an error or a failure to create the account
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	//Testing all the things...
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, arg.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

// Testing Delete Account
//func TestQueries_DeleteAccount(t *testing.T) {
//	//Creating test account
//	account1 := CreateRandomAccount(t)
//	err := testQueries.DeleteAccount(context.Background(), account1.ID)
//	require.NoError(t, err)
//
//	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
//	//Making sure no error with the delete
//	require.Error(t, err)
//	//Making sure the delete removed all items in the row
//	require.EqualError(t, err, sql.ErrNoRows.Error())
//	//Making sure the appropriate account object is empty
//	require.Empty(t, account2)
//}

func TestQueries_ListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}

	// We expect to get at least 5 records in the DB
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
