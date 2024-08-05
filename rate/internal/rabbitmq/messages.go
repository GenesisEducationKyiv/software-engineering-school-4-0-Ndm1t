package rabbitmq

import (
	"github.com/google/uuid"
	"rate-service/internal/models"
	"time"
)

type (
	RateMessage struct {
		EventID     string      `json:"eventId"`
		EventType   string      `json:"eventType"`
		AggregateID *string     `json:"aggregateId"`
		Timestamp   string      `json:"timestamp"`
		Data        models.Rate `json:"data"`
	}
)

func NewRateMessage(eventType string, rate models.Rate) RateMessage {
	return RateMessage{
		EventID:     uuid.NewString(),
		EventType:   eventType,
		AggregateID: nil,
		Timestamp:   time.Now().String(),
		Data:        rate,
	}
}
