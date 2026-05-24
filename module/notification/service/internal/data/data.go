package data

import (
	"context"

	biometricV1 "mall-go/api/biometrics/service/v1"
	departmentV1 "mall-go/api/department/service/v1"

	"mall-go/module/notification/service/internal/conf"
	"mall-go/module/notification/service/internal/data/model"
	"mall-go/module/notification/service/internal/data/model/migrate"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-redis/redis/v8"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var ProviderSet = wire.NewSet(NewData, NewEntClient, NewNotificationRepo, NewRedisClient)

// Data .
type Data struct {
	db               *model.Client
	log              *log.Helper
	BiometricClient  biometricV1.BiometricClient
	DepartmentClient departmentV1.DepartmentClient
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

	ctx := context.Background()
	// Run migration.
	err = client.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)
	if err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}

// NewData .
func NewData(entClient *model.Client, logger log.Logger, r registry.Discovery) (*Data, func(), error) {
	logHelper := log.NewHelper(logger)

	// Connect to biometric service
	connBio, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///biometric.service"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		logHelper.Fatalf("failed to connect to biometric service: %v", err)
	}

	// Connect to department service
	connDept, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///department.service"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		logHelper.Fatalf("failed to connect to department service: %v", err)
	}

	cleanup := func() {
		logHelper.Info("closing the data resources")
		_ = connBio.Close()
		_ = connDept.Close()
	}

	return &Data{
		db:               entClient,
		log:              logHelper,
		BiometricClient:  biometricV1.NewBiometricClient(connBio),
		DepartmentClient: departmentV1.NewDepartmentClient(connDept),
	}, cleanup, nil
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
