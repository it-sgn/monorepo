package data

import (
	"context"
	"encoding/json"
	"fmt"
	biometricV1 "mall-go/api/biometrics/service/v1"
	departmentV1 "mall-go/api/department/service/v1"
	orgv1 "mall-go/api/organization/service/v1"
	"mall-go/module/employers/service/internal/biz"
	"mall-go/module/employers/service/internal/data/model/employers"
	"mall-go/pkg/utils/pagination"
	"strconv"
	"strings"
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
	orgv1       orgv1.OrganizationServiceClient
	redisClient *redis.Client // Add Redis client
}

func NewEmployersRepo(data *Data, bio biometricV1.BiometricClient, dept departmentV1.DepartmentClient, org orgv1.OrganizationServiceClient, logger log.Logger, rdb *redis.Client) biz.EmployersRepo {
	return &employersRepo{
		data:        data,
		bioClient:   bio,
		deptClient:  dept,
		orgv1:       org,
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
		SetKodePerusahaan(b.KodePerusahaan).
		SetKodeCabang(b.KodeCabang).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	r.redisClient.Del(ctx, "employers:list")
	id := strconv.FormatInt(po.ID, 10)
	r.redisClient.Del(ctx, "employers:"+id) //  Invalidate specific employer cache if exists
	// r.redisClient.Del(ctx, "employers:"+string(po.ID)) //  Invalidate specific employer cache if exists

	return &biz.Employers{
		Id:             po.ID,
		NoSap:          *po.Nosap,
		Nip:            *po.Nip,
		KaryaCode:      *po.Karyacode,
		KaryaName:      po.Karyaname,
		DispName:       po.DispName,
		PassMesin:      *po.PassMesin,
		RFIDCard:       *po.RfidCard,
		Finger:         po.KodeFinger,
		Department:     po.DepartCode,
		Status:         po.Status,
		KodePerusahaan: po.KodePerusahaan,
		KodeCabang:     po.KodeCabang,
		CreatedAt:      po.CreateTime.String(),
		UpdatedAt:      po.UpdateTime.String(),
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
		SetKodePerusahaan(b.KodePerusahaan).
		SetKodeCabang(b.KodeCabang).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	// Invalidate relevant caches after update
	r.redisClient.Del(ctx, "employers:list")
	id := strconv.FormatInt(po.ID, 10)
	r.redisClient.Del(ctx, "employers:"+id)
	// r.redisClient.Del(ctx, "employers:"+string(po.ID)) // Invalidate specific employer cache

	return &biz.Employers{
		Id:             po.ID,
		NoSap:          *po.Nosap,
		Nip:            *po.Nip,
		KaryaCode:      *po.Karyacode,
		KaryaName:      po.Karyaname,
		DispName:       po.DispName,
		PassMesin:      *po.PassMesin,
		RFIDCard:       *po.RfidCard,
		Finger:         po.KodeFinger,
		Department:     po.DepartCode,
		Status:         po.Status,
		KodePerusahaan: po.KodePerusahaan,
		KodeCabang:     po.KodeCabang,
	}, nil
}

func (r *employersRepo) GetEmployersID(ctx context.Context, id int64) (*biz.Employers, error) {
	idx := strconv.FormatInt(id, 10)
	cacheKey := "employer:" + idx
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
		Id:             po.ID,
		NoSap:          *po.Nosap,
		Nip:            *po.Nip,
		KaryaCode:      *po.Karyacode,
		KaryaName:      po.Karyaname,
		DispName:       po.DispName,
		PassMesin:      *po.PassMesin,
		RFIDCard:       *po.RfidCard,
		Finger:         po.KodeFinger,
		Department:     po.DepartCode,
		Status:         po.Status,
		KodePerusahaan: po.KodePerusahaan,
		KodeCabang:     po.KodeCabang,
	}

	// Cache the result
	employerBytes, marshalErr := json.Marshal(employer)
	if marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, employerBytes, 3*time.Minute) // Cache for 5 minutes
	} else {
		r.Log.Errorf("Failed to marshal employer for caching: %v", marshalErr)
	}

	return employer, nil
}

