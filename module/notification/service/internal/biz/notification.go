package biz

import (
	"mall-go/pkg/page_token"

	"golang.org/x/sync/singleflight"

	"github.com/go-kratos/kratos/v2/log"
	// "github.com/go-kratos/beer-shop/pkg/page_token"

	biometricV1 "mall-go/api/biometrics/service/v1"
	departmentV1 "mall-go/api/department/service/v1"
)

type Notification struct {
	Id         int64
	NoSap      string
	Nip        string
	KaryaCode  string
	KaryaName  string
	DispName   string
	PassMesin  string
	RFIDCard   string
	Finger     string
	Department string
	Status     int32
	CreatedAt  string
	UpdatedAt  string
}
type DepartData struct {
	DepartCode string
	DepartName string
}
type FingerData struct {
	Fingercode string
	Finger0    string
	Finger1    string
	Finger2    string
	Finger3    string
	Finger4    string
	Finger5    string
	Finger6    string
	Finger7    string
	Finger8    string
	Finger9    string
}

type EmployerData struct {
	Id         int64
	NoSap      string
	Nip        string
	KaryaCode  string
	KaryaName  string
	DispName   string
	PassMesin  string
	RFIDCard   string
	Finger     []FingerData
	Department []DepartData
	Status     int32
	CreatedAt  string
	UpdatedAt  string
}

type NotificationRepo interface {
	// CreateNotification(ctx context.Context, c *Notification) (*Notification, error)
	// UpdateNotification(ctx context.Context, c *Notification) (*Notification, error)
	// GetNotificationID(ctx context.Context, id int64) (*Notification, error)
	// DeleteNotification(ctx context.Context, id int64) error
	// ListNotification(ctx context.Context, pageNum, pageSize int64) ([]*EmployerData, int, error)
	// ListNotificationNext(ctx context.Context, start, end int32) ([]*Notification, error)
	// Count(ctx context.Context) (int, error)
	// GetEmployerDetail(ctx context.Context, id int64) (*EmployerData, error)
	// GetFingerByKode(ctx context.Context, fkode string) (*FingerData, error)
	// GetDepartmentByKode(ctx context.Context, dkode string) (*DepartData, error)
	// GetEmployeeFullData(ctx context.Context, fingerID string, departmentID int64) (*FullEmployeeData, error)

	//
	// GetFingerByID(ctx context.Context, id int64) (*KodeFingerRecord, error)
	// GetFingerByKode(ctx context.Context, kodefinger string) (*KodeFingerRecord, error)
	// CreateFinger(ctx context.Context, f *KodeFingerRecord) (*KodeFingerRecord, error)
	// UpdateFinger(ctx context.Context, f *KodeFingerRecord) (*KodeFingerRecord, error)
	// DeleteFinger(ctx context.Context, id int64) error
}
type NotificationUsecase struct {
	repo       NotificationRepo
	pageToken  page_token.ProcessPageTokens
	log        *log.Helper
	deptClient departmentV1.DepartmentClient
	bioCilent  biometricV1.BiometricClient
	sg         singleflight.Group
}

func NewNotificationUsecase(repo NotificationRepo, logger log.Logger, dept departmentV1.DepartmentClient, bio biometricV1.BiometricClient) *NotificationUsecase {
	return &NotificationUsecase{
		repo:       repo,
		log:        log.NewHelper(log.With(logger, "module", "usecase/notification")),
		deptClient: dept,
		bioCilent:  bio,
	}
}

// func (uc *NotificationUsecase) Create(ctx context.Context, u *Notification) (*Notification, error) {
// 	return uc.repo.CreateNotification(ctx, u)
// }
// func (uc *NotificationUsecase) Update(ctx context.Context, u *Notification) (*Notification, error) {
// 	return uc.repo.UpdateNotification(ctx, u)
// }

// // get by id
// func (uc *NotificationUsecase) GetByID(ctx context.Context, id int64) (*Notification, error) {
// 	return uc.repo.GetNotificationID(ctx, id)
// }

// func (uc *NotificationUsecase) List(ctx context.Context, pageNum, pageSize int64) ([]*EmployerData, int, error) {
// 	uc.log.Infof("ListNotification called with pageNum=%d, pageSize=%d", pageNum, pageSize)

// 	// return uc.repo.ListNotification(ctx, pageNum, pageSize)
// 	list, total, err := uc.repo.ListNotification(ctx, pageNum, pageSize)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	if list == nil {
// 		return nil, 0, fmt.Errorf("repo.ListNotification returned nil list")
// 	}

// 	return list, total, nil
// }

// func (uc *NotificationUsecase) ListNext(ctx context.Context, pageSize int32, pageToken string) ([]*Notification, string, error) {
// 	total, err := uc.repo.Count(ctx)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	// log.Error("INI TOTAL :", total)

// 	start, end, nextToken, err := uc.pageToken.ProcessPageTokens(total, pageSize, pageToken)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	// log.Error("PAGE TOKEN :", pageToken)
// 	// single flight
// 	data, err, _ := uc.sg.Do(fmt.Sprintf("list_next_%d_%d", start, end), func() (interface{}, error) {
// 		return uc.repo.ListNotificationNext(ctx, start, end)
// 	})
// 	return data.([]*Notification), nextToken, err
// }

// func (uc *NotificationUsecase) Delete(ctx context.Context, id int64) error {
// 	return uc.repo.DeleteNotification(ctx, id)
// }

// // func (uc *NotificationUsecase) GetFingerByKode(ctx context.Context, fcode string) (*FingerData, error) {
// // 	return uc.repo.GetFingerByKode(ctx, fcode)
// // }
// // func (uc *NotificationUsecase) GetDepartmentByKode(ctx context.Context, dkode string) (*DepartData, error) {
// // 	return uc.repo.GetDepartmentByKode(ctx, dkode)
// // }

// func (uc *NotificationUsecase) GetDepartmentByCode(ctx context.Context, deptCode string) (*departmentV1.GetDepartmentCodeResponse, error) {
// 	return uc.deptClient.GetDepartmentCode(ctx, &departmentV1.GetDepartmentCodeRequest{DepartCode: deptCode})
// }

// func (uc *NotificationUsecase) GetFingerByKode(ctx context.Context, bioCode string) (*biometricV1.GetFingerByKodeResponse, error) {
// 	return uc.bioCilent.GetFingerByKode(ctx, &biometricV1.GetFingerByKodeRequest{Fingercode: bioCode})
// }

// func (uc *NotificationUsecase) GetDetail(ctx context.Context, id int64) (*EmployerData, error) {
// 	return uc.repo.GetEmployerDetail(ctx, id)
// }
