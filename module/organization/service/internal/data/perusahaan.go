package data

import (
	"context"
	"encoding/json" // Digunakan untuk serialisasi/deserialisasi JSON
	"fmt"
	"strconv"
	"time" // Digunakan untuk mengatur waktu kedaluwarsa cache

	"mall-go/module/organization/service/internal/biz"
	"mall-go/module/organization/service/internal/data/model/perusahaan"
	"mall-go/pkg/utils/pagination"

	v1 "mall-go/api/organization/service/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8" // Pastikan versi Redis client yang benar
)

// Pastikan antarmuka PerusahaanRepo diimpor dengan benar dari biz package
var _ biz.PerusahaanRepo = (*perusahaanRepo)(nil)

type perusahaanRepo struct {
	data        *Data
	Log         *log.Helper
	redisClient *redis.Client
}

func NewPerusahaanRepo(data *Data, logger log.Logger, rdb *redis.Client) biz.PerusahaanRepo {
	return &perusahaanRepo{
		data:        data,
		redisClient: rdb,
		Log:         log.NewHelper(log.With(logger, "module", "data/perusahaan")),
	}
}

func (r perusahaanRepo) CreatePerusahaan(ctx context.Context, b *v1.CreatePerusahaanRequest) (*biz.PerusahaanData, error) {
	po, err := r.data.db.Perusahaan.
		Create().
		SetKodePerusahaan(b.KodePerusahaan).
		SetNamaPerusahaan(b.NamaPerusahaan).
		SetKodeCabang(b.KodeCabang).
		SetCabang(b.Cabang).
		SetAlamat(b.Alamat).
		SetTelp(b.Telp).
		SetEmail(b.Email).
		SetCreateTime(time.Now()).
		SetCreatedBy("System").
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	r.redisClient.Del(ctx, "perusahaans:list")
	id := strconv.FormatInt(po.ID, 10)
	r.redisClient.Del(ctx, "perusahaan:"+id)
	r.redisClient.Del(ctx, "perusahaan_code:"+po.KodePerusahaan)
	r.redisClient.Del(ctx, "perusahaan_code:"+po.KodeCabang)

	return &biz.PerusahaanData{
		Id:             po.ID,
		KodePerusahaan: po.KodePerusahaan,
		NamaPerusahaan: po.NamaPerusahaan,
		KodeCabang:     po.KodeCabang,
		Cabang:         po.Cabang,
		Alamat:         po.Alamat,
		Telp:           *po.Telp,
		Email:          *po.Email,
		CreatedAt:      po.CreateTime,
		CreatedBy:      po.CreatedBy,
	}, nil
}

func (r *perusahaanRepo) UpdatePerusahaan(ctx context.Context, b *v1.UpdatePerusahaanRequest) (*biz.PerusahaanData, error) {
	po, err := r.data.db.Perusahaan.
		UpdateOneID(b.Id).
		SetKodePerusahaan(b.KodePerusahaan).
		SetNamaPerusahaan(b.NamaPerusahaan).
		SetKodeCabang(b.KodeCabang).
		SetCabang(b.Cabang).
		SetAlamat(b.Alamat).
		SetTelp(b.Telp).
		SetEmail(b.Email).
		SetUpdateTime(time.Now()).
		SetUpdatedBy("System").
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Invalidate cache setelah update
	r.redisClient.Del(ctx, "perusahaans:list")
	ids := strconv.FormatInt(po.ID, 10)
	r.redisClient.Del(ctx, "perusahaan:"+ids)                    // Invalidate cache spesifik berdasarkan ID
	r.redisClient.Del(ctx, "perusahaan_code:"+po.KodePerusahaan) // Invalidate cache spesifik berdasarkan kode
	r.redisClient.Del(ctx, "perusahaan_code:"+po.KodeCabang)

	return &biz.PerusahaanData{
		Id:             po.ID,
		KodePerusahaan: po.KodePerusahaan,
		NamaPerusahaan: po.NamaPerusahaan,
		KodeCabang:     po.KodeCabang,
		Cabang:         po.Cabang,
		Alamat:         po.Alamat,
		Telp:           *po.Telp,
		Email:          *po.Email,
		CreatedAt:      po.CreateTime,
		CreatedBy:      po.CreatedBy,
	}, nil
}

// func (r departmentRepo) DeleteDepartment(ctx context.Context, id int64) error {
// 	// Dapatkan detail departemen sebelum menghapus untuk menginvalidasi cache berdasarkan kode juga
// 	deptToDelete, err := r.data.db.Department.Get(ctx, id)
// 	if err != nil {
// 		return err // Atau log error dan tetap coba hapus jika tidak kritis
// 	}