func (r *employersRepo) GetEmployersKode(ctx context.Context, karya_code string) (*biz.EmployerKode, error) {
	cacheKey := "employer:" + string(karya_code)
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedEmployer biz.EmployerKode
		unmarshalErr := json.Unmarshal([]byte(val), &cachedEmployer)
		if unmarshalErr := json.Unmarshal([]byte(val), &cachedEmployer); unmarshalErr == nil {
			r.Log.Infof("Cache hit for employer kode: %d", karya_code)
			return &cachedEmployer, nil
		}
		r.Log.Warnf("Failed to unmarshal cached employer for karya_code %d: %v", karya_code, unmarshalErr)
	} else if err != redis.Nil {
		r.Log.Errorf("Redis GET error for employer karya_code %d: %v", karya_code, err)
	}

	po, err := r.data.db.Employers.Query().
		Where(employers.KaryacodeEQ(strings.TrimSpace(karya_code))).
		First(ctx)
	if err != nil {
		return nil, err
	}

	employer := &biz.EmployerKode{
		NoSap:          *po.Nosap,
		Nip:            *po.Nip,
		KaryaCode:      *po.Karyacode,
		KaryaName:      po.Karyaname,
		KodePerusahaan: po.KodePerusahaan,
		KodeCabang:     po.KodeCabang,
	}

	// Cache the result
	employerBytes, marshalErr := json.Marshal(employer)
	if marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, employerBytes, 3*time.Minute) // Cache for 5 minutes
	} else {
		r.Log.Errorf("Failed to marshal employer for caching: %v", marshalErr)
	}

	return employer, nil
}

func (r *employersRepo) ListEmployers(ctx context.Context, pageNum, pageSize int64) ([]*biz.EmployerData, int, error) {
	page := strconv.FormatInt(pageNum, 10)
	size := strconv.FormatInt(pageSize, 10)
	cacheKey := "employers:list:page:" + page + ":size:" + size
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
			Id:             po.ID,
			NoSap:          *po.Nosap,
			Nip:            *po.Nip,
			KaryaCode:      *po.Karyacode,
			KaryaName:      po.Karyaname,
			DispName:       po.DispName,
			PassMesin:      *po.PassMesin,
			RFIDCard:       *po.RfidCard,
			Status:         po.Status,
			KodePerusahaan: po.KodePerusahaan,
			KodeCabang:     po.KodeCabang,
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
					r.redisClient.Set(ctx, fingerCacheKey, fingerBytes, 3*time.Minute)
				} else {
					r.Log.Errorf("Failed to marshal biometric for caching: %v", marshalErr)
				}
			} else {
				r.Log.Warnf("biometric not found for kode: %s, err: %v", po.KodeFinger, err)
			}
		}

		// Department with caching
		r.Log.Infow("", po.DepartCode)
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
					r.redisClient.Set(ctx, deptCacheKey, deptBytes, 3*time.Minute)
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
		r.redisClient.Set(ctx, cacheKey, responseBytes, 3*time.Minute) // Cache for 5 minutes
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
			Id:             po.ID,
			NoSap:          *po.Nosap,
			Nip:            *po.Nip,
			KaryaCode:      *po.Karyacode,
			KaryaName:      po.Karyaname,
			DispName:       po.DispName,
			PassMesin:      *po.PassMesin,
			RFIDCard:       *po.RfidCard,
			Finger:         po.KodeFinger,
			Department:     po.DepartCode,
			Status:         po.Status,
			KodePerusahaan: po.KodePerusahaan,
			KodeCabang:     po.KodeCabang,
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
	ids := strconv.FormatInt(id, 10)
	r.redisClient.Del(ctx, "employers:list")
	r.redisClient.Del(ctx, "employers:filterbydepartcode")
	r.redisClient.Del(ctx, "employers:filterbykodes")
	r.redisClient.Del(ctx, "employers")
	r.redisClient.Del(ctx, "employer")
	r.redisClient.Del(ctx, "employer:"+ids)
	r.redisClient.Del(ctx, "employerDetail:"+ids) // Invalidate detail cache
	r.redisClient.Del(ctx)

	r.Log.Infof("employer with ID %d deleted", id)
	return nil
}

