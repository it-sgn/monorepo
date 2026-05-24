package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	v1 "mall-go/api/organization/service/v1"
	"mall-go/module/organization/service/internal/biz"
	"mall-go/module/organization/service/internal/data"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrganizationService struct {
	v1.UnimplementedOrganizationServiceServer

	bc   *biz.PositionUsecase
	rp   *biz.AssignmentUsecase
	pr   *biz.PerusahaanUsecase
	log  *log.Helper
	data *data.Data // ← Tambahkan ini
}

func NewOrganizationService(bc *biz.PositionUsecase, rp *biz.AssignmentUsecase, pr *biz.PerusahaanUsecase, data *data.Data, logger log.Logger) *OrganizationService {
	return &OrganizationService{
		bc:   bc,
		rp:   rp,
		pr:   pr,
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (s *OrganizationService) CreatePosition(ctx context.Context, in *v1.CreatePositionRequest) (*v1.PositionResponse, error) {
	b := &v1.CreatePositionRequest{
		PositionCode:        in.PositionCode,
		Name:                in.Name,
		RoleName:            in.RoleName,
		DepartmentCode:      in.DepartmentCode,
		ReportsToPositionId: in.ReportsToPositionId,
	}

	x, err := s.bc.Create(ctx, b)
	if err != nil {
		return nil, err
	}
	if err := s.data.Kafka.PublishEvent(ctx, &biz.Event{
		Type: "create",
		Key:  x.Position.PositionCode,
		Data: fmt.Sprintf(`{"position":"%s","name":"%s"}`, x.Position.PositionCode, x.Position.Name),
	}); err != nil {
		s.log.Warnf("failed to publish event: %v", err)
	}
	return &v1.PositionResponse{
		Position: &v1.Position{
			Id:                  x.Position.Id,
			PositionCode:        x.Position.PositionCode,
			Name:                x.Position.Name,
			RoleName:            x.Position.RoleName,
			DepartmentCode:      x.Position.DepartmentCode,
			ReportsToPositionId: x.Position.ReportsToPositionId,
			CreatedAt:           x.Position.CreatedAt,
			UpdatedAt:           x.Position.UpdatedAt,
			CreatedBy:           x.Position.CreatedBy,
			UpdatedBy:           x.Position.UpdatedBy,
		},
	}, nil
}

func (s *OrganizationService) UpdatePosition(ctx context.Context, in *v1.UpdatePositionRequest) (*v1.PositionResponse, error) {
	b := &v1.Position{
		PositionCode:        in.PositionCode,
		Name:                in.Name,
		RoleName:            in.RoleName,
		DepartmentCode:      in.DepartmentCode,
		ReportsToPositionId: in.ReportsToPositionId,
	}
	x, err := s.bc.Update(ctx, b)
	if err != nil {
		return nil, err
	}
	return &v1.PositionResponse{
		Position: &v1.Position{
			Id:                  x.Position.Id,
			PositionCode:        x.Position.PositionCode,
			Name:                x.Position.Name,
			RoleName:            x.Position.RoleName,
			DepartmentCode:      x.Position.DepartmentCode,
			ReportsToPositionId: x.Position.ReportsToPositionId,
			CreatedAt:           x.Position.CreatedAt,
			UpdatedAt:           x.Position.UpdatedAt,
			CreatedBy:           x.Position.CreatedBy,
			UpdatedBy:           x.Position.UpdatedBy,
		}}, nil
}
func (s *OrganizationService) GetPositionID(ctx context.Context, in *v1.GetPositionIDRequest) (*v1.PositionResponse, error) {
	x, err := s.bc.Get(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &v1.PositionResponse{
		Position: &v1.Position{
			Id:                  x.Position.Id,
			PositionCode:        x.Position.PositionCode,
			Name:                x.Position.Name,
			RoleName:            x.Position.RoleName,
			DepartmentCode:      x.Position.DepartmentCode,
			ReportsToPositionId: x.Position.ReportsToPositionId,
			CreatedAt:           x.Position.CreatedAt,
			UpdatedAt:           x.Position.UpdatedAt,
			CreatedBy:           x.Position.CreatedBy,
			UpdatedBy:           x.Position.UpdatedBy,
		},
	}, nil
}

func (s *OrganizationService) GetByCode(ctx context.Context, in *v1.GetPositionRequest) (*v1.PositionResponse, error) {
	x, err := s.bc.GetByCode(ctx, in.PositionCode)
	if err != nil {
		return nil, err
	}
	return &v1.PositionResponse{
		Position: &v1.Position{
			Id:                  x.Position.Id,
			PositionCode:        x.Position.PositionCode,
			Name:                x.Position.Name,
			RoleName:            x.Position.RoleName,
			DepartmentCode:      x.Position.DepartmentCode,
			ReportsToPositionId: x.Position.ReportsToPositionId,
			CreatedAt:           x.Position.CreatedAt,
			UpdatedAt:           x.Position.UpdatedAt,
			CreatedBy:           x.Position.CreatedBy,
			UpdatedBy:           x.Position.UpdatedBy,
		},
	}, nil
}
func (s *OrganizationService) DeletePosition(ctx context.Context, in *v1.DeletePositionRequest) (*v1.DeletePositionResponse, error) {
	err := s.bc.Delete(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DeletePositionResponse{Success: true}, nil
}

func (s *OrganizationService) ListPosition(ctx context.Context, in *v1.ListPositionRequest) (*v1.ListPositionsResponse, error) {
	list, total, err := s.bc.List(ctx, int64(in.Pn), int64(in.Psize))
	if err != nil {
		return nil, err
	}

	var positions []*v1.Position
	for _, x := range list {
		positions = append(positions, &v1.Position{
			Id:                  x.ID,
			PositionCode:        x.PositionCode,
			Name:                x.Name,
			RoleName:            x.RoleName,
			DepartmentCode:      x.DepartmentCode,
			ReportsToPositionId: x.ReportsToPositionId,
			CreatedAt:           timestamppb.New(x.CreatedAt),
			UpdatedAt:           timestamppb.New(x.UpdatedAt),
			CreatedBy:           x.CreatedBy,
			UpdatedBy:           x.UpdatedBy,
		})
	}

	return &v1.ListPositionsResponse{
		Total:     int64(total),
		Positions: positions,
	}, nil
}

func (s *OrganizationService) AssignPosition(ctx context.Context, req *v1.AssignPositionRequest) (*v1.AssignmentResponse, error) {
	// Validasi request awal
	if req.EmployeeId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "employee_id is required")
	}
	if req.PositionId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "position_id is required")
	}
	if req.StartDate == "" {
		return nil, status.Errorf(codes.InvalidArgument, "start_date is required")
	}

	// Parse PositionId dari string ke int64
	positionID, err := strconv.ParseInt(req.PositionId, 10, 64)
	if err != nil || positionID <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid position_id format")
	}

	// Parse StartDate dari string ke time.Time
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid start_date format, must be YYYY-MM-DD")
	}

	// Mapping ke struct biz layer
	bizReq := &biz.AssignmentData{
		EmployeeID: req.EmployeeId,
		PositionID: positionID,
		StartDate:  startDate,
	}

	// Delegasikan ke usecase (biz logic)
	assignment, err := s.rp.Assign(ctx, bizReq)
	if err != nil {
		return nil, err
	}

	// Mapping ke proto response
	resp := &v1.AssignmentResponse{
		Assignment: &v1.Assignment{
			Id:         assignment.Assignment.Id,
			EmployeeId: assignment.Assignment.EmployeeId,
			PositionId: assignment.Assignment.PositionId,
			StartDate:  assignment.Assignment.StartDate,
			CreatedAt:  assignment.Assignment.CreatedAt,
		},
	}

	return resp, nil
}

