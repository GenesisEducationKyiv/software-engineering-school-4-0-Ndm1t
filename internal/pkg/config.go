package pkg

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadConfig(paths ...string) {
	err := godotenv.Load(paths...)
	if err != nil {
		log.Fatalf("Error loading .env file : %v", err)
	}
}