func (r *employersRepo) GetEmployerDetail(ctx context.Context, id int64) (*biz.EmployerData, error) {
	ids := strconv.FormatInt(id, 10)
	cacheKey := "employerDetail:" + ids
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
		Id:             emp.ID,
		NoSap:          *emp.Nosap,
		Nip:            *emp.Nip,
		KaryaCode:      *emp.Karyacode,
		KaryaName:      emp.Karyaname,
		DispName:       emp.DispName,
		PassMesin:      *emp.PassMesin,
		RFIDCard:       *emp.RfidCard,
		Status:         emp.Status,
		KodePerusahaan: emp.KodePerusahaan,
		KodeCabang:     emp.KodeCabang,
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
				r.redisClient.Set(ctx, fingerCacheKey, fingerBytes, 3*time.Minute)
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
				r.redisClient.Set(ctx, deptCacheKey, deptBytes, 3*time.Minute)
			} else {
				r.Log.Errorf("Failed to marshal department for caching (detail): %v", marshalErr)
			}
		}
	}

	// Cache the employer detail
	employerDataBytes, marshalErr := json.Marshal(employerData)
	if marshalErr == nil {
		r.redisClient.Set(ctx, cacheKey, employerDataBytes, 3*time.Minute) // Cache for 5 minutes
	} else {
		r.Log.Errorf("Failed to marshal employer detail for caching: %v", marshalErr)
	}

	return employerData, nil
}

func (r *employersRepo) GetByFilter(ctx context.Context, karyacodes []string) ([]*biz.EmployerItem, error) {
	cacheKey := fmt.Sprintf("employers:filterbykodes:%s", strings.Join(karyacodes, ","))

	// 🔍 Coba ambil dari cache Redis
	if val, err := r.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var cachedEmployers []*biz.EmployerItem
		if err := json.Unmarshal([]byte(val), &cachedEmployers); err == nil {
			r.Log.Infof("✅ Cache hit for key: %s", cacheKey)
			return cachedEmployers, nil
		}
		r.Log.Warnf("⚠️ Failed to unmarshal cache for key %s: %v", cacheKey, err)
	} else if err != redis.Nil {
		r.Log.Errorf("❌ Redis GET error for key %s: %v", cacheKey, err)
	}

	// 🗃️ Query ke DB jika cache tidak tersedia
	query := r.data.db.Employers.Query()
	if len(karyacodes) > 0 {
		query = query.Where(employers.KaryacodeIn(karyacodes...))
	}

	entities, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	var bizEmployers []*biz.EmployerItem
	for _, e := range entities {
		emp := &biz.EmployerItem{
			KaryaCode:      deref(e.Karyacode),
			KaryaName:      deref(&e.Karyaname),
			KodePerusahaan: deref(&e.KodePerusahaan),
			KodeCabang:     deref(&e.KodeCabang),
		}

		deptCode := e.DepartCode
		deptCacheKey := "department:" + deptCode

		// 📦 Coba ambil departemen dari Redis
		if deptVal, err := r.redisClient.Get(ctx, deptCacheKey).Result(); err == nil {
			var cachedDept departmentV1.GetDepartmentCodeResponse
			if err := json.Unmarshal([]byte(deptVal), &cachedDept); err == nil {
				emp.Department = biz.DepartData{
					DepartCode: cachedDept.DepartCode,
					DepartName: cachedDept.DepartName,
				}
				r.Log.Infof("✅ Department cache hit for code: %s", deptCode)
			} else {
				r.Log.Warnf("⚠️ Failed to unmarshal department cache for code %s: %v", deptCode, err)
			}
		} else if err != redis.Nil {
			r.Log.Errorf("❌ Redis GET error for department code %s: %v", deptCode, err)
		}

		// ☎️ Jika tidak ditemukan di cache, ambil dari service
		if emp.Department.DepartCode == "" && deptCode != "" {
			deptResp, err := r.deptClient.GetDepartmentCode(ctx, &departmentV1.GetDepartmentCodeRequest{
				DepartCode: deptCode,
			})
			if err != nil {
				r.Log.Warnf("⚠️ Failed to fetch department for code=%s: %v", deptCode, err)
			} else {
				emp.Department = biz.DepartData{
					DepartCode: deptResp.DepartCode,
					DepartName: deptResp.DepartName,
				}

				// 💾 Simpan ke Redis cache
				if deptBytes, err := json.Marshal(deptResp); err == nil {
					_ = r.redisClient.Set(ctx, deptCacheKey, deptBytes, 3*time.Minute).Err()
					r.Log.Infof("📥 Department cached for code: %s", deptCode)
				} else {
					r.Log.Warnf("⚠️ Failed to marshal department for cache: %v", err)
				}
			}
		}

		bizEmployers = append(bizEmployers, emp)
	}

	// 💾 Simpan hasil employers ke Redis cache
	if data, err := json.Marshal(bizEmployers); err == nil {
		if err := r.redisClient.Set(ctx, cacheKey, data, 3*time.Minute).Err(); err != nil {
			r.Log.Warnf("⚠️ Failed to cache employers for key %s: %v", cacheKey, err)
		} else {
			r.Log.Infof("📥 Cached %d employers for key: %s", len(bizEmployers), cacheKey)
		}
	}

	return bizEmployers, nil
}

