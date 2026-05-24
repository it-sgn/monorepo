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

type PositionRepo interface {
	CreatePosition(ctx context.Context, c *v1.CreatePositionRequest) (*v1.PositionResponse, error)
	GetPositionID(ctx context.Context, id int64) (*v1.PositionResponse, error)
	GetPosition(ctx context.Context, position_code string) (*v1.PositionResponse, error)
	UpdatePosition(ctx context.Context, c *v1.Position) (*v1.PositionResponse, error)
	ListPosition(ctx context.Context, pageNum, pageSize int64) ([]*PositionData, int, error)
	DeletePosition(ctx context.Context, id int64) error
	// AssignPosition(ctx context.Context c *v1.AssignPositionRequest)(*v1.AssignmentResponse)

	// ListPositionNext(ctx context.Context, start, end int32) ([]*v1.Position, error)
	Count(ctx context.Context) (int, error)
	GetPositionCode(ctx context.Context, Dcode string) (*v1.Position, error)
}
type PositionUsecase struct {
	repo      PositionRepo
	pageToken page_token.ProcessPageTokens
	log       *log.Helper
	sg        singleflight.Group
}
type PositionData struct {
	ID                  int64
	PositionCode        string
	Name                string
	RoleName            string
	DepartmentCode      string
	ReportsToPositionId string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	CreatedBy           string
	UpdatedBy           string
}

func NewPositionUsecase(repo PositionRepo, logger log.Logger) *PositionUsecase {
	return &PositionUsecase{
		repo: repo,
		// log:  log.NewHelper(logger),
		log: log.NewHelper(log.With(logger, "module", "usecase/Position")),
	}
}
func (uc *PositionUsecase) Create(ctx context.Context, u *v1.CreatePositionRequest) (*v1.PositionResponse, error) {
	return uc.repo.CreatePosition(ctx, u)
}
func (uc *PositionUsecase) Update(ctx context.Context, u *v1.Position) (*v1.PositionResponse, error) {
	return uc.repo.UpdatePosition(ctx, u)
}

// get by id
func (uc *PositionUsecase) Get(ctx context.Context, id int64) (*v1.PositionResponse, error) {
	return uc.repo.GetPositionID(ctx, id)
}

// get by id
func (uc *PositionUsecase) GetByCode(ctx context.Context, position_code string) (*v1.PositionResponse, error) {
	return uc.repo.GetPosition(ctx, position_code)
}

func (uc *PositionUsecase) List(ctx context.Context, pageNum, pageSize int64) ([]*PositionData, int, error) {
	return uc.repo.ListPosition(ctx, pageNum, pageSize)
}
func (uc *PositionUsecase) Delete(ctx context.Context, id int64) error {
	return uc.repo.DeletePosition(ctx, id)
}
