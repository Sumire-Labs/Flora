package database

import (
	"database/sql"
	"flora/configs"
	"log"

	_ "modernc.org/sqlite" // Import the pure Go driver
)

// Connect opens a connection to the SQLite database file specified in the config.
func Connect(cfg *configs.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite", cfg.DatabasePath)
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
	CREATE TABLE IF NOT EXISTS guild_settings (
		guild_id TEXT PRIMARY KEY,
		log_channel_id TEXT
	);
	`

	_, err := db.Exec(query)
	return err
}

// SetLogChannel sets or updates the log channel for a specific guild.
func SetLogChannel(db *sql.DB, guildID, channelID string) error {
	query := `
	INSERT INTO guild_settings (guild_id, log_channel_id)
	VALUES (?, ?)
	ON CONFLICT(guild_id) DO UPDATE SET log_channel_id = excluded.log_channel_id;
	`
	_, err := db.Exec(query, guildID, channelID)
	return err
}
