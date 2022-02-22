package modules

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Help(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.ToLower(m.Content) == PreCommand+"help" {
		startMessage := "Hey there. I'm a PDF utility bot written in Golang by @tesla59.\nI'm still in my initial phase so don't expect much."
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Description: startMessage,
			URL:         "http://nishantns.xyz/help",
			Type:        "link",
			Title:       "For more help, click here",
		})
	}
}
