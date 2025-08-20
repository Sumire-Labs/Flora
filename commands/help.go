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
		// ここではまだ cmdManager を直接参照できないため、手動でコマンドリストを作成します。
		// DI導入時に、ここを動的に生成するように修正します。
		embed := ui.NewEmbed("🌿 Flora Commands", "利用可能なコマンドの一覧です。", ui.ColorPrimary)

		ui.AddField(embed, "/ping", "BOTの応答速度を測定します。", false)
		ui.AddField(embed, "/help", "このヘルプメッセージを表示します。", false)

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})

		if err != nil {
			// 仮のエラーハンドリング
			fmt.Println("Error sending help message:", err)
		}
	},
}
