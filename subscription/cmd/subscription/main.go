package main

import (
	"log"
	"subscription-service/internal/config"
	"subscription-service/internal/database"
	"subscription-service/internal/models"
	"subscription-service/internal/server"
	"subscription-service/internal/server/controllers"
	"subscription-service/internal/services"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	db := database.ConnectDatabase()

	if err := db.AutoMigrate(&models.Email{}); err != nil {
		log.Fatalf("Coulnd't migrate database: %v", err.Error())
	}

	subscriptionRepository := database.NewSubscriptionRepository(db)

	subscriptionService := services.NewSubscriptionService(subscriptionRepository)

	subscriptionController := controllers.NewSubscriptionController(subscriptionService)

	s := server.NewServer(subscriptionController)
	s.Run()
}