func (s *OrganizationService) CreatePerusahaan(ctx context.Context, req *v1.CreatePerusahaanRequest) (*v1.PerusahaanResponse, error) {
	b := &v1.CreatePerusahaanRequest{
		KodePerusahaan: req.KodePerusahaan,
		NamaPerusahaan: req.NamaPerusahaan,
		KodeCabang:     req.KodeCabang,
		Cabang:         req.KodeCabang,
		Alamat:         req.Alamat,
		Telp:           req.Telp,
		Email:          req.Email,
	}
	x, err := s.pr.Create(ctx, b)
	if err != nil {
		return nil, err
	}
	if err := s.data.Kafka.PublishEvent(ctx, &biz.Event{
		Type: "create",
		Key:  x.KodePerusahaan,
		Data: fmt.Sprintf(`{"kode_perusahaan: %s","nama_perusahaan: %s", "kode_cabang: %s", "cabang: %s","alamat: %s","telp: %s","email: %s"}`, x.KodePerusahaan, x.NamaPerusahaan, x.KodeCabang, x.Cabang, x.Alamat, x.Telp, x.Email),
	}); err != nil {
		s.log.Warnf("failed to pubish event: %v", err)
	}
	return &v1.PerusahaanResponse{
		Result: &v1.Perusahaan{
			Id:             x.Id,
			KodePerusahaan: x.KodePerusahaan,
			NamaPerusahaan: x.NamaPerusahaan,
			KodeCabang:     x.KodeCabang,
			Cabang:         x.Cabang,
			Alamat:         x.Alamat,
			Telp:           x.Telp,
			Email:          x.Email,
			// CreatedAt:      x.CreatedAt,
		},
	}, nil

}

