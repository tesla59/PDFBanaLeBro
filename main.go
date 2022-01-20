package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/go-ping/ping"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"gopkg.in/ini.v1"
)

// Variables used for command line parameters
var (
	Token      string
	PreCommand string = "soja."
	Debug      bool
)

func main() {

	// Loading config.ini
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Println("Failed to read config.ini,", err)
		return
	}
	Token = cfg.Section("").Key("botToken").String()

	// Creating a new Discord session
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(pingCreate)
	discord.AddHandler(start)

	// Debug handler: only enable in debug mode
	Debug, err = cfg.Section("").Key("app_mode").Bool()
	if err != nil {
		fmt.Println("Error: app_mode not defined in config.ini,", err)
		Debug = false // Set Debug to false if not defined
	}
	if Debug {
		discord.AddHandler(debug)
	}

	// In this example, we only care about receiving message events.
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running on", discord.State.User.Username)
	fmt.Println("Press CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func pingCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

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

func start(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == PreCommand+"help" {
		startMessage := "Hey there. I'm a PDF utility bot written in Golang by @tesla59.\nI'm still in my initial phase so don't expect much."
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Description: startMessage,
			URL:         "http://nishantns.xyz/help",
			Type:        "link",
			Title:       "For more help, click here",
		})
	}
}

func debug(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	//
	if m.Content == PreCommand+"debug" {
		imp, _ := api.Import("form:A3, pos:c, s:1.0", pdfcpu.POINTS)
		api.ImportImagesFile([]string{"test/2.png"}, "test/out.pdf", imp, nil)

		file, err := os.Open("test/out.pdf")
		if err != nil {
			fmt.Println("Error Reading output: ", err)
			return
		}
		defer file.Close()

		s.ChannelFileSendWithMessage(m.ChannelID, "Ye Le Bro", "lamma.pdf", file)

		os.Remove("test/out.pdf")
	}
}
