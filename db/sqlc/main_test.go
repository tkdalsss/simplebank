package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	// to tell go formatter to keep
	_ "github.com/lib/pq"
	// "github.com/lib/pq" -> auto removed when we save this file
	// go mod tidy -> clean up the dependencies
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	// conn, err := sql.Open(dbDriver, dbSource)
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil { // error is not null -> error exists
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