func (s *OrganizationService) UpdatePerusahaan(ctx context.Context, req *v1.UpdatePerusahaanRequest) (*v1.PerusahaanResponse, error) {
	b := &v1.UpdatePerusahaanRequest{
		KodePerusahaan: req.KodePerusahaan,
		NamaPerusahaan: req.NamaPerusahaan,
		KodeCabang:     req.KodeCabang,
		Cabang:         req.KodeCabang,
		Alamat:         req.Alamat,
		Telp:           req.Telp,
		Email:          req.Email,
	}
	x, err := s.pr.Update(ctx, b)
	if err != nil {
		return nil, err
	}
	if err := s.data.Kafka.PublishEvent(ctx, &biz.Event{
		Type: "create",
		Key:  x.KodePerusahaan,
		Data: fmt.Sprintf(`{"kode_perusahaan: %s","nama_perusahaan: %s", "kode_cabang: %s", "cabang: %s","alamat: %s","telp: %s","email: %s"}`, x.KodePerusahaan, x.NamaPerusahaan, x.KodeCabang, x.Cabang, x.Alamat, x.Telp, x.Email),
	}); err != nil {
		s.log.Warnf("failed to pubish event: %v", err)
	}
	return &v1.PerusahaanResponse{
		Result: &v1.Perusahaan{
			Id:             x.Id,
			KodePerusahaan: x.KodePerusahaan,
			NamaPerusahaan: x.NamaPerusahaan,
			KodeCabang:     x.KodeCabang,
			Cabang:         x.Cabang,
			Alamat:         x.Alamat,
			Telp:           x.Telp,
			Email:          x.Email,
			// CreatedAt:      x.CreatedAt,
		},
	}, nil

}
func (s *OrganizationService) DeletePerusahaan(ctx context.Context, in *v1.DeletePerusahaanRequest) (*v1.DeletePerusahaanResponse, error) {
	err := s.pr.Delete(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DeletePerusahaanResponse{Success: true}, nil
}

// Get Perusahaan by Kode
func (s *OrganizationService) GetPerusahaan(ctx context.Context, in *v1.GetPerusahaanRequest) (*v1.PerusahaanResponse, error) {
	x, err := s.pr.Get(ctx, in.KodePerusahaan)
	if err != nil {
		return nil, err
	}
	return &v1.PerusahaanResponse{
		Result: &v1.Perusahaan{
			Id:             x.Id,
			KodePerusahaan: x.KodePerusahaan,
			NamaPerusahaan: x.NamaPerusahaan,
			KodeCabang:     x.KodeCabang,
			Cabang:         x.Cabang,
			Alamat:         x.Alamat,
			Telp:           x.Telp,
			Email:          x.Email,
			// CreatedAt:      x.CreatedAt,
		},
	}, nil
}

func (s *OrganizationService) GetCabang(ctx context.Context, in *v1.GetCabangRequest) (*v1.PerusahaanResponse, error) {
	x, err := s.pr.GetCabang(ctx, in.KodeCabang)
	if err != nil {
		return nil, err
	}
	return &v1.PerusahaanResponse{
		Result: &v1.Perusahaan{
			Id:             x.Id,
			KodePerusahaan: x.KodePerusahaan,
			NamaPerusahaan: x.NamaPerusahaan,
			KodeCabang:     x.KodeCabang,
			Cabang:         x.Cabang,
			Alamat:         x.Alamat,
			Telp:           x.Telp,
			Email:          x.Email,
		},
	}, nil
}
