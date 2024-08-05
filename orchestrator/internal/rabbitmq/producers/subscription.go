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

type Producer struct {
	Chan  *amqp.Channel
	Queue amqp.Queue
	topic string
}

func NewSubscriptionProducer(conn *amqp.Connection, topic string) (*Producer, error) {
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

	return &Producer{
		Chan:  ch,
		Queue: q,
		topic: topic,
	}, nil
}

func (p *Producer) Publish(eventType string, subscription models.Subscription, ctx context.Context) error {
	subscriptionMessage := rabbitmq.NewSubscriptionMessage(eventType, subscription)

	body, err := json.Marshal(subscriptionMessage)

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
