package providers

import (
	"log"
	"rate-service/internal/services"
)

type (
	LoggingProvider struct {
		name         string
		rateProvider services.IRateAPIProvider
	}
)

func NewLoggingProvider(name string, rateProvider services.IRateAPIProvider) *LoggingProvider {
	return &LoggingProvider{
		name:         name,
		rateProvider: rateProvider,
	}
}
func (l *LoggingProvider) FetchRate() (*float64, error) {
	rate, err := l.rateProvider.FetchRate()
	if err != nil {
		log.Printf("%v provider returned error: %v", l.name, err)
		return nil, err
	}
	log.Printf("%v provider returned value: %v", l.name, *rate)
	return rate, err
}
