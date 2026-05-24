package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// Producer defines the interface for a Kafka producer.
type Producer interface {
	Produce(ctx context.Context, key []byte, value []byte) error
	Close() error
}

// KafkaProducer implements the Producer interface using segmentio/kafka-go.
type KafkaProducer struct {
	writer *kafka.Writer
	topic  string
}

// NewKafkaProducer creates a new KafkaProducer.
func NewKafkaProducer(brokerUrls []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokerUrls...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{}, // Or &kafka.Hash{}, or other balancers
		},
		topic: topic,
	}
}

// Produce sends a message to Kafka.
func (p *KafkaProducer) Produce(ctx context.Context, key []byte, value []byte) error {
	return p.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   key,
			Value: value,
		},
	)
}

// Close closes the Kafka producer.
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
