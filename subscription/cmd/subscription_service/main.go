package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"subscription-service/internal/config"
	"subscription-service/internal/database"
	"subscription-service/internal/models"
	"subscription-service/internal/rabbitmq/consumers"
	"subscription-service/internal/rabbitmq/producers"
	"subscription-service/internal/server"
	"subscription-service/internal/server/controllers"
	"subscription-service/internal/services"
)

const topic = "emails"

func main() {
	err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf(err.Error())
	}

	logger := zap.Must(zap.NewProduction()).Sugar()

	db := database.ConnectDatabase()

	conn, err := amqp.Dial(viper.GetString("RABBIT_URL"))
	if err != nil {
		logger.Errorf("Failed to connetct to rabbitmq: %v", err.Error())
	}
	defer func(conn *amqp.Connection) {
		err = conn.Close()
		if err != nil {
			logger.Errorf("Failed to close rabbit connection: %v", err.Error())
		}
	}(conn)

	subscriptionProducer, err := producers.NewEmailProducer(conn, topic)
	if err != nil {
		logger.Errorf("Failed to initialize message producer: %v", err.Error())
	}

	subscriptionSagaProducer, err := producers.NewEmailProducer(conn, "sagaEmails")
	if err != nil {
		logger.Errorf("Failed to initialize message producer: %v", err.Error())
	}

	if err = db.AutoMigrate(&models.Email{}); err != nil {
		logger.Errorf("Coulnd't migrate database: %v", err.Error())
	}

	subscriptionRepository := database.NewSubscriptionRepository(db)

	subscriptionService := services.NewSubscriptionService(subscriptionRepository, subscriptionProducer, subscriptionSagaProducer, logger)

	subscriptionController := controllers.NewSubscriptionController(subscriptionService)

	subscriptionConsumer, err := consumers.NewSubscriptionConsumer(conn, "sagaEmailsReply", subscriptionService, logger)
	defer subscriptionConsumer.Chan.Close()

	s := server.NewServer(subscriptionController)

	subscriptionConsumer.Listen()
	s.Run()
}
