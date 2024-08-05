package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"orchestrator/internal/models"
	"orchestrator/internal/rabbitmq"
)

const (
	subscriptionCreated string = "SubscriptionCreated"
	createCustomer      string = "CreateCustomer"
)

type (
	CustomerProducerInterface interface {
		Publish(eventType string, customer models.Customer, ctx context.Context) error
	}

	SubscriptionConsumer struct {
		Chan             *amqp.Channel
		Queue            amqp.Queue
		topic            string
		customerProducer CustomerProducerInterface
	}
)

func NewSubscriptionConsumer(conn *amqp.Connection,
	topic string,
	customerProducer CustomerProducerInterface) (*SubscriptionConsumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create rabbit channel: %v", err)
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
		return nil, fmt.Errorf("failed to create rabbit queue: %v", err)
	}

	return &SubscriptionConsumer{
		Chan:             ch,
		Queue:            q,
		topic:            topic,
		customerProducer: customerProducer,
	}, nil
}

func (c *SubscriptionConsumer) Listen(forever chan struct{}) {
	msgs, err := c.Chan.Consume(
		c.Queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("failed to consume subscriptions: %v", err.Error())
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("recovered from panic: %v", r)
			}
		}()

		for d := range msgs {
			var message rabbitmq.SubscriptionMessage
			err = json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Printf("failed to unmarshal message: %v", err)
				d.Nack(false, true)
				continue
			}
			switch message.EventType {
			case subscriptionCreated:
				c.handleCreated(d, message)
			default:
				log.Printf("unhandled event type: %s", message.EventType)
				d.Nack(false, false)
			}
		}
	}()
	<-forever
}

func (c *SubscriptionConsumer) handleCreated(delivery amqp.Delivery, message rabbitmq.SubscriptionMessage) {
	err := c.customerProducer.Publish(createCustomer, models.Customer{
		TxID:  message.Data.TxID,
		Email: message.Data.Email,
	}, context.Background())
	if err != nil {
		log.Printf("failed to create customer message: %v", err)
		delivery.Nack(false, true)
		return
	}

	delivery.Ack(false)
}
