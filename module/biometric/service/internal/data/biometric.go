package data

import (
	"context"
	"mall-go/module/biometric/service/internal/biz"
	"mall-go/module/biometric/service/internal/data/model/biometric"

	"github.com/go-kratos/kratos/v2/log"
)

// var _ biz.EmployersRepo = (*employersRepo)(nil)
var _ biz.BiometricRepo = (*biometricRepo)(nil)

type biometricRepo struct {
	data *Data
	Log  *log.Helper
}

func NewBiometricRepo(data *Data, logger log.Logger) biz.BiometricRepo {
	return &biometricRepo{
		data: data,
		Log:  log.NewHelper(log.With(logger, "module", "data/biometri")),
	}
}

func (r *biometricRepo) GetFingerByID(ctx context.Context, id int64) (*biz.Biometric, error) {
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
	r.Log.Infof("Biometric dengan ID %d terhapus", id)
	return nil
}
