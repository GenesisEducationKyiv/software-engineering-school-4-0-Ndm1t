package consumers

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"informing-service/internal/models"
	"informing-service/internal/rabbitmq"
	"log"
)

const (
	subscriptionCreated string = "SubscriptionCreated"
)

type (
	SubscriptionRepositoryInterface interface {
		Create(email string) (*models.Subscription, error)
		ListSubscribed() ([]models.Subscription, error)
		Update(subscription models.Subscription) (*models.Subscription, error)
	}

	Consumer struct {
		Chan                   *amqp.Channel
		Queue                  amqp.Queue
		topic                  string
		subscriptionRepository SubscriptionRepositoryInterface
	}
)

func NewSubscriptionConsumer(conn *amqp.Connection, topic string, subscriptionRepository SubscriptionRepositoryInterface) (*Consumer, error) {
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

	return &Consumer{
		Chan:                   ch,
		Queue:                  q,
		topic:                  topic,
		subscriptionRepository: subscriptionRepository,
	}, nil
}

func (c *Consumer) Listen() {
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
				_, err = c.subscriptionRepository.Create(message.Data.Email)
				if err != nil {
					log.Printf("failed to create subscription: %v", err)
					d.Nack(false, true)
				} else {
					d.Ack(false)
				}
			default:
				log.Printf("unhandled event type: %s", message.EventType)
				d.Nack(false, false)
			}
		}
	}()
}
