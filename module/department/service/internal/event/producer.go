package event

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type EventPublisher struct {
	producer *kafka.Producer
	topic    string
}

func NewEventPublisher(broker, topic string) (*EventPublisher, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
	})
	if err != nil {
		return nil, err
	}
	return &EventPublisher{producer: p, topic: topic}, nil
}

func (p *EventPublisher) Publish(ctx context.Context, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          data,
	}

	return p.producer.Produce(msg, nil)
}
