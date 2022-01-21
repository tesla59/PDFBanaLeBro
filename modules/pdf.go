package modules

import (
	dload "PDFBanaLeBro/downloader"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func PDF(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	var session Session

	// Opening DB
	db, err := gorm.Open(sqlite.Open("session.db"), &gorm.Config{})
	if err != nil {
		log.Println()
	}

	// Create Schema
	db.AutoMigrate(&Session{})

	//
	// Check if Session already exist
	if m.Content == PreCommand+"start" {
		result := db.Where(&Session{UserID: m.Author.ID}).First(&session)
		if result.Error != nil {
			// Add entry if user doesn't exist
			db.Create(&Session{
				UserID:       m.Author.ID,
				RState:       false,
				CurrentJPEGs: 0,
			})
			db.Where(&Session{UserID: m.Author.ID}).First(&session)
		}

		if !session.RState {
			s.ChannelMessageSend(m.ChannelID, "Never Gonna Give u up\nNever gonna let u down\nAlright You may send nudes now")
			session.RState = true
			db.Save(&session)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Hold up, you already have an active session\nSend Pictures instead")
		}
		os.Mkdir(session.UserID, 0777)
	}

	// Take in all the pictures
	if m.Content == PreCommand+"f" {
		result := db.Where(&Session{UserID: m.Author.ID}).First(&session)
		// User not in DB
		// User in DB, RState True
		// User in DB, RState False
		// User in DB, RState True, soja.f with pik
		// User in DB, RState True, soja.f without pik
		// User in DB, RState False, soja.f with pik
		// User in DB, RState False, soja.f without pik
		if result.Error != nil {
			s.ChannelMessageSend(m.ChannelID, "I dont even know who u are\nSend soja.start to send me your bank details")
		} else if len(m.Attachments) != 0 {
			filePath := session.UserID + "/" + session.UserID + "_" + fmt.Sprint(session.CurrentJPEGs) + ".jpeg" // fmt.Sprint(i) + ".jpeg"
			err := dload.DownloadFile(m.Attachments[0].ProxyURL, filePath)
			if err != nil {
				log.Println("Error downloading file: ", err)
				return
			}
			s.ChannelMessageSend(m.ChannelID, "Hippity Hoppty your nudes are now my property")
			session.CurrentJPEGs++
			db.Save(&session)
		} else if !session.RState {
			s.ChannelMessageSend(m.ChannelID, "Error: I can't convert this to PDF")
		} else {
			s.ChannelMessageSend(m.ChannelID, "You dont have an active session\nSend soja.start to enable a session")
		}
	}

	if m.Content == PreCommand+"end" {
		result := db.Where(&Session{UserID: m.Author.ID}).First(&session)
		if result.Error != nil {
			s.ChannelMessageSend(m.ChannelID, "I dont even know who u are\nSend soja.start to send me your bank details")
		} else if session.RState {
			var inputJPEGs []string
			for i := 0; i < session.CurrentJPEGs; i++ {
				inputJPEGs = append(inputJPEGs, session.UserID+"/"+session.UserID+"_"+fmt.Sprint(i)+".jpeg")
			}

			imp, _ := api.Import("form:A3, pos:c, s:1.0", pdfcpu.POINTS)
			filePDF := session.UserID + "/" + session.UserID + ".pdf"
			err = api.ImportImagesFile(inputJPEGs, filePDF, imp, nil)
			if err != nil {
				log.Println("Error Creating output PDF: ", err)
				return
			}
			file, err := os.Open(filePDF)
			if err != nil {
				log.Println("Error Reading output PDF: ", err)
				return
			}
			defer file.Close()

			s.ChannelFileSendWithMessage(m.ChannelID, "Ye Le Bro", "lamma.pdf", file)

			err = os.RemoveAll(session.UserID)
			if err != nil {
				log.Println("Error Removing temp directory: ", err)
				return
			}

			session.RState = false
			session.CurrentJPEGs = 0
			db.Save(&session)
		} else {
			s.ChannelMessageSend(m.ChannelID, "I dont even know who you are\nSend soja.start to send your bank details")
		}
	}

}
