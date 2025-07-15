package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	// "github.com/go-kratos/beer-shop/pkg/page_token"
)

type Biometric struct {
	Id         int64
	Fingercode string
	Finger0    string
	Finger1    string
	Finger2    string
	Finger3    string
	Finger4    string
	Finger5    string
	Finger6    string
	Finger7    string
	Finger8    string
	Finger9    string
}

type BiometricRepo interface {
	GetFingerByID(ctx context.Context, id int64) (*Biometric, error)
	GetFingerByKode(ctx context.Context, kodefinger string) (*Biometric, error)
	CreateFinger(ctx context.Context, f *Biometric) (*Biometric, error)
	UpdateFinger(ctx context.Context, f *Biometric) (*Biometric, error)
	DeleteFinger(ctx context.Context, id int64) error
}
type BiometricUsecase struct {
	repo BiometricRepo
	// pageToken page_token.ProcessPageTokens
	log *log.Helper
	// sg        singleflight.Group
}

func NewBiometricUsecase(repo BiometricRepo, logger log.Logger) *BiometricUsecase {
	return &BiometricUsecase{
		repo: repo,
		// log:  log.NewHelper(logger),
		log: log.NewHelper(log.With(logger, "module", "usecase/Biometric")),
	}
}

// Finger
func (uc *BiometricUsecase) CreateFinger(ctx context.Context, u *Biometric) (*Biometric, error) {
	return uc.repo.CreateFinger(ctx, u)
}
func (uc *BiometricUsecase) UpdateFinger(ctx context.Context, u *Biometric) (*Biometric, error) {
	return uc.repo.UpdateFinger(ctx, u)
}

func (uc *BiometricUsecase) DeleteFinger(ctx context.Context, id int64) error {
	return uc.repo.DeleteFinger(ctx, id)
}
func (uc *BiometricUsecase) GetFingerByID(ctx context.Context, id int64) (*Biometric, error) {
	return uc.repo.GetFingerByID(ctx, id)
}

func (uc *BiometricUsecase) GetFingerByKode(ctx context.Context, kode string) (*Biometric, error) {
	return uc.repo.GetFingerByKode(ctx, kode)
}
