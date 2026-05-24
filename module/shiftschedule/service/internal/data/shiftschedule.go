package data

import (
	"context"
	"encoding/json"
	departmentV1 "mall-go/api/department/service/v1"
	employersV1 "mall-go/api/employers/service/v1"
	"mall-go/module/shiftschedule/service/internal/biz"
	"mall-go/pkg/utils/pagination"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

// var _ biz.BeerRepo = (*beerRepo)(nil)
var _ biz.ShiftScheduleRepo = (*shiftscheduleRepo)(nil)

type shiftscheduleRepo struct {
	data            *Data
	Log             *log.Helper
	employersClient employersV1.EmployersClient
	deptClient      departmentV1.DepartmentClient
	redisClient     *redis.Client // Add Redis client
}

func NewShiftScheduleRepo(data *Data, employer employersV1.EmployersClient, dept departmentV1.DepartmentClient, logger log.Logger, rdb *redis.Client) biz.ShiftScheduleRepo {
	return &shiftscheduleRepo{
		data:            data,
		employersClient: employer,
		deptClient:      dept,
		redisClient:     rdb,                   // Initialize Redis client
		Log:             log.NewHelper(logger), // ✅ INI WAJIB ADA
	}
}

func (r *shiftscheduleRepo) CreateShiftSchedule(ctx context.Context, b *biz.ShiftSchedule) (*biz.ShiftSchedule, error) {
	po, err := r.data.db.ShiftSchedule.
		Create().
		SetScheduleCode(b.ScheduleCode).
		SetKaryaCode(b.KaryaCode).
		SetTanggal(b.Tanggal).
		SetDepartCode(b.DepartCode).
		SetCreatedBy(b.CreatedBy).
		SetShiftID(b.ShiftID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	r.redisClient.Del(ctx, "Shiftscheduler:list")
	r.redisClient.Del(ctx, "Shiftscheduler:"+string(po.ID)) //  Invalidate specific employer cache if exists

	return &biz.ShiftSchedule{
		Id:           po.ID,
		ScheduleCode: po.ScheduleCode,
		KaryaCode:    po.KaryaCode,
		Tanggal:      po.Tanggal,
		DepartCode:   po.DepartCode,
		CreatedBy:    po.CreatedBy,
		ShiftID:      po.ShiftID,
	}, nil

}

func (r *shiftscheduleRepo) UpdateShiftSchedule(ctx context.Context, b *biz.ShiftSchedule) (*biz.ShiftSchedule, error) {
	po, err := r.data.db.ShiftSchedule.
		UpdateOneID(b.Id).
		SetScheduleCode(b.ScheduleCode).
		SetKaryaCode(b.KaryaCode).
		SetTanggal(b.Tanggal).
		SetDepartCode(b.DepartCode).
		SetCreatedBy(b.CreatedBy).
		SetShiftID(b.ShiftID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	// Invalidate relevant caches after update
	r.redisClient.Del(ctx, "Shiftscheduler:list")
	r.redisClient.Del(ctx, "Shiftscheduler:"+string(po.ID)) // Invalidate specific employer cache

	return &biz.ShiftSchedule{
		Id:           po.ID,
		ScheduleCode: po.ScheduleCode,
		KaryaCode:    po.KaryaCode,
		Tanggal:      po.Tanggal,
		DepartCode:   po.DepartCode,
		CreatedBy:    po.CreatedBy,
	}, nil
}
func (r *shiftscheduleRepo) GetShiftScheduleID(ctx context.Context, id int64) (*biz.ShiftSchedule, error) {
	cacheKey := "Shiftscheduler:" + string(id)
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedShiftSchedule biz.ShiftSchedule
		unmarshalErr := json.Unmarshal([]byte(val), &cachedShiftSchedule)
		if unmarshalErr := json.Unmarshal([]byte(val), &cachedShiftSchedule); unmarshalErr == nil {
			r.Log.Infof("Cache hit for Shiftscheduler ID: %d", id)
			return &cachedShiftSchedule, nil
		}
		r.Log.Warnf("Failed to unmarshal cached cachedShiftSchedule for ID %d: %v", id, unmarshalErr)
	} else if err != redis.Nil {
		r.Log.Errorf("Redis GET error for cachedShiftSchedule ID %d: %v", id, err)
	}

	po, err := r.data.db.ShiftSchedule.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	shiftschedule := &biz.ShiftSchedule{
		Id:           po.ID,
		ScheduleCode: po.ScheduleCode,
		KaryaCode:    po.KaryaCode,
		Tanggal:      po.Tanggal,
		DepartCode:   po.DepartCode,
		CreatedBy:    po.CreatedBy,
	}

	// Cache the result
	shiftscheduleBytes, marshalErr := json.Marshal(shiftschedule)
	if marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, shiftscheduleBytes, 5*time.Minute) // Cache for 5 minutes
	} else {
		r.Log.Errorf("Failed to marshal shiftschedule for caching: %v", marshalErr)
	}

	return shiftschedule, nil
}

