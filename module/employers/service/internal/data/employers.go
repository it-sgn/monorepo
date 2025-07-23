package data

import (
	"context"
	"encoding/json"
	biometricV1 "mall-go/api/biometrics/service/v1"
	departmentV1 "mall-go/api/department/service/v1"
	"mall-go/module/employers/service/internal/biz"
	"mall-go/pkg/utils/pagination"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

// var _ biz.BeerRepo = (*beerRepo)(nil)
var _ biz.EmployersRepo = (*employersRepo)(nil)

type employersRepo struct {
	data        *Data
	Log         *log.Helper
	bioClient   biometricV1.BiometricClient
	deptClient  departmentV1.DepartmentClient
	redisClient *redis.Client // Add Redis client
}

func NewEmployersRepo(data *Data, bio biometricV1.BiometricClient, dept departmentV1.DepartmentClient, logger log.Logger, rdb *redis.Client) biz.EmployersRepo {
	return &employersRepo{
		data:        data,
		bioClient:   bio,
		deptClient:  dept,
		redisClient: rdb,                   // Initialize Redis client
		Log:         log.NewHelper(logger), // ✅ INI WAJIB ADA
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
	r.redisClient.Del(ctx, "employers:list")
	r.redisClient.Del(ctx, "employers:"+string(po.ID)) //  Invalidate specific employer cache if exists

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
	// Invalidate relevant caches after update
	r.redisClient.Del(ctx, "employers:list")
	r.redisClient.Del(ctx, "employers:"+string(po.ID)) // Invalidate specific employer cache

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
	cacheKey := "employer:" + string(id)
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedEmployer biz.Employers
		unmarshalErr := json.Unmarshal([]byte(val), &cachedEmployer)
		if unmarshalErr := json.Unmarshal([]byte(val), &cachedEmployer); unmarshalErr == nil {
			r.Log.Infof("Cache hit for employer ID: %d", id)
			return &cachedEmployer, nil
		}
		r.Log.Warnf("Failed to unmarshal cached employer for ID %d: %v", id, unmarshalErr)
	} else if err != redis.Nil {
		r.Log.Errorf("Redis GET error for employer ID %d: %v", id, err)
	}

	po, err := r.data.db.Employers.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	employer := &biz.Employers{
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
	}

	// Cache the result
	employerBytes, marshalErr := json.Marshal(employer)
	if marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, employerBytes, 5*time.Minute) // Cache for 5 minutes
	} else {
		r.Log.Errorf("Failed to marshal employer for caching: %v", marshalErr)
	}

	return employer, nil
}

// func (r *employersRepo) ListEmployers(ctx context.Context, pageNum, pageSize int64) ([]*biz.EmployerData, int, error) {
// 	query := r.data.db.Employers.Query()

