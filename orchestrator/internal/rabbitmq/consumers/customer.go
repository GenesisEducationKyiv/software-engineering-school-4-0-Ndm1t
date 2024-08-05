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
	customerCreated    string = "CustomerCreated"
	customerFailed     string = "CustomerFailed"
	verifySubscription string = "VerifySubscription"
	deleteSubscription string = "DeleteSubscription"
)

type (
	SubscriptionProducerInterface interface {
		Publish(eventType string, customer models.Subscription, ctx context.Context) error
	}

	CustomerConsumer struct {
		Chan                 *amqp.Channel
		Queue                amqp.Queue
		topic                string
		subscriptionProducer SubscriptionProducerInterface
	}
)

func NewCustomerConsumer(conn *amqp.Connection,
	topic string,
	subscriptionProducer SubscriptionProducerInterface) (*CustomerConsumer, error) {
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

	return &CustomerConsumer{
		Chan:                 ch,
		Queue:                q,
		topic:                topic,
		subscriptionProducer: subscriptionProducer,
	}, nil
}

func (c *CustomerConsumer) Listen(forever chan struct{}) {
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
		log.Printf("failed to consume customers: %v", err.Error())
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("recovered from panic: %v", r)
			}
		}()

		for d := range msgs {
			var message rabbitmq.CustomerMessage
			err = json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Printf("failed to unmarshal message: %v", err)
				d.Nack(false, true)
				continue
			}
			switch message.EventType {
			case customerCreated:
				c.handleCustomerCreated(d, message)
			case customerFailed:
				c.handleCustomerFailed(d, message)
			default:
				log.Printf("unhandled event type: %s", message.EventType)
				d.Nack(false, false)
			}
		}
	}()
	<-forever
}

func (c *CustomerConsumer) handleCustomerCreated(delivery amqp.Delivery, message rabbitmq.CustomerMessage) {
	err := c.subscriptionProducer.Publish(verifySubscription, models.Subscription{
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

func (c *CustomerConsumer) handleCustomerFailed(delivery amqp.Delivery, message rabbitmq.CustomerMessage) {
	err := c.subscriptionProducer.Publish(deleteSubscription, models.Subscription{
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
