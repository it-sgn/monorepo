package db

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"mall-go/module/log-downloader-service/internalent"
// 	"mall-go/ent/attendanceLog"
// 	"mall-go/module/log-downloader-service/internal/zk"
// )

// func sendToDatabase(ctx context.Context, client *ent.Client, logs []zk.LogData) error {
// 	if len(logs) == 0 {
// 		return nil
// 	}

// 	// Ambil semua log unik berdasarkan (user_id, device_ip, timestamp)
// 	type key struct {
// 		UserID    int64
// 		DeviceIP  string
// 		Timestamp time.Time
// 	}
// 	uniqueMap := make(map[key]zk.LogData)

// 	for _, logEntry := range logs {
// 		k := key{UserID: logEntry.UserID, DeviceIP: logEntry.DeviceIP, Timestamp: logEntry.Timestamp}
// 		uniqueMap[k] = logEntry
// 	}

// 	// Ambil key yang akan dicek apakah sudah ada di database
// 	keys := make([]key, 0, len(uniqueMap))
// 	for k := range uniqueMap {
// 		keys = append(keys, k)
// 	}

// 	// Query semua data yang sudah ada
// 	existingLogs, err := client.AttendanceLog.
// 		Query().
// 		Where(
// 			attendanceLog.Or(buildPredicates(keys)...)...,
// 		).
// 		All(ctx)

// 	if err != nil {
// 		return fmt.Errorf("query existing logs failed: %w", err)
// 	}

// 	// Buat map data yang sudah ada
// 	exists := map[key]bool{}
// 	for _, row := range existingLogs {
// 		k := key{UserID: row.UserID, DeviceIP: row.DeviceIP, Timestamp: row.Timestamp}
// 		exists[k] = true
// 	}

// 	// Siapkan batch insert untuk data yang belum ada
// 	var bulk []*ent.AttendanceLogCreate
// 	for k, logEntry := range uniqueMap {
// 		if exists[k] {
// 			continue
// 		}
// 		bulk = append(bulk, client.AttendanceLog.
// 			Create().
// 			SetUserID(logEntry.UserID).
// 			SetDeviceIP(logEntry.DeviceIP).
// 			SetAttLog(logEntry.Attendace),
// 		)
// 	}

// 	if len(bulk) == 0 {
// 		return nil // tidak ada yang perlu disimpan
// 	}

// 	// Simpan sekaligus
// 	if err := client.AttendanceLog.CreateBulk(bulk...).Exec(ctx); err != nil {
// 		return fmt.Errorf("batch insert failed: %w", err)
// 	}

// 	return nil
// }

// package db

// import (
// 	"context"
// 	"fmt"
// 	_ "mall-go/module/log-downloader-service/internal/db/model"
// 	_ "mall-go/module/log-downloader-service/internal/db/model/attendancelog"
// 	"mall-go/module/log-downloader-service/internal/zk"

// 	// 	"mall-go/module/log-downloader-service/internal/utils"
// 	// 	"mall-go/module/log-downloader-service/internal/zk"

// 	"entgo.io/ent"
// )

// func sendToDatabase(ctx context.Context, client *ent.Client, logs []zk.LogData) error {
// 	for _, logEntry := range logs {
// 		// Cek apakah data dengan kombinasi user_id + timestamp + device_ip sudah ada
// 		exist, err := client.AttendanceLog.
// 			Query().
// 			Where(
// 				attendanceLog.UserID(logEntry.UserID),
// 				attendanceLog.DeviceIP(logEntry.DeviceIP),
// 				attendanceLog.Timestamp(logEntry.Timestamp),
// 			).
// 			Exist(ctx)

// 		if err != nil {
// 			return fmt.Errorf("error checking existing log: %w", err)
// 		}
// 		if exist {
// 			continue // skip jika sudah ada
// 		}

// 		// Simpan jika belum ada
// 		_, err = client.AttendanceLog.
// 			Create().
// 			SetUserID(logEntry.UserID).
// 			SetDeviceIP(logEntry.DeviceIP).
// 			SetTimestamp(logEntry.Timestamp).
// 			Save(ctx)

// 		if err != nil {
// 			return fmt.Errorf("error inserting log: %w", err)
// 		}
// 	}
// 	return nil
// }

// // package db

// // import (
// // 	"context"
// // 	"fmt"

// // 	"mall-go/module/log-downloader-service/internal/db/model"
// // 	"mall-go/module/log-downloader-service/internal/db/model/attendancelog"
// // 	"mall-go/module/log-downloader-service/internal/utils"
// // 	"mall-go/module/log-downloader-service/internal/zk"

