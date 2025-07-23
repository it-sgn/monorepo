package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"mall-go/module/log-downloader-service/internal/zk"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(dsn string) error {
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}
	log.Println("✅ Connected to PostgreSQL")
	return nil
}

func SendToDatabase(logs []zk.LogData) error {
	if db == nil {
		return fmt.Errorf("DB not initialized")
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	stmt, err := tx.Prepare(`
        INSERT INTO attendance_log (device_ip, user_id, att_log, status, created_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (device_ip, user_id, att_log) DO NOTHING
    `)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	for _, logData := range logs {
		_, err := stmt.Exec(
			logData.DeviceIP,
			logData.UserID,
			logData.Attendace,
			logData.Status,
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("exec: %w", err)
		}
	}

	return nil
}

func SendToKafka(topic string, ip string, severity string, message string) {
	// Dummy: replace with real Kafka
	log.Printf("[Kafka][%s] %s [%s]: %s", topic, ip, severity, message)
}

// package storage

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"time"

// 	"mall-go/module/log-downloader-service/internal/zk"

// 	_ "github.com/lib/pq"
// )

// var db *sql.DB

// // InitDB opens DB connection (call this once from main)
// func InitDB(dsn string) error {
// 	var err error
// 	db, err = sql.Open("postgres", dsn)
// 	if err != nil {
// 		return fmt.Errorf("failed to open DB: %w", err)
// 	}

// 	if err := db.Ping(); err != nil {
// 		return fmt.Errorf("failed to ping DB: %w", err)
// 	}

// 	log.Println("✅ Connected to PostgreSQL")
// 	return nil
// }

// // SendToDatabase inserts attendance log to database
// func SendToDatabase(data zk.LogData) error {
// 	if db == nil {
// 		return fmt.Errorf("DB not initialized")
// 	}

// 	query := `
//         INSERT INTO attendance_log (device_ip, user_id, att_log, status, created_at)
//         VALUES ($1, $2, $3, $4, $5)
//         ON CONFLICT DO NOTHING
//     `

// 	_, err := db.Exec(query,
// 		data.DeviceIP,
// 		data.UserID,
// 		data.Attendace,
// 		data.Status,
// 		// data.Punch,
// 		time.Now(),
// 	)

// 	if err != nil {
// 		return fmt.Errorf("failed insert: %w", err)
// 	}

// 	return nil
// }

// // SendToKafka sends log message to Kafka
// func SendToKafka(topic string, ip string, severity string, message string) {
// 	// Dummy Kafka log
// 	log.Printf("[Kafka][%s] %s [%s]: %s", topic, ip, severity, message)
// 	// TODO: Replace with actual Kafka producer logic
// }