func deref(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func (r *employersRepo) GetByDepart(ctx context.Context, departCode string) ([]*biz.EmployerItem, error) {
	cacheKey := fmt.Sprintf("employers:filterbydepartcode:%s", departCode)

	// 🔍 Coba ambil dari cache Redis
	if val, err := r.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var cachedEmployers []*biz.EmployerItem
		if err := json.Unmarshal([]byte(val), &cachedEmployers); err == nil {
			r.Log.Infof("✅ Cache hit for key: %s", cacheKey)
			return cachedEmployers, nil
		}
		r.Log.Warnf("⚠️ Failed to unmarshal cache for key %s: %v", cacheKey, err)
	} else if err != redis.Nil {
		r.Log.Errorf("❌ Redis GET error for key %s: %v", cacheKey, err)
	}

	// 🗃️ Query ke DB jika cache tidak tersedia
	query := r.data.db.Employers.Query()
	if len(departCode) > 0 {
		query = query.Where(employers.DepartCodeEQ(departCode))
	}

	entities, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	var bizEmployers []*biz.EmployerItem
	for _, e := range entities {
		emp := &biz.EmployerItem{
			KaryaCode:      deref(e.Karyacode),
			KaryaName:      deref(&e.Karyaname),
			KodePerusahaan: deref(&e.KodePerusahaan),
			KodeCabang:     deref(&e.KodeCabang),
		}

		deptCode := e.DepartCode
		deptCacheKey := "department:" + deptCode

		// 📦 Coba ambil departemen dari Redis
		if deptVal, err := r.redisClient.Get(ctx, deptCacheKey).Result(); err == nil {
			var cachedDept departmentV1.GetDepartmentCodeResponse
			if err := json.Unmarshal([]byte(deptVal), &cachedDept); err == nil {
				emp.Department = biz.DepartData{
					DepartCode: cachedDept.DepartCode,
					DepartName: cachedDept.DepartName,
				}
				r.Log.Infof("✅ Department cache hit for code: %s", deptCode)
			} else {
				r.Log.Warnf("⚠️ Failed to unmarshal department cache for code %s: %v", deptCode, err)
			}
		} else if err != redis.Nil {
			r.Log.Errorf("❌ Redis GET error for department code %s: %v", deptCode, err)
		}

		// ☎️ Jika tidak ditemukan di cache, ambil dari service
		if emp.Department.DepartCode == "" && deptCode != "" {
			deptResp, err := r.deptClient.GetDepartmentCode(ctx, &departmentV1.GetDepartmentCodeRequest{
				DepartCode: deptCode,
			})
			if err != nil {
				r.Log.Warnf("⚠️ Failed to fetch department for code=%s: %v", deptCode, err)
			} else {
				emp.Department = biz.DepartData{
					DepartCode: deptResp.DepartCode,
					DepartName: deptResp.DepartName,
				}

				// 💾 Simpan ke Redis cache
				if deptBytes, err := json.Marshal(deptResp); err == nil {
					_ = r.redisClient.Set(ctx, deptCacheKey, deptBytes, 3*time.Minute).Err()
					r.Log.Infof("📥 Department cached for code: %s", deptCode)
				} else {
					r.Log.Warnf("⚠️ Failed to marshal department for cache: %v", err)
				}
			}
		}

		bizEmployers = append(bizEmployers, emp)
	}

	// 💾 Simpan hasil employers ke Redis cache
	if data, err := json.Marshal(bizEmployers); err == nil {
		if err := r.redisClient.Set(ctx, cacheKey, data, 3*time.Minute).Err(); err != nil {
			r.Log.Warnf("⚠️ Failed to cache employers for key %s: %v", cacheKey, err)
		} else {
			r.Log.Infof("📥 Cached %d employers for key: %s", len(bizEmployers), cacheKey)
		}
	}

	return bizEmployers, nil
}

