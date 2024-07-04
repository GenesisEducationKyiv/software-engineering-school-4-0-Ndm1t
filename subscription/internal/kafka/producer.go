package kafkaClient

import "github.com/confluentinc/confluent-kafka-go/kafka"

const kafkaServer = "localhost:9092"

func NewProducer() (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaServer,
	})

	if err != nil {
		return nil, err
	}

	return producer, err
}
