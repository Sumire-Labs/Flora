package ui

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// M3E-inspired color palette
const (
	ColorPrimary = 0x6750A4 // A rich, deep purple
	ColorSuccess = 0x4CAF50 // A clear, positive green
	ColorError   = 0xE53935 // A strong, attention-grabbing red
	ColorInfo    = 0x1E88E5 // A calm, informative blue
)

// BaseEmbed creates a new MessageEmbed with a consistent style.
func BaseEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Flora",
			// IconURL: "URL_TO_FLORA_ICON", // TODO: Add bot icon URL later
		},
	}
}

// NewEmbed creates a customized embed message.
func NewEmbed(title, description string, color int) *discordgo.MessageEmbed {
	embed := BaseEmbed()
	embed.Title = title
	embed.Description = description
	embed.Color = color
	return embed
}

// AddField adds a field to an existing embed.
func AddField(embed *discordgo.MessageEmbed, name, value string, inline bool) {
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
}
