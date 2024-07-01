package main

import (
	"informing-service/internal/clients"
	"informing-service/internal/config"
	"informing-service/internal/crons"
	"informing-service/internal/mailers"
	"informing-service/internal/server"
	"informing-service/internal/server/controllers"
	"informing-service/internal/services"
	"log"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}

	smtpSender := mailers.NewSMTPEmailSender()

	rateCLient := clients.NewRateClient()
	subscriptionClient := clients.NewSubscriptionClient()

	informingService := services.NewInformingService(subscriptionClient, rateCLient, smtpSender)

	cronScheduler := crons.NewCronScheduler(informingService)

	informingController := controllers.NewInformingController(informingService)

	s := server.NewServer(informingController, cronScheduler)
	s.Scheduler.Start()
	defer s.Scheduler.Stop()
	s.Run()
}
