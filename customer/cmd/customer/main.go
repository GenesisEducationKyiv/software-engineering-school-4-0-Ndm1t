package main

import (
	"customer-service/internal/config"
	"customer-service/internal/database"
	"customer-service/internal/models"
	"customer-service/internal/rabbitmq/consumers"
	"customer-service/internal/rabbitmq/producers"
	"customer-service/internal/services"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"log"
)

func main() {
	err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf(err.Error())
	}

	db := database.ConnectDatabase()

	if err = db.AutoMigrate(&models.Customer{}); err != nil {
		log.Fatalf("Coulnd't migrate database: %v", err.Error())
	}

	customerRepository := database.NewCustomerRepository(db)

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

	customerProducer, err := producers.NewCustomerProducer(conn, "sagaCustomersReply")
	if err != nil {
		log.Fatalf("Failed to initialize message producer: %v", err.Error())
	}

	customerService := services.NewCustomerService(customerProducer, customerRepository)

	customerConumer, err := consumers.NewCustomerConsumer(conn, "sagaCustomers", customerService)
	if err != nil {
		log.Fatalf("Failed to initialize message producer: %v", err.Error())
	}
	defer customerConumer.Chan.Close()

	var forever chan struct{}

	customerConumer.Listen(forever)
}
