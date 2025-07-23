package main

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"mall-go/module/log-downloader-service/internal/schedule"
	"mall-go/module/log-downloader-service/internal/storage"
	"mall-go/module/log-downloader-service/internal/zk"

	_ "mall-go/module/log-downloader-service/internal/conf"
)

type Config struct {
	Server struct {
		HTTP struct {
			Addr    string        `yaml:"addr"`
			Timeout time.Duration `yaml:"timeout"`
		} `yaml:"http"`
		GRPC struct {
			Addr    string        `yaml:"addr"`
			Timeout time.Duration `yaml:"timeout"`
		} `yaml:"grpc"`
	} `yaml:"server"`

	Devices []struct {
		IP string `yaml:"ip"`
	} `yaml:"devices"`

	Data struct {
		Database struct {
			Driver string `yaml:"driver"`
			Source string `yaml:"source"`
		} `yaml:"database"`
		Redis struct {
			Addr         string        `yaml:"addr"`
			ReadTimeout  time.Duration `yaml:"read_timeout"`
			WriteTimeout time.Duration `yaml:"write_timeout"`
		} `yaml:"redis"`
	} `yaml:"data"`
}

func main() {
	var cfg Config

	// 📖 Load config file
	data, err := os.ReadFile("../../configs/config.yaml")
	if err != nil {
		log.Fatalf("❌ Gagal membaca config.yaml: %v", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("❌ Gagal parsing YAML: %v", err)
	}
	// init memphis
	storage.InitMemphis()
	// 🔌 Inisialisasi DB
	if err := storage.InitDB(cfg.Data.Database.Source); err != nil {
		log.Fatalf("❌ Gagal konek ke DB: %v", err)
	}
	// kafka.InitKafka([]string{"localhost:9092"}, "attendance.log.success")
	// 📡 Ambil IP dari config
	var ips []string
	for _, d := range cfg.Devices {
		ips = append(ips, d.IP)
	}

	// 🕒 Start scheduler
	schedule.Start(ips, func(ip string) error {
		logsValid := zk.GetAllScannedEvents(ip, true, func(ip, severity, msg string) {
			log.Printf("[%s] %s - %s", severity, ip, msg)
			storage.SendToKafka("attendance.log.error", ip, severity, msg)
		})

		logsInvalid := zk.GetAllScannedEvents(ip, false, func(ip, severity, msg string) {
			log.Printf("[%s] %s - %s", severity, ip, msg)
			storage.SendToKafka("attendance.log.error", ip, severity, msg)
		})

		processLogs := func(logs []zk.LogData, label string) {
			if len(logs) == 0 {
				log.Printf("[INFO] Tidak ada log %s dari %s", label, ip)
				return
			}
			for _, logData := range logs {
				log.Printf("[INFO][%s] Log %s: %+v", ip, label, logData)
				// kafka.SendToKafka(logData)
				err := storage.PublishLog(storage.AttLog{
					UserID:    logData.UserID,
					DeviceIP:  logData.DeviceIP,
					Attendace: logData.Attendace,
					Status:    logData.Status,
				})
				if err != nil {
					return
				}
			}
			if err := storage.SendToDatabase(logs); err != nil {
				log.Printf("[ERROR] Gagal insert DB log %s: %v", label, err)
				// storage.SendToKafka("attendance.log.error", ip, "error", err.Error())
				// storage.PublishLog()
			}
		}

		processLogs(logsValid, "valid")
		processLogs(logsInvalid, "invalid")

		return nil
	})
}

// package main

// import (
// 	"log"
// 	"mall-go/module/log-downloader-service/internal/kafka"
// 	"mall-go/module/log-downloader-service/internal/schedule"
// 	"mall-go/module/log-downloader-service/internal/storage"
// 	"mall-go/module/log-downloader-service/internal/zk"
// 	"os"
// 	"time"

// 	_ "mall-go/module/log-downloader-service/internal/conf"

// 	"gopkg.in/yaml.v3"
// 	// _ "mall-go/module/log-downloader-service/internal/conf"
// )

// type Device struct {
// 	IP string `yaml:"ip"`
// }

// // type Config struct {
// // 	Devices []Device `yaml:"devices"`
// // }

// type Config struct {
// 	Server struct {
// 		HTTP struct {
// 			Addr    string        `yaml:"addr"`
// 			Timeout time.Duration `yaml:"timeout"`
// 		} `yaml:"http"`
// 		GRPC struct {
// 			Addr    string        `yaml:"addr"`
// 			Timeout time.Duration `yaml:"timeout"`
// 		} `yaml:"grpc"`
// 	} `yaml:"server"`

// 	Devices []struct {
// 		IP string `yaml:"ip"`
// 	} `yaml:"devices"`

// 	Data struct {
// 		Database struct {
// 			Driver string `yaml:"driver"`
// 			Source string `yaml:"source"`
// 		} `yaml:"database"`
// 		Redis struct {
// 			Addr         string        `yaml:"addr"`
// 			ReadTimeout  time.Duration `yaml:"read_timeout"`
// 			WriteTimeout time.Duration `yaml:"write_timeout"`
// 		} `yaml:"redis"`
// 	} `yaml:"data"`
// }

// func main() {
// 	var cfg Config

// 	data, err := os.ReadFile("../../configs/config.yaml")
// 	if err != nil {
// 		log.Fatalf("❌ gagal membaca config.yaml: %v", err)
// 	}
// 	if err := yaml.Unmarshal(data, &cfg); err != nil {
// 		log.Fatalf("❌ gagal parsing YAML: %v", err)
// 	}

// 	// 🔌 Init koneksi DB
// 	if err := storage.InitDB(cfg.Data.Database.Source); err != nil {
// 		log.Fatalf("❌ gagal konek ke DB: %v", err)
// 	}

// 	// Ambil list IP
// 	var ips []string
// 	for _, d := range cfg.Devices {
// 		ips = append(ips, d.IP)
// 	}

// 	// ⏱ Start scheduler
// 	schedule.Start(ips, func(ip string) error {

// 		logs := zk.GetAllScannedEvents(ip, true, func(ip, severity, msg string) {
// 			log.Printf("[%s] %s - %s", severity, ip, msg)
// 			storage.SendToKafka("attendance.log.error", ip, severity, msg)
// 		})
// 		logsx := zk.GetAllScannedEvents(ip, false, func(ip, severity, msg string) {
// 			log.Printf("[%s] %s - %s", severity, ip, msg)
// 			storage.SendToKafka("attendance.log.error", ip, severity, msg)
// 		})

// 		if len(logs) == 0 {
// 			log.Printf("[INFO] Tidak ada log dari %s", ip)
// 			return nil
// 		}
// 		if len(logsx) == 0 {
// 			log.Printf("[INFO] Tidak ada log dari %s", ip)
// 			return nil
// 		}
// 		// // 		// Kirim ke Kafka
// 		for _, logData := range logs {
// 			log.Printf("[INFO][%s] LogData: %+v", ip, logData)
// 			kafka.SendToKafka(logData)
// 		}
// 		for _, logData := range logsx {
// 			log.Printf("[INFO][%s] LogData: %+v", ip, logData)
// 			kafka.SendToKafka(logData)
// 		}
// 		if err := storage.SendToDatabase(logs); err != nil {
// 			log.Printf("[ERROR] Gagal insert DB: %v", err)
// 			storage.SendToKafka("attendance.log.error", ip, "error", err.Error())
// 		}
// 		if err := storage.SendToDatabase(logsx); err != nil {
// 			log.Printf("[ERROR] Gagal insert DB: %v", err)
// 			storage.SendToKafka("attendance.log.error", ip, "error", err.Error())
// 		}

// 		return nil
// 	})
// }

// package main

// import (
// 	"log"
// 	"mall-go/module/log-downloader-service/internal/schedule"
// 	"mall-go/module/log-downloader-service/internal/storage"
// 	"mall-go/module/log-downloader-service/internal/zk"
// 	"os"

// 	"gopkg.in/yaml.v3"
// )

// type Device struct {
// 	IP string `yaml:"ip"`
// }

// type Config struct {
// 	Devices []Device `yaml:"devices"`
// }

// func main() {
// 	cfg := &Config{}
// 	data, err := os.ReadFile("../../configs/config.yaml")
// 	if err != nil {
// 		log.Fatalf("❌ config read error: %v", err)
// 	}
// 	if err := yaml.Unmarshal(data, cfg); err != nil {
// 		log.Fatalf("❌ config parse error: %v", err)
// 	}

// 	// 🔌 Init PostgreSQL
// 	dbDSN := os.Getenv("DATABASE_DSN")
// 	if dbDSN == "" {
// 		log.Fatal("❌ DATABASE_DSN env not set")
// 	}
// 	if err := storage.InitDB(dbDSN); err != nil {
// 		log.Fatalf("❌ gagal init DB: %v", err)
// 	}

// 	ips := []string{}
// 	for _, d := range cfg.Devices {
// 		ips = append(ips, d.IP)
// 	}

// 	// 🔁 Start scheduler
// 	schedule.Start(ips, func(ip string) error {
// 		logs := zk.GetAllScannedEvents(ip, true, func(ip, severity, msg string) {
// 			log.Printf("[%s] %s - %s", severity, ip, msg)
// 			storage.SendToKafka("attendance.log.error", ip, severity, msg)
// 		})

// 		if len(logs) == 0 {
// 			log.Printf("[INFO] Tidak ada log dari %s", ip)
// 			return nil
// 		}

// 		log.Printf("[INFO] Dapat %d log dari %s", len(logs), ip)

// 		// 💾 Simpan ke database
// 		if err := storage.SendToDatabase(logs); err != nil {
// 			log.Printf("[ERROR] Gagal simpan DB: %v", err)
// 			storage.SendToKafka("attendance.log.error", ip, "error", err.Error())
// 		}

// 		// (Opsional) Kirim ke Kafka juga
// 		// storage.SendToKafka("attendance.log.processed", ip, "info", fmt.Sprintf("%d log disimpan", len(logs)))

// 		return nil
// 	})
// }

// // package main

// // import (
// // 	"log"

// // 	"mall-go/module/log-downloader-service/internal/kafka"
// // 	"mall-go/module/log-downloader-service/internal/schedule"
// // 	"mall-go/module/log-downloader-service/internal/storage"
// // 	"mall-go/module/log-downloader-service/internal/zk"
// // 	"os"

// // 	_ "github.com/lib/pq"
// // 	"gopkg.in/yaml.v2"
// // )

// // type Device struct {
// // 	IP string `yaml:"ip"`
// // }

// // type Config struct {
// // 	Devices []Device `yaml:"devices"`
// // }

// // func main() {
// // 	cfg := &Config{}
// // 	data, err := os.ReadFile("../../configs/config.yaml")
// // 	if err != nil {
// // 		log.Fatalf("❌ Gagal baca config.yaml: %v", err)
// // 	}
// // 	if err := yaml.Unmarshal(data, cfg); err != nil {
// // 		log.Fatalf("❌ Gagal parse YAML config: %v", err)
// // 	}
// // 	var ips []string
// // 	for _, d := range cfg.Devices {
// // 		ips = append(ips, d.IP)
// // 	}
// // 	kafka.InitKafka([]string{"localhost:9092"}, "attendance.log.success")

// // 	// client, err := model.Open("postgres", "host=127.0.0.1 user=admin password=admin123 dbname=attendance port=5432 sslmode=disable TimeZone=Asia/Jakarta")
// // 	// if err != nil {
// // 	// 	log.Fatalf("failed to connect to database: %v", err)
// // 	// }
// // 	// defer client.Close()

// // 	// ctx := context.Background()
// // 	// if err := client.Schema.Create(ctx); err != nil {
// // 	// 	log.Fatalf("❌ Gagal migrate schema: %v", err)
// // 	// }
// // 	const pgDSN = "postgres://admin:admin123@localhost:5432/attendance?sslmode=disable"
// // 	if err := storage.InitDB(pgDSN); err != nil {
// // 		log.Fatalf("❌ Failed to connect to DB: %v", err)
// // 	}

// // 	// ... bagian pemanggilan
// // 	// Start scheduler dan proses log dari tiap device
// // 	schedule.Start(ips, func(ip string) error {
// // 		// Callback error untuk tiap IP
// // 		onError := func(ip, severity, msg string) {
// // 			log.Printf("[%s][%s] %s", severity, ip, msg)
// // 			// Kafka jika perlu
// // 			// kafka.SendErrorKafka(ip, severity, msg)
// // 		}

// // 		// Ambil logs dari device
// // 		allLogs := append(
// // 			zk.GetAllScannedEvents(ip, false, onError),
// // 			zk.GetAllScannedEvents(ip, true, onError)...,
// // 		)

// // 		if len(allLogs) == 0 {
// // 			log.Printf("[INFO][%s] Tidak ada log baru", ip)
// // 			return nil
// // 		}

// // 		// Kirim ke Kafka
// // 		// for _, logData := range allLogs {
// // 		// 	log.Printf("[INFO][%s] LogData: %+v", ip, logData)
// // 		// 	kafka.SendToKafka(logData)
// // 		// }

// // 		// Simpan ke database
// // 		for _, logData := range allLogs {
// // 			log.Printf("[INFO] Log from %s: %+v", ip, logData)

// // 			if err := storage.SendToDatabase(logData); err != nil {
// // 				log.Printf("❌ DB Insert Error: %v", err)
// // 				storage.SendToKafka("attendance.log.error", ip, "error", err.Error())
// // 			}
// // 		}
// // 		return nil
// // 	})
// // }

// // // package main

// // // import (
// // // 	"log"
// // // 	"mall-go/module/log-downloader-service/internal/kafka"
// // // 	"mall-go/module/log-downloader-service/internal/schedule"
// // // 	"mall-go/module/log-downloader-service/internal/zk"
// // // 	"os"

// // // 	"gopkg.in/yaml.v3"
// // // )

// // // type Device struct {
// // // 	IP string `yaml:"ip"`
// // // }

// // // type Config struct {
// // // 	Devices []Device `yaml:"devices"`
// // // }

// // // func main() {
// // // 	cfg := &Config{}
// // // 	data, err := os.ReadFile("../../configs/config.yaml")
// // // 	if err != nil {
// // // 		log.Fatalf("❌ Gagal baca config.yaml: %v", err)
// // // 	}

// // // 	if err := yaml.Unmarshal(data, cfg); err != nil {
// // // 		log.Fatalf("❌ Gagal parse YAML config: %v", err)
// // // 	}

// // // 	var ips []string
// // // 	for _, d := range cfg.Devices {
// // // 		ips = append(ips, d.IP)
// // // 	}

// // // 	kafka.InitKafka([]string{"localhost:9092"}, "attendance.log.success")

// // // 	// Start scheduler dan proses log dari tiap device
// // // 	schedule.Start(ips, func(ip string) error {
// // // 		// Callback error untuk tiap IP
// // // 		onError := func(ip, severity, msg string) {
// // // 			switch severity {
// // // 			case "warning":
// // // 				log.Printf("[WARNING][%s] %s", ip, msg)
// // // 			case "error":
// // // 				log.Printf("[ERROR][%s] %s", ip, msg)
// // // 			default:
// // // 				log.Printf("[INFO][%s] %s", ip, msg)
// // // 			}
// // // 			// TODO: Kirim juga ke Kafka kalau perlu
// // // 			// sendToKafka("attendance.log.error", ip, severity, msg)
// // // 		}

// // // 		// Ambil logs dari device
// // // 		logs := zk.GetAllScannedEvents(ip, false, onError)
// // // 		logs2 := zk.GetAllScannedEvents(ip, true, onError)
// // // 		for _, logData := range logs2 {
// // // 			log.Printf("[INFO][%s] LogData: %+v", ip, logData)
// // // 			// kafka.sendToKafka(logData)
// // // 			kafka.SendToKafka(logData)
// // // 			// TODO: Simpan atau kirim logData ke sistem lain
// // // 			// sendToDatabase(logData)
// // // 		}

// // // 		// Proses log (kirim ke gRPC / Kafka / DB)
// // // 		for _, logData := range logs {
// // // 			log.Printf("[INFO][%s] LogData: %+v", ip, logData)
// // // 			// kafka.sendToKafka(logData)
// // // 			kafka.SendToKafka(logData)
// // // 			// TODO: Simpan atau kirim logData ke sistem lain
// // // 			// sendToDatabase(logData)
// // // 		}

// // // 		return nil
// // // 	})
// // // }

// // // // package main

// // // // import (
// // // // 	"log"
// // // // 	"mall-go/module/log-downloader-service/internal/schedule"
// // // // 	"mall-go/module/log-downloader-service/internal/zk"
// // // // 	"os"

// // // // 	"gopkg.in/yaml.v3"
// // // // )

// // // // type Device struct {
// // // // 	IP string `yaml:"ip"`
// // // // }
// // // // type Config struct {
// // // // 	Devices []Device `yaml:"devices"`
// // // // }

// // // // func main() {
// // // // 	cfg := &Config{}
// // // // 	data, err := os.ReadFile("../../configs/config.yaml")
// // // // 	if err != nil {
// // // // 		log.Fatalf("config read error: %v", err)
// // // // 	}
// // // // 	yaml.Unmarshal(data, cfg)

// // // // 	ips := []string{}
// // // // 	for _, d := range cfg.Devices {
// // // // 		ips = append(ips, d.IP)
// // // // 	}

// // // // 	schedule.Start(ips, func(ip string) error {
// // // // 		logs := zk.GetAllScannedEvents(ip, false, func(ip, msg string) {
// // // // 			log.Printf("[ERROR] device %s: %s", ip, msg)
// // // // 		})

// // // // 		// Lakukan sesuatu dengan logs (kirim ke gRPC/Kafka/db)
// // // // 		for _, logData := range logs {
// // // // 			log.Printf("[INFO] Log from %s: %+v", ip, logData)
// // // // 		}
// // // // 		return nil
// // // // 	})
// // // // }

// // // // package main

// // // // import (
// // // // 	"log"
// // // // 	"mall-go/module/log-downloader-service/internal/schedule"
// // // // 	"mall-go/module/log-downloader-service/internal/zk"
// // // // 	"os"

// // // // 	"gopkg.in/yaml.v3"
// // // // )

// // // // type Device struct {
// // // // 	IP string `yaml:"ip"`
// // // // }
// // // // type Config struct {
// // // // 	Devices []Device `yaml:"devices"`
// // // // }

// // // // func main() {
// // // // 	cfg := &Config{}
// // // // 	data, err := os.ReadFile("../../configs/config.yaml")
// // // // 	if err != nil {
// // // // 		log.Fatalf("config read error: %v", err)
// // // // 	}
// // // // 	yaml.Unmarshal(data, cfg)

// // // // 	ips := []string{}
// // // // 	for _, d := range cfg.Devices {
// // // // 		ips = append(ips, d.IP)
// // // // 	}

// // // // 	schedule.Start(ips, func(ip string) error {
// // // // 		return zk.GetAllScannedEvents(ip, false)
// // // // 	})
// // // // }
