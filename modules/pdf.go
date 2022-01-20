package modules

import (
	dload "PDFBanaLeBro/downloader"
	"fmt"
	"log"
	"os"
	// "strings"

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
		if result.Error != nil {
			s.ChannelMessageSend(m.ChannelID, "Niggah, I dont even know who u are\nSend soja.start to send me your bank details")
		} else {
			// for i := 0 + session.CurrentJPEGs; i < len(m.Attachments)+ session.CurrentJPEGs; i++ {
			filePath := session.UserID + "/" + session.UserID + "_" + fmt.Sprint(session.CurrentJPEGs) + ".jpeg" // fmt.Sprint(i) + ".jpeg"
			err := dload.DownloadFile(m.Attachments[0].ProxyURL, filePath)
			if err != nil {
				log.Println("Error downloading file: ", err)
				return
			}
			s.ChannelMessageSend(m.ChannelID, "Hippity Hoppty your nudes are now my property")
			session.CurrentJPEGs++
			db.Save(&session)
			// }
		}
	}

	if m.Content == PreCommand+"end" {
		result := db.Where(&Session{UserID: m.Author.ID}).First(&session)
		if result.Error != nil {
			s.ChannelMessageSend(m.ChannelID, "Niggah, I dont even know who u are\nSend soja.start to send me your bank details")
		} else {
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
			session.CurrentJPEGs=0
			db.Save(&session)
		}
	}

}

// func isPic(name string) bool {
// 	length := len(name)
// 	arr := strings.Split(name, "")
// 	ext := arr[length-4] + arr[length-3] + arr[length-2] + arr[length-1]
// 	if ext == ".png" || ext == "jpeg" || ext == ".jpg" {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// // if m.Author.ID == session.UserID
// result := db.Where(&Session{UserID: m.Author.ID}).First(&session)
// if result.Error != nil {
// 	s.ChannelMessageSend(m.ChannelID, "I don't think I know you")
// 	return
// }
// var inputJPEGs []string
// for i := 0; i < session.CurrentJPEGs; i++ {
// 	inputJPEGs = append(inputJPEGs, session.UserID + fmt.Sprint(i) + ".jpeg")
// }

// imp, _ := api.Import("form:A3, pos:c, s:1.0", pdfcpu.POINTS)
// filePDF := session.UserID + "/" + session.UserID + ".pdf"
// api.ImportImagesFile(inputJPEGs, filePDF, imp, nil)

// file, err := os.Open(filePDF)
// if err != nil {
// 	log.Println("Error Reading output PDF: ", err)
// 	return
// }
// defer file.Close()

// s.ChannelFileSendWithMessage(m.ChannelID, "Ye Le Bro", "lamma.pdf", file)

// os.Remove("test/out.pdf")
