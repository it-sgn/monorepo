package data

import (
	"context"

	departmentv1 "mall-go/api/department/service/v1"
	employerv1 "mall-go/api/employers/service/v1"
	orgv1 "mall-go/api/organization/service/v1"
	"mall-go/module/attendance-raw/service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(
	NewDB,
	NewRedisClient,
	NewKafkaProducer,
	NewEmployersClient,
	NewDepartmentClient,
	NewOrganizationClient,
	NewData,
	NewAttendanceRawRepo,
)

// Data holds all external resources.
type Data struct {
	db               *gorm.DB
	log              *log.Helper
	rdb              *redis.Client
	Kafka            *KafkaProducer
	EmployersClient  employerv1.EmployersClient
	DepartmentClient departmentv1.DepartmentClient
	OrgV1Client      orgv1.OrganizationServiceClient
}

// NewData initializes the application's shared resources.
func NewData(
	conf *conf.Data,
	logger log.Logger,
	db *gorm.DB,
	rdb *redis.Client,
	kafkaProducer *KafkaProducer,
	empClient employerv1.EmployersClient,
	deptClient departmentv1.DepartmentClient,
) (*Data, func(), error) {
	logHelper := log.NewHelper(logger)

	data := &Data{
		db:               db,
		log:              logHelper,
		rdb:              rdb,
		Kafka:            kafkaProducer,
		EmployersClient:  empClient,
		DepartmentClient: deptClient,
	}

	cleanup := func() {
		logHelper.Info("🧹 Cleaning up data resources")
		if kafkaProducer != nil && kafkaProducer.producer != nil {
			// Optional: kafkaProducer.Close() if implemented
		}
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}

	return data, cleanup, nil
}

// NewDB initializes a PostgreSQL GORM connection.
func NewDB(conf *conf.Data, logger log.Logger) (*gorm.DB, error) {
	logHelper := log.NewHelper(logger)

	dsn := conf.Database.Source
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logHelper.Errorf("❌ Failed to open GORM DB: %v", err)
		return nil, err
	}

	logHelper.Info("✅ GORM PostgreSQL connected")
	return db, nil
}

// NewRedisClient initializes a Redis client.
func NewRedisClient(conf *conf.Data, logger log.Logger) *redis.Client {
	logHelper := log.NewHelper(logger)

	if conf.Redis == nil {
		logHelper.Fatalf("❌ Redis configuration is missing")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logHelper.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	logHelper.Info("✅ Redis client connected")
	return rdb
}
