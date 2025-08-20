package database

import (
	"database/sql"
	"flora/configs"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import the driver
)

// Connect opens a connection to the SQLite database file specified in the config.
func Connect(cfg *configs.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.DatabasePath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the database.")
	return db, nil
}

// Migrate runs the initial database migrations to create tables.
func Migrate(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS settings (
		guild_id TEXT PRIMARY KEY,
		key TEXT NOT NULL,
		value TEXT NOT NULL
	);
	`

	_, err := db.Exec(query)
	return err
}
