package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	connStr := os.Getenv("DB_TEST_CONN_STR")
	if connStr == "" {
		log.Fatal("connStr doesnt exist in enviroment variables")
	}
	var err error

	testDB, err = OpenPostgre(connStr)
	if err != nil {
		log.Fatal("failed to open connection to DB:", err)
	}

	defer testDB.Close()
	os.Exit(m.Run())
}

func TestDBConnection(t *testing.T) {
	if err := testDB.Ping(); err != nil {
		t.Fatalf("expected db to be reachable, got: %v", err)
	}
}