// 	// Hitung total sebelum pagination
// 	total, err := query.Clone().Count(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	pos, err := r.data.db.Employers.Query().
// 		Offset(int(pagination.GetPageOffset(pageNum, pageSize))).
// 		Limit(int(pageSize)).
// 		All(ctx)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	if r.Log != nil {
// 		r.Log.Warnf("ListEmployers result: %+v", pos)
// 	}
// 	rv := make([]*biz.EmployerData, 0, len(pos))
// 	for _, po := range pos {
// 		empData := &biz.EmployerData{
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
func (r *employersRepo) ListEmployers(ctx context.Context, pageNum, pageSize int64) ([]*biz.EmployerData, int, error) {
	cacheKey := "employers:list:page:" + string(pageNum) + ":size:" + string(pageSize)
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedResponse struct {
			Employers []*biz.EmployerData
			Total     int
		}
		var cachedEmployer biz.Employers
		unmarshalErr := json.Unmarshal([]byte(val), &cachedEmployer)
		if unmarshalErr := json.Unmarshal([]byte(val), &cachedResponse); unmarshalErr == nil {
			r.Log.Infof("Cache hit for ListEmployers page %d, size %d", pageNum, pageSize)
			return cachedResponse.Employers, cachedResponse.Total, nil
		}

		r.Log.Warnf("Failed to unmarshal cached ListEmployers: %v", unmarshalErr)
	} else if err != redis.Nil {
		r.Log.Errorf("Redis GET error for ListEmployers: %v", err)
	}

	query := r.data.db.Employers.Query()

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
		r.Log.Warnf("ListEmployers result from DB: %+v", pos)
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

		// Get finger data with caching
		fingerCacheKey := "biometric:" + po.KodeFinger
		fingerVal, fingerErr := r.redisClient.Get(ctx, fingerCacheKey).Result()
		if fingerErr == nil {
			var cachedFinger biometricV1.GetFingerByKodeResponse
			if unmarshalErr := json.Unmarshal([]byte(fingerVal), &cachedFinger); unmarshalErr == nil {
				empData.Finger = []biz.FingerData{
					{
						Fingercode: cachedFinger.Fingercode,
						Finger0:    cachedFinger.Finger0,
						Finger1:    cachedFinger.Finger1,
						Finger2:    cachedFinger.Finger2,
						Finger3:    cachedFinger.Finger3,
						Finger4:    cachedFinger.Finger4,
						Finger5:    cachedFinger.Finger5,
						Finger6:    cachedFinger.Finger6,
						Finger7:    cachedFinger.Finger7,
						Finger8:    cachedFinger.Finger8,
						Finger9:    cachedFinger.Finger9,
					},
				}
				r.Log.Infof("Cache hit for biometric code: %s", po.KodeFinger)
			} else {
				r.Log.Warnf("Failed to unmarshal cached biometric for code %s: %v", po.KodeFinger, unmarshalErr)
			}
		} else if fingerErr != redis.Nil {
			r.Log.Errorf("Redis GET error for biometric code %s: %v", po.KodeFinger, fingerErr)
		}

		if empData.Finger == nil && po.KodeFinger != "" { // If not in cache, or unmarshal failed, fetch from service
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
				// Cache the finger data
				fingerBytes, marshalErr := json.Marshal(fingerResp)
				if marshalErr == nil {
					r.redisClient.Set(ctx, fingerCacheKey, fingerBytes, 5*time.Minute)
				} else {
					r.Log.Errorf("Failed to marshal biometric for caching: %v", marshalErr)
				}
			} else {
				r.Log.Warnf("biometric not found for kode: %s, err: %v", po.KodeFinger, err)
			}
		}

		// Department with caching
		// r.Log.Infow("", po.DepartCode)
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
		Employers []*biz.EmployerData
		Total     int
	}{
		Employers: rv,
		Total:     total,
	}
	responseBytes, marshalErr := json.Marshal(responseToCache)
	if marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, responseBytes, 5*time.Minute) // Cache for 5 minutes
	} else {
		r.Log.Errorf("Failed to marshal ListEmployers response for caching: %v", marshalErr)
	}

	return rv, total, nil
}
func (r *employersRepo) Count(ctx context.Context) (int, error) {
	dt, _ := r.data.db.Employers.Query().Count(ctx)
	// fmt.Println("INI PL DT: ")
	return dt, nil
}
func (r *employersRepo) ListEmployersNext(ctx context.Context, start, end int32) ([]*biz.Employers, error) {
	// You can apply similar caching logic here if this method is heavily used
	// and benefits from caching. The cache key would need to reflect start and end.
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

	// Invalidate caches related to the deleted employer
	r.redisClient.Del(ctx, "employers:list")
	r.redisClient.Del(ctx, "employer:"+string(id))
	r.redisClient.Del(ctx, "employerDetail:"+string(id)) // Invalidate detail cache

	r.Log.Infof("employer with ID %d deleted", id)
	return nil
}