// // 	"github.com/sirupsen/logrus"
// // )

// // func SendToDatabase(ctx context.Context, client *model.Client, logs []zk.LogData) error {
// // 	keyMap, keys := utils.BuildKeyMap(logs)
// // 	predicates := utils.BuildPredicates(keys)

// // 	existing, err := client.AttendanceLog.
// // 		Query().
// // 		Where(attendancelog.Or(predicates...)).
// // 		All(ctx)
// // 	if err != nil {
// // 		return fmt.Errorf("query existing logs: %w", err)
// // 	}

// // 	for _, e := range existing {
// // 		k := utils.LogKey{UserID: e.UserID, DeviceIP: e.DeviceIP, Time: e.AttLog, Status: e.Status}
// // 		delete(keyMap, k)
// // 	}

// // 	var bulk []*model.AttendanceLogCreate
// // 	for _, log := range keyMap {
// // 		bulk = append(bulk, client.AttendanceLog.
// // 			Create().
// // 			SetUserID(log.UserID).
// // 			SetDeviceIP(log.DeviceIP).
// // 			SetAttLog(log.Attendace).
// // 			SetStatus(log.Status),
// // 		)
// // 	}
// // 	logrus.Info(bulk)

// // 	if len(bulk) == 0 {
// // 		return nil
// // 	}
// // 	return client.AttendanceLog.CreateBulk(bulk...).Exec(ctx)
// // }

// // // package data

// // // import (
// // // 	"context"
// // // 	"fmt"
// // // 	"time"

// // // 	"mall-go/module/log-downloader-service/internal/data/model"
// // // 	"mall-go/module/log-downloader-service/internal/data/model/attendancelog"
// // // 	"mall-go/module/log-downloader-service/internal/data/model/predicate"
// // // 	"mall-go/module/log-downloader-service/internal/zk"
// // // )

// // // type logKey struct {
// // // 	UserID   int64
// // // 	DeviceIP string
// // // 	Time     time.Time
// // // 	Status   int
// // // }

// // // func sendToDatabase(ctx context.Context, client *model.Client, logs []zk.LogData) error {
// // // 	// 1. Bangun key map
// // // 	keyMap := make(map[logKey]zk.LogData)
// // // 	var keys []logKey
// // // 	for _, log := range logs {
// // // 		k := logKey{UserID: log.UserID, DeviceIP: log.DeviceIP, Time: log.Attendace, Status: log.Status}
// // // 		keys = append(keys, k)
// // // 		keyMap[k] = log
// // // 	}

// // // 	// 2. Bangun predicates
// // // 	var predicates []predicate.AttendanceLog
// // // 	for _, k := range keys {
// // // 		predicates = append(predicates, attendancelog.And(
// // // 			attendancelog.UserID(k.UserID),
// // // 			attendancelog.DeviceIP(k.DeviceIP),
// // // 			attendancelog.AttLog(k.Time),
// // // 			attendancelog.Status(k.Status),
// // // 		))
// // // 	}

// // // 	// 3. Query existing logs
// // // 	existing, err := client.AttendanceLog.
// // // 		Query().
// // // 		Where(attendancelog.Or(predicates...)).
// // // 		All(ctx)
// // // 	if err != nil {
// // // 		return fmt.Errorf("query existing logs: %w", err)
// // // 	}

// // // 	// 4. Hapus yang sudah ada
// // // 	for _, e := range existing {
// // // 		k := logKey{UserID: e.UserID, DeviceIP: e.DeviceIP, Time: e.AttLog, Status: e.Status}
// // // 		delete(keyMap, k)
// // // 	}

// // // 	// 5. Build CreateBulk
// // // 	var bulk []*model.AttendanceLogCreate
// // // 	for _, log := range keyMap {
// // // 		bulk = append(bulk, client.AttendanceLog.
// // // 			Create().
// // // 			SetUserID(log.UserID).
// // // 			SetDeviceIP(log.DeviceIP).
// // // 			SetAttLog(log.Attendace).
// // // 			SetStatus(log.Status),
// // // 		)
// // // 	}

// // // 	// 6. Insert jika ada
// // // 	if len(bulk) == 0 {
// // // 		return nil
// // // 	}
// // // 	if err := client.AttendanceLog.CreateBulk(bulk...).Exec(ctx); err != nil {
// // // 		return fmt.Errorf("bulk insert failed: %w", err)
// // // 	}

// // // 	return nil
// // // }

// // // package data

// // // import (
// // // 	"context"
// // // 	"fmt"

