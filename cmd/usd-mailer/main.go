package main

import (
	"github.com/robfig/cron/v3"
	"gses4_project/internal/crons"
	"gses4_project/internal/database"
	"gses4_project/internal/models"
	"gses4_project/internal/pkg"
	"gses4_project/internal/server"
	"log"
)

func main() {
	pkg.LoadConfig()
	database.ConnectDatabase()

	database.DB.AutoMigrate(&models.Email{})

	s := server.NewServer()

	c := cron.New()
	_, err := c.AddFunc("*/5 * * * *", crons.SendRateEmails)
	if err != nil {
		log.Fatalf("Error scheduling crons job: %v", err)
	}
	c.Start()
	defer c.Stop()

	s.Run()
}
