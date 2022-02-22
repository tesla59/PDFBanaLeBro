package modules

import (
	dload "PDFBanaLeBro/downloader"
	"log"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func PDF(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	var session Session

	// Create Schema
	DB.AutoMigrate(&Session{})

	// Start command
	if m.Content == PreCommand+"start" {

		// Fetch user's data from DB
		if (DB.Where(&Session{UserID: m.Author.ID}).First(&session)).Error != nil {
			// Add entry if user doesn't exist
			DB.Create(&Session{
				UserID:        m.Author.ID,
				RState:        false,
				CurrentImages: 0,
			})
			DB.Where(&Session{UserID: m.Author.ID}).First(&session)
		}

		if !session.RState {
			// Initiate a session
			s.ChannelMessageSend(m.ChannelID, "Never Gonna Give u up\nNever gonna let u down\nAlright You may send images now")
			session.RState = true
			DB.Save(&session)
		} else {
			// Already active session
			s.ChannelMessageSend(m.ChannelID, "Hold up, you already have an active session\nSend images instead")
		}
		// Store all temp files in temp/usersID/userID_i.png
		if Err = os.MkdirAll(filepath.Join("temp", session.UserID), 0777); Err != nil {
			log.Println("Cannot create directory", Err)
		}
	}

	// Fetch images command
	if m.Content == PreCommand+"f" {
		// Fetch user's data from DB
		if (DB.Where(&Session{UserID: m.Author.ID}).First(&session)).Error != nil {
			// User not in DB
			s.ChannelMessageSend(m.ChannelID, "I dont even know who u are\nSend "+PreCommand+"start to send me your bank details")
			return
		}
		// User not in DB
		if !session.RState {
			// User in DB + RState False
			s.ChannelMessageSend(m.ChannelID, "You dont have an active session\nSend "+PreCommand+"start to enable a session")
			return
		}
		// User in DB + RState True
		if len(m.Attachments) == 0 {
			// "+PreCommand+".f without attachments
			s.ChannelMessageSend(m.ChannelID, "Error: no file sent")
			return
		}
		// "+PreCommand+".f with attachment
		for i := range m.Attachments {
			filePath := filepath.Join("temp", session.UserID, "temp.jpeg")
			if Err = dload.DownloadFile(m.Attachments[i].ProxyURL, filePath); Err != nil {
				log.Println("Error downloading file: ", Err)
				return
			}
			if !isImage(filePath) {
				// Unsupported filetype
				s.ChannelMessageSend(m.ChannelID, "Error: This file format is not supported")
				os.RemoveAll(filePath)
				return
			}
			// File is Image
			if i == 0 {
				// Only send this once
				s.ChannelMessageSend(m.ChannelID, "Hippity Hoppty your images are now my property")
			}
			// Creating a new PDF/Appending to existing one
			imp, _ := api.Import("form:A3, pos:c, s:1.0", pdfcpu.POINTS)
			filePDF := filepath.Join("temp", session.UserID, m.Author.Username+".pdf")
			Err = api.ImportImagesFile([]string{filePath}, filePDF, imp, nil)
			if Err != nil {
				log.Println("Error Creating output PDF: ", Err)
				return
			}
			os.RemoveAll(filePath)
			session.CurrentImages++
			DB.Save(&session)
		}
	}

	// End command
	if m.Content == PreCommand+"end" {
		// Fetch user's data from DB
		if (DB.Where(&Session{UserID: m.Author.ID}).First(&session)).Error != nil {
			// User doesn't exist/New User
			s.ChannelMessageSend(m.ChannelID, "I dont even know who u are\nSend "+PreCommand+"start to send me your bank details")
			return
		}
		if !session.RState {
			// User exist + inactive session
			s.ChannelMessageSend(m.ChannelID, "You don't have any active session\nSend "+PreCommand+"start to initiate a session")
			return
		}
		// User exist + Has an active session
		if session.CurrentImages == 0 {
			// User exist + Active session + 0 images sent
			s.ChannelMessageSend(m.ChannelID, "Okay your session has been ended\nSend "+PreCommand+"start to initiate a new session again")
			session.RState = false
			DB.Save(&session)
			if Err = os.RemoveAll(filepath.Join("temp", session.UserID)); Err != nil {
				log.Println("Error Removing temp directory: ", Err)
				return
			}
		}
		// Has at least 1 image to convert
		filePDF := filepath.Join("temp", session.UserID, m.Author.Username+".pdf")
		// Create *FILE for output.pdf
		file, Err := os.Open(filePDF)
		if Err != nil {
			log.Println("Error Reading output PDF: ", Err)
			return
		}
		defer file.Close()

		stats, Err := file.Stat()
		if Err != nil {
			log.Println("Error Reading output PDF: ", Err)
			return
		}

		if stats.Size() > 8388608 {
			// 8mB is the size limit. POG Discord
			s.ChannelMessageSend(m.ChannelID, "You're file exceed Discord's file size limit of 8mB. Try recreating PDF with less images")
		} else {
			// Send the final PDF
			s.ChannelFileSendWithMessage(m.ChannelID, "Ye Le Bro", m.Author.Username+".pdf", file)
		}

		// Clean all temp directories
		if Err = os.RemoveAll(filepath.Join("temp", session.UserID)); Err != nil {
			log.Println("Error Removing temp directory: ", Err)
			return
		}

		// Reset all DB entries except userID
		session.RState = false
		session.CurrentImages = 0
		DB.Save(&session)
	}
}
