package ui

import (
	"fmt"
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

// BaseEmbed creates a new MessageEmbed with a consistent style, including footer with requester info.
func BaseEmbed(user *discordgo.User) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Requested by %s", user.Username),
			IconURL: user.AvatarURL(""),
		},
	}
}

// SuccessEmbed creates a new embed with a success style.
func SuccessEmbed(user *discordgo.User, title, description string) *discordgo.MessageEmbed {
	embed := BaseEmbed(user)
	embed.Title = "✅ " + title
	embed.Description = description
	embed.Color = ColorSuccess
	return embed
}

// ErrorEmbed creates a new embed with an error style.
func ErrorEmbed(user *discordgo.User, title, description string) *discordgo.MessageEmbed {
	embed := BaseEmbed(user)
	embed.Title = "❌ " + title
	embed.Description = description
	embed.Color = ColorError
	return embed
}

// InfoEmbed creates a new embed with an info style.
func InfoEmbed(user *discordgo.User, title, description string) *discordgo.MessageEmbed {
	embed := BaseEmbed(user)
	embed.Title = "ℹ️ " + title
	embed.Description = description
	embed.Color = ColorInfo
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

// SetThumbnail sets the thumbnail for an embed.
func SetThumbnail(embed *discordgo.MessageEmbed, url string) {
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: url}
}
