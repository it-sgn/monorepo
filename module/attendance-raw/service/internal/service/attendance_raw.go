package service

import (
	"context"

	v1 "mall-go/api/attendance-raw/service/v1"
	"mall-go/module/attendance-raw/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type AttendanceRawService struct {
	v1.UnimplementedAttendanceServiceServer
	log *log.Helper
	uc  *biz.AttendanceRawUsecase
}

func NewAttendanceRawService(uc *biz.AttendanceRawUsecase, logger log.Logger) *AttendanceRawService {
	return &AttendanceRawService{
		log: log.NewHelper(log.With(logger, "module", "service/attendanceRaw")),
		uc:  uc,
	}
}

// GetAttendanceReport retrieves the attendance report.
func (s *AttendanceRawService) GetAttendanceReport(
	ctx context.Context,
	req *v1.GetAttendanceReportRequest,
) (*v1.GetAttendanceReportResponse, error) {

	// 1️⃣ Ambil data dari usecase/repo
	report, err := s.uc.GetAttendanceReport(ctx, req.Start, req.End, req.DepartCode)
	if err != nil {
		s.log.Errorf("failed to get attendance report: %v", err)
		return nil, err
	}

	if report == nil {
		return &v1.GetAttendanceReportResponse{}, nil
	}

	// 2️⃣ Mapping ke response proto
	protoReport := &v1.AttendanceReport{
		NamaPerusahaan: report.NamaPerusahaan,
		Cabang:         report.Cabang,
		Department:     report.Department,
		Jabatan:        report.Jabatan,
		DibuatOleh:     report.DibuatOleh,
		DiperiksaOleh:  report.DiperiksaOleh,
		DisetujuiOleh:  report.DisetujuiOleh,
		Karyawan:       make([]*v1.Karyawan, 0, len(report.Karyawan)),
	}

	for _, k := range report.Karyawan {
		protoKaryawan := &v1.Karyawan{
			Karyaname: k.Karyaname,
			Periode:   k.Periode,
			Data:      make([]*v1.Absensi, 0, len(k.Data)),
		}

		for _, d := range k.Data {
			statusItems := make([]*v1.Status, 0, len(d.Status))
			for _, sItem := range d.Status {
				statusItems = append(statusItems, &v1.Status{
					ClockIn:  sItem.ClockIn,
					ClockOut: sItem.ClockOut,
				})
			}

			protoKaryawan.Data = append(protoKaryawan.Data, &v1.Absensi{
				Tanggal:  d.Tanggal,
				Evaluasi: d.Evaluasi,
				Status:   statusItems,
			})
		}

		protoReport.Karyawan = append(protoReport.Karyawan, protoKaryawan)
	}

	return &v1.GetAttendanceReportResponse{
		Report: protoReport,
	}, nil
}

// package service

// import (
// 	"context"

// 	v1 "mall-go/api/attendance-raw/service/v1"
// 	"mall-go/module/attendance-raw/service/internal/biz"

// 	"github.com/go-kratos/kratos/v2/log"
// )

// type AttendanceRawService struct {
// 	v1.UnimplementedAttendanceServiceServer
// 	log *log.Helper
// 	uc  *biz.AttendanceRawUsecase
// }

// func NewAttendanceRawService(uc *biz.AttendanceRawUsecase, logger log.Logger) *AttendanceRawService {
// 	return &AttendanceRawService{
// 		log: log.NewHelper(log.With(logger, "module", "service/attendanceRaw")),
// 		uc:  uc,
// 	}
// }

// // GetAttendanceReport retrieves the attendance report.
// func (s *AttendanceRawService) GetAttendanceReport(
// 	ctx context.Context,
// 	req *v1.GetAttendanceReportRequest,
// ) (*v1.GetAttendanceReportResponse, error) {

// 	// 1️⃣ Ambil data dari usecase/repo
// 	report, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.Depart)
// 	if err != nil {
// 		s.log.Errorf("failed to get attendance report: %v", err)
// 		return nil, err
// 	}

// 	if report == nil {
// 		return &v1.GetAttendanceReportResponse{}, nil
// 	}

// 	// 2️⃣ Mapping ke response proto
// 	resp := &v1.GetAttendanceReportResponse{
// 		NamaPerusahaan: report.NamaPerusahaan,
// 		Cabang:         report.Cabang,
// 		Department:     report.Department,
// 		Jabatan:        report.Jabatan,
// 		Periode:        report.Periode,
// 		DibuatOleh:     report.DibuatOleh,
// 		DiperiksaOleh:  report.DiperiksaOleh,
// 		DisetujuiOleh:  report.DisetujuiOleh,
// 		Karyawan:       make([]*v1.Karyawan, 0, len(report.Karyawan)),
// 	}

// 	for _, k := range report.Karyawan {
// 		dataItems := make([]*v1.AttendanceData, 0, len(k.Data))
// 		for _, d := range k.Data {
// 			statusItems := make([]*v1.StatusData, 0, len(d.Status))
// 			for _, sItem := range d.Status {
// 				statusItems = append(statusItems, &v1.StatusData{
// 					ClockIn:  sItem.ClockIn,
// 					ClockOut: sItem.ClockOut,
// 				})
// 			}

// 			dataItems = append(dataItems, &v1.AttendanceData{
// 				Tanggal: d.Tanggal,
// 				// Karyaname: "", // Kosong karena nama karyawan ada di level Karyawan
// 				Evaluasi: d.Evaluasi,
// 				Status:   statusItems,
// 			})
// 		}

// 		resp.Karyawan = append(resp.Karyawan, &v1.Karyawan{
// 			Karyaname: k.Karyaname,
// 			Data:      dataItems,
// 		})
// 	}

// 	return resp, nil
// }

// // package service

// // import (
// // 	"context"

// // 	v1 "mall-go/api/attendance-raw/service/v1"
// // 	"mall-go/module/attendance-raw/service/internal/biz"

// // 	"github.com/go-kratos/kratos/v2/log"
// // )

// // type AttendanceRawService struct {
// // 	v1.UnimplementedAttendanceServiceServer
// // 	log *log.Helper
// // 	uc  *biz.AttendanceRawUsecase
// // }

// // func NewAttendanceRawService(uc *biz.AttendanceRawUsecase, logger log.Logger) *AttendanceRawService {
// // 	return &AttendanceRawService{
// // 		log: log.NewHelper(log.With(logger, "module", "service/attendanceRaw")),
// // 		uc:  uc,
// // 	}
// // }

// // // GetAttendanceReport retrieves the attendance report.
// // func (s *AttendanceRawService) GetAttendanceReport(
// // 	ctx context.Context,
// // 	req *v1.GetAttendanceReportRequest,
// // ) (*v1.GetAttendanceReportResponse, error) {

// // 	// 1️⃣ Ambil data dari usecase/repo
// // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.Depart)
// // 	if err != nil {
// // 		s.log.Errorf("failed to get attendance report: %v", err)
// // 		return nil, err
// // 	}

// // 	if len(reports) == 0 {
// // 		return &v1.GetAttendanceReportResponse{}, nil
// // 	}

// // 	first := reports[0]

// // 	// 2️⃣ Mapping ke response proto
// // 	resp := &v1.GetAttendanceReportResponse{
// // 		NamaPerusahaan: first.NamaPerusahaan,
// // 		Cabang:         first.Cabang,
// // 		Department:     first.Department,
// // 		Jabatan:        first.Jabatan,
// // 		Periode:        first.Periode,
// // 		DibuatOleh:     first.DibuatOleh,
// // 		DiperiksaOleh:  first.DiperiksaOleh,
// // 		DisetujuiOleh:  first.DisetujuiOleh,
// // 		Karyawan:       make([]*v1.Karyawan, 0, len(first.Karyawan)),
// // 	}

// // 	for _, k := range first.Karyawan {
// // 		dataItems := make([]*v1.AttendanceData, 0, len(k.Data))
// // 		for _, d := range k.Data {
// // 			statusItems := make([]*v1.StatusData, 0, len(d.Status))
// // 			for _, sItem := range d.Status {
// // 				statusItems = append(statusItems, &v1.StatusData{
// // 					ClockIn:  sItem.ClockIn,
// // 					ClockOut: sItem.ClockOut,
// // 				})
// // 			}

// // 			dataItems = append(dataItems, &v1.AttendanceData{
// // 				Tanggal:   d.Tanggal,
// // 				Karyaname: "", // Kosong karena nama karyawan ada di level Karyawan
// // 				Evaluasi:  d.Evaluasi,
// // 				Status:    statusItems,
// // 			})
// // 		}

// // 		resp.Karyawan = append(resp.Karyawan, &v1.Karyawan{
// // 			Karyaname: k.KaryaName,
// // 			Data:      dataItems,
// // 		})
// // 	}

// // 	return resp, nil
// // }

