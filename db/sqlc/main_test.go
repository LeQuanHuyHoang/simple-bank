package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"simple-bank/utils"
	"testing"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("can't connect to database", err)
	}

	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
