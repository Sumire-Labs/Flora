package commands

import (
	"flora/pkg/ui"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// HelpCommand is the definition for the /help command.
var HelpCommand = &Command{
	Def: &discordgo.ApplicationCommand{
		Name:        "help",
		Description: "BOTのコマンド一覧を表示します。",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		user := i.Member.User

		// Note: In a real application, you would get the command list from the CommandManager.
		// Since we don't have access to it here without DI refactoring, we list them manually.
		embed := ui.InfoEmbed(user, "Flora Commands", "利用可能なコマンドの一覧です。")

		ui.AddField(embed, "/ping", "BOTの応答速度を測定します。", false)
		ui.AddField(embed, "/help", "このヘルプメッセージを表示します。", false)

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})

		if err != nil {
			fmt.Println("Error sending help message:", err)
		}
	},
}