// // // package service

// // // import (
// // // 	"context"

// // // 	v1 "mall-go/api/attendance-raw/service/v1"
// // // 	"mall-go/module/attendance-raw/service/internal/biz"

// // // 	"github.com/go-kratos/kratos/v2/log"
// // // )

// // // type AttendanceRawService struct {
// // // 	v1.UnimplementedAttendanceServiceServer
// // // 	log *log.Helper
// // // 	uc  *biz.AttendanceRawUsecase
// // // }

// // // func NewAttendanceRawService(uc *biz.AttendanceRawUsecase, logger log.Logger) *AttendanceRawService {
// // // 	return &AttendanceRawService{
// // // 		log: log.NewHelper(log.With(logger, "module", "service/attendanceRaw")),
// // // 		uc:  uc,
// // // 	}
// // // }

// // // // GetAttendanceReport retrieves the attendance report.
// // // func (s *AttendanceRawService) GetAttendanceReport(
// // // 	ctx context.Context,
// // // 	req *v1.GetAttendanceReportRequest,
// // // ) (*v1.GetAttendanceReportResponse, error) {

// // // 	// 1️⃣ Ambil data dari usecase/repo
// // // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.Depart)
// // // 	if err != nil {
// // // 		s.log.Errorf("failed to get attendance report: %v", err)
// // // 		return nil, err
// // // 	}

// // // 	if len(reports) == 0 {
// // // 		return &v1.GetAttendanceReportResponse{}, nil
// // // 	}

// // // 	first := reports[0]
// // // 	s.log.Infof("REPORT pertama: %+v", first)

// // // 	// 2️⃣ Mapping ke response proto
// // // 	resp := &v1.GetAttendanceReportResponse{
// // // 		NamaPerusahaan: first.NamaPerusahaan,
// // // 		Cabang:         first.Cabang,
// // // 		Department:     first.Department,
// // // 		Jabatan:        first.Jabatan,
// // // 		Periode:        first.Periode,
// // // 		DibuatOleh:     first.DibuatOleh,
// // // 		DiperiksaOleh:  first.DiperiksaOleh,
// // // 		DisetujuiOleh:  first.DisetujuiOleh,
// // // 		Karyawan:       make([]*v1.Karyawan, 0, len(first.Karyawan)),
// // // 	}

// // // 	for _, k := range first.Karyawan {
// // // 		dataItems := make([]*v1.AttendanceData, 0, len(k.Data))
// // // 		for _, d := range k.Data {
// // // 			statusItems := make([]*v1.StatusData, 0, len(d.Status))
// // // 			for _, s := range d.Status {
// // // 				statusItems = append(statusItems, &v1.StatusData{
// // // 					ClockIn:  s.ClockIn,
// // // 					ClockOut: s.ClockOut,
// // // 				})
// // // 			}

// // // 			dataItems = append(dataItems, &v1.AttendanceData{
// // // 				Tanggal:   d.Tanggal,
// // // 				Karyaname: "", // kosong karena nama karyawan di Karyawan.Karyaname
// // // 				Evaluasi:  d.Evaluasi,
// // // 				Status:    statusItems,
// // // 			})
// // // 		}

// // // 		resp.Karyawan = append(resp.Karyawan, &v1.Karyawan{
// // // 			Karyaname: k.KaryaName, // hanya nama karyawan
// // // 			Data:      dataItems,
// // // 		})
// // // 	}
// // // 	// resp := &v1.GetAttendanceReportResponse{
// // // 	// 	NamaPerusahaan: first.NamaPerusahaan,
// // // 	// 	Cabang:         first.Cabang,
// // // 	// 	Department:     first.Department,
// // // 	// 	Jabatan:        first.Jabatan,
// // // 	// 	Periode:        first.Periode,
// // // 	// 	DibuatOleh:     first.DibuatOleh,
// // // 	// 	DiperiksaOleh:  first.DiperiksaOleh,
// // // 	// 	DisetujuiOleh:  first.DisetujuiOleh,
// // // 	// 	Karyawan:       make([]*v1.Karyawan, 0, len(first.Karyawan)),
// // // 	// }

// // // 	// for _, k := range first.Karyawan {
// // // 	// 	dataItems := make([]*v1.AttendanceData, 0, len(k.Data))
// // // 	// 	for _, d := range k.Data {
// // // 	// 		statusItems := make([]*v1.StatusData, 0, len(d.Status))
// // // 	// 		for _, s := range d.Status {
// // // 	// 			statusItems = append(statusItems, &v1.StatusData{
// // // 	// 				ClockIn:  s.ClockIn,
// // // 	// 				ClockOut: s.ClockOut,
// // // 	// 			})
// // // 	// 		}

// // // 	// 		dataItems = append(dataItems, &v1.AttendanceData{
// // // 	// 			Tanggal:   d.Tanggal,
// // // 	// 			Karyaname: d.KaryaName,
// // // 	// 			Evaluasi:  d.Evaluasi,
// // // 	// 			Status:    statusItems,
// // // 	// 		})
// // // 	// 	}

// // // 	// 	resp.Karyawan = append(resp.Karyawan, &v1.Karyawan{
// // // 	// 		Karyaname: k.KaryaName,
// // // 	// 		Data:      dataItems,
// // // 	// 	})
// // // 	// }

// // // 	return resp, nil
// // // }

// // // package service

// // // import (
// // // 	"context"

// // // 	v1 "mall-go/api/attendance-raw/service/v1"
// // // 	"mall-go/module/attendance-raw/service/internal/biz"

// // // 	"github.com/go-kratos/kratos/v2/log"
// // // )

// // // type AttendanceRawService struct {
// // // 	v1.UnimplementedAttendanceServiceServer
// // // 	log *log.Helper
// // // 	uc  *biz.AttendanceRawUsecase
// // // }

// // // func NewAttendanceRawService(uc *biz.AttendanceRawUsecase, logger log.Logger) *AttendanceRawService {
// // // 	return &AttendanceRawService{
// // // 		log: log.NewHelper(log.With(logger, "module", "service/attendanceRaw")),
// // // 		uc:  uc,
// // // 	}
// // // }

// // // // GetAttendanceReport retrieves the attendance report.
// // // func (s *AttendanceRawService) GetAttendanceReport(
// // // 	ctx context.Context,
// // // 	req *v1.GetAttendanceReportRequest,
// // // ) (*v1.GetAttendanceReportResponse, error) {

// // // 	// 1️⃣ Ambil data dari usecase/repo (sudah pakai DISTINCT ON + window function)
// // // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.DepartCode)
// // // 	if err != nil {
// // // 		s.log.Errorf("failed to get attendance report: %v", err)
// // // 		return nil, err
// // // 	}

// // // 	if len(reports) == 0 {
// // // 		return &v1.GetAttendanceReportResponse{}, nil
// // // 	}

// // // 	first := reports[0]
// // // 	s.log.Infof("REPORT pertama: %s", reports)

// // // 	// 2️⃣ Mapping ke response proto
// // // 	resp := &v1.GetAttendanceReportResponse{
// // // 		NamaPerusahaan: first.NamaPerusahaan,
// // // 		Cabang:         first.Cabang,
// // // 		Department:     first.Department,
// // // 		Jabatan:        first.Jabatan,
// // // 		Periode:        first.Periode,
// // // 		DibuatOleh:     first.DibuatOleh,
// // // 		DiperiksaOleh:  first.DiperiksaOleh,
// // // 		DisetujuiOleh:  first.DisetujuiOleh,
// // // 		Data:           make([]*v1.DataItem, 0, len(first.Data)),
// // // 	}
// // // 	s.log.Infof("DATA pertama: %s", first.Data)

// // // 	for _, d := range first.Data {
// // // 		// Mapping DailyAttendance -> DataItem
// // // 		statusItems := make([]*v1.StatusItem, 0, len(d.Status))
// // // 		for _, s := range d.Status {
// // // 			statusItems = append(statusItems, &v1.StatusItem{
// // // 				ClockIn:  s.ClockIn,
// // // 				ClockOut: s.ClockOut,
// // // 			})
// // // 		}

// // // 		s.log.Infof("D pertama: %s", statusItems)
// // // 		resp.Data = append(resp.Data, &v1.DataItem{
// // // 			Tanggal:   d.Tanggal,
// // // 			Karyaname: d.KaryaName,
// // // 			Status:    statusItems,
// // // 			Evaluasi:  d.Evaluasi,
// // // 		})
// // // 	}
// // // 	// s.log.Infof("RESP pertama: %s", resp)

// // // 	return resp, nil
// // // }

// // // package service

// // // import (
// // // 	"context"

// // // 	v1 "mall-go/api/attendance-raw/service/v1"
// // // 	"mall-go/module/attendance-raw/service/internal/biz"
// // // 	"strings"

