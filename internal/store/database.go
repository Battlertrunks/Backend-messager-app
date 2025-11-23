package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

// Open the connection to the database
func Open() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./message-app.db")
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	fmt.Println("Connected to the Database...")
	return db, nil
}

func MigrateFS(db *sql.DB, migrationFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationFS)

	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}