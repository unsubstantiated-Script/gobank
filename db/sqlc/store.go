package db

import (
	"context"
	"database/sql"
	"fmt"
)

//Store provides all functions to execute db queries and transactions
type Store struct {
	//Composing a struct here instead of inheritance
	*Queries
	db *sql.DB
}

//This will create a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTX creates a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	// Getting a new queries object
	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input params of the transfer transaction
type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	//After amount is removed
	FromAccount Account `json:"from_account"`
	//After amount is added
	ToAccount Account `json:"to_account"`
	//Recording money going out
	FromEntry Entry `json:"from_entry"`
	//Recording money moving in
	ToEntry Entry `json:"to_entry"`
}

//Second bracket means we are creating a new empty object of that type. For debugging
//var txKey = struct{}{}

// TransferTx performs a money transfer from one account to the other.
// Creates a transfer record, adds account entries, and updates accounts' balance within a single DB transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	//Creating an empty result
	var result TransferTxResult

	//Setting up the transfer package
	err := store.execTx(ctx, func(q *Queries) error {

		var err error

		//For debugging
		//txName := ctx.Value(txKey)

		//Creating a transfer
		//For debugging
		//fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		//Account Entries
		// Removal From Account Entry
		//For debugging
		//fmt.Println(txName, "create Entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// Addition TO Account Entry
		//For debugging
		//fmt.Println(txName, "create Entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//Updating Balances
		//get account -> update its balance
		//For debugging
		//fmt.Println(txName, "Get account 1 for update")
		//account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountId)
		//if err != nil {
		//	return err
		//}

		//For debugging
		//fmt.Println(txName, "Update account 1 balance")
		//result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		//	ID:      arg.FromAccountId,
		//	Balance: account1.Balance - arg.Amount,
		//})
		//if err != nil {
		//	return err
		//}

		//Combo-ing jazz now
		//For debugging
		//fmt.Println(txName, "Update account 1 balance")
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountId,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		//get account -> update its balance
		//For debugging
		//fmt.Println(txName, "Get account 2 for update")
		//account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountId)
		//if err != nil {
		//	return err
		//}
		//
		////For debugging
		////fmt.Println(txName, "Update account 2 balance")
		//result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		//	ID:      arg.ToAccountId,
		//	Balance: account2.Balance + arg.Amount,
		//})
		//if err != nil {
		//	return err
		//}

		//For debugging
		//fmt.Println(txName, "Update account 2 balance")
		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountId,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
