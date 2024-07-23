package services

import (
	"context"
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"rate-service/internal/models"
	"time"
)

const (
	rateFetched               = "RateFetched"
	rateMessagePublishSuccess = "rate_publish_success"
	rateMessagePublishFail    = "rate_publish_fail"
)

type (
	IRateAPIProvider interface {
		FetchRate() (*float64, error)
	}

	IRateService interface {
		Get() (*float64, error)
	}

	RateProducerInterface interface {
		Publish(eventType string, rate models.Rate, ctx context.Context) error
	}

	RateService struct {
		APIProvider  IRateAPIProvider
		rateProducer RateProducerInterface
	}
)

func NewRateService(provider IRateAPIProvider, rateProducer RateProducerInterface) *RateService {
	return &RateService{
		APIProvider:  provider,
		rateProducer: rateProducer,
	}
}

func (r *RateService) Get() (*float64, error) {
	rate, err := r.APIProvider.FetchRate()
	if err != nil {
		return nil, err
	}

	rateData := models.Rate{
		Rate:      *rate,
		CreatedAt: time.Now(),
	}

	err = r.rateProducer.Publish(rateFetched, rateData, context.Background())
	if err != nil {
		metrics.GetOrCreateCounter(fmt.Sprintf(`%v{type=%q}`, rateMessagePublishFail, rateFetched)).Inc()
		return nil, err
	}
	metrics.GetOrCreateCounter(fmt.Sprintf(`%v{type=%q}`, rateMessagePublishSuccess, rateFetched)).Inc()

	return rate, nil
}