func (r *shiftscheduleRepo) ListShiftSchedule(ctx context.Context, pageNum, pageSize int64) ([]*biz.ShiftScheduleData, int, error) {
	cacheKey := "Shiftscheduler:list:page:" + string(pageNum) + ":size:" + string(pageSize)
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedResponse struct {
			ShiftSchedule []*biz.ShiftScheduleData
			Total         int
		}
		var cachedShiftSchedule biz.ShiftScheduleData
		unmarshalErr := json.Unmarshal([]byte(val), &cachedShiftSchedule)
		if unmarshalErr := json.Unmarshal([]byte(val), &cachedShiftSchedule); unmarshalErr == nil {
			r.Log.Infof("Cache hit for ListShiftSchedule page %d, size %d", pageNum, pageSize)
			return cachedResponse.ShiftSchedule, cachedResponse.Total, nil
		}

		r.Log.Warnf("Failed to unmarshal cached ListShiftSchedule: %v", unmarshalErr)
	} else if err != redis.Nil {
		r.Log.Errorf("Redis GET error for ListShiftSchedule: %v", err)
	}

	query := r.data.db.ShiftSchedule.Query()

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	pos, err := r.data.db.ShiftSchedule.Query().
		Offset(int(pagination.GetPageOffset(pageNum, pageSize))).
		Limit(int(pageSize)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	if r.Log != nil {
		r.Log.Warnf("ListShiftSchedule result from DB: %+v", pos)
	}
	rv := make([]*biz.ShiftScheduleData, 0, len(pos))
	for _, po := range pos {
		empData := &biz.ShiftScheduleData{
			Id:           po.ID,
			ScheduleCode: po.ScheduleCode,
			// Empl:    po.KaryaCode,
			Tanggal: po.Tanggal,
			// Department: po.DepartCode,
			CreatedBy: po.CreatedBy,
		}

		// Get finger data with caching
		employersCacheKey := "employers:" + po.KaryaCode
		employerVal, employerErr := r.redisClient.Get(ctx, employersCacheKey).Result()
		if employerErr == nil {
			var cachedEmployer employersV1.GetEmployersKodeResponse
			if unmarshalErr := json.Unmarshal([]byte(employerVal), &cachedEmployer); unmarshalErr == nil {
				empData.Employers = []biz.EmployerData{
					{
						NoSap:     cachedEmployer.Nosap,
						Nip:       cachedEmployer.Nip,
						KaryaCode: cachedEmployer.KaryaCode,
						KaryaName: cachedEmployer.KaryaName,
					},
				}
				r.Log.Infof("Cache hit for employers code: %s", po.KaryaCode)
			} else {
				r.Log.Warnf("Failed to unmarshal cached employer for code %s: %v", po.KaryaCode, unmarshalErr)
			}
		} else if employerErr != redis.Nil {
			r.Log.Errorf("Redis GET error for employers code %s: %v", po.KaryaCode, employerErr)
		}

		if empData.Employers == nil && po.KaryaCode != "" { // If not in cache, or unmarshal failed, fetch from service
			employersResp, err := r.employersClient.GetEmployersKode(ctx, &employersV1.GetEmployersKodeRequest{
				KaryaCode: po.KaryaCode,
			})
			if err == nil {
				empData.Employers = []biz.EmployerData{
					{
						NoSap:     employersResp.Nosap,
						Nip:       employersResp.Nip,
						KaryaCode: employersResp.KaryaCode,
						KaryaName: employersResp.KaryaName,
					},
				}
				// Cache the finger data
				employersBytes, marshalErr := json.Marshal(employersResp)
				if marshalErr == nil {
					r.redisClient.Set(ctx, employersCacheKey, employersBytes, 5*time.Minute)
				} else {
					r.Log.Errorf("Failed to marshal employers for caching: %v", marshalErr)
				}
			} else {
				r.Log.Warnf("employers not found for kode: %s, err: %v", po.KaryaCode, err)
			}
		}

		// Department with caching
		deptCacheKey := "department:" + po.DepartCode
		deptVal, deptErr := r.redisClient.Get(ctx, deptCacheKey).Result()
		if deptErr == nil {
			var cachedDept departmentV1.GetDepartmentCodeResponse
			if unmarshalErr := json.Unmarshal([]byte(deptVal), &cachedDept); unmarshalErr == nil {
				empData.Department = []biz.DepartData{
					{
						DepartCode: cachedDept.DepartCode,
						DepartName: cachedDept.DepartName,
					},
				}
				r.Log.Infof("Cache hit for department code: %s", po.DepartCode)
			} else {
				r.Log.Warnf("Failed to unmarshal cached department for code %s: %v", po.DepartCode, unmarshalErr)
			}
		} else if deptErr != redis.Nil {
			r.Log.Errorf("Redis GET error for department code %s: %v", po.DepartCode, deptErr)
		}

		if empData.Department == nil && po.DepartCode != "" { // If not in cache, or unmarshal failed, fetch from service
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
				// Cache the department data
				deptBytes, marshalErr := json.Marshal(deptResp)
				if marshalErr == nil {
					r.redisClient.Set(ctx, deptCacheKey, deptBytes, 5*time.Minute)
				} else {
					r.Log.Errorf("Failed to marshal department for caching: %v", marshalErr)
				}
			}
		}
		rv = append(rv, empData)
	}

	// Cache the entire list response
	responseToCache := struct {
		ShiftSchedule []*biz.ShiftScheduleData
		Total         int
	}{
		ShiftSchedule: rv,
		Total:         total,
	}
	responseBytes, marshalErr := json.Marshal(responseToCache)
	if marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, responseBytes, 5*time.Minute) // Cache for 5 minutes
	} else {
		r.Log.Errorf("Failed to marshal ListShiftSchedule response for caching: %v", marshalErr)
	}

	return rv, total, nil
}
func (r *shiftscheduleRepo) Count(ctx context.Context) (int, error) {
	dt, _ := r.data.db.ShiftSchedule.Query().Count(ctx)
	// fmt.Println("INI PL DT: ")
	return dt, nil
}

