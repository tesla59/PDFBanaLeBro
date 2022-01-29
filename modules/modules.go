package modules

import (
	"log"
	"net/http"
	"os"
)

var PreCommand string = "soja."
var Err error

func isImage(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("File doesn't exist: ", err)
	}

	buf := make([]byte, 512)
	file.Read(buf)
	if err != nil {
		log.Println("File doesn't exist: ", err)
	}

	fileType := http.DetectContentType(buf)
	if fileType == "image/png" || fileType == "image/jpeg" || fileType == "image/webp" {
		return true
	} else {
		return false
	}
}
