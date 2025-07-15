package data

import (
	"context"
	biometricV1 "mall-go/api/biometric/service/v1"
	departmentV1 "mall-go/api/department/service/v1"
	"mall-go/module/employers/service/internal/biz"
	"mall-go/pkg/utils/pagination"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

// var _ biz.BeerRepo = (*beerRepo)(nil)
var _ biz.EmployersRepo = (*employersRepo)(nil)

type employersRepo struct {
	data       *Data
	Log        *log.Helper
	bioClient  biometricV1.BiometricClient
	deptClient departmentV1.DepartmentClient
}

func NewEmployersRepo(data *Data, bio biometricV1.BiometricClient, dept departmentV1.DepartmentClient, logger log.Logger) biz.EmployersRepo {
	return &employersRepo{
		data:       data,
		bioClient:  bio,
		deptClient: dept,
		Log:        log.NewHelper(logger), // ✅ INI WAJIB ADA
	}
}

func (r *employersRepo) CreateEmployers(ctx context.Context, b *biz.Employers) (*biz.Employers, error) {
	po, err := r.data.db.Employers.
		Create().
		SetNosap(b.NoSap).
		SetNip(b.Nip).
		SetKaryacode(b.KaryaCode).
		SetKaryaname(b.KaryaName).
		SetDispName(b.DispName).
		SetPassMesin(b.PassMesin).
		SetRfidCard(b.RFIDCard).
		SetKodeFinger(b.Finger).
		SetDepartCode(b.Department).
		SetStatus(b.Status).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &biz.Employers{
		Id:         po.ID,
		NoSap:      *po.Nosap,
		Nip:        *po.Nip,
		KaryaCode:  *po.Karyacode,
		KaryaName:  po.Karyaname,
		DispName:   po.DispName,
		PassMesin:  *po.PassMesin,
		RFIDCard:   *po.RfidCard,
		Finger:     po.KodeFinger,
		Department: po.DepartCode,
		Status:     po.Status,
	}, nil

}

func (r *employersRepo) UpdateEmployers(ctx context.Context, b *biz.Employers) (*biz.Employers, error) {
	po, err := r.data.db.Employers.
		UpdateOneID(b.Id).
		SetNosap(b.NoSap).
		SetNip(b.Nip).
		SetKaryacode(b.KaryaCode).
		SetKaryaname(b.KaryaName).
		SetDispName(b.DispName).
		SetPassMesin(b.PassMesin).
		SetRfidCard(b.RFIDCard).
		SetKodeFinger(b.Finger).
		SetDepartCode(b.Department).
		SetStatus(b.Status).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &biz.Employers{
		Id:         po.ID,
		NoSap:      *po.Nosap,
		Nip:        *po.Nip,
		KaryaCode:  *po.Karyacode,
		KaryaName:  po.Karyaname,
		DispName:   po.DispName,
		PassMesin:  *po.PassMesin,
		RFIDCard:   *po.RfidCard,
		Finger:     po.KodeFinger,
		Department: po.DepartCode,
		Status:     po.Status,
	}, nil
}
func (r *employersRepo) GetEmployersID(ctx context.Context, id int64) (*biz.Employers, error) {
	po, err := r.data.db.Employers.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &biz.Employers{
		Id:         po.ID,
		NoSap:      *po.Nosap,
		Nip:        *po.Nip,
		KaryaCode:  *po.Karyacode,
		KaryaName:  po.Karyaname,
		DispName:   po.DispName,
		PassMesin:  *po.PassMesin,
		RFIDCard:   *po.RfidCard,
		Finger:     po.KodeFinger,
		Department: po.DepartCode,
		Status:     po.Status,
	}, nil
}

func (r *employersRepo) ListEmployers(ctx context.Context, pageNum, pageSize int64) ([]*biz.EmployerData, int, error) {
	query := r.data.db.Employers.Query()

	// Hitung total sebelum pagination
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	pos, err := r.data.db.Employers.Query().
		Offset(int(pagination.GetPageOffset(pageNum, pageSize))).
		Limit(int(pageSize)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	if r.Log != nil {
		r.Log.Warnf("ListEmployers result: %+v", pos)
	}
	rv := make([]*biz.EmployerData, 0, len(pos))
	for _, po := range pos {
		empData := &biz.EmployerData{
			Id:        po.ID,
			NoSap:     *po.Nosap,
			Nip:       *po.Nip,
			KaryaCode: *po.Karyacode,
			KaryaName: po.Karyaname,
			DispName:  po.DispName,
			PassMesin: *po.PassMesin,
			RFIDCard:  *po.RfidCard,
			Status:    po.Status,
		}

		// Get finger data
		if po.KodeFinger != "" {
			fingerResp, err := r.bioClient.GetFingerByKode(ctx, &biometricV1.GetFingerByKodeRequest{
				Fingercode: po.KodeFinger,
			})
			if err == nil {
				empData.Finger = []biz.FingerData{
					{
						Fingercode: fingerResp.Fingercode,
						Finger0:    fingerResp.Finger0,
						Finger1:    fingerResp.Finger1,
						Finger2:    fingerResp.Finger2,
						Finger3:    fingerResp.Finger3,
						Finger4:    fingerResp.Finger4,
						Finger5:    fingerResp.Finger5,
						Finger6:    fingerResp.Finger6,
						Finger7:    fingerResp.Finger7,
						Finger8:    fingerResp.Finger8,
						Finger9:    fingerResp.Finger9,
					},
				}
			} else {
				r.Log.Warnf("finger not found for kode: %s, err: %v", po.KodeFinger, err)
			}
		}

		// Department
		r.Log.Infow("", po.DepartCode)
		if po.DepartCode != "" {
			deptResp, err := r.deptClient.GetDepartmentCode(ctx, &departmentV1.GetDepartmentCodeRequest{
				DepartCode: po.DepartCode,
			})
			if err != nil {
				r.Log.Warnf("Failed to get department for depart_code=%s: %v", po.DepartCode, err)
			} else {
				empData.Department = []biz.DepartData{
					{
						DepartCode: deptResp.DepartCode,
						DepartName: deptResp.DepartName,
					},
				}
				// r.Log.Infof("DEPARTMENT: %s, %s,%s", empData.Department, deptResp.DepartCode, deptResp.DepartName)
			}
		}

		rv = append(rv, empData)
	}

	return rv, total, nil
}

func (r *employersRepo) Count(ctx context.Context) (int, error) {
	dt, _ := r.data.db.Employers.Query().Count(ctx)
	// fmt.Println("INI PL DT: ")
	return dt, nil
}
func (r *employersRepo) ListEmployersNext(ctx context.Context, start, end int32) ([]*biz.Employers, error) {
	pos, err := r.data.db.Employers.Query().
		Offset(int(start)).
		Limit(int(end - start)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	rv := make([]*biz.Employers, 0, len(pos))
	for _, po := range pos {
		rv = append(rv, &biz.Employers{
			Id:         po.ID,
			NoSap:      *po.Nosap,
			Nip:        *po.Nip,
			KaryaCode:  *po.Karyacode,
			KaryaName:  po.Karyaname,
			DispName:   po.DispName,
			PassMesin:  *po.PassMesin,
			RFIDCard:   *po.RfidCard,
			Finger:     po.KodeFinger,
			Department: po.DepartCode,
			Status:     po.Status,
		})
	}
	return rv, nil
}

func (r *employersRepo) DeleteEmployers(ctx context.Context, id int64) error {
	err := r.data.db.Employers.
		DeleteOneID(id).
		Exec(ctx)

	if err != nil {
		return err
	}

	r.Log.Infof("employer with ID %d deleted", id)
	return nil
}

func (r *employersRepo) GetEmployerDetail(ctx context.Context, id int64) (*biz.EmployerData, error) {
	emp, err := r.data.db.Employers.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get finger data
	fingerResp, err := r.bioClient.GetFingerByKode(ctx, &biometricV1.GetFingerByKodeRequest{
		Fingercode: emp.KodeFinger,
	})
	if err != nil {
		r.Log.Errorf("Failed to fetch finger: %v", err)
	}

	// Get department data
	deptResp, err := r.deptClient.GetDepartmentCode(ctx, &departmentV1.GetDepartmentCodeRequest{
		DepartCode: emp.DepartCode,
	})
	if err != nil {
		r.Log.Errorf("Failed to fetch department: %v", err)
	}

	return &biz.EmployerData{
		Id:        emp.ID,
		NoSap:     *emp.Nosap,
		Nip:       *emp.Nip,
		KaryaCode: *emp.Karyacode,
		KaryaName: emp.Karyaname,
		DispName:  emp.DispName,
		PassMesin: *emp.PassMesin,
		RFIDCard:  *emp.RfidCard,
		Status:    emp.Status,
		// CreatedAt: emp.created_at.String(),
		// UpdatedAt: emp.UpdatedAt.String(),
		Finger: []biz.FingerData{
			{
				Fingercode: fingerResp.Fingercode,
				Finger0:    fingerResp.Finger0,
				Finger1:    fingerResp.Finger1,
				Finger2:    fingerResp.Finger2,
				Finger3:    fingerResp.Finger3,
				Finger4:    fingerResp.Finger4,
				Finger5:    fingerResp.Finger5,
				Finger6:    fingerResp.Finger6,
				Finger7:    fingerResp.Finger7,
				Finger8:    fingerResp.Finger8,
				Finger9:    fingerResp.Finger9,
			},
		},
		Department: []biz.DepartData{
			{
				DepartCode: deptResp.DepartCode,
				DepartName: deptResp.DepartName,
			},
		},
	}, nil
}

func safeTrim(s *string) string {
	if s == nil {
		return ""
	}
	return strings.TrimSpace(*s)
}
