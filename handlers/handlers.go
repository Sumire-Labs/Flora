package handlers

import (
	"database/sql"
	"flora/commands"
	"flora/database"
	"flora/pkg/ui"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// EventHandler holds all dependencies for handling Discord events.
type EventHandler struct {
	DB *sql.DB
	CM *commands.Manager
}

// NewEventHandler creates a new EventHandler.
func NewEventHandler(db *sql.DB, cm *commands.Manager) *EventHandler {
	return &EventHandler{DB: db, CM: cm}
}

// OnReady is called when the bot is ready to start receiving events.
func (h *EventHandler) OnReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	s.UpdateGameStatus(0, "Flora")

	if err := h.CM.RegisterAll(s); err != nil {
		log.Fatalf("Failed to register commands: %v", err)
	}
	log.Println("Commands registered successfully.")
}

// OnInteractionCreate is called when a user interacts with the bot (slash command, button, etc.).
func (h *EventHandler) OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if cmd, exists := h.CM.Get(i.ApplicationCommandData().Name); exists {
			cmd.Handler(s, i)
		}
	case discordgo.InteractionMessageComponent:
		h.handleComponentInteraction(s, i)
	}
}

// OnMessageDelete handles logging for deleted messages.
func (h *EventHandler) OnMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	if m.GuildID == "" {
		return
	}
	logChannelID, exists, err := database.GetLogChannel(h.DB, m.GuildID)
	if err != nil || !exists {
		return
	}
	if m.BeforeDelete == nil {
		return
	}
	if m.BeforeDelete.Author.ID == s.State.User.ID || m.BeforeDelete.ChannelID == logChannelID {
		return
	}
	embed := ui.ErrorEmbed(m.BeforeDelete.Author, "Message Deleted", "")
	ui.AddField(embed, "Channel", "<#"+m.ChannelID+">", true)
	ui.AddField(embed, "Author", m.BeforeDelete.Author.Mention(), true)
	if m.BeforeDelete.Content != "" {
		ui.AddField(embed, "Content", m.BeforeDelete.Content, false)
	}
	s.ChannelMessageSendEmbeds(logChannelID, []*discordgo.MessageEmbed{embed})
}

// OnMemberAdd handles logging for new members.
func (h *EventHandler) OnMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	logChannelID, exists, err := database.GetLogChannel(h.DB, m.GuildID)
	if err != nil || !exists {
		return
	}
	embed := ui.SuccessEmbed(m.User, "Member Joined", "")
	ui.SetThumbnail(embed, m.User.AvatarURL("128"))
	ui.AddField(embed, "User", m.User.Mention(), true)
	ui.AddField(embed, "Account Created", fmt.Sprintf("<t:%d:R>", m.User.ID>>22+1420070400000/1000), true)
	s.ChannelMessageSendEmbeds(logChannelID, []*discordgo.MessageEmbed{embed})
}

// OnMemberRemove handles logging for members who leave.
func (h *EventHandler) OnMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	logChannelID, exists, err := database.GetLogChannel(h.DB, m.GuildID)
	if err != nil || !exists {
		return
	}
	embed := ui.ErrorEmbed(m.User, "Member Left", "")
	ui.SetThumbnail(embed, m.User.AvatarURL("128"))
	ui.AddField(embed, "User", m.User.Username, false)
	s.ChannelMessageSendEmbeds(logChannelID, []*discordgo.MessageEmbed{embed})
}

// OnMemberUpdate handles logging for member updates (e.g., nickname).
func (h *EventHandler) OnMemberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	if m.BeforeUpdate != nil && m.Nick != m.BeforeUpdate.Nick {
		logChannelID, exists, err := database.GetLogChannel(h.DB, m.GuildID)
		if err != nil || !exists {
			return
		}
		beforeNick := m.BeforeUpdate.Nick
		if beforeNick == "" {
			beforeNick = "(None)"
		}
		afterNick := m.Nick
		if afterNick == "" {
			afterNick = "(None)"
		}
		embed := ui.InfoEmbed(m.User, "Nickname Changed", "")
		ui.SetThumbnail(embed, m.User.AvatarURL("128"))
		ui.AddField(embed, "User", m.User.Mention(), false)
		ui.AddField(embed, "Before", beforeNick, true)
		ui.AddField(embed, "After", afterNick, true)
		s.ChannelMessageSendEmbeds(logChannelID, []*discordgo.MessageEmbed{embed})
	}
}
