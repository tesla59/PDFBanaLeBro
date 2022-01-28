package modules

import (
	dload "PDFBanaLeBro/downloader"
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
		log.Println("Failed to Open session.db: ", err)
	}
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Println("Failed to generate sqlDB: ", err)
		}
		sqlDB.Close()
	}()

	// Create Schema
	db.AutoMigrate(&Session{})

	// Start command
	if m.Content == PreCommand+"start" {

		// Fetch user's data from DB
		err := db.Where(&Session{UserID: m.Author.ID}).First(&session)
		if err.Error != nil {
			// Add entry if user doesn't exist
			db.Create(&Session{
				UserID:        m.Author.ID,
				RState:        false,
				CurrentImages: 0,
			})
			db.Where(&Session{UserID: m.Author.ID}).First(&session)
		}

		if !session.RState {
			// Initiate a session
			s.ChannelMessageSend(m.ChannelID, "Never Gonna Give u up\nNever gonna let u down\nAlright You may send images now")
			session.RState = true
			db.Save(&session)
		} else {
			// Already active session
			s.ChannelMessageSend(m.ChannelID, "Hold up, you already have an active session\nSend images instead")
		}
		// Store all temp files in usersID/userID_i.png
		os.Mkdir(session.UserID, 0777)
	}

	// Fetch images command
	if m.Content == PreCommand+"f" {
		// Fetch user's data from DB
		err := db.Where(&Session{UserID: m.Author.ID}).First(&session)

		if err.Error != nil {
			// User not in DB
			s.ChannelMessageSend(m.ChannelID, "I dont even know who u are\nSend soja.start to send me your bank details")
		} else {
			// User in DB
			if session.RState {
				// User in DB + RState True
				if len(m.Attachments) == 0 {
					// soja.f without attachments
					s.ChannelMessageSend(m.ChannelID, "Error: no file sent")
				} else {
					// soja.f with attachment
					for i := range m.Attachments {
						filePath := session.UserID + "/" + "temp.jpeg"
						err := dload.DownloadFile(m.Attachments[i].ProxyURL, filePath)
						if err != nil {
							log.Println("Error downloading file: ", err)
							return
						}
						if isImage(filePath) {
							// File is Image
							if i == 0 {
								// Only send this once
								s.ChannelMessageSend(m.ChannelID, "Hippity Hoppty your images are now my property")
							}
							// Creating a new PDF/Appending to existing one
							imp, _ := api.Import("form:A3, pos:c, s:1.0", pdfcpu.POINTS)
							filePDF := session.UserID + "/" + m.Author.Username + ".pdf"
							err := api.ImportImagesFile([]string{filePath}, filePDF, imp, nil)
							if err != nil {
								log.Println("Error Creating output PDF: ", err)
								return
							}
							os.Remove(filePath)
							session.CurrentImages++
							db.Save(&session)
						} else {
							// Unsupported filetype
							s.ChannelMessageSend(m.ChannelID, "Error: This file format is not supported")
							os.Remove(filePath)
						}
					}
				}
			} else {
				// User in DB + RState False
				s.ChannelMessageSend(m.ChannelID, "You dont have an active session\nSend soja.start to enable a session")
			}
		}
	}

	// End command
	if m.Content == PreCommand+"end" {
		// Fetch user's data from DB
		err := db.Where(&Session{UserID: m.Author.ID}).First(&session)

		if err.Error != nil {
			// User doesn't exist/New User
			s.ChannelMessageSend(m.ChannelID, "I dont even know who u are\nSend soja.start to send me your bank details")
		} else if session.RState {
			// User exist + Has an active session
			if session.CurrentImages != 0 {
				// Has at least 1 image to convert
				filePDF := session.UserID + "/" + m.Author.Username + ".pdf"

				// Create *FILE for output.pdf
				file, err := os.Open(filePDF)
				if err != nil {
					log.Println("Error Reading output PDF: ", err)
					return
				}
				defer file.Close()

				// Send the final PDF
				s.ChannelFileSendWithMessage(m.ChannelID, "Ye Le Bro", m.Author.Username+".pdf", file)

				// Clean all temp directories
				err = os.RemoveAll(session.UserID)
				if err != nil {
					log.Println("Error Removing temp directory: ", err)
					return
				}

				// Reset all DB entries except userID
				session.RState = false
				session.CurrentImages = 0
				db.Save(&session)
			} else {
				// User exist + Active session + 0 images sent
				session.RState = false
				s.ChannelMessageSend(m.ChannelID, "Okay your session has been ended\nSend soja.start to initiate a new session again")
				db.Save(&session)
			}
		} else {
			// User exist + inactive session
			s.ChannelMessageSend(m.ChannelID, "You don't have any active session\nSend soja.start to initiate a session")
		}
	}
}
