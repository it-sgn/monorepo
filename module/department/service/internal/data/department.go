package data

import (
	"context"
	"mall-go/module/department/service/internal/biz"
	"mall-go/module/department/service/internal/data/model/department"
	"mall-go/pkg/utils/pagination"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

// var _ biz.departmentRepo = (*departmentRepo)(nil)
var _ biz.DepartmentRepo = (*departmentRepo)(nil)

type departmentRepo struct {
	data *Data
	Log  *log.Helper
}

func NewDepartmentRepo(data *Data, logger log.Logger) biz.DepartmentRepo {
	return &departmentRepo{
		data: data,
		Log:  log.NewHelper(log.With(logger, "module", "data/department")),
	}
}

func (r *departmentRepo) GetDepartmentID(ctx context.Context, id int64) (*biz.Department, error) {
	po, err := r.data.db.Department.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &biz.Department{
		Id:         po.ID,
		DepartCode: po.DepartCode,
		DepartName: po.DepartName,
		Status:     po.Status,
		Ket:        po.Ket,
	}, nil
}

// func (r *departmentRepo) GetDepartmentCode(ctx context.Context, Dcodes []string) ([]*biz.Department, error) {
// 	// Bersihkan spasi dari semua kode
// 	cleaned := make([]string, 0, len(Dcodes))
// 	for _, code := range Dcodes {
// 		cleaned = append(cleaned, strings.TrimSpace(code))
// 	}

// 	// Query pakai DepartCodeIn
// 	pos, err := r.data.db.Department.Query().
// 		Where(department.DepartCodeIn(cleaned...)).
// 		All(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Mapping ke []*biz.Department
// 	result := make([]*biz.Department, 0, len(pos))
// 	for _, po := range pos {
// 		result = append(result, &biz.Department{
// 			Id:         po.ID,
// 			DepartCode: po.DepartCode,
// 			DepartName: po.DepartName,
// 			Status:     po.Status,
// 			Ket:        po.Ket,
// 		})
// 	}

// 	return result, nil
// }

func (r *departmentRepo) GetDepartmentCode(ctx context.Context, Dcode string) (*biz.Department, error) {
	po, err := r.data.db.Department.Query().
		Where(department.DepartCodeIn(strings.TrimSpace(Dcode))).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.Department{
		Id:         po.ID,
		DepartCode: po.DepartCode,
		DepartName: po.DepartName,
		Status:     po.Status,
		Ket:        po.Ket,
	}, nil
}

func (r departmentRepo) CreateDepartment(ctx context.Context, b *biz.Department) (*biz.Department, error) {
	po, err := r.data.db.Department.
		Create().
		SetDepartCode(b.DepartCode).
		SetDepartName(b.DepartName).
		SetStatus(b.Status).
		SetKet(b.Ket).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.Department{
		Id:         po.ID,
		DepartCode: po.DepartCode,
		DepartName: po.DepartName,
		Status:     po.Status,
		Ket:        po.Ket,
	}, nil
}
func (r *departmentRepo) UpdateDepartment(ctx context.Context, b *biz.Department) (*biz.Department, error) {
	po, err := r.data.db.Department.
		UpdateOneID(b.Id).
		SetDepartCode(b.DepartCode).
		SetDepartName(b.DepartName).
		SetStatus(b.Status).
		SetKet(b.Ket).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.Department{
		Id:         po.ID,
		DepartCode: po.DepartCode,
		DepartName: po.DepartName,
		Status:     po.Status,
		Ket:        po.Ket,
	}, nil
}

func (r departmentRepo) DeleteDepartment(ctx context.Context, id int64) error {
	err := r.data.db.Department.
		DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	r.Log.Infof("Department dengan ID %d terhapus", id)
	return nil
}

func (r *departmentRepo) Count(ctx context.Context) (int, error) {
	dt, _ := r.data.db.Department.Query().Count(ctx)
	// fmt.Println("INI PL DT: ")
	return dt, nil
}

func (r *departmentRepo) ListDepartment(ctx context.Context, pageNum, pageSize int64) ([]*biz.Department, int, error) {
	query := r.data.db.Department.Query()

	// Hitung total sebelum pagination
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	pos, err := r.data.db.Department.Query().
		Offset(int(pagination.GetPageOffset(pageNum, pageSize))).
		Limit(int(pageSize)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}
	rv := make([]*biz.Department, 0, len(pos))
	for _, po := range pos {
		rv = append(rv, &biz.Department{
			Id:         po.ID,
			DepartCode: po.DepartCode,
			DepartName: po.DepartName,
			Status:     po.Status,
			Ket:        po.Ket,
		})
	}
	return rv, total, nil
}

func (r *departmentRepo) ListDepartmentNext(ctx context.Context, start, end int32) ([]*biz.Department, error) {
	pos, err := r.data.db.Department.Query().
		Offset(int(start)).
		Limit(int(end - start)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	rv := make([]*biz.Department, 0, len(pos))
	for _, po := range pos {
		rv = append(rv, &biz.Department{
			Id:         po.ID,
			DepartCode: po.DepartCode,
			DepartName: po.DepartName,
			Status:     po.Status,
			Ket:        po.Ket,
		})
	}
	return rv, nil
}
