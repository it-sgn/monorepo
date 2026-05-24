package service

import (
	"log"
	"mall-go/module/log_downloader/service/internal/biz"
)

var DefaultJobs map[string]JobFunc

type JobFunc func()

type DownloadJobService struct {
	uc *biz.DownloadUseCase
}

func NewDownloadJob(uc *biz.DownloadUseCase) *DownloadJobService {
	job := &DownloadJobService{
		uc: uc,
	}
	return job
}

func (s *DownloadJobService) Init() {
	DefaultJobs = map[string]JobFunc{
		"device_1_5min":  s.Device1,
		"device_1_3hour": s.Device1,
		"device_2_5min":  s.Device2,
		"device_2_3hour": s.Device2,
		"device_3_5min":  s.Device3,
		"device_3_3hour": s.Device3,
		"device_4_5min":  s.Device4,
		"device_4_3hour": s.Device4,
		"device_5_5min":  s.Device5,
		"device_5_3hour": s.Device5,
		// "device_1":       s.Device1,
		// "device_2":       s.Device2,
		// "device_3":       s.Device3,
		// "device_4":       s.Device4,
		// "device_5":       s.Device5,
	}
}

func (s *DownloadJobService) Device1() {
	// ctx := context.Background()
	ips := []string{"192.168.80.26"}
	// ips := []string{"192.168.80.26"}

	for _, ip := range ips {
		// Ambil log valid (via TCP)
		logsValid, err := s.uc.GetAllScannedEvents(ip, true, func(ip, severity, msg string) {
			log.Printf("[%s] %s - %s", severity, ip, msg)
		})
		if err != nil {
			log.Printf("[ERROR] Gagal ambil log valid dari %s: %v", ip, err)
			continue
		}

		// Ambil log invalid (via UDP)
		logsInvalid, err := s.uc.GetAllScannedEvents(ip, false, func(ip, severity, msg string) {
			log.Printf("[%s] %s - %s", severity, ip, msg)
		})
		if err != nil {
			log.Printf("[ERROR] Gagal ambil log invalid dari %s: %v", ip, err)
			continue
		}

		// Cetak hasil log
		processLogs := func(logs []biz.LogData, label string) {
			if len(logs) == 0 {
				log.Printf("[INFO] Tidak ada log %s dari %s", label, ip)
				return
			}
			for _, logData := range logs {
				log.Printf("[INFO][%s] Log %s: %+v", ip, label, logData)
			}
		}

		processLogs(logsValid, "tcp")
		processLogs(logsInvalid, "udp")
	}
}

func (s *DownloadJobService) Device2() {
	// ctx := context.Background()
	ips := []string{"192.168.80.27"}
	// ips := []string{"192.168.80.26"}

	for _, ip := range ips {
		// Ambil log valid (via TCP)
		logsValid, err := s.uc.GetAllScannedEvents(ip, true, func(ip, severity, msg string) {
			log.Printf("[%s] %s - %s", severity, ip, msg)
		})
		if err != nil {
			log.Printf("[ERROR] Gagal ambil log valid dari %s: %v", ip, err)
			continue
		}

		// Ambil log invalid (via UDP)
		logsInvalid, err := s.uc.GetAllScannedEvents(ip, false, func(ip, severity, msg string) {
			log.Printf("[%s] %s - %s", severity, ip, msg)
		})
		if err != nil {
			log.Printf("[ERROR] Gagal ambil log invalid dari %s: %v", ip, err)
			continue
		}

		// Cetak hasil log
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
}

func (s *DownloadJobService) Device3() {
	// ctx := context.Background()
	ips := []string{"192.168.80.28"}
	// ips := []string{"192.168.80.26"}

	for _, ip := range ips {
		// Ambil log valid (via TCP)
		logsValid, err := s.uc.GetAllScannedEvents(ip, true, func(ip, severity, msg string) {
			log.Printf("[%s] %s - %s", severity, ip, msg)
		})
		if err != nil {
			log.Printf("[ERROR] Gagal ambil log valid dari %s: %v", ip, err)
			continue
		}

		// Ambil log invalid (via UDP)
		logsInvalid, err := s.uc.GetAllScannedEvents(ip, false, func(ip, severity, msg string) {
			log.Printf("[%s] %s - %s", severity, ip, msg)
		})
		if err != nil {
			log.Printf("[ERROR] Gagal ambil log invalid dari %s: %v", ip, err)
			continue
		}

		// Cetak hasil log
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
}

func (s *DownloadJobService) Device4() {
	// ctx := context.Background()
	ips := []string{"192.168.80.29"}
	// ips := []string{"192.168.80.26"}

	for _, ip := range ips {
		// Ambil log valid (via TCP)
		logsValid, err := s.uc.GetAllScannedEvents(ip, true, func(ip, severity, msg string) {
			log.Printf("[%s] %s - %s", severity, ip, msg)
		})
		if err != nil {
			log.Printf("[ERROR] Gagal ambil log valid dari %s: %v", ip, err)
			continue
		}

		// Ambil log invalid (via UDP)
		logsInvalid, err := s.uc.GetAllScannedEvents(ip, false, func(ip, severity, msg string) {
			log.Printf("[%s] %s - %s", severity, ip, msg)
		})
		if err != nil {
			log.Printf("[ERROR] Gagal ambil log invalid dari %s: %v", ip, err)
			continue
		}

		// Cetak hasil log
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
}

func (s *DownloadJobService) Device5() {
	// ctx := context.Background()
	ips := []string{"192.168.80.19"}
	// ips := []string{"192.168.80.26"}

	for _, ip := range ips {
		// Ambil log valid (via TCP)
		logsValid, err := s.uc.GetAllScannedEvents(ip, true, func(ip, severity, msg string) {
			log.Printf("[%s] %s - %s", severity, ip, msg)
		})
		if err != nil {
			log.Printf("[ERROR] Gagal ambil log valid dari %s: %v", ip, err)
			continue
		}

		// Ambil log invalid (via UDP)
		logsInvalid, err := s.uc.GetAllScannedEvents(ip, false, func(ip, severity, msg string) {
			log.Printf("[%s] %s - %s", severity, ip, msg)
		})
		if err != nil {
			log.Printf("[ERROR] Gagal ambil log invalid dari %s: %v", ip, err)
			continue
		}

		// Cetak hasil log
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
}