// func (r *shiftscheduleRepo) ListShiftScheduleNext(ctx context.Context, start, end int32) ([]*biz.ShiftSchedule, error) {
// 	// You can apply similar caching logic here if this method is heavily used
// 	// and benefits from caching. The cache key would need to reflect start and end.
// 	pos, err := r.data.db.ShiftSchedule.Query().
// 		Offset(int(start)).
// 		Limit(int(end - start)).
// 		All(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rv := make([]*biz.ShiftScheduleData, 0, len(pos))
// 	for _, po := range pos {
// 		rv = append(rv, &biz.ShiftScheduleData{
// 			Id:         po.ID,
// 			NoSap:      *po.Nosap,
// 			Nip:        *po.Nip,
// 			KaryaCode:  *po.Karyacode,
// 			KaryaName:  po.Karyaname,
// 			DispName:   po.DispName,
// 			PassMesin:  *po.PassMesin,
// 			RFIDCard:   *po.RfidCard,
// 			Finger:     po.KodeFinger,
// 			Department: po.DepartCode,
// 			Status:     po.Status,
// 		})
// 	}
// 	return rv, nil
// }

func (r *shiftscheduleRepo) DeleteShiftSchedule(ctx context.Context, id int64) error {
	err := r.data.db.ShiftSchedule.
		DeleteOneID(id).
		Exec(ctx)

	if err != nil {
		return err
	}

	// Invalidate caches related to the deleted employer
	r.redisClient.Del(ctx, "Shiftscheduler:list")
	r.redisClient.Del(ctx, "Shiftscheduler:"+string(id))
	r.redisClient.Del(ctx, "Shiftscheduler:"+string(id)) // Invalidate detail cache

	r.Log.Infof("Shiftscheduler with ID %d deleted", id)
	return nil
}

