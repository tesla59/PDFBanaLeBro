package modules

import (
	// dload "PDFBanaLeBro/downloader"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func Debug(s *discordgo.Session, m *discordgo.MessageCreate) {

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
			log.Println("Error Reading output: ", err)
			return
		}
		defer file.Close()

		s.ChannelFileSendWithMessage(m.ChannelID, "Ye Le Bro", "lamma.pdf", file)

		os.Remove("test/out.pdf")
	}

	// if m.Content == PreCommand+"test" {
		// link := m.Attachments[0].ProxyURL

		// filePath, err := dload.DownloadFile(link)

		// if err != nil {
		// 	log.Println("Error downloading file: ", err)
		// 	return
		// }

		// imp, _ := api.Import("form:A3, pos:c, s:1.0", pdfcpu.POINTS)
		// api.ImportImagesFile([]string{filePath}, filePath+".pdf", imp, nil)

		// file, err := os.Open(filePath + ".pdf")
		// if err != nil {
		// 	log.Println("Error Reading output: ", err)
		// 	return
		// }
		// defer file.Close()

		// s.ChannelFileSendWithMessage(m.ChannelID, "Ye Le Bro", "lamma.pdf", file)

		// os.Remove(filePath + ".pdf")
		// os.Remove(filePath)

	// }
}
