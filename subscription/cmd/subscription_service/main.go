package main

import (
	"log"
	"subscription-service/internal/config"
	"subscription-service/internal/crons"
	"subscription-service/internal/database"
	kafkaClient "subscription-service/internal/kafka"
	"subscription-service/internal/kafka/producers"
	"subscription-service/internal/models"
	"subscription-service/internal/server"
	"subscription-service/internal/server/controllers"
	"subscription-service/internal/services"
)

const topic = "emails"

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}

	producer, err := kafkaClient.NewProducer()
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err.Error())
	}
	defer producer.Close()

	emailProducer := producers.NewEmailProducer(producer, topic)

	db := database.ConnectDatabase()

	if err = db.AutoMigrate(&models.Email{}); err != nil {
		log.Fatalf("Coulnd't migrate database: %v", err.Error())
	}

	subscriptionRepository := database.NewSubscriptionRepository(db)

	subscriptionService := services.NewSubscriptionService(subscriptionRepository)

	cronScheduler := crons.NewCronScheduler(subscriptionService, emailProducer)

	subscriptionController := controllers.NewSubscriptionController(subscriptionService)

	s := server.NewServer(subscriptionController, cronScheduler)
	s.Scheduler.Start()
	defer s.Scheduler.Stop()
	s.Run()
}
