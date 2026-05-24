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

type AssignmentData struct {
	ID         int64
	EmployeeID string
	PositionID int64
	StartDate  time.Time
	EndDate    time.Time
	CreatedAt  string
	UpdatedAt  string
	CreatedBy  string
	UpdatedBy  string
}

type AssignmentRepo interface {
	AssignPosition(ctx context.Context, c *AssignmentData) (*v1.AssignmentResponse, error)
	GetAssignment(ctx context.Context, c *v1.GetAssignmentRequest) (*v1.AssignmentResponse, error)
	ListAssignment(ctx context.Context, pageNum, pageSize int64) (*v1.ListAssignmentsResponse, int, error)
	DeleteAssignment(ctx context.Context, id int64) error
	Count(ctx context.Context) (int, error)
}
type AssignmentUsecase struct {
	repo      AssignmentRepo
	pageToken page_token.ProcessPageTokens
	log       *log.Helper
	sg        singleflight.Group
}

func NewAssignmentUsecase(repo AssignmentRepo, logger log.Logger) *AssignmentUsecase {
	return &AssignmentUsecase{
		repo: repo,
		// log:  log.NewHelper(logger),
		log: log.NewHelper(log.With(logger, "module", "usecase/Assignment")),
	}
}

func (uc *AssignmentUsecase) Assign(ctx context.Context, u *AssignmentData) (*v1.AssignmentResponse, error) {
	return uc.repo.AssignPosition(ctx, u)
}

// func (uc *AssignmentUsecase) Assign(ctx context.Context, in *AssignmentData) (*AssignmentData, error) {
// 	in.CreatedAt = time.Now()
// 	in, err := uc.repo.AssignPosition(ctx, in)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return in, nil
// }

// get by id
func (uc *AssignmentUsecase) Get(ctx context.Context, u *v1.GetAssignmentRequest) (*v1.AssignmentResponse, error) {
	return uc.repo.GetAssignment(ctx, u)
}

func (uc *AssignmentUsecase) List(ctx context.Context, pageNum, pageSize int64) (*v1.ListAssignmentsResponse, int, error) {
	return uc.repo.ListAssignment(ctx, pageNum, pageSize)
}

func (uc *AssignmentUsecase) Delete(ctx context.Context, id int64) error {
	return uc.repo.DeleteAssignment(ctx, id)
}
