package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Open up the database for the application
func Open() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "host=localhost user=sqlite password=sqlite dbname=message-db port=5432 sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	fmt.Println("Connected to SQLite Database...")
	return db, nil
}