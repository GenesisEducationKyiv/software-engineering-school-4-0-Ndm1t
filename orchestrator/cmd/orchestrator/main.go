package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
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

	conn, err := amqp.Dial(viper.GetString("RABBIT_URL"))
	if err != nil {
		log.Fatalf("Failed to connetct to rabbitmq: %v", err.Error())
	}
	defer func(conn *amqp.Connection) {
		err = conn.Close()
		if err != nil {
			log.Fatalf("Failed to close rabbit connection: %v", err.Error())
		}
	}(conn)

	customerProducer, err := producers.NewCustomerProducer(conn, "sagaCustomers")
	if err != nil {
		log.Fatal("failed to initialize customer producer")
	}
	subscriptionProducer, err := producers.NewSubscriptionProducer(conn, "sagaEmailsReply")
	if err != nil {
		log.Fatal("failed to initialize customer producer")
	}
	customerConsumer, err := consumers.NewCustomerConsumer(conn, "sagaCustomersReply", subscriptionProducer)
	defer customerConsumer.Chan.Close()
	subscriptionConsumer, err := consumers.NewSubscriptionConsumer(conn, "sagaEmails", customerProducer)
	defer subscriptionConsumer.Chan.Close()

	var forever chan struct{}
	go subscriptionConsumer.Listen(forever)
	go customerConsumer.Listen(forever)
	<-forever
}
