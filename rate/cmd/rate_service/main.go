package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"rate-service/internal/config"
	"rate-service/internal/crons"
	"rate-service/internal/providers"
	"rate-service/internal/providers/chain"
	"rate-service/internal/rabbitmq/producers"
	"rate-service/internal/server"
	"rate-service/internal/server/controllers"
	"rate-service/internal/services"
)

const topic = "rates"

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

	subscriptionProducer, err := producers.NewRateProducer(conn, topic)
	if err != nil {
		logger.Errorf("Failed to initialize message producer: %v", err.Error())
	}

	rateService := services.NewRateService(prepareChain(logger), subscriptionProducer)

	rateController := controllers.NewRateController(rateService)

	cronScheduler := crons.NewCronScheduler(rateService, logger)

	s := server.NewServer(rateController, cronScheduler)
	s.Scheduler.Start()
	defer s.Scheduler.Stop()
	s.Run()
}

func prepareChain(logger *zap.SugaredLogger) chain.Chain {
	privatProvider := providers.NewLoggingProvider("privat", providers.NewPrivatProvider(), logger)
	exchangeAPIProvider := providers.NewLoggingProvider("exchangeAPI", providers.NewExchangeAPIProvider(), logger)

	privatChain := chain.NewBaseChain(privatProvider)
	exchangeAPIChain := chain.NewBaseChain(exchangeAPIProvider)
	exchangeAPIChain.SetNext(privatChain)
	return exchangeAPIChain
}
