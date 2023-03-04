package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	//Testing by running concurrent transfer transactions
	n := 5
	amount := int64(10)

	//This is how we test stuff running concurrently as we can't check it within the routine
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})

			//Sending things to their appropriate channel
			errs <- err
			results <- result
		}()
	}

	//Receiving the data from the channels
	for i := 0; i < n; i++ {
		//Making sure no errors
		err := <-errs
		require.NoError(t, err)

		//Making sure results aren't empty
		result := <-results
		require.NotEmpty(t, result)

		//Now we check all the things...

		//Transfer
		transfer := result.Transfer

		//Making sure the transfer is created
		require.NotEmpty(t, transfer)

		//Making sure the IDs line up
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)

		//Making sure the amount of money is correct
		require.Equal(t, amount, transfer.Amount)

		//Making sure the auto-increment ID for transfer is jiving as well as the timestamp
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// checking From Entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		//Making sure this really got created
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// checking To Entries
		toEntry := result.ToEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		//Making sure this really got created
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//	TODO: Check accounts' balance
	}

}