// // // 	"github.com/go-kratos/kratos/v2/log"
// // // )

// // // type AttendanceRawService struct {
// // // 	v1.UnimplementedAttendanceServiceServer
// // // 	log *log.Helper
// // // 	uc  *biz.AttendanceRawUsecase
// // // }

// // // func NewAttendanceRawService(uc *biz.AttendanceRawUsecase, logger log.Logger) *AttendanceRawService {
// // // 	return &AttendanceRawService{
// // // 		log: log.NewHelper(log.With(logger, "module", "service/attendanceRaw")),
// // // 		uc:  uc,
// // // 	}
// // // }

// // // // GetAttendanceReport retrieves the attendance report.
// // // func (s *AttendanceRawService) GetAttendanceReport(
// // // 	ctx context.Context,
// // // 	req *v1.GetAttendanceReportRequest,
// // // ) (*v1.GetAttendanceReportResponse, error) {

// // // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.DepartCode)
// // // 	if err != nil {
// // // 		s.log.Errorf("failed to get attendance report: %v", err)
// // // 		return nil, err
// // // 	}

// // // 	if len(reports) == 0 {
// // // 		return &v1.GetAttendanceReportResponse{}, nil
// // // 	}

// // // 	first := reports[0]
// // // 	log.Info("ini jabatan", first.Jabatan)
// // // 	resp := &v1.GetAttendanceReportResponse{
// // // 		NamaPerusahaan: first.NamaPerusahaan,
// // // 		Cabang:         first.Cabang,
// // // 		Department:     first.Department,
// // // 		Jabatan:        first.Jabatan,
// // // 		Periode:        first.Periode,
// // // 		DibuatOleh:     first.DibuatOleh,
// // // 		DiperiksaOleh:  first.DiperiksaOleh,
// // // 		DisetujuiOleh:  first.DisetujuiOleh,
// // // 		Data:           make([]*v1.DataItem, 0, len(first.Data)),
// // // 	}

// // // 	for _, d := range first.Data {
// // // 		// Mapping DailyAttendance -> DataItem
// // // 		statusItems := make([]*v1.StatusItem, 0, len(d.Status))
// // // 		for _, s := range d.Status {
// // // 			statusItems = append(statusItems, &v1.StatusItem{
// // // 				ClockIn:  s.ClockIn,
// // // 				ClockOut: s.ClockOut,
// // // 			})
// // // 		}

// // // 		resp.Data = append(resp.Data, &v1.DataItem{
// // // 			Tanggal:   d.Tanggal,
// // // 			Karyaname: d.KaryaName,
// // // 			Status:    statusItems,
// // // 			Evaluasi:  d.Evaluasi,
// // // 		})
// // // 	}

// // // 	// // Debug: cek marshal proto
// // // 	// if rawBytes, err := proto.Marshal(resp); err != nil {
// // // 	// 	s.log.Errorf("proto marshal error: %v", err)
// // // 	// } else {
// // // 	// 	s.log.Debugf("PROTO HEX DUMP: %s", hex.EncodeToString(rawBytes))
// // // 	// }

// // // 	return resp, nil
// // // }

// // // // package service

// // // // import (
// // // // 	"context"
// // // // 	"encoding/hex"
// // // // 	v1 "mall-go/api/attendance-raw/service/v1"
// // // // 	"mall-go/module/attendance-raw/service/internal/biz"
// // // // 	"strings"
// // // // 	"time"

// // // // 	"github.com/go-kratos/kratos/v2/log"
// // // // 	"google.golang.org/protobuf/proto"
// // // // )

// // // // type AttendanceRawService struct {
// // // // 	v1.UnimplementedAttendanceServiceServer
// // // // 	log *log.Helper
// // // // 	uc  *biz.AttendanceRawUsecase
// // // // }

// // // // func NewAttendanceRawService(uc *biz.AttendanceRawUsecase, logger log.Logger) *AttendanceRawService {
// // // // 	return &AttendanceRawService{
// // // // 		log: log.NewHelper(log.With(logger, "module", "service/attendanceRaw")),
// // // // 		uc:  uc,
// // // // 	}
// // // // }

// // // // // GetAttendanceReport retrieves the attendance report.
// // // // func (s *AttendanceRawService) GetAttendanceReport(
// // // // 	ctx context.Context,
// // // // 	req *v1.GetAttendanceReportRequest,
// // // // ) (*v1.GetAttendanceReportResponse, error) {

// // // // 	// Ambil data attendance dari usecase (return []*biz.DailyAttendanceReport)
// // // // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.Departcode)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("failed to get attendance report: %v", err)
// // // // 		return nil, err
// // // // 	}

// // // // 	// Kalau tidak ada data, kembalikan response kosong
// // // // 	if len(reports) == 0 {
// // // // 		return &v1.GetAttendanceReportResponse{
// // // // 			Data: []*v1.AttendanceData{},
// // // // 		}, nil
// // // // 	}

// // // // 	// Gunakan report pertama untuk header
// // // // 	first := reports[0]
// // // // 	resp := &v1.GetAttendanceReportResponse{
// // // // 		KodePerusahaan: first.KodePerusahaan,
// // // // 		NamaPerusahaan: first.NamaPerusahaan,
// // // // 		KodeCabang:     first.KodeCabang,
// // // // 		Cabang:         first.Cabang,
// // // // 		Department:     first.Department,
// // // // 		Periode:        first.Periode,
// // // // 		DibuatOleh:     first.DibuatOleh,
// // // // 		DiperiksaOleh:  first.DiperiksaOleh,
// // // // 		DisetujuiOleh:  first.DisetujuiOleh,
// // // // 		Data:           make([]*v1.AttendanceData, 0, len(first.Data)),
// // // // 	}

// // // // 	// Mapping detail harian
// // // // 	for _, d := range first.Data {
// // // // 		// Format tanggal (fallback ke original string jika gagal)
// // // // 		tanggalFormatted := d.Tanggal
// // // // 		if t, err := time.Parse("2006-01-02", d.Tanggal); err == nil {
// // // // 			tanggalFormatted = t.Format("02-Jan-2006")
// // // // 		}

// // // // 		// Format jam (fallback ke original string jika gagal)
// // // // 		jamFormatted := d.Jam
// // // // 		if t, err := time.Parse(time.RFC3339, d.Jam); err == nil {
// // // // 			jamFormatted = t.Format("15:04")
// // // // 		}

// // // // 		// Ambil clock in/out dengan aman
// // // // 		var clockIn, clockOut string
// // // // 		if len(d.Status) > 0 {
// // // // 			clockIn = d.Status[0].ClockIn
// // // // 			clockOut = d.Status[0].ClockOut
// // // // 		}

// // // // 		resp.Data = append(resp.Data, &v1.AttendanceData{
// // // // 			KodeKarya: d.KaryaCode,
// // // // 			Nama:      d.KaryaName,
// // // // 			Tanggal:   tanggalFormatted,
// // // // 			Jam:       jamFormatted,
// // // // 			Status:    &v1.AttStatus{ClockIn: clockIn, ClockOut: clockOut},
// // // // 			Evaluasi:  d.Evaluasi,
// // // // 		})
// // // // 	}

// // // // 	// Debug: cek marshal proto
// // // // 	if rawBytes, err := proto.Marshal(resp); err != nil {
// // // // 		s.log.Errorf("proto marshal error: %v", err)
// // // // 	} else {
// // // // 		s.log.Debugf("PROTO HEX DUMP: %s", hex.EncodeToString(rawBytes))
// // // // 	}

// // // // 	return resp, nil
// // // // }

// // // // func (s *AttendanceRawService) GetAttendanceReport(
// // // // 	ctx context.Context,
// // // // 	req *v1.GetAttendanceReportRequest,
// // // // ) (*v1.GetAttendanceReportResponse, error) {

// // // // 	// Ambil data attendance dari usecase (return []*biz.DailyAttendanceReport)
// // // // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.Departcode)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("failed to get attendance report: %v", err)
// // // // 		return nil, err
// // // // 	}

// // // // 	// Kalau kosong langsung return
// // // // 	if len(reports) == 0 {
// // // // 		return &v1.GetAttendanceReportResponse{
// // // // 			Data: []*v1.AttendanceData{},
// // // // 		}, nil
// // // // 	}

// // // // 	// Gunakan report pertama (sesuai repo hanya ada satu report per periode)
// // // // 	first := reports[0]

// // // // 	// Mapping header
// // // // 	resp := &v1.GetAttendanceReportResponse{
// // // // 		KodePerusahaan: first.KodePerusahaan,
// // // // 		NamaPerusahaan: first.NamaPerusahaan,
// // // // 		KodeCabang:     first.KodeCabang,
// // // // 		Cabang:         first.Cabang,
// // // // 		Department:     first.Department,
// // // // 		Periode:        first.Periode,
// // // // 		DibuatOleh:     first.DibuatOleh,
// // // // 		DiperiksaOleh:  first.DiperiksaOleh,
// // // // 		DisetujuiOleh:  first.DisetujuiOleh,
// // // // 		Data:           make([]*v1.AttendanceData, 0, len(first.Data)),
// // // // 	}
// // // // 	log.Info("INI RESP :", resp)

