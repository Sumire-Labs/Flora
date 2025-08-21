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

// Migrate runs the initial database migrations to create tables and add columns.
func Migrate(db *sql.DB) error {
	// guild_settings table
	query1 := `
	CREATE TABLE IF NOT EXISTS guild_settings (
		guild_id TEXT PRIMARY KEY,
		log_channel_id TEXT
	);
	`
	if _, err := db.Exec(query1); err != nil {
		return err
	}

	// Add columns to guild_settings for the ticket system.
	// We ignore errors here because the columns might already exist.
	db.Exec("ALTER TABLE guild_settings ADD COLUMN ticket_panel_channel_id TEXT;")
	db.Exec("ALTER TABLE guild_settings ADD COLUMN ticket_category_id TEXT;")
	db.Exec("ALTER TABLE guild_settings ADD COLUMN support_role_id TEXT;")
	db.Exec("ALTER TABLE guild_settings ADD COLUMN transcript_channel_id TEXT;")

	// tickets table
	query2 := `
	CREATE TABLE IF NOT EXISTS tickets (
		ticket_id INTEGER PRIMARY KEY AUTOINCREMENT,
		guild_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		channel_id TEXT NOT NULL,
		status TEXT NOT NULL
	);
	`
	_, err := db.Exec(query2)
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

// GetLogChannel retrieves the log channel for a specific guild.
// It returns the channel ID and a boolean indicating if it was found.
func GetLogChannel(db *sql.DB, guildID string) (string, bool, error) {
	var channelID sql.NullString
	query := `SELECT log_channel_id FROM guild_settings WHERE guild_id = ?`
	err := db.QueryRow(query, guildID).Scan(&channelID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}

	if channelID.Valid {
		return channelID.String, true, nil
	} else {
		return "", false, nil
	}
}

// SetTicketPanelChannel sets or updates the ticket panel channel for a specific guild.
func SetTicketPanelChannel(db *sql.DB, guildID, channelID string) error {
	query := `
	INSERT INTO guild_settings (guild_id, ticket_panel_channel_id)
	VALUES (?, ?)
	ON CONFLICT(guild_id) DO UPDATE SET ticket_panel_channel_id = excluded.ticket_panel_channel_id;
	`
	_, err := db.Exec(query, guildID, channelID)
	return err
}