// func (r *employersRepo) GetPerusahaan(ctx context.Context, karya_code string) ([]*biz.PerusahaanData, error) {
// 	cacheKey := fmt.Sprintf("employers:company:%s", karya_code)
// 	if val, err := r.redisClient.Get(ctx, cacheKey).Result(); err == nil {
// 		var cachedData []*biz.PerusahaanData
// 		if err := json.Unmarshal([]byte(val), &cachedData); err == nil {
// 			r.Log.Infof("✅ Cache hit for key: %s", cacheKey)
// 			return cachedData, nil
// 		}
// 		r.Log.Warnf("⚠️ Failed to unmarshal cache for key %s: %v", cacheKey, err)
// 	} else if err != redis.Nil {
// 		r.Log.Errorf("❌ Redis GET error for key %s: %v", cacheKey, err)
// 	}

// 	// 🗃️ Query ke DB jika cache tidak tersedia
// 	query := r.data.db.Employers.Query()
// 	if len(karya_code) > 0 {
// 		query = query.Where(employers.Karyacode(karya_code))
// 	}
// 	entity, err := query.First(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 🔍 Ambil detail perusahaan & cabang dari Organization Microservice
// 	orgCompany, err := r.orgv1.GetPerusahaan(ctx, &orgv1.GetPerusahaanRequest{
// 		KodePerusahaan: deref(&entity.KodePerusahaan),
// 	})
// 	if err != nil {
// 		r.Log.Errorf("❌ Gagal ambil data perusahaan dari Organization Service: %v", err)
// 		return nil, err
// 	}

// 	// 📦 Mapping ke struct biz.PerusahaanData
// 	result := &biz.PerusahaanData{
// 		KodePerusahaan: orgCompany.Result.KodePerusahaan,
// 		NamaPerusahaan: orgCompany.Result.NamaPerusahaan,
// 		KodeCabang:     orgCompany.Result.KodeCabang,
// 		Cabang:         orgCompany.Result.Cabang,
// 	}

// 	// 💾 Simpan ke Redis cache
// 	if bytes, err := json.Marshal(result); err == nil {
// 		if err := r.redisClient.Set(ctx, cacheKey, bytes, 30*time.Minute).Err(); err != nil {
// 			r.Log.Warnf("⚠️ Gagal set cache untuk key %s: %v", cacheKey, err)
// 		}
// 	}

