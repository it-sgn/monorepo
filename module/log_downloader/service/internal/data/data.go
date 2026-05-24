package data

import (
	"context"
	"fmt"

	"mall-go/module/log_downloader/service/internal/conf"
	"mall-go/module/log_downloader/service/internal/data/model"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var ProviderSet = wire.NewSet(NewData, NewEntClient, NewDownloadRepo, NewCronZKRepo, NewRedisClient, NewKafkaProducer, NewGreeterRepo)

// type KafkaProducer struct {
// 	producer *kafka.Producer
// 	topic    string
// 	logger   *log.Helper
// }

// Data .
type Data struct {
	db  *model.Client
	log *log.Helper
	// rdb *redis.Client
	Kafka *KafkaProducer
}

func NewEntClient(conf *conf.Data, logger log.Logger) *model.Client {
	l := log.NewHelper(logger)
	l.Infof("Attempting to connect to DB. Driver: %s, Source: %s", conf.Database.Driver, conf.Database.Source)
	client, err := model.Open(
		conf.Database.Driver,
		conf.Database.Source,
	)
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
	}
	client = client.Debug()
	// ctx := context.Background()
	// // Run migration.
	// // err = client.Schema.Create(
	// // 	ctx,
	// // 	migrate.WithDropIndex(true),
	// // 	migrate.WithDropColumn(true),
	// // )
	// if err != nil {
	// 	log.Fatalf("failed creating schema resources: %v", err)
	// }
	return client
}

func NewData(conf *conf.Data, entClient *model.Client, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(logger)

	// Inisialisasi Kafka producer (dari confluent-kafka-go)
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": conf.Kafka.Addrs[0],
		"security.protocol": "SASL_PLAINTEXT",
		"sasl.mechanisms":   "SCRAM-SHA-256",
		"sasl.username":     conf.Kafka.Username,
		"sasl.password":     conf.Kafka.Password,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	kafkaProducer := &KafkaProducer{
		producer: producer,
		topic:    conf.Kafka.Topic,
		log:      logHelper,
	}

	data := &Data{
		db:    entClient,
		log:   logHelper,
		Kafka: kafkaProducer,
	}

	cleanup := func() {
		logHelper.Info("closing the data resources")
		producer.Close()
	}

	return data, cleanup, nil
}
func NewRedisClient(c *conf.Data, logger log.Logger) *redis.Client {
	l := log.NewHelper(logger)
	if c.Redis == nil {
		l.Fatalf("Redis configuration is missing in conf.Data")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       0, // Assuming Db is int32 or similar in conf, cast to int
	})

	// Ping the Redis server to ensure connection is established
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		l.Fatalf("failed to connect to redis: %v", err)
	}
	l.Info("Redis client connected successfully")
	return rdb
}

func WithTx(ctx context.Context, client *model.Client, fn func(tx *model.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = errors.Wrapf(err, "rolling back transaction: %v", rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrapf(err, "committing transaction: %v", err)
	}
	return nil
}
