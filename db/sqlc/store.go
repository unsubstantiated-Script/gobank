package db

import "database/sql"

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
