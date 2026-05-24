package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	// biometricV1 "mall-go/api/biometrics/service/v1"
	// departmentV1 "mall-go/api/department/service/v1"
	"mall-go/module/shift/service/internal/biz"
	"mall-go/pkg/utils/pagination"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

var _ biz.ShiftRepo = (*shiftRepo)(nil)

type shiftRepo struct {
	data *Data
	Log  *log.Helper
	// bioClient   biometricV1.BiometricClient
	// deptClient  departmentV1.DepartmentClient
	redisClient *redis.Client // Add Redis client
}

func NewShiftRepo(data *Data, logger log.Logger, rdb *redis.Client) biz.ShiftRepo {
	return &shiftRepo{
		data:        data,
		redisClient: rdb,                   // Initialize Redis client
		Log:         log.NewHelper(logger), // ✅ INI WAJIB ADA
	}
}

//	func NewShiftRepo(data *Data, bio biometricV1.BiometricClient, dept departmentV1.DepartmentClient, logger log.Logger, rdb *redis.Client) biz.ShiftRepo {
//		return &shiftRepo{
//			data:        data,
//			bioClient:   bio,
//			deptClient:  dept,
//			redisClient: rdb,                   // Initialize Redis client
//			Log:         log.NewHelper(logger), // ✅ INI WAJIB ADA
//		}
//	}
func (r *shiftRepo) CreateShift(ctx context.Context, b *biz.Shift) (*biz.Shift, error) {
	// Parse start_time and end_time
	startTime, err := time.Parse("15:04", b.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time format: %w", err)
	}

	endTime, err := time.Parse("15:04", b.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end_time format: %w", err)
	}

	po, err := r.data.db.Shift.
		Create().
		SetName(b.Name).
		SetStartTime(startTime).
		SetEndTime(endTime).
		SetBreakDurationMinutes(int(b.BreakDurationMinutes)).
		SetCreatedBy("System").
		SetCreatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	r.redisClient.Del(ctx, "shift:list")
	r.redisClient.Del(ctx, fmt.Sprintf("shift:%d", po.ID))

	return &biz.Shift{
		Id:                   po.ID,
		Name:                 po.Name,
		StartTime:            po.StartTime.Format("15:04"),
		EndTime:              po.EndTime.Format("15:04"),
		BreakDurationMinutes: int32(po.BreakDurationMinutes),
		CreatedBy:            po.CreatedBy,
		CreatedAt:            po.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            po.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (r *shiftRepo) UpdateShift(ctx context.Context, b *biz.Shift) (*biz.Shift, error) {
	// Parse start_time and end_time
	startTime, err := time.Parse("15:04", b.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time format: %w", err)
	}

	endTime, err := time.Parse("15:04", b.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end_time format: %w", err)
	}

	// Update data di database via Ent
	po, err := r.data.db.Shift.
		UpdateOneID(b.Id).
		SetName(b.Name).
		SetStartTime(startTime).
		SetEndTime(endTime).
		SetBreakDurationMinutes(int(b.BreakDurationMinutes)).
		SetUpdatedAt(time.Now()).
		SetUpdatedBy("System"). // jika kamu punya field ini
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Invalidate Redis cache
	_ = r.redisClient.Del(ctx, "shift:list")
	_ = r.redisClient.Del(ctx, "shift:"+strconv.FormatInt(po.ID, 10))

	return &biz.Shift{
		Id:                   po.ID,
		Name:                 po.Name,
		StartTime:            po.StartTime.Format("15:04"),
		EndTime:              po.EndTime.Format("15:04"),
		BreakDurationMinutes: int32(po.BreakDurationMinutes),
		CreatedBy:            po.CreatedBy,
		CreatedAt:            po.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            po.UpdatedAt.Format(time.RFC3339),
	}, nil
}
func (r *shiftRepo) GetShiftID(ctx context.Context, id int64) (*biz.Shift, error) {
	cacheKey := fmt.Sprintf("shift:%d", id)

	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedShift biz.Shift
		unmarshalErr := json.Unmarshal([]byte(val), &cachedShift)
		if unmarshalErr := json.Unmarshal([]byte(val), &cachedShift); unmarshalErr == nil {
			r.Log.Infof("Cache hit for shift ID: %d", id)
			return &cachedShift, nil
		}
		r.Log.Warnf("Failed to unmarshal cached shift for ID %d: %v", id, unmarshalErr)
	} else if err != redis.Nil {
		r.Log.Errorf("Redis GET error for shift ID %d: %v", id, err)
	}

	// Fallback to DB
	po, err := r.data.db.Shift.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	shift := &biz.Shift{
		Id:                   po.ID,
		Name:                 po.Name,
		StartTime:            po.StartTime.Format("15:04"),
		EndTime:              po.EndTime.Format("15:04"),
		BreakDurationMinutes: int32(po.BreakDurationMinutes),
		CreatedBy:            po.CreatedBy,
		UpdatedBy:            po.UpdatedBy,
		CreatedAt:            po.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            po.UpdatedAt.Format(time.RFC3339),
	}

	// Cache the result
	if shiftBytes, marshalErr := json.Marshal(shift); marshalErr == nil {
		if err := r.redisClient.Set(ctx, cacheKey, shiftBytes, 5*time.Minute).Err(); err != nil {
			r.Log.Errorf("Redis SET error for shift ID %d: %v", id, err)
		}
	} else {
		r.Log.Errorf("Failed to marshal shift for caching: %v", marshalErr)
	}

	return shift, nil
}

// func (r *shiftRepo) GetShiftID(ctx context.Context, id int64) (*biz.Shift, error) {
// 	cacheKey := "shift:" + string(id)
// 	val, err := r.redisClient.Get(ctx, cacheKey).Result()
// 	if err == nil {
// 		var cachedShift biz.Shift
// 		unmarshalErr := json.Unmarshal([]byte(val), &cachedShift)
// 		if unmarshalErr := json.Unmarshal([]byte(val), &cachedShift); unmarshalErr == nil {
// 			r.Log.Infof("Cache hit for shift ID: %d", id)
// 			return &cachedShift, nil
// 		}
// 		r.Log.Warnf("Failed to unmarshal cached shift for ID %d: %v", id, unmarshalErr)
// 	} else if err != redis.Nil {
// 		r.Log.Errorf("Redis GET error for shift ID %d: %v", id, err)
// 	}

// 	po, err := r.data.db.Shift.Get(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	shift := &biz.Shift{
// 		Id:                   po.ID,
// 		Name:                 po.Name,
// 		StartTime:            po.StartTime.Format("15:04"),
// 		EndTime:              po.EndTime.Format("15:04"),
// 		BreakDurationMinutes: int32(po.BreakDurationMinutes),
// 		CreatedBy:            po.CreatedBy,
// 		UpdatedBy:            po.UpdatedBy,
// 		CreatedAt:            po.CreatedAt.Format(time.RFC3339),
// 		UpdatedAt:            po.UpdatedAt.Format(time.RFC3339),
// 	}

// 	// Cache the result
// 	shiftBytes, marshalErr := json.Marshal(shift)
// 	if marshalErr == nil {
// 		r.redisClient.Set(ctx, cacheKey, shiftBytes, 5*time.Minute) // Cache for 5 minutes
// 	} else {
// 		r.Log.Errorf("Failed to marshal shift for caching: %v", marshalErr)
// 	}

// 	return shift, nil
// }

// func (r *shiftRepo) ListShift(ctx context.Context, pageNum, pageSize int64) ([]*biz.ShiftData, int, error) {
// 	query := r.data.db.Shift.Query()

// 	// Hitung total sebelum pagination
// 	total, err := query.Clone().Count(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	pos, err := r.data.db.Shift.Query().
// 		Offset(int(pagination.GetPageOffset(pageNum, pageSize))).
// 		Limit(int(pageSize)).
// 		All(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	if r.Log != nil {
// 		r.Log.Warnf("ListShift result: %+v", pos)
// 	}
// 	rv := make([]*biz.ShiftData, 0, len(pos))
// 	for _, po := range pos {
// 		empData := &biz.ShiftData{
// 			Id:        po.ID,
// 			NoSap:     *po.Nosap,
// 			Nip:       *po.Nip,
// 			KaryaCode: *po.Karyacode,
// 			KaryaName: po.Karyaname,
// 			DispName:  po.DispName,
// 			PassMesin: *po.PassMesin,
// 			RFIDCard:  *po.RfidCard,
// 			Status:    po.Status,
// 		}

// 		// Get finger data
// 		if po.KodeFinger != "" {
// 			fingerResp, err := r.bioClient.GetFingerByKode(ctx, &biometricV1.GetFingerByKodeRequest{
// 				Fingercode: po.KodeFinger,
// 			})
// 			if err == nil {
// 				empData.Finger = []biz.FingerData{
// 					{
// 						Fingercode: fingerResp.Fingercode,
// 						Finger0:    fingerResp.Finger0,
// 						Finger1:    fingerResp.Finger1,
// 						Finger2:    fingerResp.Finger2,
// 						Finger3:    fingerResp.Finger3,
// 						Finger4:    fingerResp.Finger4,
// 						Finger5:    fingerResp.Finger5,
// 						Finger6:    fingerResp.Finger6,
// 						Finger7:    fingerResp.Finger7,
// 						Finger8:    fingerResp.Finger8,
// 						Finger9:    fingerResp.Finger9,
// 					},
// 				}
// 			} else {
// 				r.Log.Warnf("finger not found for kode: %s, err: %v", po.KodeFinger, err)
// 			}
// 		}

// 		// Department
// 		r.Log.Infow("", po.DepartCode)
// 		if po.DepartCode != "" {
// 			deptResp, err := r.deptClient.GetDepartmentCode(ctx, &departmentV1.GetDepartmentCodeRequest{
// 				DepartCode: po.DepartCode,
// 			})
// 			if err != nil {
// 				r.Log.Warnf("Failed to get department for depart_code=%s: %v", po.DepartCode, err)
// 			} else {
// 				empData.Department = []biz.DepartData{
// 					{
// 						DepartCode: deptResp.DepartCode,
// 						DepartName: deptResp.DepartName,
// 					},
// 				}
// 				// r.Log.Infof("DEPARTMENT: %s, %s,%s", empData.Department, deptResp.DepartCode, deptResp.DepartName)
// 			}
// 		}

// 		rv = append(rv, empData)
// 	}

//		return rv, total, nil
//	}

func (r *shiftRepo) ListShift(ctx context.Context, pageNum, pageSize int64) ([]*biz.Shift, int, error) {
	cacheKey := fmt.Sprintf("shift:list:page:%d:size:%d", pageNum, pageSize)

	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cached struct {
			Shift []*biz.Shift `json:"shift"`
			Total int          `json:"total"`
		}
		if unmarshalErr := json.Unmarshal([]byte(val), &cached); unmarshalErr == nil {
			r.Log.Infof("Cache hit for ListShift page %d, size %d", pageNum, pageSize)
			return cached.Shift, cached.Total, nil
		}
		r.Log.Warnf("Failed to unmarshal cached ListShift: %v", err)
	} else if err != redis.Nil {
		r.Log.Errorf("Redis GET error for ListShift: %v", err)
	}

	// Query total
	total, err := r.data.db.Shift.Query().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Query paginated data
	pos, err := r.data.db.Shift.Query().
		Offset(int(pagination.GetPageOffset(pageNum, pageSize))).
		Limit(int(pageSize)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	rv := make([]*biz.Shift, 0, len(pos))
	for _, po := range pos {
		rv = append(rv, &biz.Shift{
			Id:                   po.ID,
			Name:                 po.Name,
			StartTime:            po.StartTime.Format("15:04"),
			EndTime:              po.EndTime.Format("15:04"),
			BreakDurationMinutes: int32(po.BreakDurationMinutes),
			CreatedBy:            po.CreatedBy,
			UpdatedBy:            po.UpdatedBy,
			CreatedAt:            po.CreatedAt.Format(time.RFC3339),
			UpdatedAt:            po.UpdatedAt.Format(time.RFC3339),
		})
	}

	// Cache result
	responseBytes, marshalErr := json.Marshal(struct {
		Shift []*biz.Shift `json:"shift"`
		Total int          `json:"total"`
	}{
		Shift: rv,
		Total: total,
	})
	if marshalErr == nil {
		err = r.redisClient.Set(ctx, cacheKey, responseBytes, 5*time.Minute).Err()
		if err != nil {
			r.Log.Errorf("Redis SET error for ListShift cacheKey %s: %v", cacheKey, err)
		}
	} else {
		r.Log.Errorf("Marshal error ListShift cache: %v", marshalErr)
	}

	return rv, total, nil
}

// func (r *shiftRepo) ListShift(ctx context.Context, pageNum, pageSize int64) ([]*biz.Shift, int, error) {
// 	cacheKey := "shift:list:page:" + string(pageNum) + ":size:" + string(pageSize)
// 	val, err := r.redisClient.Get(ctx, cacheKey).Result()
// 	if err == nil {
// 		var cachedResponse struct {
// 			Shift []*biz.Shift
// 			Total int
// 		}
// 		var cachedShift biz.Shift
// 		unmarshalErr := json.Unmarshal([]byte(val), &cachedShift)
// 		if unmarshalErr := json.Unmarshal([]byte(val), &cachedResponse); unmarshalErr == nil {
// 			r.Log.Infof("Cache hit for ListShift page %d, size %d", pageNum, pageSize)
// 			return cachedResponse.Shift, cachedResponse.Total, nil
// 		}

// 		r.Log.Warnf("Failed to unmarshal cached ListShift: %v", unmarshalErr)
// 	} else if err != redis.Nil {
// 		r.Log.Errorf("Redis GET error for ListShift: %v", err)
// 	}

// 	query := r.data.db.Shift.Query()

// 	total, err := query.Clone().Count(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	pos, err := r.data.db.Shift.Query().
// 		Offset(int(pagination.GetPageOffset(pageNum, pageSize))).
// 		Limit(int(pageSize)).
// 		All(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	if r.Log != nil {
// 		r.Log.Warnf("ListShift result from DB: %+v", pos)
// 	}
// 	rv := make([]*biz.Shift, 0, len(pos))
// 	for _, po := range pos {
// 		shiftData := &biz.Shift{
// 			Id:                   po.ID,
// 			Name:                 po.Name,
// 			StartTime:            po.StartTime.Format("15:04"),
// 			EndTime:              po.EndTime.Format("15:04"),
// 			BreakDurationMinutes: int32(po.BreakDurationMinutes),
// 			CreatedBy:            po.CreatedBy,
// 			UpdatedBy:            po.UpdatedBy,
// 			CreatedAt:            po.CreatedAt.Format(time.RFC3339),
// 			UpdatedAt:            po.UpdatedAt.Format(time.RFC3339),
// 		}
// 		rv = append(rv, shiftData)
// 	}

// 	// Cache the entire list response
// 	responseToCache := struct {
// 		Shift []*biz.Shift
// 		Total int
// 	}{
// 		Shift: rv,
// 		Total: total,
// 	}
// 	responseBytes, marshalErr := json.Marshal(responseToCache)
// 	if marshalErr == nil {
// 		r.redisClient.Set(ctx, cacheKey, responseBytes, 5*time.Minute) // Cache for 5 minutes
// 	} else {
// 		r.Log.Errorf("Failed to marshal ListShift response for caching: %v", marshalErr)
// 	}

// 	return rv, total, nil
// }

func (r *shiftRepo) Count(ctx context.Context) (int, error) {
	dt, _ := r.data.db.Shift.Query().Count(ctx)
	// fmt.Println("INI PL DT: ")
	return dt, nil
}
func (r *shiftRepo) ListShiftNext(ctx context.Context, start, end int32) ([]*biz.Shift, error) {
	// You can apply similar caching logic here if this method is heavily used
	// and benefits from caching. The cache key would need to reflect start and end.
	pos, err := r.data.db.Shift.Query().
		Offset(int(start)).
		Limit(int(end - start)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	rv := make([]*biz.Shift, 0, len(pos))
	for _, po := range pos {
		rv = append(rv, &biz.Shift{
			Id: po.ID,
			// NoSap:      *po.Nosap,
			// Nip:        *po.Nip,
			// KaryaCode:  *po.Karyacode,
			// KaryaName:  po.Karyaname,
			// DispName:   po.DispName,
			// PassMesin:  *po.PassMesin,
			// RFIDCard:   *po.RfidCard,
			// Finger:     po.KodeFinger,
			// Department: po.DepartCode,
			// Status:     po.Status,
		})
	}
	return rv, nil
}
func (r *shiftRepo) DeleteShift(ctx context.Context, id int64) error {
	err := r.data.db.Shift.
		DeleteOneID(id).
		Exec(ctx)
	if err != nil {
		return err
	}

	// Convert ID to string properly
	idStr := strconv.FormatInt(id, 10)

	// Invalidate related cache
	r.redisClient.Del(ctx, "shift:list")
	r.redisClient.Del(ctx, "shift:"+idStr)
	r.redisClient.Del(ctx, "shiftDetail:"+idStr)

	r.Log.Infof("Shift with ID %d deleted and related cache invalidated", id)
	return nil
}

// func (r *shiftRepo) DeleteShift(ctx context.Context, id int64) error {
// 	err := r.data.db.Shift.
// 		DeleteOneID(id).
// 		Exec(ctx)

// 	if err != nil {
// 		return err
// 	}

// 	// Invalidate caches related to the deleted shift
// 	r.redisClient.Del(ctx, "shift:list")
// 	r.redisClient.Del(ctx, "shift:"+string(id))
// 	r.redisClient.Del(ctx, "shiftDetail:"+string(id)) // Invalidate detail cache

// 	r.Log.Infof("shift with ID %d deleted", id)
// 	return nil
// }

// func (r *shiftRepo) GetShiftDetail(ctx context.Context, id int64) (*biz.ShiftData, error) {
// 	cacheKey := "shiftDetail:" + string(id)
// 	val, err := r.redisClient.Get(ctx, cacheKey).Result()
// 	if err == nil {
// 		var cachedShiftData biz.ShiftData
// 		unmarshalErr := json.Unmarshal([]byte(val), &cachedShiftData)
// 		if unmarshalErr := json.Unmarshal([]byte(val), &cachedShiftData); unmarshalErr == nil {
// 			r.Log.Infof("Cache hit for shift detail ID: %d", id)
// 			return &cachedShiftData, nil
// 		}
// 		r.Log.Warnf("Failed to unmarshal cached shift detail for ID %d: %v", id, unmarshalErr)
// 	} else if err != redis.Nil {
// 		r.Log.Errorf("Redis GET error for shift detail ID %d: %v", id, err)
// 	}

// 	emp, err := r.data.db.Shift.Get(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	shiftData := &biz.ShiftData{
// 		Id: emp.ID,
// 		// NoSap:     *emp.Nosap,
// 		// Nip:       *emp.Nip,
// 		// KaryaCode: *emp.Karyacode,
// 		// KaryaName: emp.Karyaname,
// 		// DispName:  emp.DispName,
// 		// PassMesin: *emp.PassMesin,
// 		// RFIDCard:  *emp.RfidCard,
// 		// Status:    emp.Status,
// 		// CreatedAt: emp.created_at.String(),
// 		// UpdatedAt: emp.UpdatedAt.String(),
// 	}

// 	// Get finger data with caching
// 	fingerCacheKey := "biometric:" + emp.KodeFinger
// 	fingerVal, fingerErr := r.redisClient.Get(ctx, fingerCacheKey).Result()
// 	if fingerErr == nil {
// 		var cachedFinger biometricV1.GetFingerByKodeResponse
// 		if unmarshalErr := json.Unmarshal([]byte(fingerVal), &cachedFinger); unmarshalErr == nil {
// 			shiftData.Finger = []biz.FingerData{
// 				{
// 					Fingercode: cachedFinger.Fingercode,
// 					Finger0:    cachedFinger.Finger0,
// 					Finger1:    cachedFinger.Finger1,
// 					Finger2:    cachedFinger.Finger2,
// 					Finger3:    cachedFinger.Finger3,
// 					Finger4:    cachedFinger.Finger4,
// 					Finger5:    cachedFinger.Finger5,
// 					Finger6:    cachedFinger.Finger6,
// 					Finger7:    cachedFinger.Finger7,
// 					Finger8:    cachedFinger.Finger8,
// 					Finger9:    cachedFinger.Finger9,
// 				},
// 			}
// 			r.Log.Infof("Cache hit for biometric code: %s (detail)", emp.KodeFinger)
// 		} else {
// 			r.Log.Warnf("Failed to unmarshal cached biometric for code %s (detail): %v", emp.KodeFinger, unmarshalErr)
// 		}
// 	} else if fingerErr != redis.Nil {
// 		r.Log.Errorf("Redis GET error for biometric code %s (detail): %v", emp.KodeFinger, fingerErr)
// 	}

// 	if shiftData.Finger == nil && emp.KodeFinger != "" { // If not in cache, or unmarshal failed, fetch from service
// 		fingerResp, err := r.bioClient.GetFingerByKode(ctx, &biometricV1.GetFingerByKodeRequest{
// 			Fingercode: emp.KodeFinger,
// 		})
// 		if err != nil {
// 			r.Log.Errorf("Failed to fetch biometric for detail: %v", err)
// 		} else {
// 			shiftData.Finger = []biz.FingerData{
// 				{
// 					Fingercode: fingerResp.Fingercode,
// 					Finger0:    fingerResp.Finger0,
// 					Finger1:    fingerResp.Finger1,
// 					Finger2:    fingerResp.Finger2,
// 					Finger3:    fingerResp.Finger3,
// 					Finger4:    fingerResp.Finger4,
// 					Finger5:    fingerResp.Finger5,
// 					Finger6:    fingerResp.Finger6,
// 					Finger7:    fingerResp.Finger7,
// 					Finger8:    fingerResp.Finger8,
// 					Finger9:    fingerResp.Finger9,
// 				},
// 			}
// 			// Cache the finger data
// 			fingerBytes, marshalErr := json.Marshal(fingerResp)
// 			if marshalErr == nil {
// 				r.redisClient.Set(ctx, fingerCacheKey, fingerBytes, 5*time.Minute)
// 			} else {
// 				r.Log.Errorf("Failed to marshal biometric for caching (detail): %v", marshalErr)
// 			}
// 		}
// 	}

// 	// Get department data with caching
// 	deptCacheKey := "department:" + emp.DepartCode
// 	deptVal, deptErr := r.redisClient.Get(ctx, deptCacheKey).Result()
// 	if deptErr == nil {
// 		var cachedDept departmentV1.GetDepartmentCodeResponse
// 		if unmarshalErr := json.Unmarshal([]byte(deptVal), &cachedDept); unmarshalErr == nil {
// 			shiftData.Department = []biz.DepartData{
// 				{
// 					DepartCode: cachedDept.DepartCode,
// 					DepartName: cachedDept.DepartName,
// 				},
// 			}
// 			r.Log.Infof("Cache hit for department code: %s (detail)", emp.DepartCode)
// 		} else {
// 			r.Log.Warnf("Failed to unmarshal cached department for code %s (detail): %v", emp.DepartCode, unmarshalErr)
// 		}
// 	} else if deptErr != redis.Nil {
// 		r.Log.Errorf("Redis GET error for department code %s (detail): %v", emp.DepartCode, deptErr)
// 	}

// 	if shiftData.Department == nil && emp.DepartCode != "" { // If not in cache, or unmarshal failed, fetch from service
// 		deptResp, err := r.deptClient.GetDepartmentCode(ctx, &departmentV1.GetDepartmentCodeRequest{
// 			DepartCode: emp.DepartCode,
// 		})
// 		if err != nil {
// 			r.Log.Errorf("Failed to fetch department for detail: %v", err)
// 		} else {
// 			shiftData.Department = []biz.DepartData{
// 				{
// 					DepartCode: deptResp.DepartCode,
// 					DepartName: deptResp.DepartName,
// 				},
// 			}
// 			// Cache the department data
// 			deptBytes, marshalErr := json.Marshal(deptResp)
// 			if marshalErr == nil {
// 				r.redisClient.Set(ctx, deptCacheKey, deptBytes, 5*time.Minute)
// 			} else {
// 				r.Log.Errorf("Failed to marshal department for caching (detail): %v", marshalErr)
// 			}
// 		}
// 	}

// 	// Cache the shift detail
// 	shiftDataBytes, marshalErr := json.Marshal(shiftData)
// 	if marshalErr == nil {
// 		r.redisClient.Set(ctx, cacheKey, shiftDataBytes, 5*time.Minute) // Cache for 5 minutes
// 	} else {
// 		r.Log.Errorf("Failed to marshal shift detail for caching: %v", marshalErr)
// 	}

// 	return shiftData, nil
// }
