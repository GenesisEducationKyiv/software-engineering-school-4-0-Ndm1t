package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"subscription-service/internal/config"
	"subscription-service/internal/database"
	"subscription-service/internal/models"
	"subscription-service/internal/rabbitmq/producers"
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

	db := database.ConnectDatabase()

	conn, err := amqp.Dial(os.Getenv("RABBIT_URL"))
	if err != nil {
		log.Fatalf("Failed to connetct to rabbitmq: %v", err.Error())
	}
	defer func(conn *amqp.Connection) {
		err = conn.Close()
		if err != nil {
			log.Fatalf("Failed to close rabbit connection: %v", err.Error())
		}
	}(conn)

	subscriptionProducer, err := producers.NewEmailProducer(conn, topic)
	if err != nil {
		log.Fatalf("Failed to initialize message producer: %v", err.Error())
	}

	if err = db.AutoMigrate(&models.Email{}); err != nil {
		log.Fatalf("Coulnd't migrate database: %v", err.Error())
	}

	subscriptionRepository := database.NewSubscriptionRepository(db)

	subscriptionService := services.NewSubscriptionService(subscriptionRepository, subscriptionProducer)

	subscriptionController := controllers.NewSubscriptionController(subscriptionService)

	s := server.NewServer(subscriptionController)

	s.Run()
}
