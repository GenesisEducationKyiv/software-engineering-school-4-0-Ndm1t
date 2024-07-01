package main

import (
	"gateway/internal/clients"
	"gateway/internal/config"
	"gateway/internal/server"
	"gateway/internal/server/controllers"
	"log"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}

	rateClient := clients.NewRateClient()
	subscriptionClient := clients.NewSubscriptionClient()

	rateController := controllers.NewRateController(rateClient)
	subscriptionController := controllers.NewSubscriptionController(subscriptionClient)

	s := server.NewServer(rateController, subscriptionController)

	s.Run()
}
