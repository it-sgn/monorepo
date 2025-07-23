package db

// import (
// 	"context"
// 	"fmt"

// 	"mall-go/module/log-downloader-service/internal/db/model"
// 	"mall-go/module/log-downloader-service/internal/db/model/attendancelog"
// 	"mall-go/module/log-downloader-service/internal/utils"
// 	"mall-go/module/log-downloader-service/internal/zk"
// )

// func SendToDatabase(ctx context.Context, client *model.Client, logs []zk.LogData) error {
// 	keyMap, keys := utils.BuildKeyMap(logs)
// 	predicates := utils.BuildPredicates(keys)

// 	existing, err := client.AttendanceLog.
// 		Query().
// 		Where(attendancelog.Or(predicates...)).
// 		All(ctx)
// 	if err != nil {
// 		return fmt.Errorf("query existing logs: %w", err)
// 	}

// 	for _, e := range existing {
// 		k := utils.LogKey{UserID: e.UserID, DeviceIP: e.DeviceIP, Time: e.AttLog, Status: e.Status}
// 		delete(keyMap, k)
// 	}

// 	var bulk []*model.AttendanceLogCreate
// 	for _, log := range keyMap {
// 		bulk = append(bulk, client.AttendanceLog.
// 			Create().
// 			SetUserID(log.UserID).
// 			SetDeviceIP(log.DeviceIP).
// 			SetAttLog(log.Attendace).
// 			SetStatus(log.Status),
// 		)
// 	}

// 	if len(bulk) == 0 {
// 		return nil
// 	}
// 	return client.AttendanceLog.CreateBulk(bulk...).Exec(ctx)
// }

// // package data

// // import (
// // 	"context"
// // 	"fmt"
// // 	"time"

// // 	"mall-go/module/log-downloader-service/internal/data/model"
// // 	"mall-go/module/log-downloader-service/internal/data/model/attendancelog"
// // 	"mall-go/module/log-downloader-service/internal/data/model/predicate"
// // 	"mall-go/module/log-downloader-service/internal/zk"
// // )

// // type logKey struct {
// // 	UserID   int64
// // 	DeviceIP string
// // 	Time     time.Time
// // 	Status   int
// // }

// // func sendToDatabase(ctx context.Context, client *model.Client, logs []zk.LogData) error {
// // 	// 1. Bangun key map
// // 	keyMap := make(map[logKey]zk.LogData)
// // 	var keys []logKey
// // 	for _, log := range logs {
// // 		k := logKey{UserID: log.UserID, DeviceIP: log.DeviceIP, Time: log.Attendace, Status: log.Status}
// // 		keys = append(keys, k)
// // 		keyMap[k] = log
// // 	}

// // 	// 2. Bangun predicates
// // 	var predicates []predicate.AttendanceLog
// // 	for _, k := range keys {
// // 		predicates = append(predicates, attendancelog.And(
// // 			attendancelog.UserID(k.UserID),
// // 			attendancelog.DeviceIP(k.DeviceIP),
// // 			attendancelog.AttLog(k.Time),
// // 			attendancelog.Status(k.Status),
// // 		))
// // 	}

// // 	// 3. Query existing logs
// // 	existing, err := client.AttendanceLog.
// // 		Query().
// // 		Where(attendancelog.Or(predicates...)).
// // 		All(ctx)
// // 	if err != nil {
// // 		return fmt.Errorf("query existing logs: %w", err)
// // 	}

// // 	// 4. Hapus yang sudah ada
// // 	for _, e := range existing {
// // 		k := logKey{UserID: e.UserID, DeviceIP: e.DeviceIP, Time: e.AttLog, Status: e.Status}
// // 		delete(keyMap, k)
// // 	}

// // 	// 5. Build CreateBulk
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

// // 	// 6. Insert jika ada
// // 	if len(bulk) == 0 {
// // 		return nil
// // 	}
// // 	if err := client.AttendanceLog.CreateBulk(bulk...).Exec(ctx); err != nil {
// // 		return fmt.Errorf("bulk insert failed: %w", err)
// // 	}

// // 	return nil
// // }

// // package data

// // import (
// // 	"context"
// // 	"fmt"

// // 	"mall-go/module/log-downloader-service/internal/data/model"
// // 	"mall-go/module/log-downloader-service/internal/data/model/attendancelog"
// // 	"mall-go/module/log-downloader-service/internal/zk"
// // )

// // func sendToDatabase(ctx context.Context, client *model.Client, logs []zk.LogData) error {
// // 	for _, logEntry := range logs {
// // 		exist, err := client.AttendanceLog.
// // 			Query().
// // 			Where(
// // 				attendancelog.UserID(logEntry.UserID),
// // 				attendancelog.DeviceIP(logEntry.DeviceIP),
// // 				attendancelog.AttLog(logEntry.Attendace),
// // 				attendancelog.Status(logEntry.Status),
// // 			).
// // 			Exist(ctx)

// // 		if err != nil {
// // 			return fmt.Errorf("error checking existing log: %w", err)
// // 		}
// // 		if exist {
// // 			continue
// // 		}

// // 		_, err = client.AttendanceLog.
// // 			Create().
// // 			SetUserID(logEntry.UserID).
// // 			SetDeviceIP(logEntry.DeviceIP).
// // 			SetAttLog(logEntry.Attendace).
// // 			SetStatus(logEntry.Status).
// // 			Save(ctx)

// // 		if err != nil {
// // 			return fmt.Errorf("error inserting log: %w", err)
// // 		}
// // 	}
// // 	return nil
// // }