// func (r *shiftscheduleRepo) GetEmployerDetail(ctx context.Context, id int64) (*biz.EmployerData, error) {
// 	cacheKey := "employerDetail:" + string(id)
// 	val, err := r.redisClient.Get(ctx, cacheKey).Result()
// 	if err == nil {
// 		var cachedEmployerData biz.EmployerData
// 		unmarshalErr := json.Unmarshal([]byte(val), &cachedEmployerData)
// 		if unmarshalErr := json.Unmarshal([]byte(val), &cachedEmployerData); unmarshalErr == nil {
// 			r.Log.Infof("Cache hit for employer detail ID: %d", id)
// 			return &cachedEmployerData, nil
// 		}
// 		r.Log.Warnf("Failed to unmarshal cached employer detail for ID %d: %v", id, unmarshalErr)
// 	} else if err != redis.Nil {
// 		r.Log.Errorf("Redis GET error for employer detail ID %d: %v", id, err)
// 	}

// 	emp, err := r.data.db.ShiftSchedule.Get(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	employerData := &biz.EmployerData{
// 		Id:        emp.ID,
// 		NoSap:     *emp.Nosap,
// 		Nip:       *emp.Nip,
// 		KaryaCode: *emp.Karyacode,
// 		KaryaName: emp.Karyaname,
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
// 			employerData.Finger = []biz.FingerData{
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

// 	if employerData.Finger == nil && emp.KodeFinger != "" { // If not in cache, or unmarshal failed, fetch from service
// 		fingerResp, err := r.bioClient.GetFingerByKode(ctx, &biometricV1.GetFingerByKodeRequest{
// 			Fingercode: emp.KodeFinger,
// 		})
// 		if err != nil {
// 			r.Log.Errorf("Failed to fetch biometric for detail: %v", err)
// 		} else {
// 			employerData.Finger = []biz.FingerData{
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
// 			employerData.Department = []biz.DepartData{
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

// 	if employerData.Department == nil && emp.DepartCode != "" { // If not in cache, or unmarshal failed, fetch from service
// 		deptResp, err := r.deptClient.GetDepartmentCode(ctx, &departmentV1.GetDepartmentCodeRequest{
// 			DepartCode: emp.DepartCode,
// 		})
// 		if err != nil {
// 			r.Log.Errorf("Failed to fetch department for detail: %v", err)
// 		} else {
// 			employerData.Department = []biz.DepartData{
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

// 	// Cache the employer detail
// 	employerDataBytes, marshalErr := json.Marshal(employerData)
// 	if marshalErr == nil {
// 		r.redisClient.Set(ctx, cacheKey, employerDataBytes, 5*time.Minute) // Cache for 5 minutes
// 	} else {
// 		r.Log.Errorf("Failed to marshal employer detail for caching: %v", marshalErr)
// 	}

// 	return employerData, nil
// }
