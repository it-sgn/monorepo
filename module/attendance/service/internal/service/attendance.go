package service

// import (
// 	"context"
// 	v1 "mall-go/api/attendance/service/v1"
// 	"mall-go/module/attendance/service/internal/biz"

// 	"github.com/go-kratos/kratos/v2/log"
// )

// type AttendanceService struct {
// 	v1.UnimplementedAttendanceServer
// 	uc  *biz.AttendanceUsecase
// 	log *log.Helper
// }

// func NewAttendanceService(uc *biz.AttendanceUsecase, logger log.Logger) *AttendanceService {
// 	return &AttendanceService{
// 		uc:  uc,
// 		log: log.NewHelper(log.With(logger, "module", "service/attendance")),
// 	}
// }
// func (s *AttendanceService) CreateClockIn(ctx context.Context, req *v1.CreateAttendanceRequest) (*v1.CreateAttendanceResponse, error) {

// }

// func (s *AttendanceService) CreateAttendance(ctx context.Context, req *v1.CreateAttendanceRequest) (*v1.CreateAttendanceResponse, error) {
// 	emp, err := s.uc.Create(ctx, &biz.Attendance{

// 		Status: req.Status,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &v1.CreateAttendanceResponse{
// 		// Id: strconv.FormatInt(emp.Id, 10),
// 		Id: emp.Id,
// 	}, nil
// }

// func (s *AttendanceService) UpdateAttendance(ctx context.Context, req *v1.UpdateAttendanceRequest) (*v1.UpdateAttendanceResponse, error) {
// 	// id, _ := strconv.ParseInt(req.Id, 10, 64)
// 	emp, err := s.uc.Update(ctx, &biz.Attendance{
// 		Id: req.Id,

// 		Status: req.Status,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &v1.UpdateAttendanceResponse{
// 		// Id: strconv.FormatInt(emp.Id, 10),
// 		Id: emp.Id,
// 	}, nil
// }

// func (s *AttendanceService) DeleteAttendance(ctx context.Context, req *v1.DeleteAttendanceRequest) (*v1.DeleteAttendanceResponse, error) {
// 	if err := s.uc.Delete(ctx, req.Id); err != nil {
// 		return nil, err
// 	}
// 	return &v1.DeleteAttendanceResponse{Success: true}, nil
// }

// func (s *AttendanceService) GetAttendanceID(ctx context.Context, req *v1.GetAttendanceIDRequest) (*v1.GetAttendanceIDResponse, error) {
// 	emp, err := s.uc.GetByID(ctx, req.Id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &v1.GetAttendanceIDResponse{
// 		Id: emp.Id,

// 		Status: emp.Status,
// 	}, nil
// }

// // func (s *AttendanceService) GetAttendanceDetail(ctx context.Context, req *v1.GetEmployerDetailRequest) (*v1.GetAttendanceDetailResponse, error) {
// // 	// id, _ := strconv.ParseInt(req.Id, 10, 64)
// // 	data, err := s.uc.GetDetail(ctx, req.Id)
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	finger := []*v1.Finger{}
// // 	for _, f := range data.Finger {
// // 		finger = append(finger, &v1.Finger{
// // 			Fingercode: f.Fingercode,
// // 			Finger0:    f.Finger0,
// // 			Finger1:    f.Finger1,
// // 			Finger2:    f.Finger2,
// // 			Finger3:    f.Finger3,
// // 			Finger4:    f.Finger4,
// // 			Finger5:    f.Finger5,
// // 			Finger6:    f.Finger6,
// // 			Finger7:    f.Finger7,
// // 			Finger8:    f.Finger8,
// // 			Finger9:    f.Finger9,
// // 		})
// // 	}
// // 	dept := []*v1.Department{}
// // 	for _, d := range data.Department {
// // 		dept = append(dept, &v1.Department{
// // 			Departcode: d.DepartCode,
// // 			Departname: d.DepartName,
// // 		})
// // 	}
// // 	// dept := []*Depv1.GetDepartmentCodeResponse{}
// // 	// for _, d := range data.Department {
// // 	// 	dept = append(dept, &Depv1.GetDepartmentCodeResponse{
// // 	// 		DepartCode: d.DepartCode,
// // 	// 		DepartName: d.DepartName,
// // 	// 	})
// // 	// }

// // 	return &v1.GetAttendanceDetailResponse{
// // 		Id:         data.Id,
// // 		Nosap:      data.NoSap,
// // 		Nip:        data.Nip,
// // 		KaryaCode:  data.KaryaCode,
// // 		KaryaName:  data.KaryaName,
// // 		DispName:   data.DispName,
// // 		PassMesin:  data.PassMesin,
// // 		RfidCard:   data.RFIDCard,
// // 		Status:     data.Status,
// // 		Finger:     finger,
// // 		Department: dept,
// // 	}, nil
// // }

// // func (s *AttendanceService) ListAttendance(ctx context.Context, req *v1.ListAttendanceRequest) (*v1.ListAttendanceReply, error) {
// // 	emps, total, err := s.uc.List(ctx, int64(req.Pn), int64(req.PSize))
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	records := make([]*v1.AttendanceRecord, 0, len(emps))
// // 	for _, emp := range emps {
// // 		fingers := make([]*v1.Finger, 0, len(emp.Finger))
// // 		for _, f := range emp.Finger {
// // 			fingers = append(fingers, &v1.Finger{
// // 				Fingercode: f.Fingercode,
// // 				Finger0:    f.Finger0,
// // 				Finger1:    f.Finger1,
// // 				Finger2:    f.Finger2,
// // 				Finger3:    f.Finger3,
// // 				Finger4:    f.Finger4,
// // 				Finger5:    f.Finger5,
// // 				Finger6:    f.Finger6,
// // 				Finger7:    f.Finger7,
// // 				Finger8:    f.Finger8,
// // 				Finger9:    f.Finger9,
// // 			})
// // 		}

// // 		departments := make([]*v1.Department, 0, len(emp.Department))
// // 		for _, d := range emp.Department {
// // 			departments = append(departments, &v1.Department{
// // 				Departcode: d.DepartCode,
// // 				Departname: d.DepartName,
// // 			})
// // 		}

// // 		records = append(records, &v1.AttendanceRecord{
// // 			Id:         emp.Id,
// // 			Nosap:      emp.NoSap,
// // 			Nip:        emp.Nip,
// // 			KaryaCode:  emp.KaryaCode,
// // 			KaryaName:  emp.KaryaName,
// // 			DispName:   emp.DispName,
// // 			PassMesin:  emp.PassMesin,
// // 			RfidCard:   emp.RFIDCard,
// // 			Finger:     fingers,
// // 			Department: departments,
// // 			Status:     emp.Status,
// // 		})
// // 	}

// // 	return &v1.ListAttendanceReply{
// // 		Total:     int32(total),
// // 		Attendance: records,
// // 	}, nil
// // }
