package data

import (
	"context"
	"encoding/json" // Digunakan untuk serialisasi/deserialisasi JSON
	"time"          // Digunakan untuk mengatur waktu kedaluwarsa cache

	"mall-go/module/department/service/internal/biz"
	"mall-go/module/department/service/internal/data/model/department"
	"mall-go/pkg/utils/pagination"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8" // Pastikan versi Redis client yang benar
)

// Pastikan antarmuka DepartmentRepo diimpor dengan benar dari biz package
var _ biz.DepartmentRepo = (*departmentRepo)(nil)

type departmentRepo struct {
	data        *Data
	Log         *log.Helper
	redisClient *redis.Client
}

func NewDepartmentRepo(data *Data, logger log.Logger, rdb *redis.Client) biz.DepartmentRepo {
	return &departmentRepo{
		data:        data,
		redisClient: rdb,
		Log:         log.NewHelper(log.With(logger, "module", "data/department")),
	}
}

func (r *departmentRepo) GetDepartmentID(ctx context.Context, id int64) (*biz.Department, error) {
	cacheKey := "department:" + string(id)
	val, err := r.redisClient.Get(ctx, cacheKey).Result() // Coba ambil dari Redis
	if err == nil {
		var cachedDepartment biz.Department
		// Perbaikan: Deklarasi unmarshalErr di luar kondisi if
		unmarshalErr := json.Unmarshal([]byte(val), &cachedDepartment)
		if unmarshalErr == nil {
			r.Log.Infof("Cache hit for department ID: %d", id) // Perbaikan pesan log
			return &cachedDepartment, nil                      // Langsung kembalikan jika cache hit
		}
		// Log jika gagal unmarshal dari cache
		r.Log.Warnf("Failed to unmarshal cached department for ID %d: %v", id, unmarshalErr)
	} else if err != redis.Nil {
		// Log error Redis selain cache miss
		r.Log.Errorf("Redis GET error for department ID %d: %v", id, err) // Perbaikan pesan log
	}

	// Jika tidak di cache atau gagal unmarshal, ambil dari database
	po, err := r.data.db.Department.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	departmentBiz := &biz.Department{
		Id:         po.ID,
		DepartCode: po.DepartCode,
		DepartName: po.DepartName,
		Status:     po.Status,
		Ket:        po.Ket,
	}

	// Simpan hasil ke Redis cache
	departmentBytes, marshalErr := json.Marshal(departmentBiz)
	if marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, departmentBytes, 5*time.Minute) // Atur waktu kedaluwarsa (misal 5 menit)
	} else {
		r.Log.Errorf("Failed to marshal department for caching: %v", marshalErr)
	}

	return departmentBiz, nil
}

func (r *departmentRepo) GetDepartmentCode(ctx context.Context, Dcode string) (*biz.Department, error) {
	// Tambahkan caching untuk GetDepartmentCode juga jika sering diakses
	cacheKey := "department_code:" + strings.TrimSpace(Dcode)
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedDepartment biz.Department
		unmarshalErr := json.Unmarshal([]byte(val), &cachedDepartment)
		if unmarshalErr == nil {
			r.Log.Infof("Cache hit for department code: %s", Dcode)
			return &cachedDepartment, nil
		}
		r.Log.Warnf("Failed to unmarshal cached department for code %s: %v", Dcode, unmarshalErr)
	} else if err != redis.Nil {
		r.Log.Errorf("Redis GET error for department code %s: %v", Dcode, err)
	}

	po, err := r.data.db.Department.Query().
		Where(department.DepartCodeIn(strings.TrimSpace(Dcode))).
		First(ctx)
	if err != nil {
		return nil, err
	}

	departmentBiz := &biz.Department{
		Id:         po.ID,
		DepartCode: po.DepartCode,
		DepartName: po.DepartName,
		Status:     po.Status,
		Ket:        po.Ket,
	}

	// Simpan hasil ke Redis cache
	departmentBytes, marshalErr := json.Marshal(departmentBiz)
	if marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, departmentBytes, 5*time.Minute)
	} else {
		r.Log.Errorf("Failed to marshal department for caching: %v", marshalErr)
	}

	return departmentBiz, nil
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
	// Invalidate cache setelah membuat departemen baru
	r.redisClient.Del(ctx, "departments:list")
	r.redisClient.Del(ctx, "department:"+string(po.ID))      // Invalidate cache spesifik
	r.redisClient.Del(ctx, "department_code:"+po.DepartCode) // Invalidate cache by code

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

	// Invalidate cache setelah update
	r.redisClient.Del(ctx, "departments:list")
	r.redisClient.Del(ctx, "department:"+string(po.ID))      // Invalidate cache spesifik berdasarkan ID
	r.redisClient.Del(ctx, "department_code:"+po.DepartCode) // Invalidate cache spesifik berdasarkan kode

	return &biz.Department{
		Id:         po.ID,
		DepartCode: po.DepartCode,
		DepartName: po.DepartName,
		Status:     po.Status,
		Ket:        po.Ket,
	}, nil
}

