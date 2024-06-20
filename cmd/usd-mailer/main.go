package main

import (
	"gses4_project/internal/container"
	"gses4_project/internal/database"
	"gses4_project/internal/models"
	"gses4_project/internal/pkg"
	"gses4_project/internal/server"
	"log"
)

func main() {
	pkg.LoadConfig()
	db := database.ConnectDatabase()

	if err := db.AutoMigrate(&models.Email{}); err != nil {
		log.Fatalf("Coulnd't migrate database: %v", err.Error())
	}

	container := container.NewContainer(db)

	s := server.NewServer(container)
	s.Scheduler.Cron.Start()
	defer s.Scheduler.Cron.Stop()
	s.Run()
}
