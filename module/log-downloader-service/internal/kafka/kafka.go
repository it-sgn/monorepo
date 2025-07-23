package kafka

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"time"

// 	"mall-go/module/log-downloader-service/internal/zk"

// 	"github.com/segmentio/kafka-go"
// )

// var kafkaWriter *kafka.Writer

// func InitKafka(brokers []string, topic string) {
// 	kafkaWriter = kafka.NewWriter(kafka.WriterConfig{
// 		Brokers:      brokers,
// 		Topic:        topic,
// 		Balancer:     &kafka.LeastBytes{},
// 		RequiredAcks: int(kafka.RequireOne),
// 		BatchTimeout: 10 * time.Millisecond,
// 	})
// }

// func SendToKafka(data zk.LogData) error {
// 	if kafkaWriter == nil {
// 		return fmt.Errorf("Kafka writer not initialized")
// 	}

// 	// Serialize LogData to JSON
// 	value, err := json.Marshal(data)
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal LogData: %w", err)
// 	}

// 	msg := kafka.Message{
// 		Key:   []byte(fmt.Sprintf("%s-%d", data.DeviceIP, data.UserID)),
// 		Value: value,
// 		Time:  time.Now(),
// 	}

// 	err = kafkaWriter.WriteMessages(context.Background(), msg)
// 	if err != nil {
// 		log.Printf("❌ Failed to send message to Kafka: %v", err)
// 		return err
// 	}

// 	log.Printf("✅ Kafka: sent log from %s user %d", data.DeviceIP, data.UserID)
// 	return nil
// }

// // package kafka

// // import (
// // 	"encoding/json"
// // 	"log"

// // 	"github.com/Shopify/sarama"
// // )

// // var producer sarama.SyncProducer

// // func InitProducer(brokers []string) error {
// // 	config := sarama.NewConfig()
// // 	config.Producer.Return.Successes = true
// // 	config.Producer.Retry.Max = 3

// // 	var err error
// // 	producer, err = sarama.NewSyncProducer(brokers, config)
// // 	return err
// // }

// // func sendToKafka(topic string, ip, severity, message string) {
// // 	payload := map[string]interface{}{
// // 		"ip":       ip,
// // 		"severity": severity,
// // 		"message":  message,
// // 	}

// // 	data, err := json.Marshal(payload)
// // 	if err != nil {
// // 		log.Printf("[ERROR] Gagal marshal payload Kafka: %v", err)
// // 		return
// // 	}

// // 	msg := &sarama.ProducerMessage{
// // 		Topic: topic,
// // 		Value: sarama.ByteEncoder(data),
// // 	}

// // 	_, _, err = producer.SendMessage(msg)
// // 	if err != nil {
// // 		log.Printf("[ERROR] Gagal kirim ke Kafka: %v", err)
// // 	}
// // }