// // // // 	// Mapping detail harian
// // // // 	for _, d := range first.Data {
// // // // 		// Format tanggal (input YYYY-MM-DD → output DD-MMM-YYYY)
// // // // 		tanggalFormatted := d.Tanggal
// // // // 		if t, err := time.Parse("2006-01-02", d.Tanggal); err == nil {
// // // // 			tanggalFormatted = t.Format("02-Jan-2006")
// // // // 		}

// // // // 		// Format jam (kalau RFC3339)
// // // // 		jamFormatted := d.Jam
// // // // 		if t, err := time.Parse(time.RFC3339, d.Jam); err == nil {
// // // // 			jamFormatted = t.Format("15:04")
// // // // 		}

// // // // 		// Ambil status clock in/out
// // // // 		var clockIn, clockOut string
// // // // 		if len(d.Status) > 0 {
// // // // 			clockIn = d.Status[0].ClockIn
// // // // 			clockOut = d.Status[0].ClockOut
// // // // 		}

// // // // 		// Bersihkan evaluasi dari karakter aneh
// // // // 		evaluasiStr := cleanString(d.Evaluasi)

// // // // 		// Append ke response
// // // // 		resp.Data = append(resp.Data, &v1.AttendanceData{
// // // // 			KodeKarya: d.KaryaCode,
// // // // 			Nama:      d.KaryaName,
// // // // 			Tanggal:   tanggalFormatted,
// // // // 			Jam:       jamFormatted,
// // // // 			Status: &v1.AttStatus{
// // // // 				ClockIn:  clockIn,
// // // // 				ClockOut: clockOut,
// // // // 			},
// // // // 			Evaluasi: evaluasiStr,
// // // // 		})
// // // // 	}

// // // // 	// // Debug HEX dump proto
// // // // 	// if rawBytes, err := proto.Marshal(resp); err != nil {
// // // // 	// 	s.log.Errorf("proto marshal error: %v", err)
// // // // 	// } else {
// // // // 	// 	s.log.Infof("PROTO HEX DUMP: %s", hex.EncodeToString(rawBytes))
// // // // 	// }

// // // // 	return resp, nil
// // // // }

// // // // cleanString hanya menyisakan karakter printable + whitespace umum
// // // func cleanString(s string) string {
// // // 	var b strings.Builder
// // // 	for _, r := range s {
// // // 		if (r >= 32 && r <= 126) || r == '\n' || r == '\r' || r == '\t' {
// // // 			b.WriteRune(r)
// // // 		}
// // // 	}
// // // 	return b.String()
// // // }

// // // // func (s *AttendanceRawService) GetAttendanceReport(
// // // // 	ctx context.Context,
// // // // 	req *v1.GetAttendanceReportRequest,
// // // // ) (*v1.GetAttendanceReportResponse, error) {

// // // // 	// Ambil data attendance dari usecase (return []biz.AttendanceReport)
// // // // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.Departcode)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("failed to get attendance report: %v", err)
// // // // 		return nil, err
// // // // 	}

// // // // 	// Kalau kosong
// // // // 	if len(reports) == 0 {
// // // // 		return &v1.GetAttendanceReportResponse{
// // // // 			Data: make([]*v1.AttendanceData, 0),
// // // // 		}, nil
// // // // 	}

// // // // 	// Ambil report pertama
// // // // 	first := reports[0]

// // // // 	// Mapping header
// // // // 	resp := &v1.GetAttendanceReportResponse{
// // // // 		KodePerusahaan: first.KodePerusahaan,
// // // // 		NamaPerusahaan: first.NamaPerusahaan,
// // // // 		KodeCabang:     first.KodeCabang,
// // // // 		Cabang:         first.Cabang,
// // // // 		Department:     first.Department,
// // // // 		Periode:        first.Periode,
// // // // 		DibuatOleh:     first.DibuatOleh,
// // // // 		DiperiksaOleh:  first.DiperiksaOleh,
// // // // 		DisetujuiOleh:  first.DisetujuiOleh,
// // // // 		Data:           make([]*v1.AttendanceData, 0, len(first.Data)),
// // // // 	}

// // // // 	// Mapping detail harian
// // // // 	for _, d := range first.Data {
// // // // 		// Format tanggal ke "02-Jan-2006"
// // // // 		tanggalFormatted := d.Tanggal
// // // // 		if t, err := time.Parse("2006-01-02", d.Tanggal); err == nil {
// // // // 			tanggalFormatted = t.Format("02-Jan-2006")
// // // // 		}

// // // // 		// Format jam ke "15:04" kalau dari RFC3339
// // // // 		jamFormatted := d.Jam
// // // // 		if t, err := time.Parse(time.RFC3339, d.Jam); err == nil {
// // // // 			jamFormatted = t.Format("15:04")
// // // // 		}

// // // // 		// Mapping status
// // // // 		var clockIn, clockOut string
// // // // 		if len(d.Status) > 0 {
// // // // 			clockIn = d.Status[0].ClockIn
// // // // 			clockOut = d.Status[0].ClockOut
// // // // 		}

// // // // 		status := &v1.AttStatus{
// // // // 			ClockIn:  clockIn,
// // // // 			ClockOut: clockOut,
// // // // 		}

// // // // 		// // Pastikan evaluasi jadi string
// // // // 		// var evaluasiStr string
// // // // 		// switch v := any(d.Evaluasi).(type) {
// // // // 		// case string:
// // // // 		// 	evaluasiStr = v
// // // // 		// case []byte:
// // // // 		// 	evaluasiStr = string(v)
// // // // 		// default:
// // // // 		// 	evaluasiStr = fmt.Sprintf("%v", v)
// // // // 		// }
// // // // 		evaluasiStr := cleanString(d.Evaluasi)
// // // // 		// Append ke response
// // // // 		resp.Data = append(resp.Data, &v1.AttendanceData{
// // // // 			KodeKarya: d.KaryaCode,
// // // // 			Nama:      d.KaryaName,
// // // // 			Tanggal:   tanggalFormatted,
// // // // 			Jam:       jamFormatted,
// // // // 			Status:    status,
// // // // 			Evaluasi:  evaluasiStr,
// // // // 		})
// // // // 	}

// // // // 	// Debug: cek marshal proto
// // // // 	if rawBytes, err := proto.Marshal(resp); err != nil {
// // // // 		s.log.Errorf("proto marshal error: %v", err)
// // // // 	} else {
// // // // 		s.log.Infof("PROTO HEX DUMP: %s", hex.EncodeToString(rawBytes))
// // // // 	}

// // // // 	return resp, nil
// // // // }

// // // // func cleanString(s string) string {
// // // // 	var b strings.Builder
// // // // 	for _, r := range s {
// // // // 		if r >= 32 && r <= 126 || r == '\n' || r == '\r' || r == '\t' {
// // // // 			b.WriteRune(r)
// // // // 		}
// // // // 	}
// // // // 	return b.String()
// // // // }

// // // // func (s *AttendanceRawService) GetAttendanceReport(
// // // // 	ctx context.Context,
// // // // 	req *v1.GetAttendanceReportRequest,
// // // // ) (*v1.GetAttendanceReportResponse, error) {

// // // // 	// Ambil data attendance dari usecase (return []biz.AttendanceReport)
// // // // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.Departcode)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("failed to get attendance report: %v", err)
// // // // 		return nil, err
// // // // 	}

// // // // 	// Kalau kosong
// // // // 	if len(reports) == 0 {
// // // // 		return &v1.GetAttendanceReportResponse{
// // // // 			Data: make([]*v1.AttendanceData, 0),
// // // // 		}, nil
// // // // 	}

// // // // 	// Ambil report pertama
// // // // 	first := reports[0]

// // // // 	// Mapping header
// // // // 	resp := &v1.GetAttendanceReportResponse{
// // // // 		KodePerusahaan: first.KodePerusahaan,
// // // // 		NamaPerusahaan: first.NamaPerusahaan,
// // // // 		KodeCabang:     first.KodeCabang,
// // // // 		Cabang:         first.Cabang,
// // // // 		Department:     first.Department,
// // // // 		Periode:        first.Periode,
// // // // 		DibuatOleh:     first.DibuatOleh,
// // // // 		DiperiksaOleh:  first.DiperiksaOleh,
// // // // 		DisetujuiOleh:  first.DisetujuiOleh,
// // // // 		Data:           make([]*v1.AttendanceData, 0, len(first.Data)),
// // // // 	}

