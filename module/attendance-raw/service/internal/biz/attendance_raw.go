package biz

import (
	"context"
	"mall-go/pkg/page_token"

	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/sync/singleflight"

	departV1 "mall-go/api/department/service/v1"
	employersV1 "mall-go/api/employers/service/v1"
)

// Struktur Status ClockIn/ClockOut
type Status struct {
	ClockIn  string `json:"ClockIn"`
	ClockOut string `json:"ClockOut"`
}

// Data absensi per hari
type Absensi struct {
	Tanggal  string    `json:"tanggal"`
	Evaluasi string    `json:"evaluasi"`
	Status   []*Status `json:"Status"`
}

// Data per karyawan
type Karyawan struct {
	Karyaname string    `json:"karyaname"`
	Periode   string    `json:"periode"`
	Data      []Absensi `json:"Data"`
}

// Laporan keseluruhan
type AttendanceReport struct {
	NamaPerusahaan string     `json:"NamaPerusahaan"`
	Cabang         string     `json:"Cabang"`
	Department     string     `json:"Department"`
	Jabatan        string     `json:"Jabatan"`
	DibuatOleh     string     `json:"DibuatOleh"`
	DiperiksaOleh  string     `json:"DiperiksaOleh"`
	DisetujuiOleh  string     `json:"DisetujuiOleh"`
	Karyawan       []Karyawan `json:"Karyawan"`
}

// Repository interface
type AttendanceRawRepo interface {
	GetAttendanceReport(ctx context.Context, startDate, endDate, depart string) (*AttendanceReport, error)
}

// Usecase
type AttendanceRawUsecase struct {
	repo       AttendanceRawRepo
	pageToken  page_token.ProcessPageTokens
	log        *log.Helper
	sg         singleflight.Group
	deptClient departV1.DepartmentClient
	empClient  employersV1.EmployersClient
}

func NewAttendanceRawUsecase(
	repo AttendanceRawRepo,
	dept departV1.DepartmentClient,
	emp employersV1.EmployersClient,
	logger log.Logger,
) *AttendanceRawUsecase {
	return &AttendanceRawUsecase{
		repo:       repo,
		log:        log.NewHelper(log.With(logger, "module", "usecase/AttendanceRaw")),
		deptClient: dept,
		empClient:  emp,
	}
}

// GetAttendanceReport memanggil repo dan mengembalikan AttendanceReport
func (uc *AttendanceRawUsecase) GetAttendanceReport(
	ctx context.Context,
	startDate, endDate, depart string,
) (*AttendanceReport, error) {
	return uc.repo.GetAttendanceReport(ctx, startDate, endDate, depart)
}
