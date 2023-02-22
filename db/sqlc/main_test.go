package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	// to tell go formatter to keep
	_ "github.com/lib/pq"
	"github.com/techschool/simplebank/util"
	// "github.com/lib/pq" -> auto removed when we save this file
	// go mod tidy -> clean up the dependencies
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	// var err error
	// conn, err := sql.Open(dbDriver, dbSource)
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil { // error is not null -> error exists
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