// // // // 	// Mapping detail harian
// // // // 	for _, d := range first.Data {
// // // // 		// Format tanggal ke "02-Jan-2006"
// // // // 		tanggalFormatted := d.Tanggal
// // // // 		if t, err := time.Parse("2006-01-02", d.Tanggal); err == nil {
// // // // 			tanggalFormatted = t.Format("02-Jan-2006")
// // // // 		}

// // // // 		// Format jam ke "15:04" kalau dari RFC3339
// // // // 		jamFormatted := d.Jam
// // // // 		if t, err := time.Parse(time.RFC3339, d.Jam); err == nil {
// // // // 			jamFormatted = t.Format("15:04")
// // // // 		}

// // // // 		// Mapping status
// // // // 		var clockIn, clockOut string
// // // // 		if len(d.Status) > 0 {
// // // // 			clockIn = d.Status[0].ClockIn
// // // // 			clockOut = d.Status[0].ClockOut
// // // // 		}

// // // // 		status := &v1.AttStatus{
// // // // 			ClockIn:  clockIn,
// // // // 			ClockOut: clockOut,
// // // // 		}

// // // // 		// Append ke response
// // // // 		resp.Data = append(resp.Data, &v1.AttendanceData{
// // // // 			KodeKarya: d.KaryaCode,
// // // // 			Nama:      d.KaryaName,
// // // // 			Tanggal:   tanggalFormatted,
// // // // 			Jam:       jamFormatted,
// // // // 			Status:    status,
// // // // 			Evaluasi:  d.Evaluasi,
// // // // 		})
// // // // 	}

// // // // 	// Debug: cek marshal proto
// // // // 	if rawBytes, err := proto.Marshal(resp); err != nil {
// // // // 		s.log.Errorf("proto marshal error: %v", err)
// // // // 	} else {
// // // // 		s.log.Infof("PROTO HEX DUMP: %s", hex.EncodeToString(rawBytes))
// // // // 	}

// // // // 	return resp, nil
// // // // }

// // // // func (s *AttendanceRawService) GetAttendanceReport(
// // // // 	ctx context.Context,
// // // // 	req *v1.GetAttendanceReportRequest,
// // // // ) (*v1.GetAttendanceReportResponse, error) {

// // // // 	// Ambil data attendance dari usecase
// // // // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.Departcode)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("failed to get attendance report: %v", err)
// // // // 		return nil, err
// // // // 	}

// // // // 	// Kalau kosong
// // // // 	if len(reports) == 0 {
// // // // 		return &v1.GetAttendanceReportResponse{
// // // // 			Data: make([]*v1.AttendanceData, 0),
// // // // 		}, nil
// // // // 	}

// // // // 	// Untuk saat ini ambil report pertama
// // // // 	first := reports[0]

// // // // 	// Mapping header laporan
// // // // 	resp := &v1.GetAttendanceReportResponse{
// // // // 		KodePerusahaan: first.KodePerusahaan,
// // // // 		NamaPerusahaan: first.NamaPerusahaan,
// // // // 		KodeCabang:     first.KodeCabang,
// // // // 		Cabang:         first.Cabang,
// // // // 		Department:     first.Department,
// // // // 		Periode:        first.Periode,
// // // // 		DibuatOleh:     first.DibuatOleh,
// // // // 		DiperiksaOleh:  first.DiperiksaOleh,
// // // // 		DisetujuiOleh:  first.DisetujuiOleh,
// // // // 		Data:           make([]*v1.AttendanceData, 0, len(first.Data)),
// // // // 	}

// // // // 	// Mapping detail harian
// // // // 	for _, d := range first.Data {

// // // // 		// Cek apakah ini raw protobuf (misalnya []byte)
// // // // 		var srcData v1.AttendanceData
// // // // 		switch val := any(d).(type) {
// // // // 		case []byte:
// // // // 			if err := proto.Unmarshal(val, &srcData); err != nil {
// // // // 				s.log.Errorf("failed to unmarshal attendance data: %v", err)
// // // // 				continue
// // // // 			}
// // // // 		case v1.AttendanceData:
// // // // 			srcData = val
// // // // 		default:
// // // // 			s.log.Warnf("unexpected data type: %T", d)
// // // // 			continue
// // // // 		}

// // // // 		// Format tanggal
// // // // 		tanggalFormatted := srcData.Tanggal
// // // // 		if t, err := time.Parse("2006-01-02", srcData.Tanggal); err == nil {
// // // // 			tanggalFormatted = t.Format("02-Jan-2006")
// // // // 		}

// // // // 		// Format jam
// // // // 		jamFormatted := srcData.Jam
// // // // 		if t, err := time.Parse(time.RFC3339, srcData.Jam); err == nil {
// // // // 			jamFormatted = t.Format("15:04")
// // // // 		}

// // // // 		// Mapping status
// // // // 		status := &v1.AttStatus{
// // // // 			ClockIn:  srcData.Status.ClockIn,
// // // // 			ClockOut: srcData.Status.ClockOut,
// // // // 		}

// // // // 		resp.Data = append(resp.Data, &v1.AttendanceData{
// // // // 			KodeKarya: srcData.KodeKarya,
// // // // 			Nama:      srcData.Nama,
// // // // 			Tanggal:   tanggalFormatted,
// // // // 			Jam:       jamFormatted,
// // // // 			Status:    status,
// // // // 			Evaluasi:  srcData.Evaluasi,
// // // // 		})
// // // // 	}

// // // // 	// Debug: log hasil proto dalam hex
// // // // 	rawBytes, err := proto.Marshal(resp)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("proto marshal error: %v", err)
// // // // 	} else {
// // // // 		s.log.Infof("PROTO HEX DUMP: %s", hex.EncodeToString(rawBytes))
// // // // 	}

// // // // 	return resp, nil
// // // // }

// // // // func (s *AttendanceRawService) GetAttendanceReport(
// // // // 	ctx context.Context,
// // // // 	req *v1.GetAttendanceReportRequest,
// // // // ) (*v1.GetAttendanceReportResponse, error) {

// // // // 	// Ambil data attendance dari usecase
// // // // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.Departcode)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("failed to get attendance report: %v", err)
// // // // 		return nil, err
// // // // 	}

// // // // 	// Kalau kosong
// // // // 	if len(reports) == 0 {
// // // // 		return &v1.GetAttendanceReportResponse{
// // // // 			Data: make([]*v1.AttendanceData, 0),
// // // // 		}, nil
// // // // 	}

// // // // 	// Ambil report pertama
// // // // 	first := reports[0]

// // // // 	// Mapping header laporan
// // // // 	resp := &v1.GetAttendanceReportResponse{
// // // // 		KodePerusahaan: first.KodePerusahaan,
// // // // 		NamaPerusahaan: first.NamaPerusahaan,
// // // // 		KodeCabang:     first.KodeCabang,
// // // // 		Cabang:         first.Cabang,
// // // // 		Department:     first.Department,
// // // // 		Periode:        first.Periode,
// // // // 		DibuatOleh:     first.DibuatOleh,
// // // // 		DiperiksaOleh:  first.DiperiksaOleh,
// // // // 		DisetujuiOleh:  first.DisetujuiOleh,
// // // // 		Data:           make([]*v1.AttendanceData, 0),
// // // // 	}

// // // // 	// Gunakan map untuk menggabungkan ClockIn dan ClockOut per karyawan + tanggal
// // // // 	attendanceMap := make(map[string]*v1.AttendanceData)

// // // // 	for _, d := range first.Data {
// // // // 		key := d.Tanggal + "_" + d.KaryaCode // unik per tanggal + karyawan

// // // // 		// Format tanggal dari YYYY-MM-DD ke DD-Mon-YYYY
// // // // 		formattedDate := d.Tanggal
// // // // 		if t, err := time.Parse("2006-01-02", d.Tanggal); err == nil {
// // // // 			formattedDate = t.Format("02-Jan-2006")
// // // // 		}

// // // // 		// Kalau belum ada di map, buat entry baru
// // // // 		if _, exists := attendanceMap[key]; !exists {
// // // // 			attendanceMap[key] = &v1.AttendanceData{
// // // // 				KodeKarya: d.KaryaCode,
// // // // 				Nama:      d.KaryaName,
// // // // 				Tanggal:   formattedDate,
// // // // 				Jam:       d.Jam, // bisa diatur mau jam pertama atau terakhir
// // // // 				Status:    &v1.AttStatus{},
// // // // 				Evaluasi:  d.Evaluasi,
// // // // 			}
// // // // 		}

// // // // 		// Merge status ClockIn & ClockOut
// // // // 		for _, st := range d.Status {
// // // // 			if st.ClockIn != "" {
// // // // 				attendanceMap[key].Status.ClockIn = st.ClockIn
// // // // 			}
// // // // 			if st.ClockOut != "" {
// // // // 				attendanceMap[key].Status.ClockOut = st.ClockOut
// // // // 			}
// // // // 		}
// // // // 	}

