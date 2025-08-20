package commands

import (
	"flora/pkg/ui"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// PingCommand is the definition for the /ping command.
var PingCommand = &Command{
	Def: &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Pong! とレイテンシを返します。",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Get the user who initiated the command
		user := i.Member.User

		// Respond with a deferred message to show "Thinking..."
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		if err != nil {
			// In a real app, you'd send an error message back to the user.
			fmt.Println("Failed to send deferred message:", err)
			return
		}

		// Calculate latency
		latency := s.HeartbeatLatency()

		// Create the success embed
		embed := ui.SuccessEmbed(user, "Pong!", "")
		latencyStr := fmt.Sprintf("```%s```", latency.String())
		ui.AddField(embed, "API Latency", latencyStr, true)

		// Edit the original deferred message with the final embed
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		})
	},
}
