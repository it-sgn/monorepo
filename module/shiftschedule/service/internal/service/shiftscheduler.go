package service

import (
	"context"
	"fmt"
	"time"

	v1 "mall-go/api/shiftschedule/service/v1"
	"mall-go/module/shiftschedule/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type ShiftScheduleService struct {
	v1.UnimplementedShiftScheduleServer
	uc  *biz.ShiftScheduleUsecase
	log *log.Helper
}

func NewShiftScheduleService(uc *biz.ShiftScheduleUsecase, logger log.Logger) *ShiftScheduleService {
	return &ShiftScheduleService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "module", "service/ShiftSchedule")),
	}
}

// func (s *ShiftScheduleService) CreateShiftSchedule(ctx context.Context, req *v1.CreateShiftScheduleRequest) (*v1.CreateShiftScheduleResponse, error) {
// 	tanggal, _ := time.Parse("2006-01-02", req.Tanggal)

// 	sch, err := s.uc.Create(ctx, &biz.ShiftSchedule{
// 		// ScheduleKode: req.s,
// 		KaryaCode:  req.Karyacode,
// 		Tanggal:    tanggal,
// 		DepartCode: req.DepartCode,
// 		// CreatedBy: rq,

// 	})
// 	if err != nil {
// 		return nil, err
// 	}

//		return &v1.CreateShiftScheduleResponse{
//			// Id: strconv.FormatInt(emp.Id, 10),
//			Schedule: sch.Id,
//		}, nil
//	}
func (s *ShiftScheduleService) CreateShiftSchedule(ctx context.Context, req *v1.CreateShiftScheduleRequest) (*v1.CreateShiftScheduleResponse, error) {
	tanggal, err := time.Parse("2006-01-02", req.Tanggal)
	if err != nil {
		return nil, fmt.Errorf("invalid tanggal format: %w", err)
	}

	// Buat entitas domain
	sch := &biz.ShiftSchedule{
		ScheduleCode: req.ScheduleCode,
		KaryaCode:    req.Karyacode,
		Tanggal:      tanggal,
		DepartCode:   req.DepartCode,
		ShiftID:      req.ShiftId,
		CreatedBy:    "system", // bisa diganti dengan user dari context
	}

	// Simpan ke database via usecase
	created, err := s.uc.Create(ctx, sch)
	if err != nil {
		return nil, err
	}

	// Mapping ke response
	return &v1.CreateShiftScheduleResponse{
		Schedule: &v1.ShiftScheduleData{
			Id:           created.Id,
			ScheduleCode: created.ScheduleCode,
			Karyacode:    created.KaryaCode,
			Tanggal:      created.Tanggal.Format("2006-01-02"),
			DepartCode:   created.DepartCode,
			ShiftId:      created.ShiftID,
		},
	}, nil
}
func (s *ShiftScheduleService) UpdateShiftSchedule(ctx context.Context, req *v1.UpdateShiftScheduleRequest) (*v1.UpdateShiftScheduleResponse, error) {
	tanggal, err := time.Parse("2006-01-02", req.Tanggal)
	if err != nil {
		return nil, fmt.Errorf("invalid tanggal format: %w", err)
	}

	updated, err := s.uc.Update(ctx, &biz.ShiftSchedule{
		Id:           req.Id,
		ScheduleCode: req.ScheduleCode,
		KaryaCode:    req.Karyacode,
		Tanggal:      tanggal,
		DepartCode:   req.DepartCode,
		ShiftID:      req.ShiftId,
		CreatedBy:    "system", // atau dari context JWT
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateShiftScheduleResponse{
		Schedule: &v1.ShiftScheduleData{
			Id:           updated.Id,
			ScheduleCode: updated.ScheduleCode,
			Karyacode:    updated.KaryaCode,
			Tanggal:      updated.Tanggal.Format("2006-01-02"),
			DepartCode:   updated.DepartCode,
			ShiftId:      updated.ShiftID,
		},
	}, nil
}

func (s *ShiftScheduleService) DeleteShiftSchedule(ctx context.Context, req *v1.DeleteShiftScheduleRequest) (*v1.DeleteShiftScheduleResponse, error) {
	if err := s.uc.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.DeleteShiftScheduleResponse{Success: true}, nil
}
func (s *ShiftScheduleService) GetShiftScheduleID(ctx context.Context, req *v1.GetShiftScheduleRequest) (*v1.GetShiftScheduleResponse, error) {
	ss, err := s.uc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.GetShiftScheduleResponse{
		Schedule: &v1.ShiftScheduleData{
			Id:           ss.Id,
			ScheduleCode: ss.ScheduleCode,
			Karyacode:    ss.KaryaCode,
			Tanggal:      ss.Tanggal.Format("2006-01-02"),
			DepartCode:   ss.DepartCode,
			ShiftId:      ss.ShiftID,
		},
	}, nil
}

func parseTanggal(tanggalStr string) (time.Time, error) {
	return time.Parse("2006-01-02", tanggalStr) // format date-only
}
