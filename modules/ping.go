package modules

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/go-ping/ping"
)

func PingCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// '/ping' command starts here
	if m.Content == PreCommand+"ping" {
		// Respond to the command (Edit later)
		message, err := s.ChannelMessageSend(m.ChannelID, "Pinging.........")
		if err != nil {
			fmt.Println("Error sending message,", err)
			return
		}
		// Ping!!
		pinger, err := ping.NewPinger("www.google.com")
		if err != nil {
			fmt.Println("URL not reachable,", err)
			return
		}
		// Blocks until finished.
		pinger.Count = 5
		err = pinger.Run()
		if err != nil {
			fmt.Println("Error running Pinger,", err)
			return
		}
		stats := pinger.Statistics()

		s.ChannelMessageEdit(m.ChannelID, message.ID, "Ping: "+stats.AvgRtt.String()+"\nIP Addr: "+stats.IPAddr.String())
	}
}
