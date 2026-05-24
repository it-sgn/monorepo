package data

import (
	"context"
	"encoding/json" // Digunakan untuk serialisasi/deserialisasi JSON
	"fmt"
	"strconv"
	"time" // Digunakan untuk mengatur waktu kedaluwarsa cache

	"mall-go/module/organization/service/internal/biz"
	"mall-go/module/organization/service/internal/data/model/position"
	"mall-go/pkg/utils/pagination"
	"strings"

	v1 "mall-go/api/organization/service/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8" // Pastikan versi Redis client yang benar
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Pastikan antarmuka PositionRepo diimpor dengan benar dari biz package
var _ biz.PositionRepo = (*positionRepo)(nil)

type positionRepo struct {
	data        *Data
	Log         *log.Helper
	redisClient *redis.Client
}

func NewPositionRepo(data *Data, logger log.Logger, rdb *redis.Client) biz.PositionRepo {
	return &positionRepo{
		data:        data,
		redisClient: rdb,
		Log:         log.NewHelper(log.With(logger, "module", "data/position")),
	}
}

func (r *positionRepo) GetPosition(ctx context.Context, position_code string) (*v1.PositionResponse, error) {
	cacheKey := "position_code:" + position_code
	val, err := r.redisClient.Get(ctx, cacheKey).Result() // Coba ambil dari Redis
	if err == nil {
		var cachedPosition v1.PositionResponse
		// Perbaikan: Deklarasi unmarshalErr di luar kondisi if
		unmarshalErr := json.Unmarshal([]byte(val), &cachedPosition)
		if unmarshalErr == nil {
			r.Log.Infof("Cache hit for position ID: %d", position_code) // Perbaikan pesan log
			return &cachedPosition, nil                                 // Langsung kembalikan jika cache hit
		}
		// Log jika gagal unmarshal dari cache
		r.Log.Warnf("Failed to unmarshal cached position for ID %d: %v", position_code, unmarshalErr)
	} else if err != redis.Nil {
		// Log error Redis selain cache miss
		r.Log.Errorf("Redis GET error for position ID %d: %v", position_code, err) // Perbaikan pesan log
	}

	// Jika tidak di cache atau gagal unmarshal, ambil dari database
	po, err := r.data.db.Position.Query().
		Where(position.PositionCode(position_code)).First(ctx)
	if err != nil {
		return nil, err
	}

	positionBiz := &v1.PositionResponse{
		Position: &v1.Position{
			Id:                  int64(po.ID),
			PositionCode:        po.PositionCode,
			Name:                po.Name,
			RoleName:            po.RoleName,
			DepartmentCode:      po.DepartmentCode,
			ReportsToPositionId: *po.ReportsToPositionID,
			CreatedAt:           timestamppb.New(po.CreateTime),
			UpdatedAt:           nil, // or timestamppb.New(po.UpdateTime) if available
			CreatedBy:           "",  // fill if available
			UpdatedBy:           "",  // fill if available
		},
	}

	// Simpan hasil ke Redis cache
	positionBytes, marshalErr := json.Marshal(positionBiz)
	if marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, positionBytes, 3*time.Minute) // Atur waktu kedaluwarsa (misal 5 menit)
	} else {
		r.Log.Errorf("Failed to marshal position for caching: %v", marshalErr)
	}

	return positionBiz, nil
}

