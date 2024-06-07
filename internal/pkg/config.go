package pkg

import (
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
