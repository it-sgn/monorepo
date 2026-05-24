package data

import (
	"context"
	"encoding/json"
	"fmt"
	departmentv1 "mall-go/api/department/service/v1"
	empv1 "mall-go/api/employers/service/v1"
	"mall-go/module/attendance-raw/service/internal/biz"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

var _ biz.AttendanceRawRepo = (*attendanceRawRepo)(nil)

type attendanceRawRepo struct {
	data             *Data
	log              *log.Helper
	rdb              *redis.Client
	employersClient  empv1.EmployersClient
	departmentClient departmentv1.DepartmentClient
}

func NewAttendanceRawRepo(data *Data, logger log.Logger) biz.AttendanceRawRepo {
	return &attendanceRawRepo{
		data:             data,
		log:              log.NewHelper(log.With(logger, "module", "data/attendance-raw")),
		rdb:              data.rdb,
		employersClient:  data.EmployersClient,
		departmentClient: data.DepartmentClient,
	}
}

func (r *attendanceRawRepo) GetAttendanceReport(
	ctx context.Context,
	startDate, endDate, depart string,
) (*biz.AttendanceReport, error) {

	cacheKey := fmt.Sprintf("attendance:report:%s:%s-%s", depart, startDate, endDate)

	// 1️⃣ Cek Redis Cache
	if val, err := r.rdb.Get(ctx, cacheKey).Result(); err == nil {
		var cached biz.AttendanceReport
		jsonErr := json.Unmarshal([]byte(val), &cached)
		if jsonErr == nil {
			r.log.Infof("✅ Redis cache hit: %s", cacheKey)
			return &cached, nil
		}
		r.log.Warnf("⚠️ Redis unmarshal failed for key [%s]: %v", cacheKey, jsonErr)
	} else if err != redis.Nil {
		r.log.Errorf("❌ Redis GET error for key [%s]: %v", cacheKey, err)
	}

	// 2️⃣ Ambil Employers
	empsResp, err := r.employersClient.GetEmployersFilterDepartCode(ctx, &empv1.GetEmployersFilterDepartCodeRequest{
		Departcode: depart,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get employers: %w", err)
	}
	if len(empsResp.Result) == 0 {
		return &biz.AttendanceReport{}, nil
	}

	// 3️⃣ Ambil Perusahaan
	companyResp, err := r.employersClient.GetPerusahaan(ctx, &empv1.GetPerusahaanRequest{
		Departcode: depart,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get perusahaan: %w", err)
	}
	if len(companyResp.Perusahaan) == 0 {
		return &biz.AttendanceReport{}, nil
	}
	comp := companyResp.Perusahaan[0]

	// 4️⃣ Ambil Department
	departResp, err := r.departmentClient.GetDepartmentCode(ctx, &departmentv1.GetDepartmentCodeRequest{
		DepartCode: depart,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get department: %w", err)
	}
	if departResp.DepartName == "" {
		return &biz.AttendanceReport{}, nil
	}
	depName := departResp.DepartName

	// 5️⃣ Query Attendance dari DB
	type Result struct {
		UserID   string `gorm:"column:user_id"`
		Jam      string `gorm:"column:jam"`
		Tgl      string `gorm:"column:tgl"`
		ClockIn  string `gorm:"column:clock_in"`
		ClockOut string `gorm:"column:clock_out"`
	}
	var results []Result

	userIDs := extractKaryaCodes(empsResp.Result)
	rawSQL := fmt.Sprintf(`
	SELECT DISTINCT ON (user_id, att_log::date)
	    user_id,
	    att_log AS jam,
	    TO_CHAR(att_log::date, 'YYYY-MM-DD') AS tgl,
	    TO_CHAR(MIN(att_log) FILTER (WHERE status = 0) OVER (PARTITION BY user_id, att_log::date), 'HH24:MI') AS clock_in,
	    TO_CHAR(MAX(att_log) FILTER (WHERE status = 1) OVER (PARTITION BY user_id, att_log::date), 'HH24:MI') AS clock_out
	FROM attendance_log
	WHERE user_id IN ('%s')
	AND att_log::date BETWEEN '%s' AND '%s'
	ORDER BY user_id, att_log::date, att_log
	`, strings.Join(userIDs, "','"), startDate, endDate)
	r.log.Infof("Query SQL: %s", rawSQL)
	err = r.data.db.Raw(rawSQL).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("gorm raw query error: %w", err)
	}

	// 6️⃣ Mapping ke AttendanceReport (proto-compatible)
	report := &biz.AttendanceReport{
		NamaPerusahaan: comp.NamaPerusahaan,
		Cabang:         comp.Cabang,
		Department:     depName,
		Jabatan:        "Operator",
		DibuatOleh:     "Budi",
		DiperiksaOleh:  "Andi",
		DisetujuiOleh:  "Sari",
		Karyawan:       []biz.Karyawan{},
	}

	// Map employer untuk akses cepat
	empMap := make(map[string]*empv1.EmployerItem)
	for _, e := range empsResp.Result {
		empMap[e.KaryaCode] = e
	}

	// Group hasil attendance per karyawan
	karyawanMap := make(map[string][]biz.Absensi)
	for _, row := range results {
		if _, exists := empMap[row.UserID]; exists {
			karyawanMap[row.UserID] = append(karyawanMap[row.UserID], biz.Absensi{
				Tanggal:  row.Tgl,
				Status:   []*biz.Status{{ClockIn: row.ClockIn, ClockOut: row.ClockOut}},
				Evaluasi: parseEvaluasi(row.ClockIn, row.ClockOut),
			})
		}
	}

	// Buat list Karyawan (loop sekali, tanpa duplikasi)
	for _, emp := range empsResp.Result {
		dataList := karyawanMap[emp.KaryaCode] // bisa kosong jika tidak ada data
		report.Karyawan = append(report.Karyawan, biz.Karyawan{
			Karyaname: emp.KaryaName,
			Periode:   fmt.Sprintf("%s - %s", startDate, endDate),
			Data:      dataList,
		})
	}

	// 7️⃣ Simpan ke Redis
	if dataJSON, err := json.Marshal(report); err == nil {
		if setErr := r.rdb.Set(ctx, cacheKey, dataJSON, 3*time.Minute).Err(); setErr != nil {
			r.log.Warnf("⚠️ Failed to set Redis key [%s]: %v", cacheKey, setErr)
		}
	}

	return report, nil
}

// parseEvaluasi ubah jam menjadi evaluasi sederhana
func parseEvaluasi(clockIn, clockOut string) string {
	if clockIn == "" {
		return "Tidak Hadir"
	}
	if clockIn > "08:00" {
		return "Terlambat"
	}
	return "Hadir"
}

// extractKaryaCodes ambil list kode karyawan dari employer
func extractKaryaCodes(emps []*empv1.EmployerItem) []string {
	codes := make([]string, 0, len(emps))
	for _, e := range emps {
		codes = append(codes, e.KaryaCode)
	}
	return codes
}

// package data

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	departmentv1 "mall-go/api/department/service/v1"
// 	empv1 "mall-go/api/employers/service/v1"
// 	"mall-go/module/attendance-raw/service/internal/biz"
// 	"strings"
// 	"time"

// 	"github.com/go-kratos/kratos/v2/log"
// 	"github.com/go-redis/redis/v8"
// )

// var _ biz.AttendanceRawRepo = (*attendanceRawRepo)(nil)

// type attendanceRawRepo struct {
// 	data             *Data
// 	log              *log.Helper
// 	rdb              *redis.Client
// 	employersClient  empv1.EmployersClient
// 	departmentClient departmentv1.DepartmentClient
// }

// func NewAttendanceRawRepo(data *Data, logger log.Logger) biz.AttendanceRawRepo {
// 	return &attendanceRawRepo{
// 		data:             data,
// 		log:              log.NewHelper(log.With(logger, "module", "data/attendance-raw")),
// 		rdb:              data.rdb,
// 		employersClient:  data.EmployersClient,
// 		departmentClient: data.DepartmentClient,
// 	}
// }

// func (r *attendanceRawRepo) GetAttendanceReport(
// 	ctx context.Context,
// 	startDate, endDate, depart string,
// ) (*biz.GetAttendanceReportResponse, error) {

// 	cacheKey := fmt.Sprintf("attendance:report:%s:%s-%s", depart, startDate, endDate)
// 	// 1️⃣ Cek Redis Cache
// 	if val, err := r.rdb.Get(ctx, cacheKey).Result(); err == nil {
// 		var cached biz.GetAttendanceReportResponse
// 		jsonErr := json.Unmarshal([]byte(val), &cached)
// 		if jsonErr == nil {
// 			r.log.Infof("✅ Redis cache hit: %s", cacheKey)
// 			return &cached, nil
// 		}
// 		r.log.Warnf("⚠️ Redis unmarshal failed for key [%s]: %v", cacheKey, jsonErr)
// 	} else if err != redis.Nil {
// 		r.log.Errorf("❌ Redis GET error for key [%s]: %v", cacheKey, err)
// 	}

// 	// 2️⃣ Ambil Employers
// 	empsResp, err := r.employersClient.GetEmployersFilterDepartCode(ctx, &empv1.GetEmployersFilterDepartCodeRequest{
// 		Departcode: depart,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get employers: %w", err)
// 	}
// 	if len(empsResp.Result) == 0 {
// 		return &biz.GetAttendanceReportResponse{}, nil
// 	}

// 	// 3️⃣ Ambil Perusahaan
// 	companyResp, err := r.employersClient.GetPerusahaan(ctx, &empv1.GetPerusahaanRequest{
// 		Departcode: depart,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get perusahaan: %w", err)
// 	}
// 	if len(companyResp.Perusahaan) == 0 {
// 		return &biz.GetAttendanceReportResponse{}, nil
// 	}
// 	comp := companyResp.Perusahaan[0]

// 	// 4️⃣ Ambil Department
// 	departResp, err := r.departmentClient.GetDepartmentCode(ctx, &departmentv1.GetDepartmentCodeRequest{
// 		DepartCode: depart,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get department: %w", err)
// 	}
// 	if departResp.DepartName == "" {
// 		return &biz.GetAttendanceReportResponse{}, nil
// 	}
// 	depName := departResp.DepartName

// 	// 5️⃣ Query Attendance dari DB
// 	type Result struct {
// 		UserID   string `gorm:"column:user_id"`
// 		Jam      string `gorm:"column:jam"`
// 		Tgl      string `gorm:"column:tgl"`
// 		ClockIn  string `gorm:"column:clock_in"`
// 		ClockOut string `gorm:"column:clock_out"`
// 	}
// 	var results []Result

// 	userIDs := extractKaryaCodes(empsResp.Result)
// 	rawSQL := fmt.Sprintf(`
// 	SELECT DISTINCT ON (user_id, att_log::date)
// 	    user_id,
// 	    att_log AS jam,
// 	    TO_CHAR(att_log::date, 'YYYY-MM-DD') AS tgl,
// 	    TO_CHAR(MIN(att_log) FILTER (WHERE status = 0) OVER (PARTITION BY user_id, att_log::date), 'HH24:MI') AS clock_in,
// 	    TO_CHAR(MAX(att_log) FILTER (WHERE status = 1) OVER (PARTITION BY user_id, att_log::date), 'HH24:MI') AS clock_out
// 	FROM attendance_log
// 	WHERE user_id IN ('%s')
// 	AND att_log::date BETWEEN '%s' AND '%s'
// 	ORDER BY user_id, att_log::date, att_log
// 	`, strings.Join(userIDs, "','"), startDate, endDate)
// 	r.log.Infof("Query SQL: %s", rawSQL)
// 	err = r.data.db.Raw(rawSQL).Scan(&results).Error
// 	if err != nil {
// 		return nil, fmt.Errorf("gorm raw query error: %w", err)
// 	}

// 	// 6️⃣ Mapping ke GetAttendanceReportResponse (proto-compatible)
// 	report := &biz.GetAttendanceReportResponse{
// 		NamaPerusahaan: comp.NamaPerusahaan,
// 		Cabang:         comp.Cabang,
// 		Department:     depName,
// 		Jabatan:        "Operator",
// 		Periode:        fmt.Sprintf("%s - %s", startDate, endDate),
// 		DibuatOleh:     "Budi",
// 		DiperiksaOleh:  "Andi",
// 		DisetujuiOleh:  "Sari",
// 		Karyawan:       []biz.Karyawan{},
// 	}

// 	// Map employer untuk akses cepat
// 	empMap := make(map[string]*empv1.EmployerItem)
// 	for _, e := range empsResp.Result {
// 		empMap[e.KaryaCode] = e
// 	}

// 	// Group hasil attendance per karyawan
// 	karyawanMap := make(map[string][]biz.AttendanceData)
// 	for _, row := range results {
// 		if _, exists := empMap[row.UserID]; exists {
// 			karyawanMap[row.UserID] = append(karyawanMap[row.UserID], biz.AttendanceData{
// 				Tanggal:  row.Tgl,
// 				Status:   []*biz.StatusData{{ClockIn: row.ClockIn, ClockOut: row.ClockOut}},
// 				Evaluasi: parseEvaluasi(row.ClockIn, row.ClockOut),
// 			})
// 		}
// 	}

// 	// Buat list Karyawan (loop sekali, tanpa duplikasi)
// 	for _, emp := range empsResp.Result {
// 		dataList := karyawanMap[emp.KaryaCode] // bisa kosong jika tidak ada data
// 		report.Karyawan = append(report.Karyawan, biz.Karyawan{
// 			Karyaname: emp.KaryaName,
// 			Data:      dataList,
// 		})
// 	}

// 	// 7️⃣ Simpan ke Redis
// 	if dataJSON, err := json.Marshal(report); err == nil {
// 		if setErr := r.rdb.Set(ctx, cacheKey, dataJSON, 3*time.Minute).Err(); setErr != nil {
// 			r.log.Warnf("⚠️ Failed to set Redis key [%s]: %v", cacheKey, setErr)
// 		}
// 	}

// 	return report, nil
// }

// // parseEvaluasi ubah jam menjadi evaluasi sederhana
// func parseEvaluasi(clockIn, clockOut string) string {
// 	if clockIn == "" {
// 		return "Tidak Hadir"
// 	}
// 	if clockIn > "08:00" {
// 		return "Terlambat"
// 	}
// 	return "Hadir"
// }

// // extractKaryaCodes ambil list kode karyawan dari employer
// func extractKaryaCodes(emps []*empv1.EmployerItem) []string {
// 	codes := make([]string, 0, len(emps))
// 	for _, e := range emps {
// 		codes = append(codes, e.KaryaCode)
// 	}
// 	return codes
// }

// // package data

// // import (
// // 	"context"
// // 	"encoding/json"
// // 	"fmt"
// // 	departmentv1 "mall-go/api/department/service/v1"
// // 	empv1 "mall-go/api/employers/service/v1"
// // 	"mall-go/module/attendance-raw/service/internal/biz"
// // 	"strings"
// // 	"time"

// // 	"github.com/go-kratos/kratos/v2/log"
// // 	"github.com/go-redis/redis/v8"
// // )

// // var _ biz.AttendanceRawRepo = (*attendanceRawRepo)(nil)

// // type attendanceRawRepo struct {
// // 	data             *Data
// // 	log              *log.Helper
// // 	rdb              *redis.Client
// // 	employersClient  empv1.EmployersClient
// // 	departmentClient departmentv1.DepartmentClient
// // }

// // func NewAttendanceRawRepo(data *Data, logger log.Logger) biz.AttendanceRawRepo {
// // 	return &attendanceRawRepo{
// // 		data:             data,
// // 		log:              log.NewHelper(log.With(logger, "module", "data/attendance-raw")),
// // 		rdb:              data.rdb,
// // 		employersClient:  data.EmployersClient,
// // 		departmentClient: data.DepartmentClient,
// // 	}
// // }

// // func (r *attendanceRawRepo) GetAttendanceReport(
// // 	ctx context.Context,
// // 	startDate, endDate, depart string,
// // ) (*biz.GetAttendanceReportResponse, error) {

// // 	cacheKey := fmt.Sprintf("attendance:report:%s:%s-%s", depart, startDate, endDate)
// // 	// 1️⃣ Cek Redis Cache
// // 	if val, err := r.rdb.Get(ctx, cacheKey).Result(); err == nil {
// // 		var cached biz.GetAttendanceReportResponse
// // 		jsonErr := json.Unmarshal([]byte(val), &cached)
// // 		if jsonErr == nil {
// // 			r.log.Infof("✅ Redis cache hit: %s", cacheKey)
// // 			return &cached, nil
// // 		}
// // 		r.log.Warnf("⚠️ Redis unmarshal failed for key [%s]: %v", cacheKey, jsonErr)
// // 	} else if err != redis.Nil {
// // 		r.log.Errorf("❌ Redis GET error for key [%s]: %v", cacheKey, err)
// // 	}

// // 	// 2️⃣ Ambil Employers
// // 	empsResp, err := r.employersClient.GetEmployersFilterDepartCode(ctx, &empv1.GetEmployersFilterDepartCodeRequest{
// // 		Departcode: depart,
// // 	})
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to get employers: %w", err)
// // 	}
// // 	if len(empsResp.Result) == 0 {
// // 		return &biz.GetAttendanceReportResponse{}, nil
// // 	}

// // 	// 3️⃣ Ambil Perusahaan
// // 	companyResp, err := r.employersClient.GetPerusahaan(ctx, &empv1.GetPerusahaanRequest{
// // 		Departcode: depart,
// // 	})
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to get perusahaan: %w", err)
// // 	}
// // 	if len(companyResp.Perusahaan) == 0 {
// // 		return &biz.GetAttendanceReportResponse{}, nil
// // 	}
// // 	comp := companyResp.Perusahaan[0]

// // 	// 4️⃣ Ambil Department
// // 	departResp, err := r.departmentClient.GetDepartmentCode(ctx, &departmentv1.GetDepartmentCodeRequest{
// // 		DepartCode: depart,
// // 	})
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to get department: %w", err)
// // 	}
// // 	if departResp.DepartName == "" {
// // 		return &biz.GetAttendanceReportResponse{}, nil
// // 	}
// // 	depName := departResp.DepartName

// // 	// 5️⃣ Query Attendance dari DB
// // 	type Result struct {
// // 		UserID   string `gorm:"column:user_id"`
// // 		Jam      string `gorm:"column:jam"`
// // 		Tgl      string `gorm:"column:tgl"`
// // 		ClockIn  string `gorm:"column:clock_in"`
// // 		ClockOut string `gorm:"column:clock_out"`
// // 	}
// // 	var results []Result

// // 	userIDs := extractKaryaCodes(empsResp.Result)
// // 	rawSQL := fmt.Sprintf(`
// // 	SELECT DISTINCT ON (user_id, att_log::date)
// // 	    user_id,
// // 	    att_log AS jam,
// // 	    TO_CHAR(att_log::date, 'YYYY-MM-DD') AS tgl,
// // 	    TO_CHAR(MIN(att_log) FILTER (WHERE status = 0) OVER (PARTITION BY user_id, att_log::date), 'HH24:MI') AS clock_in,
// // 	    TO_CHAR(MAX(att_log) FILTER (WHERE status = 1) OVER (PARTITION BY user_id, att_log::date), 'HH24:MI') AS clock_out
// // 	FROM attendance_log
// // 	WHERE user_id IN ('%s')
// // 	AND att_log::date BETWEEN '%s' AND '%s'
// // 	ORDER BY user_id, att_log::date, att_log
// // 	`, strings.Join(userIDs, "','"), startDate, endDate)
// // 	r.log.Infof("Query SQL: %s", rawSQL)
// // 	err = r.data.db.Raw(rawSQL).Scan(&results).Error
// // 	if err != nil {
// // 		return nil, fmt.Errorf("gorm raw query error: %w", err)
// // 	}

// // 	// 6️⃣ Mapping ke GetAttendanceReportResponse (proto-compatible)
// // 	report := &biz.GetAttendanceReportResponse{
// // 		NamaPerusahaan: comp.NamaPerusahaan,
// // 		Cabang:         comp.Cabang,
// // 		Department:     depName,
// // 		Jabatan:        "Operator",
// // 		Periode:        fmt.Sprintf("%s - %s", startDate, endDate),
// // 		DibuatOleh:     "Budi",
// // 		DiperiksaOleh:  "Andi",
// // 		DisetujuiOleh:  "Sari",
// // 		Karyawan:       []biz.Karyawan{},
// // 	}

// // 	// Map employer untuk akses cepat
// // 	empMap := make(map[string]*empv1.EmployerItem)
// // 	for _, e := range empsResp.Result {
// // 		empMap[e.KaryaCode] = e
// // 	}

// // 	// Group hasil attendance per karyawan
// // 	karyawanMap := make(map[string][]biz.AttendanceData)
// // 	for _, e := range empsResp.Result {
// // 		dataList := karyawanMap[e.KaryaCode]
// // 		report.Karyawan = append(report.Karyawan, biz.Karyawan{
// // 			Karyaname: e.KaryaName,
// // 			Data:      dataList,
// // 		})
// // 	}

// // 	// Buat list Karyawan sesuai proto
// // 	for _, emp := range empsResp.Result {
// // 		if dataList, ok := karyawanMap[emp.KaryaCode]; ok {
// // 			report.Karyawan = append(report.Karyawan, biz.Karyawan{
// // 				Karyaname: emp.KaryaName,
// // 				Data:      dataList,
// // 			})
// // 		} else {
// // 			// Jika tidak ada data, tetap buat karyawan dengan data kosong
// // 			report.Karyawan = append(report.Karyawan, biz.Karyawan{
// // 				Karyaname: emp.KaryaName,
// // 				Data:      []biz.AttendanceData{},
// // 			})
// // 		}
// // 	}

// // 	// 7️⃣ Simpan ke Redis
// // 	if dataJSON, err := json.Marshal(report); err == nil {
// // 		if setErr := r.rdb.Set(ctx, cacheKey, dataJSON, 3*time.Minute).Err(); setErr != nil {
// // 			r.log.Warnf("⚠️ Failed to set Redis key [%s]: %v", cacheKey, setErr)
// // 		}
// // 	}

// // 	return report, nil
// // }

// // // parseEvaluasi ubah jam menjadi evaluasi sederhana
// // func parseEvaluasi(clockIn, clockOut string) string {
// // 	if clockIn == "" {
// // 		return "Tidak Hadir"
// // 	}
// // 	if clockIn > "08:00" {
// // 		return "Terlambat"
// // 	}
// // 	return "Hadir"
// // }

// // // extractKaryaCodes ambil list kode karyawan dari employer
// // func extractKaryaCodes(emps []*empv1.EmployerItem) []string {
// // 	codes := make([]string, 0, len(emps))
// // 	for _, e := range emps {
// // 		codes = append(codes, e.KaryaCode)
// // 	}
// // 	return codes
// // }

// // package data

// // import (
// // 	"context"
// // 	"encoding/json"
// // 	"fmt"
// // 	departmentv1 "mall-go/api/department/service/v1"
// // 	empv1 "mall-go/api/employers/service/v1"
// // 	"mall-go/module/attendance-raw/service/internal/biz"
// // 	"strings"
// // 	"time"

// // 	"github.com/go-kratos/kratos/v2/log"
// // 	"github.com/go-redis/redis/v8"
// // )

// // var _ biz.AttendanceRawRepo = (*attendanceRawRepo)(nil)

// // type attendanceRawRepo struct {
// // 	data             *Data
// // 	log              *log.Helper
// // 	rdb              *redis.Client
// // 	employersClient  empv1.EmployersClient
// // 	departmentClient departmentv1.DepartmentClient
// // }

// // func NewAttendanceRawRepo(data *Data, logger log.Logger) biz.AttendanceRawRepo {
// // 	return &attendanceRawRepo{
// // 		data:             data,
// // 		log:              log.NewHelper(log.With(logger, "module", "data/attendance-raw")),
// // 		rdb:              data.rdb,
// // 		employersClient:  data.EmployersClient,
// // 		departmentClient: data.DepartmentClient,
// // 	}
// // }

// // func (r *attendanceRawRepo) GetAttendanceReport(
// // 	ctx context.Context,
// // 	startDate, endDate, depart string,
// // ) ([]*biz.DailyAttendanceReport, error) {

// // 	cacheKey := fmt.Sprintf("attendance:report:%s:%s-%s", depart, startDate, endDate)
// // 	// 1️⃣ Cek Redis Cache
// // 	if val, err := r.rdb.Get(ctx, cacheKey).Result(); err == nil {
// // 		var cached []*biz.DailyAttendanceReport
// // 		if jsonErr := json.Unmarshal([]byte(val), &cached); jsonErr == nil {
// // 			r.log.Infof("✅ Redis cache hit: %s", cacheKey)
// // 			return cached, nil
// // 		}
// // 		r.log.Warnf("⚠️ Redis unmarshal failed for key [%s]: %v", cacheKey, err)
// // 	} else if err != redis.Nil {
// // 		r.log.Errorf("❌ Redis GET error for key [%s]: %v", cacheKey, err)
// // 	}

// // 	// 2️⃣ Ambil Employers
// // 	empsResp, err := r.employersClient.GetEmployersFilterDepartCode(ctx, &empv1.GetEmployersFilterDepartCodeRequest{
// // 		Departcode: depart,
// // 	})
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to get employers: %w", err)
// // 	}
// // 	if len(empsResp.Result) == 0 {
// // 		return []*biz.DailyAttendanceReport{}, nil
// // 	}

// // 	// 3️⃣ Ambil Perusahaan
// // 	companyResp, err := r.employersClient.GetPerusahaan(ctx, &empv1.GetPerusahaanRequest{
// // 		Departcode: depart,
// // 	})
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to get perusahaan: %w", err)
// // 	}
// // 	if len(companyResp.Perusahaan) == 0 {
// // 		return []*biz.DailyAttendanceReport{}, nil
// // 	}
// // 	comp := companyResp.Perusahaan[0]

// // 	// 4️⃣ Ambil Department
// // 	departResp, err := r.departmentClient.GetDepartmentCode(ctx, &departmentv1.GetDepartmentCodeRequest{
// // 		DepartCode: depart,
// // 	})
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to get department: %w", err)
// // 	}
// // 	if departResp.DepartName == "" {
// // 		return []*biz.DailyAttendanceReport{}, nil
// // 	}
// // 	depName := departResp.DepartName

// // 	// 5️⃣ Query Attendance dari DB
// // 	type Result struct {
// // 		UserID   string `gorm:"column:user_id"`
// // 		Jam      string `gorm:"column:jam"`
// // 		Tgl      string `gorm:"column:tgl"`
// // 		ClockIn  string `gorm:"column:clock_in"`
// // 		ClockOut string `gorm:"column:clock_out"`
// // 	}
// // 	var results []Result

// // 	userIDs := extractKaryaCodes(empsResp.Result)
// // 	rawSQL := fmt.Sprintf(`
// // 	SELECT DISTINCT ON (user_id, att_log::date)
// // 	    user_id,
// // 	    att_log AS jam,
// // 	    TO_CHAR(att_log::date, 'YYYY-MM-DD') AS tgl,
// // 	    TO_CHAR(MIN(att_log) FILTER (WHERE status = 0) OVER (PARTITION BY user_id, att_log::date), 'HH24:MI') AS clock_in,
// // 	    TO_CHAR(MAX(att_log) FILTER (WHERE status = 1) OVER (PARTITION BY user_id, att_log::date), 'HH24:MI') AS clock_out
// // 	FROM attendance_log
// // 	WHERE user_id IN ('%s')
// // 	AND att_log::date BETWEEN '%s' AND '%s'
// // 	ORDER BY user_id, att_log::date, att_log
// // 	`, strings.Join(userIDs, "','"), startDate, endDate)
// // 	r.log.Infof("Query SQL: %s", rawSQL)
// // 	err = r.data.db.Raw(rawSQL).Scan(&results).Error
// // 	if err != nil {
// // 		return nil, fmt.Errorf("gorm raw query error: %w", err)
// // 	}

// // 	// 6️⃣ Mapping ke DailyAttendanceReport (proto-compatible)
// // 	report := &biz.DailyAttendanceReport{
// // 		NamaPerusahaan: comp.NamaPerusahaan,
// // 		Cabang:         comp.Cabang,
// // 		Department:     depName,
// // 		Jabatan:        "Operator",
// // 		Periode:        fmt.Sprintf("%s - %s", startDate, endDate),
// // 		DibuatOleh:     "Budi",
// // 		DiperiksaOleh:  "Andi",
// // 		DisetujuiOleh:  "Sari",
// // 		Karyawan:       []biz.DailyAttendanceKaryawan{},
// // 	}

// // 	// Map employer untuk akses cepat
// // 	empMap := make(map[string]*empv1.EmployerItem)
// // 	for _, e := range empsResp.Result {
// // 		empMap[e.KaryaCode] = e
// // 	}

// // 	// Group hasil attendance per karyawan
// // 	karyawanMap := make(map[string][]biz.DailyAttendance)
// // 	for _, row := range results {
// // 		if emp, exists := empMap[row.UserID]; exists {
// // 			karyawanMap[row.UserID] = append(karyawanMap[row.UserID], biz.DailyAttendance{
// // 				KaryaName: emp.KaryaName,
// // 				Tanggal:   row.Tgl,
// // 				Status: []*biz.AttStatus{
// // 					{ClockIn: row.ClockIn, ClockOut: row.ClockOut},
// // 				},
// // 				Evaluasi: parseEvaluasi(row.ClockIn, row.ClockOut),
// // 			})
// // 		}
// // 	}

// // 	// Buat list Karyawan sesuai proto
// // 	for _, emp := range empsResp.Result {
// // 		if dataList, ok := karyawanMap[emp.KaryaCode]; ok {
// // 			report.Karyawan = append(report.Karyawan, biz.DailyAttendanceKaryawan{
// // 				KaryaName: emp.KaryaName,
// // 				Data:      dataList,
// // 			})
// // 		} else {
// // 			// Jika tidak ada data, tetap buat karyawan dengan data kosong
// // 			report.Karyawan = append(report.Karyawan, biz.DailyAttendanceKaryawan{
// // 				KaryaName: emp.KaryaName,
// // 				Data:      []biz.DailyAttendance{},
// // 			})
// // 		}
// // 	}

// // 	// 7️⃣ Simpan ke Redis
// // 	if dataJSON, err := json.Marshal([]*biz.DailyAttendanceReport{report}); err == nil {
// // 		if setErr := r.rdb.Set(ctx, cacheKey, dataJSON, 3*time.Minute).Err(); setErr != nil {
// // 			r.log.Warnf("⚠️ Failed to set Redis key [%s]: %v", cacheKey, setErr)
// // 		}
// // 	}

// // 	return []*biz.DailyAttendanceReport{report}, nil
// // }

// // // parseEvaluasi ubah jam menjadi evaluasi sederhana
// // func parseEvaluasi(clockIn, clockOut string) string {
// // 	if clockIn == "" {
// // 		return "Tidak Hadir"
// // 	}
// // 	if clockIn > "08:00" {
// // 		return "Terlambat"
// // 	}
// // 	return "Hadir"
// // }

// // // extractKaryaCodes ambil list kode karyawan dari employer
// // func extractKaryaCodes(emps []*empv1.EmployerItem) []string {
// // 	codes := make([]string, 0, len(emps))
// // 	for _, e := range emps {
// // 		codes = append(codes, e.KaryaCode)
// // 	}
// // 	return codes
// // }

// // package data

// // import (
// // 	"context"
// // 	"encoding/json"
// // 	"fmt"
// // 	departmentv1 "mall-go/api/department/service/v1"
// // 	empv1 "mall-go/api/employers/service/v1"
// // 	"mall-go/module/attendance-raw/service/internal/biz"
// // 	"strings"
// // 	"time"

// // 	"github.com/go-kratos/kratos/v2/log"
// // 	"github.com/go-redis/redis/v8"
// // )

// // var _ biz.AttendanceRawRepo = (*attendanceRawRepo)(nil)

// // type attendanceRawRepo struct {
// // 	data             *Data
// // 	log              *log.Helper
// // 	rdb              *redis.Client
// // 	employersClient  empv1.EmployersClient
// // 	departmentClient departmentv1.DepartmentClient
// // }

// // func NewAttendanceRawRepo(data *Data, logger log.Logger) biz.AttendanceRawRepo {
// // 	return &attendanceRawRepo{
// // 		data:             data,
// // 		log:              log.NewHelper(log.With(logger, "module", "data/attendance-raw")),
// // 		rdb:              data.rdb,
// // 		employersClient:  data.EmployersClient,
// // 		departmentClient: data.DepartmentClient,
// // 	}
// // }

// // func (r *attendanceRawRepo) GetAttendanceReport(
// // 	ctx context.Context,
// // 	startDate, endDate, depart string,
// // ) ([]*biz.DailyAttendanceReport, error) {

// // 	cacheKey := fmt.Sprintf("attendance:report:%s:%s-%s", depart, startDate, endDate)
// // 	// 1️⃣ Cek Redis Cache
// // 	if val, err := r.rdb.Get(ctx, cacheKey).Result(); err == nil {
// // 		var cached []*biz.DailyAttendanceReport
// // 		if jsonErr := json.Unmarshal([]byte(val), &cached); jsonErr == nil {
// // 			r.log.Infof("✅ Redis cache hit: %s", cacheKey)
// // 			return cached, nil
// // 		}
// // 		r.log.Warnf("⚠️ Redis unmarshal failed for key [%s]: %v", cacheKey, err)
// // 	} else if err != redis.Nil {
// // 		r.log.Errorf("❌ Redis GET error for key [%s]: %v", cacheKey, err)
// // 	}

// // 	// 2️⃣ Ambil Employers
// // 	empsResp, err := r.employersClient.GetEmployersFilterDepartCode(ctx, &empv1.GetEmployersFilterDepartCodeRequest{
// // 		Departcode: depart,
// // 	})
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to get employers: %w", err)
// // 	}
// // 	if len(empsResp.Result) == 0 {
// // 		return []*biz.DailyAttendanceReport{}, nil
// // 	}

// // 	// 3️⃣ Ambil Perusahaan
// // 	companyResp, err := r.employersClient.GetPerusahaan(ctx, &empv1.GetPerusahaanRequest{
// // 		Departcode: depart,
// // 	})
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to get perusahaan: %w", err)
// // 	}
// // 	if len(companyResp.Perusahaan) == 0 {
// // 		return []*biz.DailyAttendanceReport{}, nil
// // 	}
// // 	comp := companyResp.Perusahaan[0]

// // 	// 4️⃣ Ambil Department
// // 	departResp, err := r.departmentClient.GetDepartmentCode(ctx, &departmentv1.GetDepartmentCodeRequest{
// // 		DepartCode: depart,
// // 	})
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to get department: %w", err)
// // 	}
// // 	if departResp.DepartName == "" {
// // 		return []*biz.DailyAttendanceReport{}, nil
// // 	}
// // 	depName := departResp.DepartName

// // 	// 5️⃣ Query Attendance dari DB dengan DISTINCT ON + WINDOW FUNCTION
// // 	type Result struct {
// // 		UserID   string `gorm:"column:user_id"`
// // 		Jam      string `gorm:"column:jam"`
// // 		Tgl      string `gorm:"column:tgl"`
// // 		ClockIn  string `gorm:"column:clock_in"`
// // 		ClockOut string `gorm:"column:clock_out"`
// // 	}
// // 	var results []Result

// // 	userIDs := extractKaryaCodes(empsResp.Result)
// // 	rawSQL := fmt.Sprintf(`
// // 	SELECT DISTINCT ON (user_id, att_log::date)
// // 	    user_id,
// // 	    att_log AS jam,
// // 	    TO_CHAR(att_log::date, 'YYYY-MM-DD') AS tgl,
// // 	    TO_CHAR(MIN(att_log) FILTER (WHERE status = 0) OVER (PARTITION BY user_id, att_log::date), 'HH24:MI') AS clock_in,
// // 	    TO_CHAR(MAX(att_log) FILTER (WHERE status = 1) OVER (PARTITION BY user_id, att_log::date), 'HH24:MI') AS clock_out
// // 	FROM attendance_log
// // 	WHERE user_id IN ('%s')
// // 	AND att_log::date BETWEEN '%s' AND '%s'
// // 	ORDER BY user_id, att_log::date, att_log
// // 	`, strings.Join(userIDs, "','"), startDate, endDate)
// // 	r.log.Infof("Query SQL: %s", rawSQL)
// // 	err = r.data.db.Raw(rawSQL).Scan(&results).Error
// // 	if err != nil {
// // 		return nil, fmt.Errorf("gorm raw query error: %w", err)
// // 	}

// // 	// 6️⃣ Mapping ke DailyAttendanceReport (proto-compatible)
// // 	report := &biz.DailyAttendanceReport{
// // 		NamaPerusahaan: comp.NamaPerusahaan,
// // 		Cabang:         comp.Cabang,
// // 		Department:     depName,
// // 		Jabatan:        "Operator", // default / bisa ambil dari data employer
// // 		Periode:        fmt.Sprintf("%s - %s", startDate, endDate),
// // 		DibuatOleh:     "Budi",
// // 		DiperiksaOleh:  "Andi",
// // 		DisetujuiOleh:  "Sari",
// // 		Data:           make([]biz.DailyAttendance, 0, len(results)),
// // 	}

// // 	// Map employer untuk akses cepat
// // 	empMap := make(map[string]*empv1.EmployerItem)
// // 	for _, e := range empsResp.Result {
// // 		empMap[e.KaryaCode] = e
// // 	}

// // 	for _, row := range results {
// // 		if emp, exists := empMap[row.UserID]; exists {
// // 			report.Data = append(report.Data, biz.DailyAttendance{
// // 				KaryaName: emp.KaryaName,
// // 				Tanggal:   row.Tgl,
// // 				Status: []*biz.AttStatus{
// // 					{ClockIn: row.ClockIn, ClockOut: row.ClockOut},
// // 				},
// // 				Evaluasi: parseEvaluasi(row.ClockIn, row.ClockOut),
// // 			})
// // 		}
// // 	}

// // 	// 7️⃣ Simpan ke Redis
// // 	if dataJSON, err := json.Marshal([]*biz.DailyAttendanceReport{report}); err == nil {
// // 		if setErr := r.rdb.Set(ctx, cacheKey, dataJSON, 3*time.Minute).Err(); setErr != nil {
// // 			r.log.Warnf("⚠️ Failed to set Redis key [%s]: %v", cacheKey, setErr)
// // 		}
// // 	}

// // 	return []*biz.DailyAttendanceReport{report}, nil
// // }

// // // parseEvaluasi ubah jam menjadi evaluasi sederhana
// // func parseEvaluasi(clockIn, clockOut string) string {
// // 	if clockIn == "" {
// // 		return "Tidak Hadir"
// // 	}
// // 	if clockIn > "08:00" {
// // 		return "Terlambat"
// // 	}
// // 	return "Hadir"
// // }

// // // extractKaryaCodes ambil list kode karyawan dari employer
// // func extractKaryaCodes(emps []*empv1.EmployerItem) []string {
// // 	codes := make([]string, 0, len(emps))
// // 	for _, e := range emps {
// // 		codes = append(codes, e.KaryaCode)
// // 	}
// // 	return codes
// // }

// // func (r *attendanceRawRepo) GetAttendanceReport(
// // 	ctx context.Context,
// // 	startDate, endDate, depart string,
// // ) ([]*biz.DailyAttendanceReport, error) {

// // 	cacheKey := fmt.Sprintf("attendance:report:%s:%s-%s", depart, startDate, endDate)
// // 	// 1️⃣ Cek Redis Cache
// // 	if val, err := r.rdb.Get(ctx, cacheKey).Result(); err == nil {
// // 		var cached []*biz.DailyAttendanceReport
// // 		if jsonErr := json.Unmarshal([]byte(val), &cached); jsonErr == nil {
// // 			r.log.Infof("✅ Redis cache hit: %s", cacheKey)
// // 			return cached, nil
// // 		}
// // 		r.log.Warnf("⚠️ Redis unmarshal failed for key [%s]: %v", cacheKey, err)
// // 	} else if err != redis.Nil {
// // 		r.log.Errorf("❌ Redis GET error for key [%s]: %v", cacheKey, err)
// // 	}
// // 	// 2️⃣ Ambil Employers
// // 	empsResp, err := r.employersClient.GetEmployersFilterDepartCode(ctx, &empv1.GetEmployersFilterDepartCodeRequest{
// // 		Departcode: depart,
// // 	})
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to get employers: %w", err)
// // 	}
// // 	if len(empsResp.Result) == 0 {
// // 		return []*biz.DailyAttendanceReport{}, nil
// // 	}

// // 	// 3️⃣ Ambil Perusahaan
// // 	companyResp, err := r.employersClient.GetPerusahaan(ctx, &empv1.GetPerusahaanRequest{
// // 		Departcode: depart,
// // 	})
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to get perusahaan: %w", err)
// // 	}
// // 	if len(companyResp.Perusahaan) == 0 {
// // 		return []*biz.DailyAttendanceReport{}, nil
// // 	}
// // 	comp := companyResp.Perusahaan[0]

// // 	// 4️⃣ Ambil Department
// // 	departResp, err := r.departmentClient.GetDepartmentCode(ctx, &departmentv1.GetDepartmentCodeRequest{
// // 		DepartCode: depart,
// // 	})
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to get department: %w", err)
// // 	}
// // 	if departResp.DepartName == "" {
// // 		return []*biz.DailyAttendanceReport{}, nil
// // 	}
// // 	depName := departResp.DepartName

// // 	// 5️⃣ Query Attendance dari DB
// // 	type Result struct {
// // 		UserID   string `gorm:"column:user_id"`
// // 		Jam      string `gorm:"column:jam"`
// // 		Tgl      string `gorm:"column:tgl"`
// // 		ClockIn  string `gorm:"column:clock_in"`
// // 		ClockOut string `gorm:"column:clock_out"`
// // 	}

// // 	var results []Result

// // 	// Ambil userIDs dari employer
// // 	userIDs := extractKaryaCodes(empsResp.Result)

// // 	// Buat raw SQL dengan DISTINCT ON + window function
// // 	rawSQL := fmt.Sprintf(`
// // SELECT DISTINCT ON (user_id, att_log::date)
// //     user_id,
// //     att_log AS jam,
// //     TO_CHAR(att_log::date, 'YYYY-MM-DD') AS tgl,
// //     TO_CHAR(MIN(att_log) FILTER (WHERE status = 1) OVER (PARTITION BY user_id, att_log::date), 'HH24:MI') AS clock_in,
// //     TO_CHAR(MAX(att_log) FILTER (WHERE status = 0) OVER (PARTITION BY user_id, att_log::date), 'HH24:MI') AS clock_out
// // FROM attendance_log
// // WHERE user_id IN ('%s')
// // AND att_log::date BETWEEN '%s' AND '%s'
// // ORDER BY user_id, att_log::date, att_log
// // `, strings.Join(userIDs, "','"), startDate, endDate)

// // 	// Eksekusi query
// // 	err = r.data.db.Raw(rawSQL).Scan(&results).Error
// // 	if err != nil {
// // 		return nil, fmt.Errorf("gorm raw query error: %w", err)
// // 	}

// // 	// 6️⃣ Mapping ke DailyAttendanceReport (proto-compatible)
// // 	report := &biz.DailyAttendanceReport{
// // 		NamaPerusahaan: comp.NamaPerusahaan,
// // 		Cabang:         comp.Cabang,
// // 		Department:     depName,
// // 		Jabatan:        "Operator", // default / bisa ambil dari data employer
// // 		Periode:        fmt.Sprintf("%s - %s", startDate, endDate),
// // 		DibuatOleh:     "Budi",
// // 		DiperiksaOleh:  "Andi",
// // 		DisetujuiOleh:  "Sari",
// // 		Data:           make([]biz.DailyAttendance, 0, len(results)),
// // 	}

// // 	// Map employer untuk akses cepat
// // 	empMap := make(map[string]*empv1.EmployerItem)
// // 	for _, e := range empsResp.Result {
// // 		empMap[e.KaryaCode] = e
// // 	}

// // 	for _, row := range results {
// // 		if emp, exists := empMap[row.UserID]; exists {
// // 			report.Data = append(report.Data, biz.DailyAttendance{
// // 				KaryaName: emp.KaryaName,
// // 				Tanggal:   row.Tgl,
// // 				Status: []*biz.AttStatus{
// // 					{ClockIn: row.ClockIn, ClockOut: row.ClockOut},
// // 				},
// // 				Evaluasi: parseEvaluasi(row.ClockIn, row.ClockOut),
// // 			})
// // 		}
// // 	}

// // 	// 7️⃣ Simpan ke Redis
// // 	if dataJSON, err := json.Marshal([]*biz.DailyAttendanceReport{report}); err == nil {
// // 		if setErr := r.rdb.Set(ctx, cacheKey, dataJSON, 3*time.Minute).Err(); setErr != nil {
// // 			r.log.Warnf("⚠️ Failed to set Redis key [%s]: %v", cacheKey, setErr)
// // 		}
// // 	}

// // 	return []*biz.DailyAttendanceReport{report}, nil
// // }

// // // parseEvaluasi ubah jam menjadi evaluasi sederhana
// // func parseEvaluasi(clockIn, clockOut string) string {
// // 	if clockIn == "" {
// // 		return "Tidak Hadir"
// // 	}
// // 	if clockIn > "08:00" {
// // 		return "Terlambat"
// // 	}
// // 	return "Hadir"
// // }

// // // extractKaryaCodes ambil list kode karyawan dari employer
// // func extractKaryaCodes(emps []*empv1.EmployerItem) []string {
// // 	codes := make([]string, 0, len(emps))
// // 	for _, e := range emps {
// // 		codes = append(codes, e.KaryaCode)
// // 	}
// // 	return codes
// // }

// // package data

// // import (
// // 	"context"
// // 	"encoding/json"
// // 	"fmt"
// // 	"time"

// // 	departmentv1 "mall-go/api/department/service/v1"
// // 	empv1 "mall-go/api/employers/service/v1"
// // 	orgv1 "mall-go/api/organization/service/v1"
// // 	"mall-go/module/attendance-raw/service/internal/biz"

// // 	"github.com/go-kratos/kratos/v2/log"
// // 	"github.com/go-redis/redis/v8"
// // )

// // var _ biz.AttendanceRawRepo = (*attendanceRawRepo)(nil)

// // type attendanceRawRepo struct {
// // 	data             *Data
// // 	log              *log.Helper
// // 	rdb              *redis.Client
// // 	employersClient  empv1.EmployersClient
// // 	departmentClient departmentv1.DepartmentClient
// // 	orgV1Client      orgv1.OrganizationServiceClient
// // }

// // func NewAttendanceRawRepo(data *Data, logger log.Logger) biz.AttendanceRawRepo {
// // 	return &attendanceRawRepo{
// // 		data:             data,
// // 		log:              log.NewHelper(log.With(logger, "module", "data/attendance-raw")),
// // 		rdb:              data.rdb,
// // 		employersClient:  data.EmployersClient,
// // 		departmentClient: data.DepartmentClient,
// // 		orgV1Client:      data.OrgV1Client,
// // 	}
// // }

// // func findEmployer(karyaCode string, emps []*empv1.EmployerItem) *empv1.EmployerItem {
// // 	for _, e := range emps {
// // 		if e.KaryaCode == karyaCode {
// // 			return e
// // 		}
// // 	}
// // 	return &empv1.EmployerItem{}
// // }

// // func getDayName(dateStr string) string {
// // 	t, err := time.Parse("2006-01-02", dateStr)
// // 	if err != nil {
// // 		return ""
// // 	}
// // 	return t.Weekday().String()
// // }

// // func (r *attendanceRawRepo) GetAttendanceReport(
// // 	ctx context.Context,
// // 	startDate, endDate, depart string,
// // ) ([]*biz.DailyAttendanceReport, error) {

// // 	cacheKey := fmt.Sprintf("attendance:report:%s:%s-%s", depart, startDate, endDate)

// // 	// 1️⃣ Cek Redis Cache
// // 	if val, err := r.rdb.Get(ctx, cacheKey).Result(); err == nil {
// // 		var cached []*biz.DailyAttendanceReport
// // 		if jsonErr := json.Unmarshal([]byte(val), &cached); jsonErr == nil {
// // 			r.log.Infof("✅ Redis cache hit: %s", cacheKey)
// // 			return cached, nil
// // 		}
// // 		r.log.Warnf("⚠️ Redis unmarshal failed for key [%s]: %v", cacheKey, err)
// // 	} else if err != redis.Nil {
// // 		r.log.Errorf("❌ Redis GET error for key [%s]: %v", cacheKey, err)
// // 	}

// // 	// 2️⃣ Ambil Employers
// // 	empsResp, err := r.employersClient.GetEmployersFilterDepartCode(ctx, &empv1.GetEmployersFilterDepartCodeRequest{
// // 		Departcode: depart,
// // 	})
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to get employers: %w", err)
// // 	}
// // 	if len(empsResp.Result) == 0 {
// // 		return []*biz.DailyAttendanceReport{}, nil
// // 	}
// // 	karyaCodes := extractKaryaCodes(empsResp.Result)

// // 	// 3️⃣ Ambil Perusahaan
// // 	companyResp, err := r.employersClient.GetPerusahaan(ctx, &empv1.GetPerusahaanRequest{
// // 		Departcode: depart,
// // 	})
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to get perusahaan: %w", err)
// // 	}
// // 	if len(companyResp.Perusahaan) == 0 {
// // 		return []*biz.DailyAttendanceReport{}, nil
// // 	}
// // 	comp := companyResp.Perusahaan[0]

// // 	// 4️⃣ Ambil Department
// // 	departResp, err := r.departmentClient.GetDepartmentCode(ctx, &departmentv1.GetDepartmentCodeRequest{
// // 		DepartCode: depart,
// // 	})
// // 	if err != nil {
// // 		return nil, fmt.Errorf("failed to get department: %w", err)
// // 	}
// // 	if departResp.DepartName == "" {
// // 		return []*biz.DailyAttendanceReport{}, nil
// // 	}
// // 	depName := departResp.DepartName

// // 	// 5️⃣ Query Attendance dari DB
// // 	type Result struct {
// // 		Tgl      string
// // 		Jam      string
// // 		ClockIn  string
// // 		ClockOut string
// // 		UserID   string
// // 	}

// // 	var results []Result
// // 	startTime := time.Now()

// // 	err = r.data.db.
// // 		Table("attendance_log").
// // 		Select([]string{
// // 			"user_id",
// // 			"device_ip",
// // 			"att_log AS jam",
// // 			"TO_CHAR(att_log::date, 'YYYY-MM-DD') AS tgl",
// // 			"TO_CHAR(MIN(att_log) FILTER (WHERE status = 0), 'HH24:MI') AS clock_in",
// // 			"TO_CHAR(MAX(att_log) FILTER (WHERE status = 1), 'HH24:MI') AS clock_out",
// // 		}).
// // 		Where("user_id IN ?", karyaCodes).
// // 		Where("att_log::date BETWEEN ? AND ?", startDate, endDate).
// // 		Group("user_id, device_ip, att_log::date").
// // 		Order("user_id, att_log::date").
// // 		Scan(&results).Error
// // 	r.log.Infof("🐘 Query executed in %s", time.Since(startTime))

// // 	if err != nil {
// // 		return nil, fmt.Errorf("gorm query error: %w", err)
// // 	}

// // 	// 6️⃣ Mapping hasil ke laporan
// // 	report := &biz.DailyAttendanceReport{
// // 		KodePerusahaan: comp.KodePerusahaan,
// // 		NamaPerusahaan: comp.NamaPerusahaan,
// // 		KodeCabang:     comp.KodeCabang,
// // 		Cabang:         comp.Cabang,
// // 		Department:     depName,
// // 		Periode:        fmt.Sprintf("%s - %s", startDate, endDate),
// // 		DibuatOleh:     "Suroso",
// // 		DiperiksaOleh:  "Heri",
// // 		DisetujuiOleh:  "Agus",
// // 		Data:           make([]biz.DailyAttendance, 0, len(results)),
// // 	}

// // 	// Map employer untuk akses cepat
// // 	empMap := make(map[string]*empv1.EmployerItem)
// // 	for _, e := range empsResp.Result {
// // 		empMap[e.KaryaCode] = e
// // 	}

// // 	for _, row := range results {
// // 		if emp, exists := empMap[row.UserID]; exists {
// // 			report.Data = append(report.Data, biz.DailyAttendance{
// // 				KaryaCode: emp.KaryaCode,
// // 				KaryaName: emp.KaryaName,
// // 				Tanggal:   row.Tgl,
// // 				Jam:       row.Jam,
// // 				Status: []*biz.AttStatus{
// // 					{ClockIn: row.ClockIn, ClockOut: row.ClockOut},
// // 				},
// // 				Evaluasi: parseEvaluasi(row.ClockOut),
// // 			})
// // 		}
// // 	}

// // 	r.log.Infof("📄 Report generated with %d records", len(report.Data))

// // 	// 7️⃣ Simpan ke Redis
// // 	if dataJSON, err := json.Marshal([]*biz.DailyAttendanceReport{report}); err == nil {
// // 		if setErr := r.rdb.Set(ctx, cacheKey, dataJSON, 3*time.Minute).Err(); setErr != nil {
// // 			r.log.Warnf("⚠️ Failed to set Redis key [%s]: %v", cacheKey, setErr)
// // 		}
// // 	}

// // 	return []*biz.DailyAttendanceReport{report}, nil
// // }

// // // parseEvaluasi ubah jam menjadi evaluasi
// // func parseEvaluasi(clockOut string) string {
// // 	if clockOut == "" {
// // 		return "Tidak Pulang"
// // 	}
// // 	return "Sesuai Jadwal"
// // }
