package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// PingCommand is the definition for the /ping command.
var PingCommand = &Command{
	Def: &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Pong! とレイテンシを返します。",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		startTime := time.Now()

		// "Thinking..." を表示
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		if err != nil {
			// エラー処理は後で改善
			return
		}

		latency := time.Since(startTime)

		// メッセージを編集して結果を表示
		content := fmt.Sprintf("Pong! 🏓\nLatency: %s", latency)
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
	},
}
