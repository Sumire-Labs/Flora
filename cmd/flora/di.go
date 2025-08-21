//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	"flora/commands"
	"flora/configs"
	"flora/database"
	"flora/handlers"

	"github.com/bwmarrin/discordgo"
	"github.com/google/wire"
)

// App holds the core components of the application.
type App struct {
	Session *discordgo.Session
	DB      *sql.DB
}

// InitializeApp creates the main application with all its dependencies.
func InitializeApp(cfg *configs.Config) (*App, error) {
			wire.Build(
			database.Connect,
			commands.NewManager,
			handlers.NewEventHandler,
			NewDiscordSession,
			NewApp,
		)
	return nil, nil
}

// NewApp creates the final App object.
func NewApp(s *discordgo.Session, db *sql.DB) *App {
	return &App{Session: s, DB: db}
}

// NewDiscordSession creates a new Discord session and registers all handlers.
func NewDiscordSession(cfg *configs.Config, eh *handlers.EventHandler) (*discordgo.Session, error) {
	s, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, err
	}

	// Add necessary intents.
	s.Identify.Intents |= discordgo.IntentGuildMessages | discordgo.IntentMessageContent | discordgo.IntentGuildMembers

	// Register event handlers from the EventHandler
	s.AddHandler(eh.OnReady)
	s.AddHandler(eh.OnInteractionCreate)
	s.AddHandler(eh.OnMessageDelete)
	s.AddHandler(eh.OnMemberAdd)
	s.AddHandler(eh.OnMemberRemove)
	s.AddHandler(eh.OnMemberUpdate)

	return s, nil
}
















