package data

import (
	"context"
	"encoding/json"
	"fmt"
	"mall-go/module/log_downloader/service/internal/biz"
	"mall-go/module/log_downloader/service/internal/conf"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-kratos/kratos/v2/log"
)

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
	log      *log.Helper
}

func NewKafkaProducer(cfg *conf.Data, logger log.Logger) *KafkaProducer {
	logHelper := log.NewHelper(log.With(logger, "module", "kafka-producer"))

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Addrs[0],
		"security.protocol": "SASL_PLAINTEXT",
		"sasl.mechanisms":   "SCRAM-SHA-256",
		"sasl.username":     cfg.Kafka.Username,
		"sasl.password":     cfg.Kafka.Password,
	})
	if err != nil {
		logHelper.Fatalf("failed to create kafka producer: %v", err)
	}

	return &KafkaProducer{
		producer: p,
		topic:    cfg.Kafka.Topic,
		log:      logHelper,
	}
}

func (k *KafkaProducer) PublishEvent(ctx context.Context, e *biz.Event) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return err
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &k.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(e.Key),
		Value: payload,
	}

	// Produce message asynchronously
	deliveryChan := make(chan kafka.Event, 1)
	err = k.producer.Produce(msg, deliveryChan)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	// Optionally wait for delivery report
	ev := <-deliveryChan
	m, ok := ev.(*kafka.Message)
	if !ok || m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %v", m.TopicPartition.Error)
	}

	k.log.Infof("delivered message to %v", m.TopicPartition)
	close(deliveryChan)
	return nil
}

func StartKafkaConsumer(ctx context.Context, cfg *conf.Data, logger log.Logger, handler func(*biz.Event) error) {
	logHelper := log.NewHelper(log.With(logger, "module", "kafka-consumer"))

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Addrs[0],
		"group.id":          cfg.Kafka.GroupId,
		"auto.offset.reset": "earliest",

		"security.protocol": "SASL_PLAINTEXT",
		"sasl.mechanisms":   "SCRAM-SHA-256",
		"sasl.username":     cfg.Kafka.Username,
		"sasl.password":     cfg.Kafka.Password,
	})
	if err != nil {
		logHelper.Fatalf("failed to create kafka consumer: %v", err)
	}

	err = consumer.SubscribeTopics([]string{cfg.Kafka.Topic}, nil)
	if err != nil {
		logHelper.Fatalf("failed to subscribe to topic: %v", err)
	}

	go func() {
		defer consumer.Close()
		for {
			msg, err := consumer.ReadMessage(-1)
			if err != nil {
				logHelper.Errorf("error reading message: %v", err)
				continue
			}

			var evt biz.Event
			if err := json.Unmarshal(msg.Value, &evt); err != nil {
				logHelper.Errorf("failed to unmarshal message: %v", err)
				continue
			}

			logHelper.Infof("consumed event: %v", evt)

			if err := handler(&evt); err != nil {
				logHelper.Errorf("event handler error: %v", err)
			}
		}
	}()
}
