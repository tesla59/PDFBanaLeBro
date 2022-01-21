package modules

import (
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

func Debug(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	//
	if m.Content == PreCommand+"debug" {
		file, err := os.Open("test.jpeg")
		if err != nil {
			log.Println(err)
		}
		buf := make([]byte, 512)
		file.Read(buf)
		log.Println(http.DetectContentType(buf))
	}
}