// // // // 	// Convert map ke slice dan masukkan ke response
// // // // 	for _, v := range attendanceMap {
// // // // 		resp.Data = append(resp.Data, v)
// // // // 	}

// // // // 	// Optional: urutkan berdasarkan tanggal
// // // // 	sort.Slice(resp.Data, func(i, j int) bool {
// // // // 		ti, _ := time.Parse("02-Jan-2006", resp.Data[i].Tanggal)
// // // // 		tj, _ := time.Parse("02-Jan-2006", resp.Data[j].Tanggal)
// // // // 		return ti.Before(tj)
// // // // 	})

// // // // 	return resp, nil
// // // // }

// // // // func (s *AttendanceRawService) GetAttendanceReport(
// // // // 	ctx context.Context,
// // // // 	req *v1.GetAttendanceReportRequest,
// // // // ) (*v1.GetAttendanceReportResponse, error) {

// // // // 	// Ambil data attendance dari usecase
// // // // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.Departcode)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("failed to get attendance report: %v", err)
// // // // 		return nil, err
// // // // 	}

// // // // 	// Kalau kosong
// // // // 	if len(reports) == 0 {
// // // // 		return &v1.GetAttendanceReportResponse{
// // // // 			Data: make([]*v1.AttendanceData, 0),
// // // // 		}, nil
// // // // 	}

// // // // 	// Ambil report pertama
// // // // 	first := reports[0]

// // // // 	// Mapping header
// // // // 	resp := &v1.GetAttendanceReportResponse{
// // // // 		KodePerusahaan: first.KodePerusahaan,
// // // // 		NamaPerusahaan: first.NamaPerusahaan,
// // // // 		KodeCabang:     first.KodeCabang,
// // // // 		Cabang:         first.Cabang,
// // // // 		Department:     first.Department,
// // // // 		Periode:        first.Periode,
// // // // 		DibuatOleh:     first.DibuatOleh,
// // // // 		DiperiksaOleh:  first.DiperiksaOleh,
// // // // 		DisetujuiOleh:  first.DisetujuiOleh,
// // // // 		Data:           make([]*v1.AttendanceData, 0),
// // // // 	}

// // // // 	// Map untuk gabung ClockIn/ClockOut
// // // // 	attendanceMap := make(map[string]*v1.AttendanceData)

// // // // 	for _, d := range first.Data {
// // // // 		key := d.Tanggal + "_" + d.KaryaCode
// // // // 		if _, exists := attendanceMap[key]; !exists {
// // // // 			attendanceMap[key] = &v1.AttendanceData{
// // // // 				KodeKarya: d.KaryaCode,
// // // // 				Nama:      d.KaryaName,
// // // // 				Tanggal:   d.Tanggal,
// // // // 				Status:    &v1.AttStatus{},
// // // // 				Evaluasi:  d.Evaluasi,
// // // // 			}
// // // // 		}

// // // // 		for _, st := range d.Status {
// // // // 			if st.ClockIn != "" {
// // // // 				attendanceMap[key].Status.ClockIn = st.ClockIn
// // // // 			}
// // // // 			if st.ClockOut != "" {
// // // // 				attendanceMap[key].Status.ClockOut = st.ClockOut
// // // // 			}
// // // // 		}

// // // // 		// Set jam sesuai ClockIn kalau ada, kalau tidak ambil ClockOut
// // // // 		if attendanceMap[key].Status.ClockIn != "" {
// // // // 			attendanceMap[key].Jam = attendanceMap[key].Status.ClockIn
// // // // 		} else if attendanceMap[key].Status.ClockOut != "" {
// // // // 			attendanceMap[key].Jam = attendanceMap[key].Status.ClockOut
// // // // 		}
// // // // 	}

// // // // 	// Pindahkan hasil dari map ke slice
// // // // 	for _, v := range attendanceMap {
// // // // 		resp.Data = append(resp.Data, v)
// // // // 	}

// // // // 	// Urutkan berdasarkan tanggal & kode_karya
// // // // 	sort.Slice(resp.Data, func(i, j int) bool {
// // // // 		if resp.Data[i].Tanggal == resp.Data[j].Tanggal {
// // // // 			return resp.Data[i].KodeKarya < resp.Data[j].KodeKarya
// // // // 		}
// // // // 		return resp.Data[i].Tanggal < resp.Data[j].Tanggal
// // // // 	})

// // // // 	// Debug
// // // // 	rawBytes, err := proto.Marshal(resp)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("proto marshal error: %v", err)
// // // // 	} else {
// // // // 		s.log.Infof("PROTO HEX DUMP: %s", hex.EncodeToString(rawBytes))
// // // // 	}

// // // // 	return resp, nil
// // // // }

// // // // func (s *AttendanceRawService) GetAttendanceReport(
// // // // 	ctx context.Context,
// // // // 	req *v1.GetAttendanceReportRequest,
// // // // ) (*v1.GetAttendanceReportResponse, error) {

// // // // 	// Ambil data attendance dari usecase
// // // // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.Departcode)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("failed to get attendance report: %v", err)
// // // // 		return nil, err
// // // // 	}

// // // // 	// Kalau kosong
// // // // 	if len(reports) == 0 {
// // // // 		return &v1.GetAttendanceReportResponse{
// // // // 			Data: make([]*v1.AttendanceData, 0),
// // // // 		}, nil
// // // // 	}

// // // // 	// Untuk saat ini ambil report pertama
// // // // 	first := reports[0]

// // // // 	// Mapping header laporan
// // // // 	resp := &v1.GetAttendanceReportResponse{
// // // // 		KodePerusahaan: first.KodePerusahaan,
// // // // 		NamaPerusahaan: first.NamaPerusahaan,
// // // // 		KodeCabang:     first.KodeCabang,
// // // // 		Cabang:         first.Cabang,
// // // // 		Department:     first.Department,
// // // // 		Periode:        first.Periode,
// // // // 		DibuatOleh:     first.DibuatOleh,
// // // // 		DiperiksaOleh:  first.DiperiksaOleh,
// // // // 		DisetujuiOleh:  first.DisetujuiOleh,
// // // // 		Data:           make([]*v1.AttendanceData, 0, len(first.Data)),
// // // // 	}

// // // // 	// // Mapping detail harian
// // // // 	// for _, d := range first.Data {
// // // // 	// 	// Konversi []biz.AttStatus ke []*v1.AttStatus
// // // // 	// 	statusList := make([]*v1.AttStatus, 0, len(d.Status))
// // // // 	// 	for _, st := range d.Status {
// // // // 	// 		statusList = append(statusList, &v1.AttStatus{
// // // // 	// 			ClockIn:  st.ClockIn,
// // // // 	// 			ClockOut: st.ClockOut,
// // // // 	// 		})
// // // // 	// 	}

// // // // 	// 	resp.Data = append(resp.Data, &v1.AttendanceData{
// // // // 	// 		KodeKarya: d.KaryaCode,
// // // // 	// 		Nama:      d.KaryaName,
// // // // 	// 		Tanggal:   d.Tanggal,
// // // // 	// 		Jam:       d.Jam,
// // // // 	// 		Status:    statusList,
// // // // 	// 		Evaluasi:  d.Evaluasi,
// // // // 	// 	})
// // // // 	// }
// // // // 	// Mapping detail harian
// // // // 	for _, d := range first.Data {
// // // // 		var statusObj *v1.AttStatus
// // // // 		if len(d.Status) > 0 {
// // // // 			st := d.Status[0] // ambil yang pertama
// // // // 			statusObj = &v1.AttStatus{
// // // // 				ClockIn:  st.ClockIn,
// // // // 				ClockOut: st.ClockOut,
// // // // 			}
// // // // 		}

// // // // 		resp.Data = append(resp.Data, &v1.AttendanceData{
// // // // 			KodeKarya: d.KaryaCode,
// // // // 			Nama:      d.KaryaName,
// // // // 			Tanggal:   d.Tanggal,
// // // // 			Jam:       d.Jam,
// // // // 			Status:    statusObj,
// // // // 			Evaluasi:  d.Evaluasi,
// // // // 		})
// // // // 	}
// // // // 	// log.Info("INI RESP", resp)
// // // // 	rawBytes, err := proto.Marshal(resp)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("proto marshal error: %v", err)
// // // // 	} else {
// // // // 		s.log.Infof("PROTO HEX DUMP: %s", hex.EncodeToString(rawBytes))
// // // // 	}
// // // // 	return resp, nil
// // // // }

// // // // func (s *AttendanceRawService) GetAttendanceReport(
// // // // 	ctx context.Context,
// // // // 	req *v1.GetAttendanceReportRequest,
// // // // ) (*v1.GetAttendanceReportResponse, error) {

