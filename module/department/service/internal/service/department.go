package service

import (
	"context"
	// v1 "mall-go/api/service/employers/service"
	// v1 "mall-go/api/service/employers/service"

	v1 "mall-go/api/department/service/v1"
	"mall-go/module/department/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type DepartmentService struct {
	v1.UnimplementedDepartmentServer

	bc  *biz.DepartmentUsecase
	log *log.Helper
}

func NewDepartmentService(bc *biz.DepartmentUsecase, logger log.Logger) *DepartmentService {
	return &DepartmentService{
		bc:  bc,
		log: log.NewHelper(log.With(logger, "module", "service/employers"))}
}

// func (s *EmployerService) CreateUser(ctx context.Context, in *CreateUserRequest) (*UserVO, error) {
// return &UserVO{}, nil
// }
func (s *DepartmentService) CreateDepartment(ctx context.Context, in *v1.CreateDepartmentRequest) (*v1.CreateDepartmentResponse, error) {
	b := &biz.Department{
		DepartName: in.DepartName,
		Status:     in.Status,
		Ket:        in.Ket,
	}

	x, err := s.bc.Create(ctx, b)
	if err != nil {
		return nil, err
	}

	return &v1.CreateDepartmentResponse{
		Id:         x.Id,
		DepartName: x.DepartName,
		Status:     x.Status,
		Ket:        x.Status,
	}, nil
}

func (s *DepartmentService) UpdateDepartment(ctx context.Context, in *v1.UpdateDepartmentRequest) (*v1.UpdateDepartmentResponse, error) {
	b := &biz.Department{
		Id:         in.Id,
		DepartName: in.DepartName,
		Status:     in.Status,
		Ket:        in.Ket,
	}
	x, err := s.bc.Update(ctx, b)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateDepartmentResponse{
		Id:         x.Id,
		DepartName: x.DepartName,
		Status:     x.Status,
		Ket:        x.Ket,
	}, nil
}
func (s *DepartmentService) GetDepartmentID(ctx context.Context, in *v1.GetDepartmentIDRequest) (*v1.GetDepartmentIDResponse, error) {
	x, err := s.bc.GetByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &v1.GetDepartmentIDResponse{
		Id:         x.Id,
		DepartName: x.DepartName,
		Status:     x.Status,
		Ket:        x.Ket,
	}, nil
}

func (s *DepartmentService) GetDepartmentCode(ctx context.Context, in *v1.GetDepartmentCodeRequest) (*v1.GetDepartmentCodeResponse, error) {
	x, err := s.bc.GetByCode(ctx, in.DepartCode)
	if err != nil {
		return nil, err
	}
	return &v1.GetDepartmentCodeResponse{
		Id:         x.Id,
		DepartCode: x.DepartCode,
		DepartName: x.DepartName,
		Status:     x.Status,
	}, nil
}
func (s *DepartmentService) DeleteDepartment(ctx context.Context, in *v1.DeleteDepartmentRequest) (*v1.DeleteDepartmentResponse, error) {
	err := s.bc.Delete(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteDepartmentResponse{Success: true}, nil
}

func (s *DepartmentService) ListDepartment(ctx context.Context, in *v1.ListDepartmentRequest) (*v1.ListDepartmentResponse, error) {
	list, total, err := s.bc.List(ctx, int64(in.Pn), int64(in.PSize))
	if err != nil {
		return nil, err
	}

	var employers []*v1.DepartmentRecord
	for _, x := range list {
		employers = append(employers, &v1.DepartmentRecord{
			Id:         x.Id,
			DepartName: x.DepartName,
			Status:     x.Status,
			Ket:        x.Ket,
		})
	}

	return &v1.ListDepartmentResponse{
		Total:      int32(total),
		Department: employers,
	}, nil
}
func (s *DepartmentService) ListDepartmentNext(ctx context.Context, in *v1.ListDepartmentTokenReq) (*v1.ListDepartmentResponseNextToken, error) {
	list, nextToken, err := s.bc.ListNext(ctx, in.PageSize, in.PageToken)
	if err != nil {
		return nil, err
	}

	var employers []*v1.DepartmentRecord
	for _, x := range list {
		employers = append(employers, &v1.DepartmentRecord{
			Id:         x.Id,
			DepartName: x.DepartName,
			Status:     x.Status,
			Ket:        x.Ket,
		})
	}

	return &v1.ListDepartmentResponseNextToken{
		Department: employers,
		NextToken:  nextToken,
	}, nil
}
