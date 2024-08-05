package rabbitmq

import (
	"github.com/google/uuid"
	"orchestrator/internal/models"
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

	CustomerMessage struct {
		EventID     string          `json:"eventId"`
		EventType   string          `json:"eventType"`
		AggregateID string          `json:"aggregateId"`
		Timestamp   string          `json:"timestamp"`
		Data        models.Customer `json:"data"`
	}
)

func NewSubscriptionMessage(eventType string, data models.Subscription) SubscriptionMessage {
	return SubscriptionMessage{
		EventID:     uuid.NewString(),
		EventType:   eventType,
		AggregateID: data.Email,
		Timestamp:   time.Now().String(),
		Data:        data,
	}
}

func NewCustomerMessage(eventType string, data models.Customer) CustomerMessage {
	return CustomerMessage{
		EventID:     uuid.NewString(),
		EventType:   eventType,
		AggregateID: data.Email,
		Timestamp:   time.Now().String(),
		Data:        data,
	}
}
