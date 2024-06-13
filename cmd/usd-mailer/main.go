package main

import (
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
	s.Scheduler.Cron.Start()
	defer s.Scheduler.Cron.Stop()
	s.Run()
}
