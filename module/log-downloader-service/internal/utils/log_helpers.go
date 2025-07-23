package utils

import (
	"time"

	"mall-go/module/log-downloader-service/internal/db/model/attendancelog"
	"mall-go/module/log-downloader-service/internal/db/model/predicate"
	"mall-go/module/log-downloader-service/internal/zk"
)

type LogKey struct {
	UserID   int64
	DeviceIP string
	Time     time.Time
	Status   int
}

func BuildKeyMap(logs []zk.LogData) (map[LogKey]zk.LogData, []LogKey) {
	keyMap := make(map[LogKey]zk.LogData)
	var keys []LogKey

	for _, log := range logs {
		k := LogKey{
			UserID:   log.UserID,
			DeviceIP: log.DeviceIP,
			Time:     log.Attendace,
			Status:   log.Status,
		}
		keys = append(keys, k)
		keyMap[k] = log
	}

	return keyMap, keys
}

func BuildPredicates(keys []LogKey) []predicate.AttendanceLog {
	var preds []predicate.AttendanceLog
	for _, k := range keys {
		preds = append(preds, attendancelog.And(
			attendancelog.UserID(k.UserID),
			attendancelog.DeviceIP(k.DeviceIP),
			attendancelog.AttLog(k.Time),
			attendancelog.Status(k.Status),
		))

	}
	return preds
}
