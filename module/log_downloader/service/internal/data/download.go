package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mall-go/module/log_downloader/service/internal/biz"
	"mall-go/module/log_downloader/service/internal/data/model/attendance"
	"mall-go/pkg/gozk"

	"github.com/go-kratos/kratos/v2/log"
)

type downloadRepo struct {
	data   *Data
	logger *log.Helper
}

func NewDownloadRepo(data *Data, logger log.Logger) biz.DownloadRepo {
	return &downloadRepo{
		data:   data,
		logger: log.NewHelper(log.With(logger, "module", "data/download")),
	}
}

const (
	maxRetries = 3
	retryDelay = 2 * time.Second
)

// GetAllScannedEvents menghubungkan ke mesin fingerprint dan mengambil seluruh log
func (r *downloadRepo) GetAllScannedEvents(ip string, tcp bool, onError func(ip, severity, msg string)) ([]biz.LogData, error) {
	var (
		zkConn *gozk.ZK
		err    error
	)

	// Coba konek beberapa kali
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

	// Ambil properti dari device (opsional)
	properties, err := zkConn.GetProperties()
	if err != nil {
		onError(ip, "warning", fmt.Sprintf("⚠️ Gagal ambil properties dari %s: %v", ip, err))
	} else {
		r.logger.Infof("[Device %s] Properties loaded", ip)
		properties.Println()
	}

	// Ambil semua log absensi
	events, err := zkConn.GetAllScannedEvents()
	if err != nil {
		msg := fmt.Sprintf("❌ Gagal ambil scanned events dari %s: %v", ip, err)
		onError(ip, "error", msg)
		return nil, err
	}

	r.logger.Infof("[Device %s] Jumlah event: %d", ip, len(events))

	var logs []biz.LogData
	for i, event := range events {
		r.logger.Debugf("[Device %s] Raw Event #%d: %+v", ip, i, event)

		logs = append(logs, biz.LogData{
			UserID:     event.UserID,
			DeviceIP:   ip,
			Attendance: event.Timestamp,
			Status:     event.Status,
			CreatedBy:  "system",
		})
	}

	return logs, nil
}

func (r *downloadRepo) InsertKeDB(ctx context.Context, dt *biz.LogData) (*biz.LogData, error) {
	if dt == nil {
		return nil, errors.New("LogData is nil")
	}

	// Cek keberadaan data dengan composite key (lebih efisien)
	exist, err := r.data.db.Attendance.
		Query().
		Where(
			attendance.UserID(int(dt.UserID)),
			// attendance.DeviceIP(dt.DeviceIP),
			attendance.AttLog(dt.Attendance),
			// attendance.Status(dt.Status),
		).
		Exist(ctx)
	if err != nil {
		r.logger.Errorf("Gagal mengecek keberadaan data: %v", err)
		return nil, err
	}

	if exist {
		r.logger.Infof("Duplikat ditemukan, lewati insert: user_id=%d, ip=%s, att_log=%v, status=%d",
			dt.UserID, dt.DeviceIP, dt.Attendance, dt.Status)
		return dt, nil
	}

	// Mulai proses insert
	create := r.data.db.Attendance.Create().
		SetDeviceIP(dt.DeviceIP).
		SetUserID(int(dt.UserID)).
		SetStatus(dt.Status).
		SetAttLog(dt.Attendance).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now())

	// Jika kamu ingin logika tambahan: clock_in / clock_out dari status
	// if dt.Status == 0 {
	// 	create = create.SetClockIn(dt.Attendance)
	// } else if dt.Status == 1 {
	// 	create = create.SetClockOut(dt.Attendance)
	// }

	// Simpan ke DB
	res, err := create.Save(ctx)
	if err != nil {
		r.logger.Errorf("Gagal insert data ke DB: %v", err)
		return nil, err
	}

	r.logger.Infof("Insert berhasil: user_id=%d, ip=%s, att_log=%v, status=%d",
		dt.UserID, dt.DeviceIP, dt.Attendance, dt.Status)

	return &biz.LogData{
		DeviceIP:   res.DeviceIP,
		UserID:     int64(res.UserID),
		Attendance: res.AttLog,
		Status:     res.Status,
	}, nil
}

// func (r *downloadRepo) InsertKeDB(ctx context.Context, dt *biz.LogData) (*biz.LogData, error) {
// 	if dt == nil {
// 		return nil, errors.New("downloadRepo is nil")
// 	}