func (r *positionRepo) GetPositionID(ctx context.Context, id int64) (*v1.PositionResponse, error) {
	ids := strconv.FormatInt(id, 10)
	cacheKey := "position:" + ids
	val, err := r.redisClient.Get(ctx, cacheKey).Result() // Coba ambil dari Redis
	if err == nil {
		var cachedPosition v1.PositionResponse
		// Perbaikan: Deklarasi unmarshalErr di luar kondisi if
		unmarshalErr := json.Unmarshal([]byte(val), &cachedPosition)
		if unmarshalErr == nil {
			r.Log.Infof("Cache hit for position ID: %d", id) // Perbaikan pesan log
			return &cachedPosition, nil                      // Langsung kembalikan jika cache hit
		}
		// Log jika gagal unmarshal dari cache
		r.Log.Warnf("Failed to unmarshal cached position for ID %d: %v", id, unmarshalErr)
	} else if err != redis.Nil {
		// Log error Redis selain cache miss
		r.Log.Errorf("Redis GET error for position ID %d: %v", id, err) // Perbaikan pesan log
	}

	// Jika tidak di cache atau gagal unmarshal, ambil dari database
	po, err := r.data.db.Position.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	positionBiz := &v1.PositionResponse{
		Position: &v1.Position{
			Id:                  int64(po.ID),
			PositionCode:        po.PositionCode,
			Name:                po.Name,
			RoleName:            po.RoleName,
			DepartmentCode:      po.DepartmentCode,
			ReportsToPositionId: *po.ReportsToPositionID,
			CreatedAt:           timestamppb.New(po.CreateTime),
			UpdatedAt:           nil, // or timestamppb.New(po.UpdateTime) if available
			CreatedBy:           "",  // fill if available
			UpdatedBy:           "",  // fill if available
		},
	}

	// Simpan hasil ke Redis cache
	positionBytes, marshalErr := json.Marshal(positionBiz)
	if marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, positionBytes, 3*time.Minute) // Atur waktu kedaluwarsa (misal 5 menit)
	} else {
		r.Log.Errorf("Failed to marshal position for caching: %v", marshalErr)
	}

	return positionBiz, nil
}

