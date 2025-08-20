package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	DiscordToken string
	DatabasePath string
}

// NewConfig creates a new Config struct from environment variables.
func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("DISCORD_BOT_TOKEN")
	dbPath := os.Getenv("DATABASE_PATH")

	if token == "" || dbPath == "" {
		log.Fatal("DISCORD_BOT_TOKEN or DATABASE_PATH is not set in .env file")
	}

	return &Config{
		DiscordToken: token,
		DatabasePath: dbPath,
	}
}
