package main

import (
	"customer-service/internal/config"
	"customer-service/internal/database"
	"customer-service/internal/models"
	"customer-service/internal/rabbitmq/consumers"
	"customer-service/internal/rabbitmq/producers"
	"customer-service/internal/server"
	"customer-service/internal/services"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
)

func main() {
	err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf(err.Error())
	}

	logger := zap.Must(zap.NewProduction()).Sugar()

	db := database.ConnectDatabase()

	if err = db.AutoMigrate(&models.Customer{}); err != nil {
		logger.Errorf("couldn't migrate database: %v", err.Error())
	}

	customerRepository := database.NewCustomerRepository(db)

	conn, err := amqp.Dial(viper.GetString("RABBIT_URL"))
	if err != nil {
		logger.Errorf("couldn't migrate database: %v", err.Error())
	}
	defer func(conn *amqp.Connection) {
		err = conn.Close()
		if err != nil {
			logger.Errorf("Failed to close rabbit connection: %v", err.Error())
		}
	}(conn)

	customerProducer, err := producers.NewCustomerProducer(conn, "sagaCustomersReply", logger)
	if err != nil {
		logger.Errorf("Failed to initialize message producer: %v", err.Error())
	}

	customerService := services.NewCustomerService(customerProducer, customerRepository)

	customerConsumer, err := consumers.NewCustomerConsumer(conn, "sagaCustomers", customerService, logger)
	if err != nil {
		logger.Errorf("Failed to initialize message producer: %v", err.Error())
	}
	defer customerConsumer.Chan.Close()

	s := server.NewServer()
	s.Run()

	var forever chan struct{}

	customerConsumer.Listen(forever)
}
