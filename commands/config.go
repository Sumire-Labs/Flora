package commands

import (
	"flora/pkg/ui"

	"github.com/bwmarrin/discordgo"
)

var ConfigCommand = &Command{
	Def: &discordgo.ApplicationCommand{
		Name:        "config",
		Description: "BOTの設定パネルを開きます。",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		user := i.Member.User

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

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds:     []*discordgo.MessageEmbed{embed},
				Components: components,
				Flags:      discordgo.MessageFlagsEphemeral, // Only the user can see this
			},
		})

		if err != nil {
			// Handle error
		}
	},
}