// 	err = r.data.db.Department.
// 		DeleteOneID(id).Exec(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	// Invalidate cache setelah menghapus
// 	r.redisClient.Del(ctx, "departments:list")
// 	r.redisClient.Del(ctx, "department:"+string(id))
// 	// Asumsi ada cache 'departmentDetail' jika digunakan di tempat lain
// 	r.redisClient.Del(ctx, "departmentDetail:"+string(id))
// 	if deptToDelete != nil {
// 		r.redisClient.Del(ctx, "department_code:"+deptToDelete.DepartCode) // Invalidate cache by code
// 	}
// 	r.Log.Infof("Department dengan ID %d terhapus", id)
// 	return nil
// }

func (r perusahaanRepo) DeletePerusahaan(ctx context.Context, id int64) error {
	// Dapatkan detail departemen sebelum menghapus untuk menginvalidasi cache berdasarkan kode juga
	perToDelete, err := r.data.db.Perusahaan.Get(ctx, id)
	if err != nil {
		return err // Atau log error dan tetap coba hapus jika tidak kritis
	}

	err = r.data.db.Perusahaan.
		DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	// Invalidate cache setelah menghapus
	r.redisClient.Del(ctx, "perusahaans:list")
	ids := strconv.FormatInt(id, 10)
	r.redisClient.Del(ctx, "perusahaan:"+ids)
	// r.redisClient.Del(ctx, "perusahaan_code:"+po.KodePerusahaan)
	// r.redisClient.Del(ctx, "perusahaan_code:"+po.KodeCabang)
	if perToDelete != nil {
		r.redisClient.Del(ctx, "perusahaan:"+perToDelete.KodePerusahaan) // Invalidate cache by code
	}

	r.Log.Infof("Perusahaan dengan ID %d terhapus", id)
	return nil
}

func (r *perusahaanRepo) Count(ctx context.Context) (int, error) {
	// Pertimbangkan caching untuk count jika sering diakses dan tidak perlu sangat up-to-date
	dt, _ := r.data.db.Perusahaan.Query().Count(ctx)
	return dt, nil
}

func (r *perusahaanRepo) GetPerusahaan(ctx context.Context, kode_perusahaan string) (*biz.PerusahaanData, error) {
	cacheKey := "perusahaan_code:" + kode_perusahaan
	val, err := r.redisClient.Get(ctx, cacheKey).Result() // Coba ambil dari Redis
	if err == nil {
		var cachedPerusahaan biz.PerusahaanData
		// Perbaikan: Deklarasi unmarshalErr di luar kondisi if
		unmarshalErr := json.Unmarshal([]byte(val), &cachedPerusahaan)
		if unmarshalErr == nil {
			r.Log.Infof("Cache hit for perusahaan ID: %d", kode_perusahaan) // Perbaikan pesan log
			return &cachedPerusahaan, nil                                   // Langsung kembalikan jika cache hit
		}
		// Log jika gagal unmarshal dari cache
		r.Log.Warnf("Failed to unmarshal cached perusahaan for ID %d: %v", kode_perusahaan, unmarshalErr)
	} else if err != redis.Nil {
		// Log error Redis selain cache miss
		r.Log.Errorf("Redis GET error for perusahaan ID %d: %v", kode_perusahaan, err) // Perbaikan pesan log
	}

	// Jika tidak di cache atau gagal unmarshal, ambil dari database
	po, err := r.data.db.Perusahaan.Query().
		Where(perusahaan.KodePerusahaan(kode_perusahaan)).First(ctx)
	if err != nil {
		return nil, err
	}

	perusahaanBiz := &biz.PerusahaanData{
		Id:             po.ID,
		KodePerusahaan: po.KodePerusahaan,
		NamaPerusahaan: po.NamaPerusahaan,
		KodeCabang:     po.KodeCabang,
		Cabang:         po.Cabang,
		Alamat:         po.Alamat,
		Telp:           *po.Telp,
		Email:          *po.Email,
		CreatedAt:      po.CreateTime,
		CreatedBy:      po.CreatedBy,
	}

	// Simpan hasil ke Redis cache
	perusahaanBytes, marshalErr := json.Marshal(perusahaanBiz)
	if marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, perusahaanBytes, 3*time.Minute) // Atur waktu kedaluwarsa (misal 5 menit)
	} else {
		r.Log.Errorf("Failed to marshal perusahaan for caching: %v", marshalErr)
	}

	return perusahaanBiz, nil
}

