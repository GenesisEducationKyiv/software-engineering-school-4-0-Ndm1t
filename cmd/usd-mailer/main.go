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

	if err := database.DB.AutoMigrate(&models.Email{}); err != nil {
		log.Fatalf("Coulnd't migrate database: %v", err.Error())
	}
	s := server.NewServer()
	c := cron.New()
	_, err := c.AddFunc("0 9 * * *", crons.SendRateEmails)
	if err != nil {
		log.Fatalf("Error scheduling crons job: %v", err)
	}
	c.Start()
	defer c.Stop()

	s.Run()
}