// // // // 	// Ambil data attendance dari usecase
// // // // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.Departcode)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("failed to get attendance report: %v", err)
// // // // 		return nil, err
// // // // 	}

// // // // 	// Kalau kosong
// // // // 	if len(reports) == 0 {
// // // // 		return &v1.GetAttendanceReportResponse{}, nil
// // // // 	}

// // // // 	first := reports[0]

// // // // 	// Mapping header laporan
// // // // 	resp := &v1.GetAttendanceReportResponse{
// // // // 		KodePerusahaan: first.KodePerusahaan,
// // // // 		NamaPerusahaan: first.NamaPerusahaan,
// // // // 		KodeCabang:     first.KodeCabang,
// // // // 		Cabang:         first.Cabang,
// // // // 		Department:     first.Department,
// // // // 		Periode:        first.Periode,
// // // // 		DibuatOleh:     first.DibuatOleh,
// // // // 		DiperiksaOleh:  first.DiperiksaOleh,
// // // // 		DisetujuiOleh:  first.DisetujuiOleh,
// // // // 		Data:           make([]*v1.AttendanceData, 0, len(first.Data)),
// // // // 	}

// // // // 	// Mapping detail harian
// // // // 	for _, d := range first.Data {
// // // // 		// Konversi []biz.AttStatus ke []*v1.AttStatus
// // // // 		statusList := make([]*v1.AttStatus, 0, len(d.Status))
// // // // 		for _, st := range d.Status {
// // // // 			statusList = append(statusList, &v1.AttStatus{
// // // // 				ClockIn:  st.ClockIn,
// // // // 				ClockOut: st.ClockOut,
// // // // 			})
// // // // 		}

// // // // 		resp.Data = append(resp.Data, &v1.AttendanceData{
// // // // 			KodeKarya: d.KaryaCode,
// // // // 			Nama:      d.KaryaName,
// // // // 			Tanggal:   d.Tanggal,
// // // // 			Jam:       d.Jam,
// // // // 			Status:    statusList,
// // // // 			Evaluasi:  d.Evaluasi,
// // // // 		})
// // // // 	}

// // // // 	return resp, nil
// // // // }

// // // // func (s *AttendanceRawService) GetAttendanceReport(
// // // // 	ctx context.Context,
// // // // 	req *v1.GetAttendanceReportRequest,
// // // // ) (*v1.GetAttendanceReportResponse, error) {

// // // // 	// Ambil data attendance dari usecase
// // // // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.Departcode)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("failed to get attendance report: %v", err)
// // // // 		return nil, err
// // // // 	}

// // // // 	// Kalau kosong
// // // // 	if len(reports) == 0 {
// // // // 		return &v1.GetAttendanceReportResponse{}, nil
// // // // 	}

// // // // 	first := reports[0]

// // // // 	// Mapping header laporan
// // // // 	resp := &v1.GetAttendanceReportResponse{
// // // // 		KodePerusahaan: first.KodePerusahaan,
// // // // 		NamaPerusahaan: first.NamaPerusahaan,
// // // // 		KodeCabang:     first.KodeCabang,
// // // // 		Cabang:         first.Cabang,
// // // // 		Department:     first.Department,
// // // // 		Periode:        first.Periode,
// // // // 		DibuatOleh:     first.DibuatOleh,
// // // // 		DiperiksaOleh:  first.DiperiksaOleh,
// // // // 		DisetujuiOleh:  first.DisetujuiOleh,
// // // // 		Data:           make([]*v1.AttendanceData, 0, len(first.Data)),
// // // // 	}

// // // // 	// Mapping detail harian
// // // // 	for _, d := range first.Data {
// // // // 		resp.Data = append(resp.Data, &v1.AttendanceData{
// // // // 			KodeKarya: d.KaryaCode,
// // // // 			Nama:      d.KaryaName,
// // // // 			Tanggal:   d.Tanggal,
// // // // 			Jam:       "-", // kalau mau jam gabung clockin/clockout, tambahkan logic
// // // // 			Status:    d.Status,
// // // // 			Evaluasi:  d.Evaluasi,
// // // // 		})
// // // // 	}

// // // // 	return resp, nil
// // // // }

// // // // func (s *AttendanceRawService) GetAttendanceReport(
// // // // 	ctx context.Context,
// // // // 	req *v1.GetAttendanceReportRequest,
// // // // ) (*v1.GetAttendanceReportResponse, error) {

// // // // 	// Panggil usecase untuk ambil data attendance
// // // // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.Departcode)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("failed to get attendance report: %v", err)
// // // // 		return nil, err
// // // // 	}

// // // // 	// Kalau data kosong, balikin kosong
// // // // 	if len(reports) == 0 {
// // // // 		return &v1.GetAttendanceReportResponse{}, nil
// // // // 	}

// // // // 	// Ambil report pertama (proto saat ini mewakili 1 laporan)
// // // // 	first := reports[0]

// // // // 	// Mapping ke proto response
// // // // 	resp := &v1.GetAttendanceReportResponse{
// // // // 		KodePerusahaan: first.KodePerusahaan,
// // // // 		NamaPerusahaan: first.NamaPerusahaan,
// // // // 		KodeCabang:     first.KodeCabang,
// // // // 		Cabang:         first.Cabang,
// // // // 		Department:     first.Department,
// // // // 		Periode:        first.Periode,
// // // // 		DibuatOleh:     first.DibuatOleh,
// // // // 		DiperiksaOleh:  first.DiperiksaOleh,
// // // // 		DisetujuiOleh:  first.DisetujuiOleh,
// // // // 		Data:           make([]*v1.AttendanceData, 0, len(first.Data)),
// // // // 	}

// // // // 	// Map data attendance harian
// // // // 	for _, d := range first.Data {
// // // // 		resp.Data = append(resp.Data, &v1.AttendanceData{
// // // // 			KodeKarya:  d.KaryaCode,
// // // // 			Nama:       d.KaryaName,
// // // // 			Department: d.Department,
// // // // 			Tanggal:    d.Tgl,
// // // // 			Hari:       d.Hari,     // nama hari, contoh: "Senin"
// // // // 			ClockIn:    d.ClockIn,  // jam clock-in
// // // // 			ClockOut:   d.ClockOut, // jam clock-out
// // // // 			DeviceIp:   d.DeviceIP, // IP device jika diperlukan
// // // // 		})
// // // // 	}

// // // // 	return resp, nil
// // // // }

// // // // func (s *AttendanceRawService) GetAttendanceReport(
// // // // 	ctx context.Context,
// // // // 	req *v1.GetAttendanceReportRequest,
// // // // ) (*v1.GetAttendanceReportResponse, error) {

// // // // 	// Panggil repo / usecase untuk ambil data attendance
// // // // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.Departcode)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("failed to get attendance report: %v", err)
// // // // 		return nil, err
// // // // 	}

// // // // 	// Kalau data kosong, balikin kosong
// // // // 	if len(reports) == 0 {
// // // // 		return &v1.GetAttendanceReportResponse{}, nil
// // // // 	}

// // // // 	// Ambil salah satu (asumsi response proto ini hanya mewakili satu laporan gabungan)
// // // // 	// Kalau mau banyak, perlu ubah proto jadi repeated
// // // // 	first := reports[0]

// // // // 	// Mapping ke proto response
// // // // 	resp := &v1.GetAttendanceReportResponse{
// // // // 		KodePerusahaan: first.KodePerusahaan,
// // // // 		NamaPerusahaan: first.NamaPerusahaan,
// // // // 		KodeCabang:     first.KodeCabang,
// // // // 		Cabang:         first.Cabang,
// // // // 		Department:     first.Department,
// // // // 		// Jabatan:       first.Jabatan, // kalau field ini belum ada di struct biz, tambahkan
// // // // 		Periode:       first.Periode,
// // // // 		DibuatOleh:    first.DibuatOleh,
// // // // 		DiperiksaOleh: first.DiperiksaOleh,
// // // // 		DisetujuiOleh: first.DisetujuiOleh,
// // // // 		Data:          make([]*v1.AttendanceData, 0, len(first.Data)),
// // // // 	}

// // // // 	// Map data attendance harian
// // // // 	for _, d := range first.Data {
// // // // 		resp.Data = append(resp.Data, &v1.AttendanceData{
// // // // 			// karya_code:  d.KaryaCode,
// // // // 			// KaryaName:  d.KaryaName,
// // // // 			// Department: d.Department,
// // // // 			KodeKarya: d.KaryaCode,
// // // // 			Nama:      d.KaryaName,
// // // // 			// Department: d.Department,
// // // // 			Tanggal:  d.Tgl,
// // // // 			Jam:      d.Hari,
// // // // 			Status:   d.ClockIn,
// // // // 			Evaluasi: d.ClockOut,
// // // // 			// DeviceIp: d.DeviceIP,
// // // // 		})
// // // // 	}