func (r *perusahaanRepo) GetCabang(ctx context.Context, kode_cabang string) (*biz.PerusahaanData, error) {
	cacheKey := "perusahaan:" + kode_cabang
	val, err := r.redisClient.Get(ctx, cacheKey).Result() // Coba ambil dari Redis
	if err == nil {
		var cachedPerusahaan biz.PerusahaanData
		// Perbaikan: Deklarasi unmarshalErr di luar kondisi if
		unmarshalErr := json.Unmarshal([]byte(val), &cachedPerusahaan)
		if unmarshalErr == nil {
			r.Log.Infof("Cache hit for perusahaan ID: %d", kode_cabang) // Perbaikan pesan log
			return &cachedPerusahaan, nil                               // Langsung kembalikan jika cache hit
		}
		// Log jika gagal unmarshal dari cache
		r.Log.Warnf("Failed to unmarshal cached perusahaan for ID %d: %v", kode_cabang, unmarshalErr)
	} else if err != redis.Nil {
		// Log error Redis selain cache miss
		r.Log.Errorf("Redis GET error for perusahaan ID %d: %v", kode_cabang, err) // Perbaikan pesan log
	}

	// Jika tidak di cache atau gagal unmarshal, ambil dari database
	po, err := r.data.db.Perusahaan.Query().
		Where(perusahaan.KodeCabang(kode_cabang)).First(ctx)
	if err != nil {
		return nil, err
	}

	return &biz.PerusahaanData{
		Id:             po.ID,
		KodePerusahaan: po.KodePerusahaan,
		NamaPerusahaan: po.NamaPerusahaan,
		KodeCabang:     po.KodeCabang,
		Cabang:         po.Cabang,
		Alamat:         po.Alamat,
		Telp:           *po.Telp,
		Email:          *po.Email,
		CreatedAt:      po.CreateTime,
		CreatedBy:      po.CreatedBy,
	}, nil
}

func (r *perusahaanRepo) ListPerusahaan(ctx context.Context, pageNum, pageSize int64) ([]*biz.PerusahaanData, int, error) {
	cacheKey := fmt.Sprintf("perusahaans:list:page:%d:size:%d", pageNum, pageSize)
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedResponse struct {
			Perusahaans []*biz.PerusahaanData
			Total       int
		}
		if unmarshalErr := json.Unmarshal([]byte(val), &cachedResponse); unmarshalErr == nil {
			r.Log.Infof("Cache hit for ListPerusahaans page %d, size %d", pageNum, pageSize)
			return cachedResponse.Perusahaans, cachedResponse.Total, nil
		} else {
			r.Log.Warnf("Failed to unmarshal cached ListPerusahaans: %v", unmarshalErr)
		}
	} else if err != redis.Nil {
		r.Log.Errorf("Redis GET error for ListPerusahaans: %v", err)
	}

	query := r.data.db.Perusahaan.Query()
	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	pos, err := query.
		Offset(int(pagination.GetPageOffset(pageNum, pageSize))).
		Limit(int(pageSize)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*biz.PerusahaanData, 0, len(pos))
	for _, po := range pos {
		// var reportsTo string
		// if po.ReportsToPerusahaanID != nil {
		// 	reportsTo = *po.ReportsToPerusahaanID
		// }
		result = append(result, &biz.PerusahaanData{
			Id:             po.ID,
			KodePerusahaan: po.KodePerusahaan,
			NamaPerusahaan: po.NamaPerusahaan,
			KodeCabang:     po.KodeCabang,
			Cabang:         po.Cabang,
			Alamat:         po.Alamat,
			Telp:           *po.Telp,
			Email:          *po.Email,
			CreatedAt:      po.CreateTime,
			UpdatedAt:      po.UpdateTime,
			CreatedBy:      po.CreatedBy,
			UpdatedBy:      po.UpdatedBy,
		})
	}

	// Cache result
	responseToCache := struct {
		Perusahaans []*biz.PerusahaanData
		Total       int
	}{
		Perusahaans: result,
		Total:       total,
	}
	if responseBytes, marshalErr := json.Marshal(responseToCache); marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, responseBytes, 3*time.Minute)
	} else {
		r.Log.Errorf("Failed to marshal ListPerusahaans response for caching: %v", marshalErr)
	}

	return result, total, nil
}
