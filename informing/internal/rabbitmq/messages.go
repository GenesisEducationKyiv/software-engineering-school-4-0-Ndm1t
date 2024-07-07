package rabbitmq

import (
	"github.com/google/uuid"
	"informing-service/internal/models"
	"time"
)

type (
	SubscriptionMessage struct {
		EventID     string              `json:"eventId"`
		EventType   string              `json:"eventType"`
		AggregateID string              `json:"aggregateId"`
		Timestamp   string              `json:"timestamp"`
		Data        models.Subscription `json:"data"`
	}

	RateMessage struct {
		EventID     string      `json:"eventId"`
		EventType   string      `json:"eventType"`
		AggregateID *string     `json:"aggregateId"`
		Timestamp   string      `json:"timestamp"`
		Data        models.Rate `json:"data"`
	}
)

func NewSubscriptionMessage(eventType string, subscription models.Subscription) SubscriptionMessage {
	return SubscriptionMessage{
		EventID:     uuid.NewString(),
		EventType:   eventType,
		AggregateID: subscription.Email,
		Timestamp:   time.Now().String(),
		Data:        subscription,
	}
}
