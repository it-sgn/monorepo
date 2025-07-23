package biz

import (
	"context"
	"time"

	"mall-go/pkg/page_token"

	"golang.org/x/sync/singleflight"

	"github.com/go-kratos/kratos/v2/log"
	// "github.com/go-kratos/beer-shop/pkg/page_token"

	// biometricV1 "mall-go/api/biometric/service/v1"
	departmentV1 "mall-go/api/department/service/v1"
	employersV1 "mall-go/api/employers/service/v1"
)

type Attendance struct {
	Id                           int64
	KaryaCode                    string
	ClockInTime                  time.Time
	ClockOutTime                 string
	Status                       string
	EffectiveWorkDurationMinutes int
	OvertimeMinutes              int
	Notes                        string
	Location                     string
	CreatedAt                    string
	UpdatedAt                    string
}

type AttendanceRepo interface {
	CreateClockIn(ctx context.Context, c *Attendance) (*Attendance, error)
	UpdateClockOut(ctx context.Context, karyacode string) (*Attendance, error)
	SetStatus(ctx context.Context, karyacode string) (*Attendance, error)
	CountEffectiveWorkDurationMinutes(ctx context.Context, karyacode string) (*Attendance, error)
	CountOvertimeMinutes(ctx context.Context, karyacode string) (*Attendance, error)
	SetNotesOvertime(ctx context.Context, karyacode string) (*Attendance, error)
	// Count(ctx context.Context) (int, error)
	// GetAttendanceDetail(ctx context.Context, id int64) (*Attendance, error)
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
type AttendanceUsecase struct {
	repo       AttendanceRepo
	pageToken  page_token.ProcessPageTokens
	log        *log.Helper
	deptClient departmentV1.DepartmentClient
	empCilent  employersV1.EmployersClient
	sg         singleflight.Group
}

func NewAttendanceUsecase(repo AttendanceRepo, logger log.Logger, dept departmentV1.DepartmentClient, emp employersV1.EmployersClient) *AttendanceUsecase {
	return &AttendanceUsecase{
		repo:       repo,
		log:        log.NewHelper(log.With(logger, "module", "usecase/employers")),
		deptClient: dept,
		empCilent:  emp,
	}
}
func (uc *AttendanceUsecase) CreateClockIn(ctx context.Context, u *Attendance) (*Attendance, error) {
	return uc.repo.CreateClockIn(ctx, u)
}

func (uc *AttendanceUsecase) UpdateClockOut(ctx context.Context, karyacode string) (*Attendance, error) {
	return uc.repo.UpdateClockOut(ctx, karyacode)
}

// get by id
func (uc *AttendanceUsecase) SetStatus(ctx context.Context, karyacode string) (*Attendance, error) {
	return uc.repo.SetStatus(ctx, karyacode)
}

func (uc *AttendanceUsecase) CountEffectiveWorkDurationMinutes(ctx context.Context, karyacode string) (*Attendance, error) {
	return uc.repo.CountEffectiveWorkDurationMinutes(ctx, karyacode)
}
func (uc *AttendanceUsecase) CountOvertimeMinutes(ctx context.Context, karyacode string) (*Attendance, error) {
	return uc.repo.CountOvertimeMinutes(ctx, karyacode)
}

func (uc *AttendanceUsecase) SetNotesOvertime(ctx context.Context, karyacode string) (*Attendance, error) {
	return uc.repo.SetNotesOvertime(ctx, karyacode)
}
