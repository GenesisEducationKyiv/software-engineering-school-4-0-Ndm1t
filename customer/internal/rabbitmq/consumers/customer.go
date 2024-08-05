package consumers

import (
	"customer-service/internal/rabbitmq"
	"customer-service/internal/services"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

const (
	createCustomer string = "CreateCustomer"
)

type (
	CustomerConsumer struct {
		Chan            *amqp.Channel
		Queue           amqp.Queue
		topic           string
		customerService services.CustomerServiceInterface
	}
)

func NewCustomerConsumer(conn *amqp.Connection,
	topic string,
	customerService services.CustomerServiceInterface) (*CustomerConsumer, error) {
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
		Chan:            ch,
		Queue:           q,
		topic:           topic,
		customerService: customerService,
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
			var message rabbitmq.CustomerMessage
			err = json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Printf("failed to unmarshal message: %v", err)
				d.Nack(false, true)
				continue
			}
			switch message.EventType {
			case createCustomer:
				c.handleCreateCustomer(d, message)
			default:
				log.Printf("unhandled event type: %s", message.EventType)
				d.Nack(false, false)
			}
		}
	}()
	<-forever
}

func (c *CustomerConsumer) handleCreateCustomer(delivery amqp.Delivery, message rabbitmq.CustomerMessage) {
	c.customerService.Create(message.Data.TxID, message.Data.Email)
	delivery.Ack(false)
}
