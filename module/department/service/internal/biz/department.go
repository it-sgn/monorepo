package biz

import (
	"context"
	"fmt"

	"mall-go/pkg/page_token"

	"golang.org/x/sync/singleflight"

	"github.com/go-kratos/kratos/v2/log"
	// "github.com/go-kratos/beer-shop/pkg/page_token"
)

type Department struct {
	Id         int64
	DepartCode string
	DepartName string
	Status     string
	Ket        string
	CreatedAt  string
	UpdatedAt  string
}

type DepartmentRepo interface {
	CreateDepartment(ctx context.Context, c *Department) (*Department, error)
	UpdateDepartment(ctx context.Context, c *Department) (*Department, error)
	GetDepartmentID(ctx context.Context, id int64) (*Department, error)
	DeleteDepartment(ctx context.Context, id int64) error
	ListDepartment(ctx context.Context, pageNum, pageSize int64) ([]*Department, int, error)
	ListDepartmentNext(ctx context.Context, start, end int32) ([]*Department, error)
	Count(ctx context.Context) (int, error)
	GetDepartmentCode(ctx context.Context, Dcode string) (*Department, error)
}
type DepartmentUsecase struct {
	repo      DepartmentRepo
	pageToken page_token.ProcessPageTokens
	log       *log.Helper
	sg        singleflight.Group
}

func NewDepartmentUsecase(repo DepartmentRepo, logger log.Logger) *DepartmentUsecase {
	return &DepartmentUsecase{
		repo: repo,
		// log:  log.NewHelper(logger),
		log: log.NewHelper(log.With(logger, "module", "usecase/Department")),
	}
}
func (uc *DepartmentUsecase) Create(ctx context.Context, u *Department) (*Department, error) {
	return uc.repo.CreateDepartment(ctx, u)
}
func (uc *DepartmentUsecase) Update(ctx context.Context, u *Department) (*Department, error) {
	return uc.repo.UpdateDepartment(ctx, u)
}

// get by id
func (uc *DepartmentUsecase) GetByID(ctx context.Context, id int64) (*Department, error) {
	return uc.repo.GetDepartmentID(ctx, id)
}

// get by id
func (uc *DepartmentUsecase) GetByCode(ctx context.Context, Dcode string) (*Department, error) {
	return uc.repo.GetDepartmentCode(ctx, Dcode)
}

func (uc *DepartmentUsecase) List(ctx context.Context, pageNum, pageSize int64) ([]*Department, int, error) {
	return uc.repo.ListDepartment(ctx, pageNum, pageSize)
}

func (uc *DepartmentUsecase) ListNext(ctx context.Context, pageSize int32, pageToken string) ([]*Department, string, error) {
	total, err := uc.repo.Count(ctx)
	if err != nil {
		return nil, "", err
	}
	// log.Error("INI TOTAL :", total)

	start, end, nextToken, err := uc.pageToken.ProcessPageTokens(total, pageSize, pageToken)
	if err != nil {
		return nil, "", err
	}
	// log.Error("PAGE TOKEN :", pageToken)
	// single flight
	data, err, _ := uc.sg.Do(fmt.Sprintf("list_next_%d_%d", start, end), func() (interface{}, error) {
		return uc.repo.ListDepartmentNext(ctx, start, end)
	})
	return data.([]*Department), nextToken, err
}

func (uc *DepartmentUsecase) Delete(ctx context.Context, id int64) error {
	return uc.repo.DeleteDepartment(ctx, id)
}
