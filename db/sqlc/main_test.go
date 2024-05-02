package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"simple-bank/utils"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't connect to database", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