//		return result, nil
//	}

func (r *employersRepo) GetPerusahaan(ctx context.Context, depart string) ([]*biz.PerusahaanData, error) {
	cacheKey := fmt.Sprintf("employers:company:%s", depart)
	if val, err := r.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var cachedData []*biz.PerusahaanData
		if err := json.Unmarshal([]byte(val), &cachedData); err == nil {
			r.Log.Infof("✅ Cache hit for key: %s", cacheKey)
			return cachedData, nil
		}
		r.Log.Warnf("⚠️ Failed to unmarshal cache for key %s: %v", cacheKey, err)
	} else if err != redis.Nil {
		r.Log.Errorf("❌ Redis GET error for key %s: %v", cacheKey, err)
	}

	// Query ke DB
	query := r.data.db.Employers.Query()
	if depart != "" {
		query = query.Where(employers.DepartCode(depart))
	}
	entities, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	var results []*biz.PerusahaanData
	for _, e := range entities {
		orgCompany, err := r.orgv1.GetPerusahaan(ctx, &orgv1.GetPerusahaanRequest{
			KodePerusahaan: deref(&e.KodePerusahaan),
		})
		if err != nil {
			r.Log.Errorf("❌ Gagal ambil data perusahaan dari Organization Service: %v", err)
			continue
		}

		results = append(results, &biz.PerusahaanData{
			KodePerusahaan: orgCompany.Result.KodePerusahaan,
			NamaPerusahaan: orgCompany.Result.NamaPerusahaan,
			KodeCabang:     orgCompany.Result.KodeCabang,
			Cabang:         orgCompany.Result.Cabang,
		})
	}

	// Simpan ke cache
	if bytes, err := json.Marshal(results); err == nil {
		if err := r.redisClient.Set(ctx, cacheKey, bytes, 30*time.Minute).Err(); err != nil {
			r.Log.Warnf("⚠️ Gagal set cache untuk key %s: %v", cacheKey, err)
		}
	}

	return results, nil
}

func (r *employersRepo) GetCabang(ctx context.Context, depart string) ([]*biz.CabangData, error) {
	cacheKey := fmt.Sprintf("employers:cabang:%s", depart)
	if val, err := r.redisClient.Get(ctx, cacheKey).Result(); err == nil {
		var cachedData []*biz.CabangData
		if err := json.Unmarshal([]byte(val), &cachedData); err == nil {
			r.Log.Infof("✅ Cache hit for key: %s", cacheKey)
			return cachedData, nil
		}
		r.Log.Warnf("⚠️ Failed to unmarshal cache for key %s: %v", cacheKey, err)
	} else if err != redis.Nil {
		r.Log.Errorf("❌ Redis GET error for key %s: %v", cacheKey, err)
	}

	// Query ke DB
	query := r.data.db.Employers.Query()
	if depart != "" {
		query = query.Where(employers.DepartCode(depart))
	}
	entities, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	var results []*biz.CabangData
	for _, e := range entities {
		orgCompany, err := r.orgv1.GetCabang(ctx, &orgv1.GetCabangRequest{
			KodeCabang: deref(&e.KodeCabang),
		})
		if err != nil {
			r.Log.Errorf("❌ Gagal ambil data perusahaan dari Organization Service: %v", err)
			continue
		}

		results = append(results, &biz.CabangData{
			KodePerusahaan: orgCompany.Result.KodePerusahaan,
			NamaPerusahaan: orgCompany.Result.NamaPerusahaan,
			KodeCabang:     orgCompany.Result.KodeCabang,
			Cabang:         orgCompany.Result.Cabang,
		})
	}

	// Simpan ke cache
	if bytes, err := json.Marshal(results); err == nil {
		if err := r.redisClient.Set(ctx, cacheKey, bytes, 30*time.Minute).Err(); err != nil {
			r.Log.Warnf("⚠️ Gagal set cache untuk key %s: %v", cacheKey, err)
		}
	}

	return results, nil
}
