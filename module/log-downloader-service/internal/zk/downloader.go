package zk

import (
	"fmt"
	"mall-go/pkg/gozk"
	"time"
	// "github.com/canhlinh/gozk"
)

const (
	maxRetries = 3
	retryDelay = 2 * time.Second
)

type LogData struct {
	DeviceIP  string
	UserID    int64
	Attendace time.Time
	Status    int

	// State     int // aktifkan jika tersedia di event
	// Punch     int
}

// GetAllScannedEvents connects to a fingerprint device and returns log events
func GetAllScannedEvents(ip string, tcp bool, onError func(ip, severity, msg string)) []LogData {
	var zkConn *gozk.ZK
	var err error

	for i := 1; i <= maxRetries; i++ {
		zkConn = gozk.NewZK(ip, gozk.WithTCP(tcp), gozk.WithTimezone(gozk.DefaultTimezone))
		err = zkConn.Connect()
		if err == nil {
			break
		}
		msg := fmt.Sprintf("⚠️ Retry %d/%d: gagal koneksi ke %s: %v", i, maxRetries, ip, err)
		onError(ip, "warning", msg)
		time.Sleep(retryDelay)
	}

	if err != nil {
		msg := fmt.Sprintf("❌ Gagal konek ke device %s setelah %d percobaan: %v", ip, maxRetries, err)
		onError(ip, "error", msg)
		return nil
	}
	defer zkConn.Disconnect()

	properties, err := zkConn.GetProperties()
	if err != nil {
		onError(ip, "warning", fmt.Sprintf("⚠️ Gagal ambil properties dari %s: %v", ip, err))
	} else {
		fmt.Printf("[Device %s] Properties:\n", ip)
		properties.Println()
	}

	events, err := zkConn.GetAllScannedEvents()
	if err != nil {
		msg := fmt.Sprintf("❌ Gagal ambil scanned events dari %s: %v", ip, err)
		onError(ip, "error", msg)
		return nil
	}
	for i, e := range events {
		fmt.Printf("Raw Event #%d: %+v\n", i, e)
	}

	fmt.Printf("[Device %s] Jumlah event: %d\n", ip, len(events))
	var logs []LogData

	for _, event := range events {
		fmt.Printf("Event: %+v\n", event)
		logs = append(logs, LogData{
			DeviceIP:  ip,
			UserID:    event.UserID,
			Attendace: event.Timestamp,
			Status:    event.Status,
			// State:     event.State, // aktifkan jika ada
			// Punch:     event.Punch,
		})
	}

	return logs
}

// package zk

// import (
// 	"fmt"
// 	"time"

// 	"github.com/canhlinh/gozk"
// )

// type LogData struct {
// 	DeviceIP  string
// 	UserID    int64
// 	Timestamp time.Time
// 	State     int
// 	Punch     int
// }

// // GetAllScannedEvents connects to a fingerprint device and returns log events
// func GetAllScannedEvents(ip string, tcp bool, onError func(ip, msg string)) []LogData {
// 	zk := gozk.NewZK(ip, gozk.WithTCP(tcp), gozk.WithTimezone(gozk.DefaultTimezone))

// 	if err := zk.Connect(); err != nil {
// 		onError(ip, fmt.Sprintf("failed to connect: %v", err))
// 		return nil
// 	}
// 	defer zk.Disconnect()

// 	properties, err := zk.GetProperties()
// 	if err != nil {
// 		onError(ip, fmt.Sprintf("failed to get properties: %v", err))
// 		return nil
// 	}

// 	fmt.Printf("[Device %s] Properties:\n", ip)
// 	properties.Println()

// 	events, err := zk.GetAllScannedEvents()
// 	if err != nil {
// 		onError(ip, fmt.Sprintf("failed to get scanned events: %v", err))
// 		return nil
// 	}

// 	fmt.Printf("[Device %s] Number of events: %d\n", ip, len(events))
// 	var logs []LogData

// 	for _, event := range events {
// 		fmt.Printf("Event: %+v\n", event)
// 		logs = append(logs, LogData{
// 			DeviceIP:  ip,
// 			UserID:    event.UserID,
// 			Timestamp: event.Timestamp,
// 			// State:     event.State,
// 			// Punch:     event.Punch,
// 		})
// 	}

// 	return logs
// }

// // func GetAllScannedEvents(ip string, tcp bool, onError func(ip, msg string)) []LogData {
// // 	zk := gozk.NewZK(ip, gozk.WithTCP(tcp), gozk.WithTimezone(gozk.DefaultTimezone))

// // 	if err := zk.Connect(); err != nil {
// // 		return fmt.Errorf("failed to connect to device %s: %w", ip, err)
// // 	}
// // 	defer zk.Disconnect()

// // 	properties, err := zk.GetProperties()
// // 	if err != nil {
// // 		return fmt.Errorf("failed to get properties from %s: %w", ip, err)
// // 	}

// // 	fmt.Printf("[Device %s] Properties:\n", ip)
// // 	properties.Println()

// // 	events, err := zk.GetAllScannedEvents()
// // 	if err != nil {
// // 		return fmt.Errorf("failed to get scanned events from %s: %w", ip, err)
// // 	}

// // 	fmt.Printf("[Device %s] Number of events: %d\n", ip, len(events))
// // 	for _, event := range events {
// // 		fmt.Printf("Event: %+v\n ip: %s\n ", event, ip)
// // 	}

// // 	return nil
// // }