func (r *positionRepo) GetPositionCode(ctx context.Context, position_code string) (*v1.Position, error) {
	// Tambahkan caching untuk GetPositionCode juga jika sering diakses
	cacheKey := "position_code:" + strings.TrimSpace(position_code)
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedPosition v1.Position
		unmarshalErr := json.Unmarshal([]byte(val), &cachedPosition)
		if unmarshalErr == nil {
			r.Log.Infof("Cache hit for position code: %s", position_code)
			return &cachedPosition, nil
		}
		r.Log.Warnf("Failed to unmarshal cached position for code %s: %v", position_code, unmarshalErr)
	} else if err != redis.Nil {
		r.Log.Errorf("Redis GET error for position code %s: %v", position_code, err)
	}

	po, err := r.data.db.Position.Query().
		Where(position.PositionCode(strings.TrimSpace(position_code))).
		First(ctx)
	if err != nil {
		return nil, err
	}

	positionBiz := &v1.Position{
		Id:                  int64(po.ID),
		PositionCode:        po.PositionCode,
		Name:                po.Name,
		RoleName:            po.RoleName,
		DepartmentCode:      po.DepartmentCode,
		ReportsToPositionId: *po.ReportsToPositionID,
		CreatedAt:           timestamppb.New(po.CreateTime),
	}

	// Simpan hasil ke Redis cache
	positionBytes, marshalErr := json.Marshal(positionBiz)
	if marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, positionBytes, 3*time.Minute)
	} else {
		r.Log.Errorf("Failed to marshal position for caching: %v", marshalErr)
	}

	return positionBiz, nil
}
func (r positionRepo) CreatePosition(ctx context.Context, b *v1.CreatePositionRequest) (*v1.PositionResponse, error) {
	po, err := r.data.db.Position.
		Create().
		SetPositionCode(b.PositionCode).
		SetName(b.Name).
		SetRoleName(b.RoleName).
		SetDepartmentCode(b.DepartmentCode).
		SetNillableReportsToPositionID(&b.ReportsToPositionId).
		SetCreateTime(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	ids := strconv.FormatInt(po.ID, 10)
	r.redisClient.Del(ctx, "positions:list")
	r.redisClient.Del(ctx, "position:"+ids)
	r.redisClient.Del(ctx, "position_code:"+po.PositionCode)

	var reportsTo string
	if po.ReportsToPositionID != nil {
		reportsTo = *po.ReportsToPositionID
	}

	return &v1.PositionResponse{
		Position: &v1.Position{
			Id:                  int64(po.ID),
			PositionCode:        po.PositionCode,
			Name:                po.Name,
			RoleName:            po.RoleName,
			DepartmentCode:      po.DepartmentCode,
			ReportsToPositionId: reportsTo,
			CreatedAt:           timestamppb.New(po.CreateTime),
			UpdatedAt:           nil, // or timestamppb.New(po.UpdateTime) if available
			CreatedBy:           "",  // fill if available
			UpdatedBy:           "",  // fill if available
		},
	}, nil
}

func (r *positionRepo) UpdatePosition(ctx context.Context, b *v1.Position) (*v1.PositionResponse, error) {
	po, err := r.data.db.Position.
		UpdateOneID(b.Id).
		SetPositionCode(b.PositionCode).
		SetName(b.Name).
		SetRoleName(b.RoleName).
		SetDepartmentCode(b.DepartmentCode).
		SetReportsToPositionID(b.ReportsToPositionId).
		SetUpdateTime(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Invalidate cache setelah update
	r.redisClient.Del(ctx, "positions:list")
	ids := strconv.FormatInt(po.ID, 10)
	r.redisClient.Del(ctx, "position:"+ids)                  // Invalidate cache spesifik berdasarkan ID
	r.redisClient.Del(ctx, "position_code:"+po.PositionCode) // Invalidate cache spesifik berdasarkan kode

	return &v1.PositionResponse{
		Position: &v1.Position{
			Id:                  int64(po.ID),
			PositionCode:        po.PositionCode,
			Name:                po.Name,
			RoleName:            po.RoleName,
			DepartmentCode:      po.DepartmentCode,
			ReportsToPositionId: *po.ReportsToPositionID,
			CreatedAt:           timestamppb.New(po.CreateTime),
			UpdatedAt:           nil, // or timestamppb.New(po.UpdateTime) if available
			CreatedBy:           "",  // fill if available
			UpdatedBy:           "",  // fill if available
		},
	}, nil
}

func (r positionRepo) DeletePosition(ctx context.Context, id int64) error {
	// Dapatkan detail departemen sebelum menghapus untuk menginvalidasi cache berdasarkan kode juga
	posToDelete, err := r.data.db.Position.Get(ctx, id)
	if err != nil {
		return err // Atau log error dan tetap coba hapus jika tidak kritis
	}

	err = r.data.db.Position.
		DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}
	// Invalidate cache setelah menghapus
	r.redisClient.Del(ctx, "positions:list")
	ids := strconv.FormatInt(id, 10)
	r.redisClient.Del(ctx, "position:"+ids)
	// Asumsi ada cache 'positionDetail' jika digunakan di tempat lain
	r.redisClient.Del(ctx, "positionDetail:"+ids)
	if posToDelete != nil {
		r.redisClient.Del(ctx, "position_code:"+posToDelete.PositionCode) // Invalidate cache by code
	}
	r.Log.Infof("Position dengan ID %d terhapus", id)
	return nil
}

func (r *positionRepo) Count(ctx context.Context) (int, error) {
	// Pertimbangkan caching untuk count jika sering diakses dan tidak perlu sangat up-to-date
	dt, _ := r.data.db.Position.Query().Count(ctx)
	return dt, nil
}

// func (r *positionRepo) ListPosition(ctx context.Context, pageNum, pageSize int64) ([]*v1.ListPositionsResponse, int, error) {
// 	cacheKey := fmt.Sprintf("positions:list:page:%d:size:%d", pageNum, pageSize)
// 	val, err := r.redisClient.Get(ctx, cacheKey).Result()
// 	if err == nil {
// 		var cachedResponse struct {
// 			Positions []*v1.ListPositionsResponse
// 			Total     int
// 		}
// 		if unmarshalErr := json.Unmarshal([]byte(val), &cachedResponse); unmarshalErr == nil {
// 			r.Log.Infof("Cache hit for ListPositions page %d, size %d", pageNum, pageSize)
// 			return &v1.ListPositionsResponse{
// 				Positions: cachedResponse.Positions,
// 			}, cachedResponse.Total, nil
// 		} else {
// 			r.Log.Warnf("Failed to unmarshal cached ListPositions: %v", unmarshalErr)
// 		}
// 	} else if err != redis.Nil {
// 		r.Log.Errorf("Redis GET error for ListPositions: %v", err)
// 	}

// 	// DB query
// 	query := r.data.db.Position.Query()

// 	total, err := query.Clone().Count(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	pos, err := query.
// 		Offset(int(pagination.GetPageOffset(pageNum, pageSize))).
// 		Limit(int(pageSize)).
// 		All(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	result := make([]*v1.Position, 0, len(pos))
// 	for _, po := range pos {
// 		var reportsTo string
// 		if po.ReportsToPositionID != nil {
// 			reportsTo = *po.ReportsToPositionID
// 		}
// 		result = append(result, &v1.Position{
// 			Id:                  int64(po.ID),
// 			PositionCode:        po.PositionCode,
// 			Name:                po.Name,
// 			RoleName:            po.RoleName,
// 			DepartmentCode:      po.DepartmentCode,
// 			ReportsToPositionId: reportsTo,
// 			CreatedAt:           timestamppb.New(po.CreateTime),
// 			UpdatedAt:           nil, // or timestamppb.New(po.UpdateTime)
// 			CreatedBy:           "",  // optional
// 			UpdatedBy:           "",  // optional
// 		})
// 	}

// 	// Cache result
// 	responseToCache := struct {
// 		Positions []*v1.Position
// 		Total     int
// 	}{
// 		Positions: result,
// 		Total:     total,
// 	}
// 	if responseBytes, marshalErr := json.Marshal(responseToCache); marshalErr == nil {
// 		r.redisClient.Set(ctx, cacheKey, responseBytes, 3*time.Minute)
// 	} else {
// 		r.Log.Errorf("Failed to marshal ListPositions response for caching: %v", marshalErr)
// 	}

//		return &v1.ListPositionsResponse{
//			Positions: result,
//		}, total, nil
//	}
func (r *positionRepo) ListPosition(ctx context.Context, pageNum, pageSize int64) ([]*biz.PositionData, int, error) {
	cacheKey := fmt.Sprintf("positions:list:page:%d:size:%d", pageNum, pageSize)
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedResponse struct {
			Positions []*biz.PositionData
			Total     int
		}
		if unmarshalErr := json.Unmarshal([]byte(val), &cachedResponse); unmarshalErr == nil {
			r.Log.Infof("Cache hit for ListPositions page %d, size %d", pageNum, pageSize)
			return cachedResponse.Positions, cachedResponse.Total, nil
		} else {
			r.Log.Warnf("Failed to unmarshal cached ListPositions: %v", unmarshalErr)
		}
	} else if err != redis.Nil {
		r.Log.Errorf("Redis GET error for ListPositions: %v", err)
	}

	query := r.data.db.Position.Query()
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

	result := make([]*biz.PositionData, 0, len(pos))
	for _, po := range pos {
		var reportsTo string
		if po.ReportsToPositionID != nil {
			reportsTo = *po.ReportsToPositionID
		}
		result = append(result, &biz.PositionData{
			ID:                  int64(po.ID),
			PositionCode:        po.PositionCode,
			Name:                po.Name,
			RoleName:            po.RoleName,
			DepartmentCode:      po.DepartmentCode,
			ReportsToPositionId: reportsTo,
			CreatedAt:           po.CreateTime,
			UpdatedAt:           po.UpdateTime,
			CreatedBy:           "", // Optional, sesuaikan kalau ada
			UpdatedBy:           "", // Optional
		})
	}

	// Cache result
	responseToCache := struct {
		Positions []*biz.PositionData
		Total     int
	}{
		Positions: result,
		Total:     total,
	}
	if responseBytes, marshalErr := json.Marshal(responseToCache); marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, responseBytes, 3*time.Minute)
	} else {
		r.Log.Errorf("Failed to marshal ListPositions response for caching: %v", marshalErr)
	}

	return result, total, nil
}
