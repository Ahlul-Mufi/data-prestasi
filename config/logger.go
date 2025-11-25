package config

import (
	"log"
	"os"
)

var Logger *log.Logger

func InitLogger() {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	Logger = log.New(file, "ALUMNI-APP: ", log.Ldate|log.Ltime|log.Lshortfile)
}