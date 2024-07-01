package main

import (
	"log"
	"rate-service/internal/config"
	"rate-service/internal/providers"
	"rate-service/internal/providers/chain"
	"rate-service/internal/server"
	"rate-service/internal/server/controllers"
	"rate-service/internal/services"
)

func main() {

	err := config.LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}
	rateService := services.NewRateService(prepareChain())

	rateController := controllers.NewRateController(rateService)

	s := server.NewServer(rateController)
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
