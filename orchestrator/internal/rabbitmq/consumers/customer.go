package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
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
	Logger interface {
		Warnf(template string, arguments ...interface{})
	}

	SubscriptionProducerInterface interface {
		Publish(eventType string, customer models.Subscription, ctx context.Context) error
	}

	CustomerConsumer struct {
		Chan                 *amqp.Channel
		Queue                amqp.Queue
		topic                string
		subscriptionProducer SubscriptionProducerInterface
		logger               Logger
	}
)

func NewCustomerConsumer(conn *amqp.Connection,
	topic string,
	subscriptionProducer SubscriptionProducerInterface,
	logger Logger) (*CustomerConsumer, error) {
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
		logger:               logger,
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
		c.logger.Warnf("failed to consume customers: %v", err.Error())
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				c.logger.Warnf("recovered from panic: %v", r)
			}
		}()

		for d := range msgs {
			var message rabbitmq.CustomerMessage
			err = json.Unmarshal(d.Body, &message)
			if err != nil {
				c.logger.Warnf("failed to unmarshal message: %v", err)
				d.Nack(false, false)
				continue
			}
			switch message.EventType {
			case customerCreated:
				c.handleCustomerCreated(d, message)
			case customerFailed:
				c.handleCustomerFailed(d, message)
			default:
				c.logger.Warnf("unhandled event type: %s", message.EventType)
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
		c.logger.Warnf("failed to create customer message: %v", err)
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
		c.logger.Warnf("failed to create customer message: %v", err)
		delivery.Nack(false, true)
		return
	}

	delivery.Ack(false)
}
