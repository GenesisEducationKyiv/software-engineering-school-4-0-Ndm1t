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
	sendEmail           string = "SendEmail"
	subscriptionDeleted string = "SubscriptionDeleted"
)

type (
	SubscriptionRepositoryInterface interface {
		Create(email string) (*models.Subscription, error)
		Delete(subscription models.Subscription) error
	}

	EmailSenderInterface interface {
		SendInforming(email string, rate float64) error
	}

	SubscriptionConsumer struct {
		Chan                   *amqp.Channel
		Queue                  amqp.Queue
		topic                  string
		subscriptionRepository SubscriptionRepositoryInterface
		rateRepository         RateRepositoryInterface
		emailSender            EmailSenderInterface
	}
)

func NewSubscriptionConsumer(conn *amqp.Connection,
	topic string,
	subscriptionRepository SubscriptionRepositoryInterface,
	rateRepository RateRepositoryInterface,
	emailSender EmailSenderInterface) (*SubscriptionConsumer, error) {
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
		Chan:                   ch,
		Queue:                  q,
		topic:                  topic,
		subscriptionRepository: subscriptionRepository,
		rateRepository:         rateRepository,
		emailSender:            emailSender,
	}, nil
}

func (c *SubscriptionConsumer) Listen() {
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
			case sendEmail:
				c.handleSendEmail(d, message.Data.Email)
			case subscriptionDeleted:
				c.handleSubscriptionDeleted(d, message)
			default:
				log.Printf("unhandled event type: %s", message.EventType)
				d.Nack(false, false)
			}
		}
	}()
}

func (c *SubscriptionConsumer) handleCreated(delivery amqp.Delivery, message rabbitmq.SubscriptionMessage) {
	_, err := c.subscriptionRepository.Create(message.Data.Email)
	if err != nil {
		log.Printf("failed to create subscription: %v", err)
		delivery.Nack(false, true)
		return
	}

	delivery.Ack(false)
}

func (c *SubscriptionConsumer) handleSendEmail(delivery amqp.Delivery, email string) {
	rate, err := c.rateRepository.GetLatest()
	if err != nil {
		log.Printf("failed to get rate from database: %v", err)
		delivery.Nack(false, true)
		return
	}

	err = c.emailSender.SendInforming(email, rate.Rate)
	if err != nil {
		log.Printf("failed to send email: %v", err)
		delivery.Nack(false, true)
		return
	}

	delivery.Ack(false)
}

func (c *SubscriptionConsumer) handleSubscriptionDeleted(delivery amqp.Delivery, message rabbitmq.SubscriptionMessage) {
	err := c.subscriptionRepository.Delete(message.Data)
	if err != nil {
		log.Printf("failed to delete subscription: %v", err)
		delivery.Nack(false, true)
		return
	}

	delivery.Ack(false)
}
