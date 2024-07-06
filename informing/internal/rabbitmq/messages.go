package rabbitmq

import (
	"informing-service/internal/models"
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
