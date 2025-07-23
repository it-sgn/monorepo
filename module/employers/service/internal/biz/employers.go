package biz

import (
	"context"
	"fmt"

	"mall-go/pkg/page_token"

	"golang.org/x/sync/singleflight"

	"github.com/go-kratos/kratos/v2/log"
	// "github.com/go-kratos/beer-shop/pkg/page_token"

	biometricV1 "mall-go/api/biometrics/service/v1"
	departmentV1 "mall-go/api/department/service/v1"
)

type Employers struct {
	Id         int64
	NoSap      string
	Nip        string
	KaryaCode  string
	KaryaName  string
	DispName   string
	PassMesin  string
	RFIDCard   string
	Finger     string
	Department string
	Status     int32
	CreatedAt  string
	UpdatedAt  string
}
type DepartData struct {
	DepartCode string
	DepartName string
}
type FingerData struct {
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

type EmployerData struct {
	Id         int64
	NoSap      string
	Nip        string
	KaryaCode  string
	KaryaName  string
	DispName   string
	PassMesin  string
	RFIDCard   string
	Finger     []FingerData
	Department []DepartData
	Status     int32
	CreatedAt  string
	UpdatedAt  string
}

type EmployersRepo interface {
	CreateEmployers(ctx context.Context, c *Employers) (*Employers, error)
	UpdateEmployers(ctx context.Context, c *Employers) (*Employers, error)
	GetEmployersID(ctx context.Context, id int64) (*Employers, error)
	DeleteEmployers(ctx context.Context, id int64) error
	ListEmployers(ctx context.Context, pageNum, pageSize int64) ([]*EmployerData, int, error)
	ListEmployersNext(ctx context.Context, start, end int32) ([]*Employers, error)
	Count(ctx context.Context) (int, error)
	GetEmployerDetail(ctx context.Context, id int64) (*EmployerData, error)
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
type EmployersUsecase struct {
	repo       EmployersRepo
	pageToken  page_token.ProcessPageTokens
	log        *log.Helper
	deptClient departmentV1.DepartmentClient
	bioCilent  biometricV1.BiometricClient
	sg         singleflight.Group
}

func NewEmployersUsecase(repo EmployersRepo, logger log.Logger, dept departmentV1.DepartmentClient, bio biometricV1.BiometricClient) *EmployersUsecase {
	return &EmployersUsecase{
		repo:       repo,
		log:        log.NewHelper(log.With(logger, "module", "usecase/employers")),
		deptClient: dept,
		bioCilent:  bio,
	}
}
func (uc *EmployersUsecase) Create(ctx context.Context, u *Employers) (*Employers, error) {
	return uc.repo.CreateEmployers(ctx, u)
}
func (uc *EmployersUsecase) Update(ctx context.Context, u *Employers) (*Employers, error) {
	return uc.repo.UpdateEmployers(ctx, u)
}

// get by id
func (uc *EmployersUsecase) GetByID(ctx context.Context, id int64) (*Employers, error) {
	return uc.repo.GetEmployersID(ctx, id)
}

func (uc *EmployersUsecase) List(ctx context.Context, pageNum, pageSize int64) ([]*EmployerData, int, error) {
	uc.log.Infof("ListEmployers called with pageNum=%d, pageSize=%d", pageNum, pageSize)

	// return uc.repo.ListEmployers(ctx, pageNum, pageSize)
	list, total, err := uc.repo.ListEmployers(ctx, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}

	if list == nil {
		return nil, 0, fmt.Errorf("repo.ListEmployers returned nil list")
	}

	return list, total, nil
}

func (uc *EmployersUsecase) ListNext(ctx context.Context, pageSize int32, pageToken string) ([]*Employers, string, error) {
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
		return uc.repo.ListEmployersNext(ctx, start, end)
	})
	return data.([]*Employers), nextToken, err
}

func (uc *EmployersUsecase) Delete(ctx context.Context, id int64) error {
	return uc.repo.DeleteEmployers(ctx, id)
}

// func (uc *EmployersUsecase) GetFingerByKode(ctx context.Context, fcode string) (*FingerData, error) {
// 	return uc.repo.GetFingerByKode(ctx, fcode)
// }
// func (uc *EmployersUsecase) GetDepartmentByKode(ctx context.Context, dkode string) (*DepartData, error) {
// 	return uc.repo.GetDepartmentByKode(ctx, dkode)
// }

func (uc *EmployersUsecase) GetDepartmentByCode(ctx context.Context, deptCode string) (*departmentV1.GetDepartmentCodeResponse, error) {
	return uc.deptClient.GetDepartmentCode(ctx, &departmentV1.GetDepartmentCodeRequest{DepartCode: deptCode})
}

func (uc *EmployersUsecase) GetFingerByKode(ctx context.Context, bioCode string) (*biometricV1.GetFingerByKodeResponse, error) {
	return uc.bioCilent.GetFingerByKode(ctx, &biometricV1.GetFingerByKodeRequest{Fingercode: bioCode})
}

func (uc *EmployersUsecase) GetDetail(ctx context.Context, id int64) (*EmployerData, error) {
	return uc.repo.GetEmployerDetail(ctx, id)
}