func (r departmentRepo) DeleteDepartment(ctx context.Context, id int64) error {
	// Dapatkan detail departemen sebelum menghapus untuk menginvalidasi cache berdasarkan kode juga
	deptToDelete, err := r.data.db.Department.Get(ctx, id)
	if err != nil {
		return err // Atau log error dan tetap coba hapus jika tidak kritis
	}

	err = r.data.db.Department.
		DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	// Invalidate cache setelah menghapus
	r.redisClient.Del(ctx, "departments:list")
	r.redisClient.Del(ctx, "department:"+string(id))
	// Asumsi ada cache 'departmentDetail' jika digunakan di tempat lain
	r.redisClient.Del(ctx, "departmentDetail:"+string(id))
	if deptToDelete != nil {
		r.redisClient.Del(ctx, "department_code:"+deptToDelete.DepartCode) // Invalidate cache by code
	}
	r.Log.Infof("Department dengan ID %d terhapus", id)
	return nil
}

func (r *departmentRepo) Count(ctx context.Context) (int, error) {
	// Pertimbangkan caching untuk count jika sering diakses dan tidak perlu sangat up-to-date
	dt, _ := r.data.db.Department.Query().Count(ctx)
	return dt, nil
}

func (r *departmentRepo) ListDepartment(ctx context.Context, pageNum, pageSize int64) ([]*biz.Department, int, error) {
	cacheKey := "departments:list:page:" + string(pageNum) + ":size:" + string(pageSize)
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedResponse struct {
			Departments []*biz.Department
			Total       int
		}
		// Perbaikan: Hapus deklarasi ganda `unmarshalErr` dan `cachedDepartment` yang tidak relevan
		unmarshalErr := json.Unmarshal([]byte(val), &cachedResponse)
		if unmarshalErr == nil {
			r.Log.Infof("Cache hit for ListDepartments page %d, size %d", pageNum, pageSize)
			return cachedResponse.Departments, cachedResponse.Total, nil
		}
		// Log jika gagal unmarshal dari cache
		r.Log.Warnf("Failed to unmarshal cached ListDepartments: %v", unmarshalErr)
	} else if err != redis.Nil {
		// Log error Redis selain cache miss
		r.Log.Errorf("Redis GET error for ListDepartments: %v", err)
	}

	// Jika tidak di cache atau gagal unmarshal, ambil dari database
	query := r.data.db.Department.Query()

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
	if r.Log != nil {
		r.Log.Warnf("ListDepartments result from DB: %+v", pos)
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

	// Simpan hasil ke Redis cache
	responseToCache := struct {
		Departments []*biz.Department
		Total       int
	}{
		Departments: rv,
		Total:       total,
	}
	responseBytes, marshalErr := json.Marshal(responseToCache)
	if marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, responseBytes, 5*time.Minute) // Atur waktu kedaluwarsa
	} else {
		r.Log.Errorf("Failed to marshal ListDepartments response for caching: %v", marshalErr)
	}

	return rv, total, nil
}

func (r *departmentRepo) ListDepartmentNext(ctx context.Context, start, end int32) ([]*biz.Department, error) {
	// Pertimbangkan caching di sini juga jika sering digunakan
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
