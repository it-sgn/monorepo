package service

import (
	"context"

	v1 "mall-go/api/employers/service/v1"
	"mall-go/module/employers/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type EmployersService struct {
	v1.UnimplementedEmployersServer
	uc  *biz.EmployersUsecase
	log *log.Helper
}

func NewEmployersService(uc *biz.EmployersUsecase, logger log.Logger) *EmployersService {
	return &EmployersService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "module", "service/employers")),
	}
}

func (s *EmployersService) CreateEmployers(ctx context.Context, req *v1.CreateEmployersRequest) (*v1.CreateEmployersReply, error) {
	emp, err := s.uc.Create(ctx, &biz.Employers{
		NoSap:      req.Nosap,
		Nip:        req.Nip,
		KaryaCode:  req.KaryaCode,
		KaryaName:  req.KaryaName,
		DispName:   req.DispName,
		PassMesin:  req.PassMesin,
		RFIDCard:   req.RfidCard,
		Finger:     req.KodeFinger,
		Department: req.Departcode,
		Status:     req.Status,
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateEmployersReply{
		// Id: strconv.FormatInt(emp.Id, 10),
		Id: emp.Id,
	}, nil
}

func (s *EmployersService) UpdateEmployers(ctx context.Context, req *v1.UpdateEmployersRequest) (*v1.UpdateEmployersReply, error) {
	// id, _ := strconv.ParseInt(req.Id, 10, 64)
	emp, err := s.uc.Update(ctx, &biz.Employers{
		Id:         req.Id,
		NoSap:      req.Nosap,
		Nip:        req.Nip,
		KaryaCode:  req.KaryaCode,
		KaryaName:  req.KaryaName,
		DispName:   req.DispName,
		PassMesin:  req.PassMesin,
		RFIDCard:   req.RfidCard,
		Finger:     req.KodeFinger,
		Department: req.Departcode,
		Status:     req.Status,
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateEmployersReply{
		// Id: strconv.FormatInt(emp.Id, 10),
		Id: emp.Id,
	}, nil
}

func (s *EmployersService) DeleteEmployers(ctx context.Context, req *v1.DeleteEmployersRequest) (*v1.DeleteEmployersReply, error) {
	if err := s.uc.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.DeleteEmployersReply{Success: true}, nil
}

func (s *EmployersService) GetEmployersID(ctx context.Context, req *v1.GetEmployersIDRequest) (*v1.GetEmployersIDReply, error) {
	emp, err := s.uc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.GetEmployersIDReply{
		Id:         emp.Id,
		Nosap:      emp.NoSap,
		Nip:        emp.Nip,
		KaryaCode:  emp.KaryaCode,
		KaryaName:  emp.KaryaName,
		DispName:   emp.DispName,
		PassMesin:  emp.PassMesin,
		RfidCard:   emp.RFIDCard,
		KodeFinger: emp.Finger,
		Departcode: emp.Department,
		Status:     emp.Status,
	}, nil
}

func (s *EmployersService) GetEmployersDetail(ctx context.Context, req *v1.GetEmployerDetailRequest) (*v1.GetEmployersDetailResponse, error) {
	// id, _ := strconv.ParseInt(req.Id, 10, 64)
	data, err := s.uc.GetDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	finger := []*v1.Finger{}
	for _, f := range data.Finger {
		finger = append(finger, &v1.Finger{
			Fingercode: f.Fingercode,
			Finger0:    f.Finger0,
			Finger1:    f.Finger1,
			Finger2:    f.Finger2,
			Finger3:    f.Finger3,
			Finger4:    f.Finger4,
			Finger5:    f.Finger5,
			Finger6:    f.Finger6,
			Finger7:    f.Finger7,
			Finger8:    f.Finger8,
			Finger9:    f.Finger9,
		})
	}
	dept := []*v1.Department{}
	for _, d := range data.Department {
		dept = append(dept, &v1.Department{
			Departcode: d.DepartCode,
			Departname: d.DepartName,
		})
	}
	// dept := []*Depv1.GetDepartmentCodeResponse{}
	// for _, d := range data.Department {
	// 	dept = append(dept, &Depv1.GetDepartmentCodeResponse{
	// 		DepartCode: d.DepartCode,
	// 		DepartName: d.DepartName,
	// 	})
	// }

	return &v1.GetEmployersDetailResponse{
		Id:         data.Id,
		Nosap:      data.NoSap,
		Nip:        data.Nip,
		KaryaCode:  data.KaryaCode,
		KaryaName:  data.KaryaName,
		DispName:   data.DispName,
		PassMesin:  data.PassMesin,
		RfidCard:   data.RFIDCard,
		Status:     data.Status,
		Finger:     finger,
		Department: dept,
	}, nil
}

func (s *EmployersService) ListEmployers(ctx context.Context, req *v1.ListEmployersRequest) (*v1.ListEmployersReply, error) {
	emps, total, err := s.uc.List(ctx, int64(req.Pn), int64(req.PSize))
	if err != nil {
		return nil, err
	}

	records := make([]*v1.EmployersRecord, 0, len(emps))
	for _, emp := range emps {
		fingers := make([]*v1.Finger, 0, len(emp.Finger))
		for _, f := range emp.Finger {
			fingers = append(fingers, &v1.Finger{
				Fingercode: f.Fingercode,
				Finger0:    f.Finger0,
				Finger1:    f.Finger1,
				Finger2:    f.Finger2,
				Finger3:    f.Finger3,
				Finger4:    f.Finger4,
				Finger5:    f.Finger5,
				Finger6:    f.Finger6,
				Finger7:    f.Finger7,
				Finger8:    f.Finger8,
				Finger9:    f.Finger9,
			})
		}

		departments := make([]*v1.Department, 0, len(emp.Department))
		for _, d := range emp.Department {
			departments = append(departments, &v1.Department{
				Departcode: d.DepartCode,
				Departname: d.DepartName,
			})
		}

		records = append(records, &v1.EmployersRecord{
			Id:         emp.Id,
			Nosap:      emp.NoSap,
			Nip:        emp.Nip,
			KaryaCode:  emp.KaryaCode,
			KaryaName:  emp.KaryaName,
			DispName:   emp.DispName,
			PassMesin:  emp.PassMesin,
			RfidCard:   emp.RFIDCard,
			Finger:     fingers,
			Department: departments,
			Status:     emp.Status,
		})
	}

	return &v1.ListEmployersReply{
		Total:     int32(total),
		Employers: records,
	}, nil
}
