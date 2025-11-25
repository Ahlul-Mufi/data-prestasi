package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Tidak bisa load .env file, pakai environment variable")
	}
	return os.Getenv(key)
}
