package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"PDFBanaLeBro/modules"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/ini.v1"
)

// Variables used for command line parameters
var (
	Token string
	Debug bool
)

func main() {

	// Loading config.ini
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Println("Failed to read config.ini,", err)
		return
	}
	Token = cfg.Section("").Key("botToken").String()

	// Creating a new Discord session
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Println("Error creating discord session,", err)
		return
	}

	// Connect to DB
	err = modules.ConnectDB()
	if err != nil {
		log.Fatal("Connection to database failed,", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(modules.PingCreate)
	discord.AddHandler(modules.Help)
	discord.AddHandler(modules.PDF)

	// Debug handler: only enable in debug mode
	Debug, err = cfg.Section("").Key("app_mode").Bool()
	if err != nil {
		log.Println("Error: app_mode not defined in config.ini,", err)
		Debug = false // Set Debug to false if not defined
	}
	if Debug {
		discord.AddHandler(modules.Debug)
	}

	// In this example, we only care about receiving message events.
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		log.Println("Error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running on", discord.State.User.Username)
	log.Println("Press CTRL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}
