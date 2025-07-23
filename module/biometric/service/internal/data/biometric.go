package data

import (
	"context"
	"encoding/json"
	"mall-go/module/biometric/service/internal/biz"
	"mall-go/module/biometric/service/internal/data/model/biometric"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

// var _ biz.EmployersRepo = (*employersRepo)(nil)
var _ biz.BiometricRepo = (*biometricRepo)(nil)

type biometricRepo struct {
	data        *Data
	Log         *log.Helper
	redisClient *redis.Client
}

func NewBiometricRepo(data *Data, logger log.Logger, rdb *redis.Client) biz.BiometricRepo {
	return &biometricRepo{
		data:        data,
		redisClient: rdb,
		Log:         log.NewHelper(log.With(logger, "module", "data/biometri")),
	}
}

func (r *biometricRepo) GetFingerByID(ctx context.Context, id int64) (*biz.Biometric, error) {
	cacheKey := "biometric:" + string(id)
	val, err := r.redisClient.Get(ctx, cacheKey).Result() // Coba ambil dari Redis
	if err == nil {
		var cachedBiometric biz.Biometric
		// Perbaikan: Deklarasi unmarshalErr di luar kondisi if
		unmarshalErr := json.Unmarshal([]byte(val), &cachedBiometric)
		if unmarshalErr == nil {
			r.Log.Infof("Cache hit for biometric ID: %d", id) // Perbaikan pesan log
			return &cachedBiometric, nil                      // Langsung kembalikan jika cache hit
		}
		// Log jika gagal unmarshal dari cache
		r.Log.Warnf("Failed to unmarshal cached biometric for ID %d: %v", id, unmarshalErr)
	} else if err != redis.Nil {
		// Log error Redis selain cache miss
		r.Log.Errorf("Redis GET error for biometric ID %d: %v", id, err) // Perbaikan pesan log
	}

	po, err := r.data.db.Biometric.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &biz.Biometric{
		Id:         po.ID,
		Fingercode: po.Fingercode,
		Finger0:    po.Finger0,
		Finger1:    po.Finger1,
		Finger2:    po.Finger2,
		Finger3:    po.Finger3,
		Finger4:    po.Finger4,
		Finger5:    po.Finger5,
		Finger6:    po.Finger6,
		Finger7:    po.Finger7,
		Finger8:    po.Finger8,
		Finger9:    po.Finger9,
	}, nil
}
func (r *biometricRepo) GetFingerByKode(ctx context.Context, kode string) (*biz.Biometric, error) {
	cacheKey := "biometric_code:" + strings.TrimSpace(kode)
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedDepartment biz.Biometric
		unmarshalErr := json.Unmarshal([]byte(val), &cachedDepartment)
		if unmarshalErr == nil {
			r.Log.Infof("Cache hit for biometric code: %s", kode)
			return &cachedDepartment, nil
		}
		r.Log.Warnf("Failed to unmarshal cached biometric for code %s: %v", kode, unmarshalErr)
	} else if err != redis.Nil {
		r.Log.Errorf("Redis GET error for biometric code %s: %v", kode, err)
	}

	po, err := r.data.db.Biometric.Query().
		Where(biometric.FingercodeEQ(kode)). // LIKE %kode%
		Limit(1).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return &biz.Biometric{
		Id:         po.ID,
		Fingercode: po.Fingercode,
		Finger0:    po.Finger0,
		Finger1:    po.Finger1,
		Finger2:    po.Finger2,
		Finger3:    po.Finger3,
		Finger4:    po.Finger4,
		Finger5:    po.Finger5,
		Finger6:    po.Finger6,
		Finger7:    po.Finger7,
		Finger8:    po.Finger8,
		Finger9:    po.Finger9,
	}, nil
}
func (r biometricRepo) CreateFinger(ctx context.Context, b *biz.Biometric) (*biz.Biometric, error) {
	po, err := r.data.db.Biometric.
		Create().
		SetFingercode(b.Fingercode).
		SetFinger0(b.Finger0).
		SetFinger1(b.Finger1).
		SetFinger2(b.Finger2).
		SetFinger3(b.Finger3).
		SetFinger4(b.Finger4).
		SetFinger5(b.Finger5).
		SetFinger6(b.Finger6).
		SetFinger7(b.Finger7).
		SetFinger8(b.Finger8).
		SetFinger9(b.Finger9).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	r.redisClient.Del(ctx, "biometrics:list")
	r.redisClient.Del(ctx, "biometric:"+string(po.ID))      // Invalidate cache spesifik
	r.redisClient.Del(ctx, "biometric_code:"+po.Fingercode) // Invalidate cache by code

	return &biz.Biometric{
		Id:         po.ID,
		Fingercode: po.Fingercode,
		Finger0:    po.Finger0,
		Finger1:    po.Finger1,
		Finger2:    po.Finger2,
		Finger3:    po.Finger3,
		Finger4:    po.Finger4,
		Finger5:    po.Finger5,
		Finger6:    po.Finger6,
		Finger7:    po.Finger7,
		Finger8:    po.Finger8,
		Finger9:    po.Finger9,
	}, nil
}
func (r *biometricRepo) UpdateFinger(ctx context.Context, b *biz.Biometric) (*biz.Biometric, error) {
	po, err := r.data.db.Biometric.
		UpdateOneID(b.Id).
		SetFingercode(b.Fingercode).
		SetFinger0(b.Finger0).
		SetFinger1(b.Finger1).
		SetFinger2(b.Finger2).
		SetFinger3(b.Finger3).
		SetFinger4(b.Finger4).
		SetFinger5(b.Finger5).
		SetFinger6(b.Finger6).
		SetFinger7(b.Finger7).
		SetFinger8(b.Finger8).
		SetFinger9(b.Finger9).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	r.redisClient.Del(ctx, "biometrics:list")
	r.redisClient.Del(ctx, "biometric:"+string(po.ID))      // Invalidate cache spesifik
	r.redisClient.Del(ctx, "biometric_code:"+po.Fingercode) // Invalidate cache by code

	return &biz.Biometric{
		Id:         po.ID,
		Fingercode: po.Fingercode,
		Finger0:    po.Finger0,
		Finger1:    po.Finger1,
		Finger2:    po.Finger2,
		Finger3:    po.Finger3,
		Finger4:    po.Finger4,
		Finger5:    po.Finger5,
		Finger6:    po.Finger6,
		Finger7:    po.Finger7,
		Finger8:    po.Finger8,
		Finger9:    po.Finger9,
	}, nil
}

func (r biometricRepo) DeleteFinger(ctx context.Context, id int64) error {
	err := r.data.db.Biometric.
		DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	r.redisClient.Del(ctx, "biometrics:list")
	r.redisClient.Del(ctx, "biometric:"+string(id))      // Invalidate cache spesifik
	r.redisClient.Del(ctx, "biometric_code:"+string(id)) // Invalidate cache by code

	r.Log.Infof("Biometric dengan ID %d terhapus", id)
	return nil
}
