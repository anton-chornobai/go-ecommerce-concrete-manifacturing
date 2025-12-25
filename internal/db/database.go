package db

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func Connect(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		log.Fatalf("Failde to open db, %v", err)
	}

	return db
}