// 	// Cek apakah data sudah ada (skip jika sudah ada)
// 	exist, err := r.data.db.Attendance.
// 		Query().
// 		Where(
// 			attendance.UserID(int(dt.UserID)),
// 			attendance.DeviceIP(dt.DeviceIP),
// 			attendance.AttLog(dt.Attendance),
// 			attendance.Status(dt.Status),
// 		).
// 		Exist(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if exist {
// 		r.logger.Infof("Data sudah ada, skip insert: user_id=%d, ip=%s, att_log=%v, status=%d",
// 			dt.UserID, dt.DeviceIP, dt.Attendance, dt.Status)
// 		return dt, nil // return existing log, atau nil sesuai kebutuhan
// 	}

// 	// Insert data baru
// 	res, err := r.data.db.Attendance.Create().
// 		SetDeviceIP(dt.DeviceIP).
// 		SetUserID(int(dt.UserID)).
// 		SetStatus(dt.Status).
// 		SetAttLog(dt.Attendance).
// 		// SetAttDate(dt.Attendance). // kalau tanggal == att_log
// 		// SetClockIn(dt.Attendance). // hanya jika status = 0
// 		Save(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &biz.LogData{
// 		DeviceIP:   res.DeviceIP,
// 		UserID:     int64(res.UserID),
// 		Attendance: res.AttLog,
// 		Status:     res.Status,
// 	}, nil
// }

// func (r *downloadRepo) InsertKeDB(ctx context.Context, dt *biz.LogData) (*biz.LogData, error) {
// 	if r == nil {
// 		return nil, errors.New("downloadRepo is nil")
// 	}
// 	if dt.Status == 0 {
// 		// Insert clock_in
// 		resIn, err := r.data.db.Attendance.Create().
// 			SetDeviceIP(dt.DeviceIP).
// 			SetUserID(int(dt.UserID)).
// 			SetStatus(dt.Status).
// 			SetAttLog(dt.Attendance).
// 			SetClockIn(dt.Attendance).
// 			Save(ctx)
// 		if err != nil {
// 			return nil, fmt.Errorf("gagal insert clock_in: %w", err)
// 		}

// 		return &biz.LogData{
// 			DeviceIP:   resIn.DeviceIP,
// 			UserID:     int64(resIn.UserID),
// 			Attendance: resIn.AttLog,
// 			Status:     resIn.Status,
// 			CreatedBy:  "System",
// 		}, nil
// 	}

// 	if dt.Status == 1 {
// 		// Update clock_out berdasarkan user_id + ip + tanggal yang sama
// 		updated, err := r.data.db.Attendance.
// 			Update().
// 			Where(
// 				attendance.UserID(int(dt.UserID)),
// 				attendance.DeviceIP(dt.DeviceIP),
// 				attendance.ClockInNotNil(),
// 				attendance.ClockOutIsNil(), // Pastikan belum pernah diisi clock_out
// 				attendance.TanggalEQ(dt.Attendance.Truncate(24*time.Hour)),
// 			).
// 			SetClockOut(dt.Attendance).
// 			Save(ctx)
// 		if err != nil {
// 			return nil, fmt.Errorf("gagal update clock_out: %w", err)
// 		}

// 		// Jika tidak ada row yang diupdate
// 		if updated == 0 {
// 			log.Infof("[WARN] Tidak ada record clock_in yang cocok untuk update clock_out: user_id=%d, ip=%s", dt.UserID, dt.DeviceIP)
// 		}

// 		// Return original log (karena tidak ambil record yang diupdate)
// 		return dt, nil
// 	}

// 	return nil, fmt.Errorf("status tidak dikenal: %d", dt.Status)
// }

// func (r *downloadRepo) InsertKeDB(ctx context.Context, dt *biz.LogData) (*biz.LogData, error) {
// 	if dt.Status==0{
// 	resIn, err := r.data.db.Attendance.Create().
// 		SetDeviceIP(dt.DeviceIP).
// 		SetUserID(int(dt.UserID)).
// 		SetStatus(dt.Status).
// 		// SetClockIn()
// 		SetAttLog(dt.Attendance).
// 		SetClockIn(dt.Attendance)
// 		Save(ctx)
// 	if err != nil {
// 		return nil, nil
// 	}
// } if dt.Status==1{
// 	resUp, err:= r.data.db.Attendance.Update().
// 	SetClockOut(dt.Attendance).
// 	Save(ctx)
// }
// 	// if dt.Status=1
// 	// resUp, err := r.data.db.Attendance.Update().

// 	return &biz.LogData{
// 		DeviceIP:   resIn.ip,
// 		UserID:     resIn.UserID,
// 		Attendance: resIn.Attendance,
// 		Status:     resIn.status,
// 		CreatedBy:  resIn.CreatedBy,
// 	}, nil

// }
