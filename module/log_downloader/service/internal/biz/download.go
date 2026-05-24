package biz

import (
	"context"
	"fmt"
	"time"

	"mall-go/pkg/gozk"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	maxRetries = 3
	retryDelay = 2 * time.Second
)

type LogData struct {
	UserID     int64
	DeviceIP   string
	Attendance time.Time
	Status     int
	CreatedBy  string
}

// Logger khusus untuk download log
type dlLogger struct {
	log *log.Helper
}

func (cl *dlLogger) Info(msg string, keysAndValues ...any) {
	cl.log.Infow(append([]any{"msg", msg}, keysAndValues...)...)
}

func (cl *dlLogger) Error(err error, msg string, keysAndValues ...any) {
	cl.log.Errorw(append([]any{"msg", msg, "error", err}, keysAndValues...)...)
}

type DownloadRepo interface {
	GetAllScannedEvents(ip string, tcp bool, onError func(ip, severity, msg string)) ([]LogData, error)
	InsertKeDB(context.Context, *LogData) (*LogData, error)
}

type DownloadUseCase struct {
	repo DownloadRepo
	log  *log.Helper
}

func NewDownload(repo DownloadRepo, logger log.Logger) *DownloadUseCase {
	return &DownloadUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "biz/cronjob")),
	}
}

func (uc *DownloadUseCase) GetAllScannedEvents(ip string, tcp bool, onError func(ip, severity, msg string)) ([]LogData, error) {
	var (
		zkConn *gozk.ZK
		err    error
	)

	// Retry koneksi
	for i := 1; i <= maxRetries; i++ {
		zkConn = gozk.NewZK(ip, gozk.WithTCP(tcp), gozk.WithTimezone(gozk.DefaultTimezone))
		err = zkConn.Connect()
		if err == nil {
			break
		}
		msg := fmt.Sprintf("⚠️ Retry %d/%d: Gagal koneksi ke %s: %v", i, maxRetries, ip, err)
		onError(ip, "warning", msg)
		time.Sleep(retryDelay)
	}

	if err != nil {
		msg := fmt.Sprintf("❌ Gagal konek ke device %s setelah %d percobaan: %v", ip, maxRetries, err)
		onError(ip, "error", msg)
		return nil, err
	}
	defer zkConn.Disconnect()

	// Properties
	if properties, err := zkConn.GetProperties(); err != nil {
		onError(ip, "warning", fmt.Sprintf("⚠️ Gagal ambil properties dari %s: %v", ip, err))
	} else {
		fmt.Printf("[Device %s] Properties:\n", ip)
		properties.Println()
	}

	// Get Events
	events, err := zkConn.GetAllScannedEvents()
	if err != nil {
		msg := fmt.Sprintf("❌ Gagal ambil scanned events dari %s: %v", ip, err)
		onError(ip, "error", msg)
		return nil, err
	}

	fmt.Printf("[Device %s] Jumlah event: %d\n", ip, len(events))

	ctx := context.Background()
	var logs []LogData
	for i, event := range events {
		log := LogData{
			UserID:     event.UserID,
			DeviceIP:   ip,
			Attendance: event.Timestamp,
			Status:     event.Status,
			CreatedBy:  "system",
		}
		logs = append(logs, log)

		// Simpan ke database
		if _, err := uc.repo.InsertKeDB(ctx, &log); err != nil {
			msg := fmt.Sprintf("❌ Gagal simpan data event #%d (user %d): %v", i, event.UserID, err)
			onError(ip, "error", msg)
			continue
		}
	}

	return logs, nil
}

func (uc *DownloadUseCase) Insert(ctx context.Context, dt *LogData) (*LogData, error) {
	return uc.repo.InsertKeDB(ctx, dt)
}
