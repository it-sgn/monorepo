package data

import (
	"context"

	// biometricV1 "mall-go/api/biometric/service/v1"

	employersV1 "mall-go/api/employers/service/v1"

	"mall-go/module/attendance/service/internal/conf"
	"mall-go/module/attendance/service/internal/data/model"
	"mall-go/module/attendance/service/internal/data/model/migrate"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var ProviderSet = wire.NewSet(NewData, NewEntClient, NewAttendanceRepo)

// Data .
type Data struct {
	db  *model.Client
	log *log.Helper
	// rdb *redis.Client
	EmployersClient employersV1.EmployersClient
	// DepartmentClient departmentV1.DepartmentClient
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

// // NewData .
// func NewData(entClient *model.Client, logger log.Logger) (*Data, func(), error) {
// 	cleanup := func() {
// 		log.NewHelper(logger).Info("closing the data resources")
// 	}
// 	return &Data{
// 		db:  entClient,
// 		log: log.NewHelper(logger),
// 	}, cleanup, nil
// }

// NewData .
func NewData(entClient *model.Client, logger log.Logger, r registry.Discovery) (*Data, func(), error) {
	logHelper := log.NewHelper(logger)

	// // Connect to biometric service
	// connBio, err := grpc.DialInsecure(
	// 	context.Background(),
	// 	grpc.WithEndpoint("discovery:///biometric.service"),
	// 	grpc.WithDiscovery(r),
	// )
	// if err != nil {
	// 	logHelper.Fatalf("failed to connect to biometric service: %v", err)
	// }

	// // Connect to biometric service
	connEmp, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///employers.service"),
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
		_ = connEmp.Close()
		_ = connDept.Close()
	}

	return &Data{
		db:              entClient,
		log:             logHelper,
		EmployersClient: employersV1.NewEmployersClient(connEmp),
		// DepartmentClient: departmentV1.NewDepartmentClient(connDept),
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