func (r *employersRepo) GetEmployerDetail(ctx context.Context, id int64) (*biz.EmployerData, error) {
	cacheKey := "employerDetail:" + string(id)
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedEmployerData biz.EmployerData
		unmarshalErr := json.Unmarshal([]byte(val), &cachedEmployerData)
		if unmarshalErr := json.Unmarshal([]byte(val), &cachedEmployerData); unmarshalErr == nil {
			r.Log.Infof("Cache hit for employer detail ID: %d", id)
			return &cachedEmployerData, nil
		}
		r.Log.Warnf("Failed to unmarshal cached employer detail for ID %d: %v", id, unmarshalErr)
	} else if err != redis.Nil {
		r.Log.Errorf("Redis GET error for employer detail ID %d: %v", id, err)
	}

	emp, err := r.data.db.Employers.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	employerData := &biz.EmployerData{
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
	}

	// Get finger data with caching
	fingerCacheKey := "biometric:" + emp.KodeFinger
	fingerVal, fingerErr := r.redisClient.Get(ctx, fingerCacheKey).Result()
	if fingerErr == nil {
		var cachedFinger biometricV1.GetFingerByKodeResponse
		if unmarshalErr := json.Unmarshal([]byte(fingerVal), &cachedFinger); unmarshalErr == nil {
			employerData.Finger = []biz.FingerData{
				{
					Fingercode: cachedFinger.Fingercode,
					Finger0:    cachedFinger.Finger0,
					Finger1:    cachedFinger.Finger1,
					Finger2:    cachedFinger.Finger2,
					Finger3:    cachedFinger.Finger3,
					Finger4:    cachedFinger.Finger4,
					Finger5:    cachedFinger.Finger5,
					Finger6:    cachedFinger.Finger6,
					Finger7:    cachedFinger.Finger7,
					Finger8:    cachedFinger.Finger8,
					Finger9:    cachedFinger.Finger9,
				},
			}
			r.Log.Infof("Cache hit for biometric code: %s (detail)", emp.KodeFinger)
		} else {
			r.Log.Warnf("Failed to unmarshal cached biometric for code %s (detail): %v", emp.KodeFinger, unmarshalErr)
		}
	} else if fingerErr != redis.Nil {
		r.Log.Errorf("Redis GET error for biometric code %s (detail): %v", emp.KodeFinger, fingerErr)
	}

	if employerData.Finger == nil && emp.KodeFinger != "" { // If not in cache, or unmarshal failed, fetch from service
		fingerResp, err := r.bioClient.GetFingerByKode(ctx, &biometricV1.GetFingerByKodeRequest{
			Fingercode: emp.KodeFinger,
		})
		if err != nil {
			r.Log.Errorf("Failed to fetch biometric for detail: %v", err)
		} else {
			employerData.Finger = []biz.FingerData{
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
			// Cache the finger data
			fingerBytes, marshalErr := json.Marshal(fingerResp)
			if marshalErr == nil {
				r.redisClient.Set(ctx, fingerCacheKey, fingerBytes, 5*time.Minute)
			} else {
				r.Log.Errorf("Failed to marshal biometric for caching (detail): %v", marshalErr)
			}
		}
	}

	// Get department data with caching
	deptCacheKey := "department:" + emp.DepartCode
	deptVal, deptErr := r.redisClient.Get(ctx, deptCacheKey).Result()
	if deptErr == nil {
		var cachedDept departmentV1.GetDepartmentCodeResponse
		if unmarshalErr := json.Unmarshal([]byte(deptVal), &cachedDept); unmarshalErr == nil {
			employerData.Department = []biz.DepartData{
				{
					DepartCode: cachedDept.DepartCode,
					DepartName: cachedDept.DepartName,
				},
			}
			r.Log.Infof("Cache hit for department code: %s (detail)", emp.DepartCode)
		} else {
			r.Log.Warnf("Failed to unmarshal cached department for code %s (detail): %v", emp.DepartCode, unmarshalErr)
		}
	} else if deptErr != redis.Nil {
		r.Log.Errorf("Redis GET error for department code %s (detail): %v", emp.DepartCode, deptErr)
	}

	if employerData.Department == nil && emp.DepartCode != "" { // If not in cache, or unmarshal failed, fetch from service
		deptResp, err := r.deptClient.GetDepartmentCode(ctx, &departmentV1.GetDepartmentCodeRequest{
			DepartCode: emp.DepartCode,
		})
		if err != nil {
			r.Log.Errorf("Failed to fetch department for detail: %v", err)
		} else {
			employerData.Department = []biz.DepartData{
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
				r.Log.Errorf("Failed to marshal department for caching (detail): %v", marshalErr)
			}
		}
	}

	// Cache the employer detail
	employerDataBytes, marshalErr := json.Marshal(employerData)
	if marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, employerDataBytes, 5*time.Minute) // Cache for 5 minutes
	} else {
		r.Log.Errorf("Failed to marshal employer detail for caching: %v", marshalErr)
	}

	return employerData, nil
}
