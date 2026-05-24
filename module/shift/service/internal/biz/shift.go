package biz

import (
	"context"
	"fmt"

	"mall-go/pkg/page_token"

	"golang.org/x/sync/singleflight"

	"github.com/go-kratos/kratos/v2/log"
	// "github.com/go-kratos/beer-shop/pkg/page_token"
	// biometricV1 "mall-go/api/biometrics/service/v1"
	// departmentV1 "mall-go/api/department/service/v1"
)

type Shift struct {
	Id                   int64
	Name                 string
	StartTime            string
	EndTime              string
	BreakDurationMinutes int32
	CreatedBy            string
	UpdatedBy            string
	CreatedAt            string
	UpdatedAt            string
}

type DepartData struct {
	DepartCode string
	DepartName string
}

// type FingerData struct {
// 	Fingercode string
// 	Finger0    string
// 	Finger1    string
// 	Finger2    string
// 	Finger3    string
// 	Finger4    string
// 	Finger5    string
// 	Finger6    string
// 	Finger7    string
// 	Finger8    string
// 	Finger9    string
// }

// type EmployerData struct {
// 	Id         int64
// 	NoSap      string
// 	Nip        string
// 	KaryaCode  string
// 	KaryaName  string
// 	DispName   string
// 	PassMesin  string
// 	RFIDCard   string
// 	Finger     []FingerData
// 	Department []DepartData
// 	Status     int32
// 	CreatedAt  string
// 	UpdatedAt  string
// }

type ShiftRepo interface {
	CreateShift(ctx context.Context, c *Shift) (*Shift, error)
	UpdateShift(ctx context.Context, c *Shift) (*Shift, error)
	GetShiftID(ctx context.Context, id int64) (*Shift, error)
	DeleteShift(ctx context.Context, id int64) error
	ListShift(ctx context.Context, pageNum, pageSize int64) ([]*Shift, int, error)
	ListShiftNext(ctx context.Context, start, end int32) ([]*Shift, error)
	Count(ctx context.Context) (int, error)
	// GetEmployerDetail(ctx context.Context, id int64) (*Shift, error)
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
type ShiftUsecase struct {
	repo      ShiftRepo
	pageToken page_token.ProcessPageTokens
	log       *log.Helper
	// deptClient departmentV1.DepartmentClient
	// bioCilent  biometricV1.BiometricClient
	sg singleflight.Group
}

func NewShiftUsecase(repo ShiftRepo, logger log.Logger) *ShiftUsecase {
	return &ShiftUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "usecase/employers")),
	}
}

//	func NewShiftUsecase(repo ShiftRepo, logger log.Logger, dept departmentV1.DepartmentClient, bio biometricV1.BiometricClient) *ShiftUsecase {
//		return &ShiftUsecase{
//			repo:       repo,
//			log:        log.NewHelper(log.With(logger, "module", "usecase/employers")),
//			// deptClient: dept,
//			// bioCilent:  bio,
//		}
//	}
func (uc *ShiftUsecase) Create(ctx context.Context, u *Shift) (*Shift, error) {
	return uc.repo.CreateShift(ctx, u)
}
func (uc *ShiftUsecase) Update(ctx context.Context, u *Shift) (*Shift, error) {
	return uc.repo.UpdateShift(ctx, u)
}

// get by id
func (uc *ShiftUsecase) GetByID(ctx context.Context, id int64) (*Shift, error) {
	return uc.repo.GetShiftID(ctx, id)
}

func (uc *ShiftUsecase) Delete(ctx context.Context, id int64) error {
	return uc.repo.DeleteShift(ctx, id)
}

func (uc *ShiftUsecase) List(ctx context.Context, pageNum, pageSize int64) ([]*Shift, int, error) {
	uc.log.Infof("ListShift called with pageNum=%d, pageSize=%d", pageNum, pageSize)

	// return uc.repo.ListShift(ctx, pageNum, pageSize)
	list, total, err := uc.repo.ListShift(ctx, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}

	if list == nil {
		return nil, 0, fmt.Errorf("repo.ListShift returned nil list")
	}

	return list, total, nil
}

func (uc *ShiftUsecase) ListNext(ctx context.Context, pageSize int32, pageToken string) ([]*Shift, string, error) {
	total, err := uc.repo.Count(ctx)
	if err != nil {
		return nil, "", err
	}
	// log.Error("INI TOTAL :", total)

	start, end, nextToken, err := uc.pageToken.ProcessPageTokens(total, pageSize, pageToken)
	if err != nil {
		return nil, "", err
	}

	data, err, _ := uc.sg.Do(fmt.Sprintf("list_next_%d_%d", start, end), func() (interface{}, error) {
		return uc.repo.ListShiftNext(ctx, start, end)
	})
	return data.([]*Shift), nextToken, err
}

// func (uc *ShiftUsecase) GetDepartmentByCode(ctx context.Context, deptCode string) (*departmentV1.GetDepartmentCodeResponse, error) {
// 	return uc.deptClient.GetDepartmentCode(ctx, &departmentV1.GetDepartmentCodeRequest{DepartCode: deptCode})
// }

// func (uc *ShiftUsecase) GetFingerByKode(ctx context.Context, bioCode string) (*biometricV1.GetFingerByKodeResponse, error) {
// 	return uc.bioCilent.GetFingerByKode(ctx, &biometricV1.GetFingerByKodeRequest{Fingercode: bioCode})
// }

// func (uc *ShiftUsecase) GetDetail(ctx context.Context, id int64) (*EmployerData, error) {
// 	return uc.repo.GetEmployerDetail(ctx, id)
// }
