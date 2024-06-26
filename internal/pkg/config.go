package pkg

import (
	"fmt"
	"github.com/joho/godotenv"
)

func LoadConfig(paths ...string) error {
	err := godotenv.Load(paths...)
	if err != nil {
		return fmt.Errorf("error loading .env file : %v", err)
	}
	return err
}
