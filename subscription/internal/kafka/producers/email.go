package producers

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type (
	EmailProducer struct {
		producer *kafka.Producer
		topic    string
	}

	EmailMessage struct {
		Email string `json:"email"`
	}
)

func NewEmailProducer(producer *kafka.Producer, topic string) *EmailProducer {
	return &EmailProducer{
		producer: producer,
		topic:    topic,
	}
}

func (p *EmailProducer) Produce(email string) error {
	emailData := EmailMessage{
		Email: email,
	}

	value, err := json.Marshal(emailData)

	if err != nil {
		return err
	}

	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, nil)

	if err != nil {
		return err
	}

	return nil
}
