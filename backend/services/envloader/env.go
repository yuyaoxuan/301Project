package envloader

import (
	"log"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("❌ Failed to load .env file")
	}
}
