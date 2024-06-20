package main

import (
	"gses4_project/internal/container"
	"gses4_project/internal/crons"
	"gses4_project/internal/database"
	"gses4_project/internal/models"
	"gses4_project/internal/pkg"
	"gses4_project/internal/pkg/providers"
	"gses4_project/internal/pkg/providers/chain"
	"gses4_project/internal/server"
	"gses4_project/internal/server/controllers"
	"gses4_project/internal/services"
	"log"
)

func main() {
	pkg.LoadConfig()
	db := database.ConnectDatabase()

	if err := db.AutoMigrate(&models.Email{}); err != nil {
		log.Fatalf("Coulnd't migrate database: %v", err.Error())
	}

	appContainer := container.NewContainer(db)

	subscriptionRepository := database.NewSubscriptionRepository(db)

	smtpSender := pkg.NewSMTPEmailSender()

	rateService := services.NewRateService(prepareChain(), appContainer)
	subscriptionService := services.NewSubscriptionService(appContainer, subscriptionRepository)
	informingService := services.NewInformingService(appContainer, smtpSender, subscriptionRepository, rateService)

	cronScheduler := crons.NewCronScheduler(appContainer, informingService)

	rateController := controllers.NewRateController(appContainer, rateService)
	subscriptionController := controllers.NewSubscriptionController(appContainer, subscriptionService)

	s := server.NewServer(appContainer, rateController, subscriptionController, cronScheduler)
	s.Scheduler.Start()
	defer s.Scheduler.Stop()
	s.Run()
}

func prepareChain() chain.Chain {
	privatProvider := providers.NewLoggingProvider("privat", providers.NewPrivatProvider())
	exchangeAPIProvider := providers.NewLoggingProvider("exchangeAPI", providers.NewExchangeAPIProvider())

	privatChain := chain.NewBaseChain(privatProvider)
	exchangeAPIChain := chain.NewBaseChain(exchangeAPIProvider)
	exchangeAPIChain.SetNext(privatChain)
	return exchangeAPIChain
}
