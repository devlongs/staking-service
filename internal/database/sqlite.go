package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// NewSQLiteDB opens a SQLite connection and creates the stakes table if it does not exit
func NewSQLiteDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS stakes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		wallet_address TEXT NOT NULL,
		amount REAL NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(query); err != nil {
		log.Printf("Failed to create stakes table: %v", err)
		return nil, err
	}

	return db, nil
}
