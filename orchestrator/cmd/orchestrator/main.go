package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"orchestrator/internal/config"
	"orchestrator/internal/rabbitmq/consumers"
	"orchestrator/internal/rabbitmq/producers"
)

const topic = "emails"

func main() {
	err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf(err.Error())
	}

	logger := zap.Must(zap.NewProduction()).Sugar()

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

	customerProducer, err := producers.NewCustomerProducer(conn, "sagaCustomers")
	if err != nil {
		logger.Errorf("failed to initialize customer producer")
	}
	subscriptionProducer, err := producers.NewSubscriptionProducer(conn, "sagaEmailsReply")
	if err != nil {
		logger.Errorf("failed to initialize customer producer")
	}
	customerConsumer, err := consumers.NewCustomerConsumer(conn, "sagaCustomersReply", subscriptionProducer, logger)
	defer customerConsumer.Chan.Close()
	subscriptionConsumer, err := consumers.NewSubscriptionConsumer(conn, "sagaEmails", customerProducer, logger)
	defer subscriptionConsumer.Chan.Close()

	var forever chan struct{}
	go subscriptionConsumer.Listen(forever)
	go customerConsumer.Listen(forever)
	<-forever
}
