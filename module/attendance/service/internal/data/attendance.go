package data

import (
	"context"
	employersV1 "mall-go/api/employers/service/v1"

	// biometricV1 "mall-go/api/biometric/service/v1"
	// departmentV1 "mall-go/api/department/service/v1"
	"mall-go/module/attendance/service/internal/biz"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

// var _ biz.BeerRepo = (*beerRepo)(nil)
var _ biz.AttendanceRepo = (*attendanceRepo)(nil)

type attendanceRepo struct {
	data      *Data
	Log       *log.Helper
	empClient employersV1.EmployersClient
	// deptClient departmentV1.DepartmentClient
}

func NewAttendanceRepo(data *Data, emp employersV1.EmployersClient, logger log.Logger) biz.AttendanceRepo {
	return &attendanceRepo{
		data:      data,
		empClient: emp,
		// deptClient: dept,
		Log: log.NewHelper(logger), // ✅ INI WAJIB ADA
	}
}

func safeTrim(s *string) string {
	if s == nil {
		return ""
	}
	return strings.TrimSpace(*s)
}
func (r attendanceRepo) CreateClockIn(ctx context.Context, b *biz.Attendance) (*biz.Attendance, error) {
	po, err := r.data.db.Attendance.
		Create().
		SetKaryacode(b.KaryaCode).
		SetClockInTime(b.ClockInTime).
		// Set(b.KaryaCode).

		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.Attendance{
		KaryaCode:   po.Karyacode,
		ClockInTime: po.ClockInTime,
	}, nil
}
func (r attendanceRepo) UpdateClockOut(ctx context.Context, karyacode string) (*biz.Attendance, error) {
	// return panic("IL"), nil
	panic("implement me")
}

func (r attendanceRepo) SetStatus(ctx context.Context, karyacode string) (*biz.Attendance, error) {
	panic("implement me")
}
func (r attendanceRepo) CountEffectiveWorkDurationMinutes(ctx context.Context, karyacode string) (*biz.Attendance, error) {
	panic("ok")
}
func (r attendanceRepo) CountOvertimeMinutes(ctx context.Context, karyacode string) (*biz.Attendance, error) {
	panic("error")
}
func (r attendanceRepo) SetNotesOvertime(ctx context.Context, karyacode string) (*biz.Attendance, error) {
	panic("implement me")
}
