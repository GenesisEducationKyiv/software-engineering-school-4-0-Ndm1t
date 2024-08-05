package rabbitmq

import (
	"customer-service/internal/models"
	"github.com/google/uuid"
	"time"
)

type (
	CustomerMessage struct {
		EventID     string          `json:"eventId"`
		EventType   string          `json:"eventType"`
		AggregateID string          `json:"aggregateId"`
		Timestamp   string          `json:"timestamp"`
		Data        models.Customer `json:"data"`
	}
)

func NewCustomerMessage(eventType string, customer models.Customer) CustomerMessage {
	return CustomerMessage{
		EventID:     uuid.NewString(),
		EventType:   eventType,
		AggregateID: customer.Email,
		Timestamp:   time.Now().String(),
		Data:        customer,
	}
}
