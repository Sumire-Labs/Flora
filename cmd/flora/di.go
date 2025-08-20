//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	"flora/commands"
	"flora/configs"
	"flora/database"
	"log"

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
		NewDiscordSession, // This now provides a fully configured session
		NewApp,
	)
	return nil, nil
}

// NewApp creates the final App object.
func NewApp(s *discordgo.Session, db *sql.DB) *App {
	return &App{Session: s, DB: db}
}

// NewDiscordSession creates a new Discord session and registers all handlers.
func NewDiscordSession(cfg *configs.Config, cm *commands.Manager) (*discordgo.Session, error) {
	s, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, err
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
		s.UpdateGameStatus(0, "Flora")

		if err := cm.RegisterAll(s); err != nil {
			log.Fatalf("Failed to register commands: %v", err)
		}
		log.Println("Commands registered successfully.")
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if cmd, exists := cm.Get(i.ApplicationCommandData().Name); exists {
				cmd.Handler(s, i)
			}
		}
	})

	return s, nil
}
