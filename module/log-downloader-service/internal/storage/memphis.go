package storage

import (
	"encoding/json"
	"log"
	"time"

	"github.com/memphisdev/memphis.go"
)

var memphisConn *memphis.Conn
var producer *memphis.Producer

func InitMemphis() {
	var err error
	memphisConn, err = memphis.Connect(
		"localhost",                 // host
		"root",                      // username
		memphis.Password("memphis"), // password
	)
	if err != nil {
		log.Fatal("Failed to connect to Memphis: ", err)
	}

	producer, err = memphisConn.CreateProducer("attendance.log", "log-downloader-producer")
	if err != nil {
		log.Fatal("Failed to create producer: ", err)
	}
}

// DeviceIP  string
// UserID    int64
// Attendace time.Time
// Status    int
type AttLog struct {
	DeviceIP  string    `json:"device_ip"`
	UserID    int64     `json:"user_id"`
	Attendace time.Time `json:"attendace"`
	Status    int       `json:"status"` // 0 = clock_in, 1 = clock_out
}

func PublishLog(logData AttLog) error {
	payload, err := json.Marshal(logData)
	if err != nil {
		return err
	}

	return producer.Produce(payload)
}
