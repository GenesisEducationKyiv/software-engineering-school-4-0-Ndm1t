package consumers

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"subscription-service/internal/models"
	"subscription-service/internal/rabbitmq"
	"subscription-service/internal/services"
)

const (
	verifySubscription = "VerifySubscription"
	deleteSubscription = "DeleteSubscription"
)

type (
	SubscriptionConsumer struct {
		Chan                *amqp.Channel
		Queue               amqp.Queue
		topic               string
		subscriptionService services.ISubscriptionService
		logger              *zap.SugaredLogger
	}
)

func NewSubscriptionConsumer(conn *amqp.Connection,
	topic string,
	subscriptionService services.ISubscriptionService,
	logger *zap.SugaredLogger) (*SubscriptionConsumer, error) {
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
		Chan:                ch,
		Queue:               q,
		topic:               topic,
		subscriptionService: subscriptionService,
		logger:              logger,
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
		c.logger.Warnf("failed to consume subscriptions: %v", err.Error())
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				c.logger.Warnf("recovered from panic: %v", r)
			}
		}()

		for d := range msgs {
			var message rabbitmq.EmailMessage
			err = json.Unmarshal(d.Body, &message)
			if err != nil {
				c.logger.Warnf("failed to unmarshal message: %v", err)
				d.Nack(false, false)
				continue
			}
			switch message.EventType {
			case verifySubscription:
				c.handleVerify(d, message)
			case deleteSubscription:
				c.handleDelete(d, message.Data.Email)
			default:
				c.logger.Warnf("unhandled event type: %s", message.EventType)
				d.Nack(false, false)
			}
		}
	}()
}

func (c *SubscriptionConsumer) handleVerify(delivery amqp.Delivery, message rabbitmq.EmailMessage) {
	err := c.subscriptionService.UpdateSate(message.Data.Email, models.Verified)
	if err != nil {
		c.logger.Warnf("failed to create subscription: %v", err)
		delivery.Nack(false, true)
		return
	}

	delivery.Ack(false)
}

func (c *SubscriptionConsumer) handleDelete(delivery amqp.Delivery, email string) {
	err := c.subscriptionService.Delete(email)
	if err != nil {
		c.logger.Warnf("failed to get rate from database: %v", err)
		delivery.Nack(false, true)
		return
	}
	delivery.Ack(false)
}
