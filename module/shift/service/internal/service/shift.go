package service

import (
	"context"

	v1 "mall-go/api/shift/service/v1"
	"mall-go/module/shift/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type ShiftService struct {
	v1.UnimplementedShiftServer
	uc  *biz.ShiftUsecase
	log *log.Helper
}

func NewShiftService(uc *biz.ShiftUsecase, logger log.Logger) *ShiftService {
	return &ShiftService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "module", "service/shift")),
	}
}
func (s *ShiftService) CreateShift(ctx context.Context, req *v1.CreateShiftRequest) (*v1.CreateShiftResponse, error) {
	// // Parse start_time and end_time dari string ke time.Time
	// startTime, err := time.Parse("15:04", req.StartTime)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid start_time format: %w", err)
	// }
	// endTime, err := time.Parse("15:04", req.EndTime)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid end_time format: %w", err)
	// }

	// Panggil usecase
	shift, err := s.uc.Create(ctx, &biz.Shift{
		Name:                 req.Name,
		StartTime:            req.StartTime,
		EndTime:              req.EndTime,
		BreakDurationMinutes: req.BreakDurationMinutes,
		CreatedBy:            req.CreatedBy,
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateShiftResponse{
		Shift: &v1.ShiftData{
			Id:                   shift.Id,
			Name:                 shift.Name,
			StartTime:            shift.StartTime,
			EndTime:              shift.EndTime,
			BreakDurationMinutes: shift.BreakDurationMinutes,
			CreatedBy:            shift.CreatedBy,
			CreatedAt:            shift.CreatedAt,
			UpdatedAt:            shift.UpdatedAt,
		},
	}, nil
}
func (s *ShiftService) UpdateShift(ctx context.Context, req *v1.UpdateShiftRequest) (*v1.UpdateShiftResponse, error) {
	// startTime, err := time.Parse("15:04", req.StartTime)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid start_time: %w", err)
	// }
	// endTime, err := time.Parse("15:04", req.EndTime)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid end_time: %w", err)
	// }

	shift, err := s.uc.Update(ctx, &biz.Shift{
		Id:                   req.Id,
		Name:                 req.Name,
		StartTime:            req.StartTime,
		EndTime:              req.EndTime,
		BreakDurationMinutes: req.BreakDurationMinutes,
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateShiftResponse{
		Shift: &v1.ShiftData{
			Id:                   shift.Id,
			Name:                 shift.Name,
			StartTime:            shift.StartTime,
			EndTime:              shift.EndTime,
			BreakDurationMinutes: shift.BreakDurationMinutes,
			CreatedBy:            shift.CreatedBy,
			CreatedAt:            shift.CreatedAt,
			UpdatedAt:            shift.UpdatedAt,
		},
	}, nil
}

func (s *ShiftService) DeleteShift(ctx context.Context, req *v1.DeleteShiftRequest) (*v1.DeleteShiftResponse, error) {
	if err := s.uc.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.DeleteShiftResponse{Success: true}, nil
}

func (s *ShiftService) GetShiftID(ctx context.Context, req *v1.GetShiftIDRequest) (*v1.GetShiftIDResponse, error) {
	shift, err := s.uc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.GetShiftIDResponse{
		Shift: &v1.ShiftData{
			Id:                   shift.Id,
			Name:                 shift.Name,
			StartTime:            shift.StartTime,
			EndTime:              shift.EndTime,
			BreakDurationMinutes: shift.BreakDurationMinutes,
			CreatedBy:            shift.CreatedBy,
			CreatedAt:            shift.CreatedAt,
			UpdatedAt:            shift.UpdatedAt,
		},
	}, nil
}

// func (s *ShiftService) GetShiftDetail(ctx context.Context, req *v1.GetEmployerDetailRequest) (*v1.GetShiftDetailResponse, error) {
// 	// id, _ := strconv.ParseInt(req.Id, 10, 64)
// 	data, err := s.uc.GetDetail(ctx, req.Id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	finger := []*v1.Finger{}
// 	for _, f := range data.Finger {
// 		finger = append(finger, &v1.Finger{
// 			Fingercode: f.Fingercode,
// 			Finger0:    f.Finger0,
// 			Finger1:    f.Finger1,
// 			Finger2:    f.Finger2,
// 			Finger3:    f.Finger3,
// 			Finger4:    f.Finger4,
// 			Finger5:    f.Finger5,
// 			Finger6:    f.Finger6,
// 			Finger7:    f.Finger7,
// 			Finger8:    f.Finger8,
// 			Finger9:    f.Finger9,
// 		})
// 	}
// 	dept := []*v1.Department{}
// 	for _, d := range data.Department {
// 		dept = append(dept, &v1.Department{
// 			Departcode: d.DepartCode,
// 			Departname: d.DepartName,
// 		})
// 	}
// 	// dept := []*Depv1.GetDepartmentCodeResponse{}
// 	// for _, d := range data.Department {
// 	// 	dept = append(dept, &Depv1.GetDepartmentCodeResponse{
// 	// 		DepartCode: d.DepartCode,
// 	// 		DepartName: d.DepartName,
// 	// 	})
// 	// }

// 	return &v1.GetShiftDetailResponse{
// 		Id:         data.Id,
// 		Nosap:      data.NoSap,
// 		Nip:        data.Nip,
// 		KaryaCode:  data.KaryaCode,
// 		KaryaName:  data.KaryaName,
// 		DispName:   data.DispName,
// 		PassMesin:  data.PassMesin,
// 		RfidCard:   data.RFIDCard,
// 		Status:     data.Status,
// 		Finger:     finger,
// 		Department: dept,
// 	}, nil
// }

func (s *ShiftService) ListShift(ctx context.Context, req *v1.ListShiftRequest) (*v1.ListShiftResponse, error) {
	shifts, total, err := s.uc.List(ctx, int64(req.Pn), int64(req.PSize))
	if err != nil {
		return nil, err
	}

	shiftData := make([]*v1.ShiftData, 0, len(shifts))
	for _, s := range shifts {
		shiftData = append(shiftData, &v1.ShiftData{
			Id:                   s.Id,
			Name:                 s.Name,
			StartTime:            s.StartTime,
			EndTime:              s.EndTime,
			BreakDurationMinutes: s.BreakDurationMinutes,
			CreatedBy:            s.CreatedBy,
			UpdatedBy:            s.UpdatedBy,
			CreatedAt:            s.CreatedAt,
			UpdatedAt:            s.UpdatedAt,
		})
	}

	return &v1.ListShiftResponse{
		Total:  int32(total),
		Shifts: shiftData,
	}, nil
}

// func (s *ShiftService) ListShift(ctx context.Context, req *v1.ListShiftRequest) (*v1.ListShiftResponse, error) {
// 	shfs, total, err := s.uc.List(ctx, int64(req.Pn), int64(req.PSize))
// 	if err != nil {
// 		return nil, err
// 	}

// 	records := make([]*v1.ShiftData, 0, len(shfs))
// 	for _, shf := range shfs {
// 		// fingers := make([]*v1.Finger, 0, len(emp.Finger))
// 		// for _, f := range emp.Finger {
// 		// 	fingers = append(fingers, &v1.Finger{
// 		// 		Fingercode: f.Fingercode,
// 		// 		Finger0:    f.Finger0,
// 		// 		Finger1:    f.Finger1,
// 		// 		Finger2:    f.Finger2,
// 		// 		Finger3:    f.Finger3,
// 		// 		Finger4:    f.Finger4,
// 		// 		Finger5:    f.Finger5,
// 		// 		Finger6:    f.Finger6,
// 		// 		Finger7:    f.Finger7,
// 		// 		Finger8:    f.Finger8,
// 		// 		Finger9:    f.Finger9,
// 		// 	})
// 		// }

// 		// departments := make([]*v1.Department, 0, len(emp.Department))
// 		// for _, d := range emp.Department {
// 		// 	departments = append(departments, &v1.Department{
// 		// 		Departcode: d.DepartCode,
// 		// 		Departname: d.DepartName,
// 		// 	})
// 		// }

// 		records = append(records, &v1.ShiftData{
// 			Id:                   shf.Id,
// 			Name:                 shf.Name,
// 			StartTime:            shf.StartTime,
// 			EndTime:              shf.EndTime,
// 			BreakDurationMinutes: shf.BreakDurationMinutes,
// 			CreatedBy:            shf.CreatedBy,
// 			UpdatedBy:            shf.UpdatedBy,
// 			CreatedAt:            shf.CreatedAt,
// 			UpdatedAt:            shf.UpdatedAt,
// 		})
// 	}

// 	return &v1.ListShiftResponse{
// 		Total: int32(total),
// 		Shift: records,
// 	}, nil
// }
