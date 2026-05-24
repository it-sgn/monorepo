package biz

import (
	"context"
	"time"

	"mall-go/pkg/page_token"

	"golang.org/x/sync/singleflight"

	"github.com/go-kratos/kratos/v2/log"
	// "github.com/go-kratos/beer-shop/pkg/page_token"
	v1 "mall-go/api/organization/service/v1"
)

type PerusahaanData struct {
	Id             int64
	KodePerusahaan string
	NamaPerusahaan string
	KodeCabang     string
	Cabang         string
	Alamat         string
	Telp           string
	Email          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      string
	UpdatedBy      string
}
type PerusahaanRepo interface {
	CreatePerusahaan(ctx context.Context, c *v1.CreatePerusahaanRequest) (*PerusahaanData, error)
	UpdatePerusahaan(ctx context.Context, c *v1.UpdatePerusahaanRequest) (*PerusahaanData, error)
	DeletePerusahaan(ctx context.Context, id int64) error
	GetPerusahaan(ctx context.Context, kode_perusahaan string) (*PerusahaanData, error)
	GetCabang(ctx context.Context, kode_cabang string) (*PerusahaanData, error)
	ListPerusahaan(ctx context.Context, pageNum, pageSize int64) ([]*PerusahaanData, int, error)
	Count(ctx context.Context) (int, error)
	// GetPerusahaanCode(ctx context.Context, Dcode string) (*v1.Perusahaan, error)
}
type PerusahaanUsecase struct {
	repo      PerusahaanRepo
	pageToken page_token.ProcessPageTokens
	log       *log.Helper
	sg        singleflight.Group
}

func NewPerusahaanUsecase(repo PerusahaanRepo, logger log.Logger) *PerusahaanUsecase {
	return &PerusahaanUsecase{
		repo: repo,
		// log:  log.NewHelper(logger),
		log: log.NewHelper(log.With(logger, "module", "usecase/Perusahaan")),
	}
}
func (uc *PerusahaanUsecase) Create(ctx context.Context, u *v1.CreatePerusahaanRequest) (*PerusahaanData, error) {
	return uc.repo.CreatePerusahaan(ctx, u)
}
func (uc *PerusahaanUsecase) Update(ctx context.Context, u *v1.UpdatePerusahaanRequest) (*PerusahaanData, error) {
	return uc.repo.UpdatePerusahaan(ctx, u)
}

func (uc *PerusahaanUsecase) Delete(ctx context.Context, id int64) error {
	return uc.repo.DeletePerusahaan(ctx, id)
}

// get by id
func (uc *PerusahaanUsecase) Get(ctx context.Context, kode_perusahaan string) (*PerusahaanData, error) {
	return uc.repo.GetPerusahaan(ctx, kode_perusahaan)
}

func (uc *PerusahaanUsecase) GetCabang(ctx context.Context, kode_cabang string) (*PerusahaanData, error) {
	return uc.repo.GetCabang(ctx, kode_cabang)
}

func (uc *PerusahaanUsecase) List(ctx context.Context, pageNum, pageSize int64) ([]*PerusahaanData, int, error) {
	return uc.repo.ListPerusahaan(ctx, pageNum, pageSize)
}
