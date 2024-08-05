package main

import (
	"gateway/internal/clients"
	"gateway/internal/config"
	"gateway/internal/server"
	"gateway/internal/server/controllers"
	"go.uber.org/zap"
	"log"
)

func main() {
	err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf(err.Error())
	}

	logger := zap.Must(zap.NewProduction()).Sugar()

	rateClient := clients.NewRateClient(logger)
	subscriptionClient := clients.NewSubscriptionClient(logger)

	rateController := controllers.NewRateController(rateClient)
	subscriptionController := controllers.NewSubscriptionController(subscriptionClient)

	s := server.NewServer(rateController, subscriptionController)

	s.Run()
}