// // // // 	return resp, nil
// // // // }

// // // // func (s *AttendanceRawService) GetAttendanceReport(
// // // // 	ctx context.Context,
// // // // 	req *v1.GetAttendanceReportRequest,
// // // // ) (*v1.GetAttendanceReportResponse, error) {

// // // // 	// Panggil repository
// // // // 	reports, err := s.uc.GetAttendanceReport(ctx, req.StartDate, req.EndDate, req.DepartCode)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("failed to get attendance report: %v", err)
// // // // 		return nil, err
// // // // 	}

// // // // 	// Mapping ke response proto
// // // // 	resp := &v1.GetAttendanceReportResponse{
// // // // 		Reports: make([]*v1.DailyAttendanceReport, 0, len(reports)),
// // // // 	}

// // // // 	for _, r := range reports {
// // // // 		reportItem := &v1.DailyAttendanceReport{
// // // // 			KodePerusahaan: r.KodePerusahaan,
// // // // 			NamaPerusahaan: r.NamaPerusahaan,
// // // // 			KodeCabang:     r.KodeCabang,
// // // // 			Cabang:         r.Cabang,
// // // // 			Nama:           r.Nama,
// // // // 			Periode:        r.Periode,
// // // // 			DibuatOleh:     r.DibuatOleh,
// // // // 			DiperiksaOleh:  r.DiperiksaOleh,
// // // // 			DisetujuiOleh:  r.DisetujuiOleh,
// // // // 			Data:           make([]*v1.DailyAttendance, 0, len(r.Data)),
// // // // 		}

// // // // 		for _, d := range r.Data {
// // // // 			reportItem.Data = append(reportItem.Data, &v1.DailyAttendance{
// // // // 				KaryaCode:  d.KaryaCode,
// // // // 				KaryaName:  d.KaryaName,
// // // // 				Department: d.Department,
// // // // 				Tgl:        d.Tgl,
// // // // 				Hari:       d.Hari,
// // // // 				ClockIn:    d.ClockIn,
// // // // 				ClockOut:   d.ClockOut,
// // // // 				DeviceIP:   d.DeviceIP,
// // // // 			})
// // // // 		}

// // // // 		resp.Reports = append(resp.Reports, reportItem)
// // // // 	}

// // // // 	return resp, nil
// // // // }

// // // // func (s *AttendanceRawService) GetAttendanceReport(ctx context.Context, req *v1.GetAttendanceReportRequest) (*v1.GetAttendanceReportResponse, error) {

// // // // }

// // // // func (s *AttendanceRawService) GetAttendanceRaw(ctx context.Context, req *v1.GetAttendanceRequest) (*v1.GetAttendanceResponse, error) {
// // // // 	if len(req.KodeKarya) == 0 {
// // // // 		return nil, status.Errorf(codes.InvalidArgument, "kode_karya list tidak boleh kosong")
// // // // 	}

// // // // 	results, err := s.uc.GetAttendanceRaw(ctx, req.KodeKarya, req.StartDate, req.EndDate)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("❌ Gagal ambil attendance raw: %v", err)
// // // // 		return nil, status.Errorf(codes.Internal, "gagal mengambil data presensi: %v", err)
// // // // 	}

// // // // 	var attendances []*v1.AttendanceData
// // // // 	for _, r := range results {
// // // // 		attendances = append(attendances, &v1.AttendanceData{
// // // // 			Tanggal: r.Tanggal,
// // // // 			Status:  r.State,
// // // // 			Jam:     r.Jam,
// // // // 			// ClockIn:  r.ClockIn,
// // // // 			// ClockOut: r.ClockOut,
// // // // 			// DeviceIp: r.DeviceIP,
// // // // 		})
// // // // 	}

// // // // 	return &v1.GetAttendanceResponse{
// // // // 		Attendance: attendances,
// // // // 	}, nil
// // // // }

// // // // func (s *AttendanceRawService) GetAttendanceDayli(ctx context.Context, req *v1.GetAttendanceDayliRequest) (*v1.GetAttendanceDayliResponse, error) {
// // // // 	s.log.Infof("📥 Received GetAttendanceDayli request depart_code: %s start_date: %s end_date: %s",
// // // // 		req.DepartCode, req.StartDate, req.EndDate)

// // // // 	if strings.TrimSpace(req.DepartCode) == "" {
// // // // 		return nil, status.Error(codes.InvalidArgument, "depart code tidak boleh kosong")
// // // // 	}

// // // // 	// Ambil data dari use case
// // // // 	reports, err := s.uc.GetAttendanceDayli(ctx, req.DepartCode, req.StartDate, req.EndDate)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("❌ Gagal ambil attendance dayli: %v", err)
// // // // 		return nil, status.Errorf(codes.Internal, "gagal mengambil data presensi: %v", err)
// // // // 	}

// // // // 	var attendances []*v1.AttendanceDayli
// // // // 	for _, report := range reports {
// // // // 		for _, d := range report.Data {
// // // // 			attendances = append(attendances, &v1.AttendanceDayli{
// // // // 				// Data detail harian
// // // // 				KodeKarya:  d.KaryaCode,
// // // // 				KaryaName:  d.KaryaName,
// // // // 				Departname: d.Department,
// // // // 				Tanggal:    d.Tgl,
// // // // 				Hari:       d.Hari,
// // // // 				ClockIn:    d.ClockIn,
// // // // 				ClockOut:   d.ClockOut,
// // // // 				DeviceIp:   d.DeviceIP,

// // // // 				// Data header laporan
// // // // 				KodePerusahaan: report.KodePerusahaan,
// // // // 				NamaPerusahaan: report.NamaPerusahaan,
// // // // 				KodeCabang:     report.KodeCabang,
// // // // 				Cabang:         report.Cabang,
// // // // 				Nama:           report.Nama,
// // // // 				Jabatan:        report.Jabatan,
// // // // 				Periode:        report.Periode,
// // // // 				DibuatOleh:     report.DibuatOleh,
// // // // 				DiperiksaOleh:  report.DiperiksaOleh,
// // // // 				DisetujuiOleh:  report.DisetujuiOleh,
// // // // 			})
// // // // 		}
// // // // 	}

// // // // 	return &v1.GetAttendanceDayliResponse{
// // // // 		Result: attendances,
// // // // 	}, nil
// // // // }

// // // // func (s *AttendanceRawService) GetAttendanceDayli(ctx context.Context, req *v1.GetAttendanceDayliRequest) (*v1.GetAttendanceDayliResponse, error) {
// // // // 	s.log.Infof("Received GetAttendanceDayli request depart_code: %s start_date: %s end_date: %s",
// // // // 		req.DepartCode, req.StartDate, req.EndDate)

// // // // 	if strings.TrimSpace(req.DepartCode) == "" {
// // // // 		return nil, status.Error(codes.InvalidArgument, "depart code tidak boleh kosong")
// // // // 	}

// // // // 	// Ambil data dari use case
// // // // 	reports, err := s.uc.GetAttendanceDayli(ctx, req.DepartCode, req.StartDate, req.EndDate)
// // // // 	if err != nil {
// // // // 		s.log.Errorf("❌ Gagal ambil attendance raw: %v", err)
// // // // 		return nil, status.Errorf(codes.Internal, "gagal mengambil data presensi: %v", err)
// // // // 	}

// // // // 	var attendances []*v1.AttendanceDayli
// // // // 	for _, report := range reports {
// // // // 		for _, d := range report.Data {
// // // // 			attendances = append(attendances, &v1.AttendanceDayli{
// // // // 				KodeKarya:  d.KaryaCode,
// // // // 				KaryaName:  d.KaryaName,
// // // // 				Departname: d.Department,
// // // // 				Tanggal:    d.Tgl,
// // // // 				Hari:       d.Hari,
// // // // 				ClockIn:    d.ClockIn,
// // // // 				ClockOut:   d.ClockOut,
// // // // 				DeviceIp:   d.DeviceIP,
// // // // 				// Meta tambahan bisa dimasukkan jika di proto tersedia
// // // // 				KodePerusahaan: report.KodePerusahaan,
// // // // 				NamaPerusahaan: report.NamaPerusahaan,
// // // // 				KodeCabang:     report.KodeCabang,
// // // // 				Cabang:         report.Cabang,
// // // // 				// Jabatan:        report.Jabatan,
// // // // 				// Periode:        report.Periode,
// // // // 				// DibuatOleh:     report.DibuatOleh,
// // // // 				// DiperiksaOleh:  report.DiperiksaOleh,
// // // // 				// DisetujuiOleh:  report.DisetujuiOleh,
// // // // 			})
// // // // 		}
// // // // 	}

// // // // 	return &v1.GetAttendanceDayliResponse{
// // // // 		Result: attendances,
// // // // 	}, nil
// // // // }
