package job111

import (
	"fmt"
	"log"
	"mall-go/module/log_downloader/service/internal/biz"
)

// type JobFunc func()

type DownloadJob struct {
	uc *biz.DownloadUseCase
}

func NewDownloadJob(uc *biz.DownloadUseCase) *DownloadJob {
	return &DownloadJob{uc: uc}
}

func (s *DownloadJob) Init() {
	defaultJobs := map[string]JobFunc{
		"device_1": s.Device1,
		"device_2": s.Device2,
		"device_3": s.Device3,
	}

	for name := range defaultJobs {
		log.Printf("[INFO] Job registered: %s", name)
	}
}

func (s *DownloadJob) Device1() {
	// ctx := context.Background()
	ip := "192.168.80.26"

	// Ambil log valid (via TCP)
	logsValid, err := s.uc.GetAllScannedEvents(ip, true, func(ip, severity, msg string) {
		log.Printf("[%s] %s - %s", severity, ip, msg)
	})
	if err != nil {
		log.Printf("[ERROR] Gagal ambil log valid dari %s: %v", ip, err)
	}

	// Ambil log invalid (via UDP)
	logsInvalid, err := s.uc.GetAllScannedEvents(ip, false, func(ip, severity, msg string) {
		log.Printf("[%s] %s - %s", severity, ip, msg)
	})
	if err != nil {
		log.Printf("[ERROR] Gagal ambil log invalid dari %s: %v", ip, err)
	}

	processLogs := func(logs []biz.LogData, label string) {
		if len(logs) == 0 {
			log.Printf("[INFO] Tidak ada log %s dari %s", label, ip)
			return
		}
		for _, logData := range logs {
			log.Printf("[INFO][%s] Log %s: %+v", ip, label, logData)
		}
	}

	processLogs(logsValid, "valid")
	processLogs(logsInvalid, "invalid")
}

func (s *DownloadJob) Device2() {
	// // err := s.uc.Create(context.Background(), &biz.CronZK{})
	// if err != nil {
	// 	log.Printf("[ERROR] Device2 job gagal: %v", err)
	// }
	fmt.Println("[Device2] Job selesai.")
}

func (s *DownloadJob) Device3() {
	// err := s.uc.Create(context.Background(), &biz.CronZK{})
	// if err != nil {
	// 	log.Printf("[ERROR] Device3 job gagal: %v", err)
	// }
	fmt.Println("[Device3] Job selesai.")
}
