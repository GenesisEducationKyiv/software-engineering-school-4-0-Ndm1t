package consumers

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"informing-service/internal/models"
	"informing-service/internal/rabbitmq"
)

const (
	rateFetched string = "RateFetched"
)

type (
	RateRepositoryInterface interface {
		Create(rate models.Rate) (*models.Rate, error)
		GetLatest() (*models.Rate, error)
	}

	RateConsumer struct {
		Chan           *amqp.Channel
		Queue          amqp.Queue
		topic          string
		rateRepository RateRepositoryInterface
		logger         *zap.SugaredLogger
	}
)

func NewRateConsumer(conn *amqp.Connection, topic string, rateRepository RateRepositoryInterface, logger *zap.SugaredLogger) (*RateConsumer, error) {
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

	return &RateConsumer{
		Chan:           ch,
		Queue:          q,
		topic:          topic,
		rateRepository: rateRepository,
		logger:         logger,
	}, nil
}

func (c *RateConsumer) Listen() {
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
			var message rabbitmq.RateMessage
			err = json.Unmarshal(d.Body, &message)
			if err != nil {
				c.logger.Warnf("failed to unmarshal message: %v", err)
				d.Nack(false, false)
				continue
			}

			switch message.EventType {
			case rateFetched:
				_, err = c.rateRepository.Create(message.Data)
				if err != nil {
					c.logger.Warnf("failed to create rate row: %v", err)
					d.Nack(false, true)
				} else {
					d.Ack(false)
				}
			default:
				c.logger.Warnf("unhandled event type: %s", message.EventType)
				d.Nack(false, false)
			}
		}
	}()
}
