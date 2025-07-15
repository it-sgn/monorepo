package data

import (
	"context"

	"mall-go/module/biometric/service/internal/conf"
	"mall-go/module/biometric/service/internal/data/model"
	"mall-go/module/biometric/service/internal/data/model/migrate"

	"github.com/go-kratos/kratos/v2/log"
	// _ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var ProviderSet = wire.NewSet(NewData, NewEntClient, NewBiometricRepo)

// Data .
type Data struct {
	db  *model.Client
	log *log.Helper
	// rdb *redis.Client
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
	// if err := client.Schema.Create(context.Background(), migrate.WithForeignKeys(false)); err != nil {
	// 	l.Fatalf("failed creating schema resources: %v", err)
	// }
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
func NewData(entClient *model.Client, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db:  entClient,
		log: log.NewHelper(logger),
	}, cleanup, nil
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
