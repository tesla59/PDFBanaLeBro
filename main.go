package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/go-ping/ping"
	"gopkg.in/ini.v1"
)

// Variables used for command line parameters
var (
	Token string
)

func main() {

	// Loading config.ini
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Println("Fail to read file,", err)
		return
	}
	Token = cfg.Section("").Key("botToken").String()

	// Creating a new Discord session
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// '/ping' command starts here
	if m.Content == "/ping" {
		// Respond to the command (Edit later)
		message, err := s.ChannelMessageSend(m.ChannelID, "Pinging.........")
		if err != nil {
			fmt.Println("Cannot send message: ", err)
			return
		}
		// Ping!!
		pinger, err := ping.NewPinger("www.google.com")
		if err != nil {
			fmt.Println("URL not reachable: ", err)
			return
		}
		// Blocks until finished.
		pinger.Count = 5
		err = pinger.Run()
		if err != nil {
			fmt.Println("Pinger couldn't run: ",err)
			return
		}
		stats := pinger.Statistics()

		if err != nil {
			fmt.Println("Cannot send message: ", err)
			return
		}
		s.ChannelMessageEdit(m.ChannelID, message.ID, "Ping: "+stats.AvgRtt.String()+"\nIP Addr: "+stats.IPAddr.String())
	}
}