// // // 	"mall-go/module/log-downloader-service/internal/data/model"
// // // 	"mall-go/module/log-downloader-service/internal/data/model/attendancelog"
// // // 	"mall-go/module/log-downloader-service/internal/zk"
// // // )

// // // func sendToDatabase(ctx context.Context, client *model.Client, logs []zk.LogData) error {
// // // 	for _, logEntry := range logs {
// // // 		exist, err := client.AttendanceLog.
// // // 			Query().
// // // 			Where(
// // // 				attendancelog.UserID(logEntry.UserID),
// // // 				attendancelog.DeviceIP(logEntry.DeviceIP),
// // // 				attendancelog.AttLog(logEntry.Attendace),
// // // 				attendancelog.Status(logEntry.Status),
// // // 			).
// // // 			Exist(ctx)

// // // 		if err != nil {
// // // 			return fmt.Errorf("error checking existing log: %w", err)
// // // 		}
// // // 		if exist {
// // // 			continue
// // // 		}

// // // 		_, err = client.AttendanceLog.
// // // 			Create().
// // // 			SetUserID(logEntry.UserID).
// // // 			SetDeviceIP(logEntry.DeviceIP).
// // // 			SetAttLog(logEntry.Attendace).
// // // 			SetStatus(logEntry.Status).
// // // 			Save(ctx)

// // // 		if err != nil {
// // // 			return fmt.Errorf("error inserting log: %w", err)
// // // 		}
// // // 	}
// // // 	return nil
// // // }

// // // import (
// // // 	"context"

// // // 	"mall-go/module/log-downloader-service/internal/conf"
// // // 	"mall-go/module/log-downloader-service/internal/data/model"
// // // 	"mall-go/module/log-downloader-service/internal/data/model/migrate"

// // // 	"github.com/go-kratos/kratos/v2/log"
// // // 	// _ "github.com/go-sql-driver/mysql"
// // // 	"github.com/google/wire"

// // // 	_ "github.com/lib/pq"
// // // 	"github.com/pkg/errors"
// // // )

// // // var ProviderSet = wire.NewSet(NewData, NewEntClient, sendToDatabase)

// // // // Data .
// // // type Data struct {
// // // 	db  *model.Client
// // // 	log *log.Helper
// // // 	// rdb *redis.Client
// // // }

// // // func NewEntClient(conf *conf.Data, logger log.Logger) *model.Client {
// // // 	l := log.NewHelper(logger)
// // // 	l.Infof("Attempting to connect to DB. Driver: %s, Source: %s", conf.Database.Driver, conf.Database.Source)
// // // 	client, err := model.Open(
// // // 		conf.Database.Driver,
// // // 		conf.Database.Source,
// // // 	)
// // // 	if err != nil {
// // // 		l.Fatalf("failed opening connection to db: %v", err)
// // // 	}
// // // 	client = client.Debug()
// // // 	// if err := client.Schema.Create(context.Background(), migrate.WithForeignKeys(false)); err != nil {
// // // 	// 	l.Fatalf("failed creating schema resources: %v", err)
// // // 	// }
// // // 	ctx := context.Background()
// // // 	// Run migration.
// // // 	err = client.Schema.Create(
// // // 		ctx,
// // // 		migrate.WithDropIndex(true),
// // // 		migrate.WithDropColumn(true),
// // // 	)
// // // 	if err != nil {
// // // 		log.Fatalf("failed creating schema resources: %v", err)
// // // 	}
// // // 	return client
// // // }

// // // // NewData .
// // // func NewData(entClient *model.Client, logger log.Logger) (*Data, func(), error) {
// // // 	cleanup := func() {
// // // 		log.NewHelper(logger).Info("closing the data resources")
// // // 	}
// // // 	return &Data{
// // // 		db:  entClient,
// // // 		log: log.NewHelper(logger),
// // // 	}, cleanup, nil
// // // }

// // // func WithTx(ctx context.Context, client *model.Client, fn func(tx *model.Tx) error) error {
// // // 	tx, err := client.Tx(ctx)
// // // 	if err != nil {
// // // 		return err
// // // 	}
// // // 	defer func() {
// // // 		if v := recover(); v != nil {
// // // 			tx.Rollback()
// // // 			panic(v)
// // // 		}
// // // 	}()
// // // 	if err := fn(tx); err != nil {
// // // 		if rerr := tx.Rollback(); rerr != nil {
// // // 			err = errors.Wrapf(err, "rolling back transaction: %v", rerr)
// // // 		}
// // // 		return err
// // // 	}
// // // 	if err := tx.Commit(); err != nil {
// // // 		return errors.Wrapf(err, "committing transaction: %v", err)
// // // 	}
// // // 	return nil
// // // }
