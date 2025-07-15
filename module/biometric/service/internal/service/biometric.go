package service

import (

	// v1 "mall-go/api/service/employers/service"
	// v1 "mall-go/api/service/employers/service"

	"context"
	v1 "mall-go/api/biometric/service/v1"
	"mall-go/module/biometric/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type BiometricService struct {
	v1.UnimplementedBiometricServer

	bc  *biz.BiometricUsecase
	log *log.Helper
}

func NewBiometricService(bc *biz.BiometricUsecase, logger log.Logger) *BiometricService {
	return &BiometricService{
		bc:  bc,
		log: log.NewHelper(log.With(logger, "module", "service/biometric"))}
}

func (s *BiometricService) CreateFinger(ctx context.Context, in *v1.CreateFingerRequest) (*v1.CreateFingerResponse, error) {
	b := &biz.Biometric{
		Fingercode: in.Fingercode,
		Finger0:    in.Finger0,
		Finger1:    in.Finger1,
		Finger2:    in.Finger2,
		Finger3:    in.Finger3,
		Finger4:    in.Finger4,
		Finger5:    in.Finger5,
		Finger6:    in.Finger6,
		Finger7:    in.Finger7,
		Finger8:    in.Finger8,
		Finger9:    in.Finger9,
	}
	x, err := s.bc.CreateFinger(ctx, b)
	if err != nil {
		return nil, err
	}
	return &v1.CreateFingerResponse{
		Id:         x.Id,
		Fingercode: x.Fingercode,
		Finger0:    x.Finger0,
		Finger1:    x.Finger1,
		Finger2:    x.Finger2,
		Finger3:    x.Finger3,
		Finger4:    x.Finger4,
		Finger5:    x.Finger5,
		Finger6:    x.Finger6,
		Finger7:    x.Finger7,
		Finger8:    x.Finger8,
		Finger9:    x.Finger9,
	}, nil
}
func (s *BiometricService) UpdateFinger(ctx context.Context, in *v1.UpdateFingerRequest) (*v1.UpdateFingerResponse, error) {
	b := &biz.Biometric{
		Fingercode: in.Fingercode,
		Finger0:    in.Finger0,
		Finger1:    in.Finger1,
		Finger2:    in.Finger2,
		Finger3:    in.Finger3,
		Finger4:    in.Finger4,
		Finger5:    in.Finger5,
		Finger6:    in.Finger6,
		Finger7:    in.Finger7,
		Finger8:    in.Finger8,
		Finger9:    in.Finger9,
	}
	x, err := s.bc.UpdateFinger(ctx, b)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateFingerResponse{
		Id:         x.Id,
		Fingercode: x.Fingercode,
		Finger0:    x.Finger0,
		Finger1:    x.Finger1,
		Finger2:    x.Finger2,
		Finger3:    x.Finger3,
		Finger4:    x.Finger4,
		Finger5:    x.Finger5,
		Finger6:    x.Finger6,
		Finger7:    x.Finger7,
		Finger8:    x.Finger8,
		Finger9:    x.Finger9,
	}, nil
}
func (s *BiometricService) GetFingerByID(ctx context.Context, in *v1.GetFingerByIDRequest) (*v1.GetFingerByIDResponse, error) {
	x, err := s.bc.GetFingerByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetFingerByIDResponse{
		Id:         x.Id,
		Fingercode: x.Fingercode,
		Finger0:    x.Finger0,
		Finger1:    x.Finger1,
		Finger2:    x.Finger2,
		Finger3:    x.Finger3,
		Finger4:    x.Finger4,
		Finger5:    x.Finger5,
		Finger6:    x.Finger6,
		Finger7:    x.Finger7,
		Finger8:    x.Finger8,
		Finger9:    x.Finger9,
	}, nil
}

func (s *BiometricService) GetFingerByKode(ctx context.Context, in *v1.GetFingerByKodeRequest) (*v1.GetFingerByKodeResponse, error) {
	x, err := s.bc.GetFingerByKode(ctx, in.Fingercode)
	if err != nil {
		return nil, err
	}
	return &v1.GetFingerByKodeResponse{
		Id:         x.Id,
		Fingercode: x.Fingercode,
		Finger0:    x.Finger0,
		Finger1:    x.Finger1,
		Finger2:    x.Finger2,
		Finger3:    x.Finger3,
		Finger4:    x.Finger4,
		Finger5:    x.Finger5,
		Finger6:    x.Finger6,
		Finger7:    x.Finger7,
		Finger8:    x.Finger8,
		Finger9:    x.Finger9,
	}, nil

}

func (s *BiometricService) DeleteFinger(ctx context.Context, in *v1.DeleteFingerRequest) (*v1.DeleteFingerResponse, error) {
	err := s.bc.DeleteFinger(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteFingerResponse{Success: true}, nil
}

// func (s *BiometricService) DeleteFinger(ctx context.Context, in *v1.CreateFingerRequest) (*v1.DeleteFingerResponse, error) {
// 	err := s.bc.DeleteFinger(ctx, in.Id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &v1.DeleteFingerResponse{Success: true}, nil
// }
