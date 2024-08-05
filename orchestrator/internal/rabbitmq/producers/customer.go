package producers

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"orchestrator/internal/models"
	"orchestrator/internal/rabbitmq"
	"time"
)

type (
	CustomerProducer struct {
		Chan  *amqp.Channel
		Queue amqp.Queue
		topic string
	}
)

func NewCustomerProducer(conn *amqp.Connection, topic string) (*CustomerProducer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create rabbit chanel: %v", err)
	}

	q, err := ch.QueueDeclare(
		topic,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create rabbit chanel: %v", err)
	}

	return &CustomerProducer{
		Chan:  ch,
		Queue: q,
		topic: topic,
	}, nil
}

func (p *CustomerProducer) Publish(eventType string, customer models.Customer, ctx context.Context) error {
	customerMessage := rabbitmq.NewCustomerMessage(eventType, customer)

	body, err := json.Marshal(customerMessage)

	if err != nil {
		return fmt.Errorf("failed to create message body: %v", err.Error())
	}

	newCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err = p.Chan.PublishWithContext(
		newCtx,
		"",
		p.topic,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err.Error())
	}

	return nil

}
