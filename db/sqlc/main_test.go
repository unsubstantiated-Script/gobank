package db

import (
	"database/sql"
	"github.com/techschool/simplebank/db/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

//const (
//	dbDriver = "postgres"
//	dbSource = "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable"
//)

var testQueries *Queries

//Adding this to access the creation of a DB for testing
var testDB *sql.DB

//Setting up our testing connection to the DB
func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the DB:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
