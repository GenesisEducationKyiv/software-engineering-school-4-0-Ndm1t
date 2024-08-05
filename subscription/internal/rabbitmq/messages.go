package rabbitmq

import (
	"github.com/google/uuid"
	"subscription-service/internal/models"
	"time"
)

type (
	EmailMessage struct {
		EventID     string       `json:"eventId"`
		EventType   string       `json:"eventType"`
		AggregateID string       `json:"aggregateId"`
		Timestamp   string       `json:"timestamp"`
		Data        models.Email `json:"data"`
	}
)

func NewEmailMessage(eventType string, email models.Email) EmailMessage {
	return EmailMessage{
		EventID:     uuid.NewString(),
		EventType:   eventType,
		AggregateID: email.Email,
		Timestamp:   time.Now().String(),
		Data:        email,
	}
}
