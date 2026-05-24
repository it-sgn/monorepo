package biz

import (
	"context"
	"fmt"
	"time"

	"mall-go/pkg/page_token"

	"golang.org/x/sync/singleflight"

	"github.com/go-kratos/kratos/v2/log"
	// "github.com/go-kratos/beer-shop/pkg/page_token"

	departmentV1 "mall-go/api/department/service/v1"
	employerV1 "mall-go/api/employers/service/v1"
)

type ShiftSchedule struct {
	Id           int64
	ScheduleCode string
	KaryaCode    string
	Tanggal      time.Time
	DepartCode   string
	CreatedBy    string
	ShiftID      string
}

type ShiftScheduleData struct {
	Id           int64
	ScheduleCode string
	Employers    []EmployerData
	Tanggal      time.Time
	Department   []DepartData
	CreatedBy    string
	ShiftID      string
}

type DepartData struct {
	DepartCode string
	DepartName string
}
type EmployerData struct {
	Id        int64
	NoSap     string
	Nip       string
	KaryaCode string
	KaryaName string
}

type ShiftScheduleRepo interface {
	CreateShiftSchedule(ctx context.Context, c *ShiftSchedule) (*ShiftSchedule, error)
	UpdateShiftSchedule(ctx context.Context, c *ShiftSchedule) (*ShiftSchedule, error)
	GetShiftScheduleID(ctx context.Context, id int64) (*ShiftSchedule, error)
	DeleteShiftSchedule(ctx context.Context, id int64) error
	ListShiftSchedule(ctx context.Context, pageNum, pageSize int64) ([]*ShiftScheduleData, int, error)
	// ListShiftScheduleNext(ctx context.Context, start, end int32) ([]*ShiftSchedule, error)
	Count(ctx context.Context) (int, error)
	// GetEmployerDetail(ctx context.Context, id int64) (*ShiftSchedule, error)
	// GetFingerByKode(ctx context.Context, fkode string) (*FingerData, error)
	// GetDepartmentByKode(ctx context.Context, dkode string) (*DepartData, error)
	// GetEmployeeFullData(ctx context.Context, fingerID string, departmentID int64) (*FullEmployeeData, error)

	//
	// GetFingerByID(ctx context.Context, id int64) (*KodeFingerRecord, error)
	// GetFingerByKode(ctx context.Context, kodefinger string) (*KodeFingerRecord, error)
	// CreateFinger(ctx context.Context, f *KodeFingerRecord) (*KodeFingerRecord, error)
	// UpdateFinger(ctx context.Context, f *KodeFingerRecord) (*KodeFingerRecord, error)
	// DeleteFinger(ctx context.Context, id int64) error
}
type ShiftScheduleUsecase struct {
	repo           ShiftScheduleRepo
	pageToken      page_token.ProcessPageTokens
	log            *log.Helper
	deptClient     departmentV1.DepartmentClient
	employerCilent employerV1.EmployersClient
	sg             singleflight.Group
}

func NewShiftScheduleUsecase(repo ShiftScheduleRepo, logger log.Logger, dept departmentV1.DepartmentClient, employer employerV1.EmployersClient) *ShiftScheduleUsecase {
	return &ShiftScheduleUsecase{
		repo:           repo,
		log:            log.NewHelper(log.With(logger, "module", "usecase/employers")),
		deptClient:     dept,
		employerCilent: employer,
	}
}
func (uc *ShiftScheduleUsecase) Create(ctx context.Context, u *ShiftSchedule) (*ShiftSchedule, error) {
	return uc.repo.CreateShiftSchedule(ctx, u)
}
func (uc *ShiftScheduleUsecase) Update(ctx context.Context, u *ShiftSchedule) (*ShiftSchedule, error) {
	return uc.repo.UpdateShiftSchedule(ctx, u)
}

// get by id
func (uc *ShiftScheduleUsecase) GetByID(ctx context.Context, id int64) (*ShiftSchedule, error) {
	return uc.repo.GetShiftScheduleID(ctx, id)
}

func (uc *ShiftScheduleUsecase) List(ctx context.Context, pageNum, pageSize int64) ([]*ShiftScheduleData, int, error) {
	uc.log.Infof("ListShiftScheduleData called with pageNum=%d, pageSize=%d", pageNum, pageSize)

	list, total, err := uc.repo.ListShiftSchedule(ctx, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}

	if list == nil {
		return nil, 0, fmt.Errorf("repo.ListShiftScheduleData returned nil list")
	}

	return list, total, nil
}

// func (uc *ShiftScheduleUsecase) ListNext(ctx context.Context, pageSize int32, pageToken string) ([]*ShiftSchedule, string, error) {
// 	total, err := uc.repo.Count(ctx)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	start, end, nextToken, err := uc.pageToken.ProcessPageTokens(total, pageSize, pageToken)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	data, err, _ := uc.sg.Do(fmt.Sprintf("list_next_%d_%d", start, end), func() (interface{}, error) {
// 		return uc.repo.ListShiftScheduleNext(ctx, start, end)
// 	})
// 	return data.([]*ShiftSchedule), nextToken, err
// }

func (uc *ShiftScheduleUsecase) Delete(ctx context.Context, id int64) error {
	return uc.repo.DeleteShiftSchedule(ctx, id)
}

func (uc *ShiftScheduleUsecase) GetDepartmentByKode(ctx context.Context, deptCode string) (*departmentV1.GetDepartmentCodeResponse, error) {
	return uc.deptClient.GetDepartmentCode(ctx, &departmentV1.GetDepartmentCodeRequest{DepartCode: deptCode})
}

func (uc *ShiftScheduleUsecase) GetEmployerByKode(ctx context.Context, KaryaCode string) (*employerV1.GetEmployersKodeResponse, error) {
	return uc.employerCilent.GetEmployersKode(ctx, &employerV1.GetEmployersKodeRequest{KaryaCode: KaryaCode})
}
