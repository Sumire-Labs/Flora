package handlers

import (
	"flora/database"
	"flora/pkg/ui"
	"log"

	"github.com/bwmarrin/discordgo"
)

// This file contains the handlers for UI interactions (buttons, selects, etc.)

func (h *EventHandler) handleComponentInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	switch customID {
	case "config_log_btn":
		h.respondWithLogMenu(s, i)
	case "config_main_menu_btn":
		user := i.Member.User
		h.respondWithMainMenu(s, i, user)
	case "config_log_channel_btn":
		h.respondWithChannelSelectMenu(s, i)
	case "log_channel_select":
		h.handleLogChannelSelect(s, i)
	case "config_ticket_btn":
		h.respondWithTicketMenu(s, i)
	}
}

func (h *EventHandler) handleLogChannelSelect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	channelID := data.Values[0]
	guildID := i.GuildID

	if err := database.SetLogChannel(h.DB, guildID, channelID); err != nil {
		log.Printf("Failed to set log channel for guild %s: %v", guildID, err)
		// TODO: Respond with an error message to the user
		return
	}

	// Respond with a success message and show the log menu again
	h.respondWithLogMenu(s, i)

	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "✅ ログチャンネルを設定しました。",
		Flags:   discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		log.Printf("Failed to send followup message: %v", err)
	}
}

func (h *EventHandler) respondWithMainMenu(s *discordgo.Session, i *discordgo.InteractionCreate, user *discordgo.User) {
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
				},
			},
		},
	}

	var err error
	if i.Message != nil {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{embed},
				Components: components,
				Flags:      discordgo.MessageFlagsEphemeral,
			},
		})
	} else {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{embed},
				Components: components,
				Flags:      discordgo.MessageFlagsEphemeral,
			},
		})
	}
	if err != nil {
		log.Printf("Failed to respond to main menu interaction: %v", err)
	}
}

func (h *EventHandler) respondWithLogMenu(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
	if err != nil {
		log.Printf("Failed to respond to log menu interaction: %v", err)
	}
}

func (h *EventHandler) respondWithChannelSelectMenu(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User
	embed := ui.InfoEmbed(user, "✍️ 記録チャンネル設定", "ログを投稿するチャンネルを下のメニューから選択してください。")
	components := []discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.SelectMenu{
					CustomID:     "log_channel_select",
					Placeholder:  "テキストチャンネルを選択...",
					ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
					Options:      []discordgo.SelectMenuOption{},
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

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
	if err != nil {
		log.Printf("Failed to respond to channel select interaction: %v", err)
	}
}

func (h *EventHandler) respondWithTicketMenu(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User
	embed := ui.InfoEmbed(user, "🎫 チケット機能設定", "チケット機能に関する設定を行います。")

	components := []discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					Label:    "パネル送信先設定",
					Style:    discordgo.PrimaryButton,
					CustomID: "config_ticket_panel_channel_btn",
					Disabled: true, // TODO: Implement
				},
				&discordgo.Button{
					Label:    "サポートロール設定",
					Style:    discordgo.PrimaryButton,
					CustomID: "config_ticket_support_role_btn",
					Disabled: true, // TODO: Implement
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

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
	if err != nil {
		log.Printf("Failed to respond to ticket menu interaction: %v", err)
	}
}
