//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	"flora/commands"
	"flora/configs"
	"flora/database"
	"flora/pkg/ui"
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
func NewDiscordSession(cfg *configs.Config, cm *commands.Manager, db *sql.DB) (*discordgo.Session, error) {
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
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if cmd, exists := cm.Get(i.ApplicationCommandData().Name); exists {
				cmd.Handler(s, i)
			}
		case discordgo.InteractionMessageComponent:
			// This is where we handle button clicks and selects
			customID := i.MessageComponentData().CustomID
			switch customID {
			case "config_log_btn":
				respondWithLogMenu(s, i)
			case "config_main_menu_btn":
				// This is not a command, so we need to get the user differently
				user := i.Member.User
				respondWithMainMenu(s, i, user)
			case "config_log_channel_btn":
				respondWithChannelSelectMenu(s, i)
			case "log_channel_select":
				handleLogChannelSelect(s, i, db)
			}
		}
	})

	return s, nil
}

// handleLogChannelSelect saves the selected log channel to the database.
func handleLogChannelSelect(s *discordgo.Session, i *discordgo.InteractionCreate, db *sql.DB) {
	// Get the selected channel ID from the interaction data
	data := i.MessageComponentData()
	channelID := data.Values[0]
	guildID := i.GuildID

	// Save to database
	err := database.SetLogChannel(db, guildID, channelID)
	if err != nil {
		log.Printf("Failed to set log channel for guild %s: %v", guildID, err)
		// TODO: Respond with an error message to the user
		return
	}

	// Respond with a success message and show the log menu again
	respondWithLogMenu(s, i)

	// Send a temporary confirmation message
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "✅ ログチャンネルを設定しました。",
		Flags:   discordgo.MessageFlagsEphemeral,
	})
}

// respondWithMainMenu sends or updates a message to show the main config menu.
func respondWithMainMenu(s *discordgo.Session, i *discordgo.InteractionCreate, user *discordgo.User) {
	embed := ui.InfoEmbed(user, "⚙️ Configuration Panel", "設定したい項目をボタンで選択してください。")
	components := []discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					Label:    "ログ機能設定",
					Style:    discordgo.SecondaryButton,
					CustomID: "config_log_btn",
					Emoji:    &discordgo.ComponentEmoji{Name: "📝"},
				},
				&discordgo.Button{
					Label:    "チケット機能設定",
					Style:    discordgo.SecondaryButton,
					CustomID: "config_ticket_btn",
					Emoji:    &discordgo.ComponentEmoji{Name: "🎫"},
					Disabled: true, // Not yet implemented
				},
			},
		},
	}

	// This is a bit of a trick. The original command handler for /config
	// is now just a wrapper around this function.
	// For button clicks, we use InteractionResponseUpdateMessage.
	// We can tell which it is by checking if the interaction has a message attached.
	if i.Message != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{embed},
				Components: components,
				Flags:      discordgo.MessageFlagsEphemeral,
			},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{embed},
				Components: components,
				Flags:      discordgo.MessageFlagsEphemeral,
			},
		})
	}
}

// respondWithLogMenu sends or updates a message to show the log config menu.
func respondWithLogMenu(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User
	embed := ui.InfoEmbed(user, "📝 ログ機能設定", "ログを記録するチャンネルや、記録するイベントの種類を設定します。")

	components := []discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					Label:    "記録チャンネル設定",
					Style:    discordgo.PrimaryButton,
					CustomID: "config_log_channel_btn",
				},
				&discordgo.Button{
					Label:    "記録イベント設定",
					Style:    discordgo.SecondaryButton,
					CustomID: "config_log_events_btn",
					Disabled: true,
				},
			},
		},
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					Label:    "戻る",
					Style:    discordgo.DangerButton,
					CustomID: "config_main_menu_btn",
				},
			},
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
}

// respondWithChannelSelectMenu shows a channel select menu to the user.
func respondWithChannelSelectMenu(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User
	embed := ui.InfoEmbed(user, "✍️ 記録チャンネル設定", "ログを投稿するチャンネルを下のメニューから選択してください。")

	components := []discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.SelectMenu{
					Type:         discordgo.SelectMenuTypeChannel,
					CustomID:     "log_channel_select",
					Placeholder:  "テキストチャンネルを選択...",
					ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
				},
			},
		},
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					Label:    "戻る (ログ設定)",
					Style:    discordgo.SecondaryButton,
					CustomID: "config_log_btn",
				},
			},
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
}
