package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"informing-service/internal/config"
	"informing-service/internal/crons"
	"informing-service/internal/database"
	"informing-service/internal/mailers"
	"informing-service/internal/models"
	"informing-service/internal/rabbitmq/consumers"
	"informing-service/internal/rabbitmq/producers"
	"informing-service/internal/server"
	"informing-service/internal/server/controllers"
	"informing-service/internal/services"
	"log"
)

const (
	emailsTopic = "emails"
	rateTopic   = "rates"
)

func main() {
	err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf(err.Error())
	}

	logger := zap.Must(zap.NewProduction()).Sugar()

	db := database.ConnectDatabase()

	if err = db.AutoMigrate(&models.Subscription{}, &models.Rate{}); err != nil {
		logger.Errorf("Coulnd't migrate database: %v", err.Error())
	}

	subscriptionRepository := database.NewSubscriptionRepository(db)
	rateRepository := database.NewRateRepository(db)

	conn, err := amqp.Dial(viper.Get("RABBIT_URL").(string))
	if err != nil {
		logger.Errorf("Failed to connetct to rabbitmq: %v", err.Error())
	}
	defer func(conn *amqp.Connection) {
		err = conn.Close()
		if err != nil {
			logger.Errorf("Failed to close rabbit connection: %v", err.Error())
		}
	}(conn)

	smtpSender := mailers.NewSMTPEmailSender()

	subscriptionProducer, err := producers.NewEmailProducer(conn, emailsTopic)
	if err != nil {
		logger.Errorf("failed to initialize subscription producer: %v", err)
	}

	subscriptionConsumer, err := consumers.NewSubscriptionConsumer(
		conn,
		emailsTopic,
		subscriptionRepository,
		rateRepository,
		smtpSender,
		logger)
	if err != nil {
		logger.Errorf("Failed to initialize message producer: %v", err.Error())
	}
	defer subscriptionConsumer.Chan.Close()

	rateConsumer, err := consumers.NewRateConsumer(conn, rateTopic, rateRepository, logger)
	if err != nil {
		logger.Errorf("Failed to initialize message producer: %v", err.Error())
	}
	defer rateConsumer.Chan.Close()

	subscriptionConsumer.Listen()
	rateConsumer.Listen()

	informingService := services.NewInformingService(subscriptionRepository, subscriptionProducer, logger)

	cronScheduler := crons.NewCronScheduler(informingService)

	informingController := controllers.NewInformingController(informingService)

	s := server.NewServer(informingController, cronScheduler)
	s.Scheduler.Start()
	defer s.Scheduler.Stop()
	s.Run()
}
