package data

import (

	// Digunakan untuk serialisasi/deserialisasi JSON
	// Digunakan untuk mengatur waktu kedaluwarsa cache

	"context"
	"fmt"
	"mall-go/module/log_downloader/service/internal/biz"
	"mall-go/module/log_downloader/service/internal/data/model"
	"mall-go/module/log_downloader/service/internal/data/model/cronzk"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	// Pastikan versi Redis client yang benar
)

// Pastikan antarmuka DepartmentRepo diimpor dengan benar dari biz package
// var _ biz.DepartmentRepo = (*departmentRepo)(nil)
var _ biz.CronZKRepo = (*cronZKRepo)(nil)

type cronZKRepo struct {
	data *Data
	log  *log.Helper
}

// NewCronZKRepo .
func NewCronZKRepo(data *Data, logger log.Logger) biz.CronZKRepo {
	return &cronZKRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/cronjob")),
	}
}

func (r *cronZKRepo) CreateCronZK(ctx context.Context, cj *biz.CronZK) (*biz.CronZK, error) {
	res, err := r.data.db.CronZK.Create().
		SetName(cj.Name).
		SetSpec(cj.Spec).
		SetCommand(cj.Command).
		SetEnabled(cj.Enabled).
		Save(ctx)
	if err != nil {
		if model.IsConstraintError(err) {
			return nil, errors.BadRequest("CRONJOB_EXISTS", fmt.Sprintf("cron job with name %s already exists", cj.Name))
		}
		return nil, errors.InternalServer("CREATE_CRONJOB_FAILED", fmt.Sprintf("failed to create cron job: %v", err))
	}
	return &biz.CronZK{
		ID:        res.ID,
		Name:      res.Name,
		Spec:      res.Spec,
		Command:   res.Command,
		Enabled:   res.Enabled,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}

func (r *cronZKRepo) GetCronZK(ctx context.Context, id int64) (*biz.CronZK, error) {
	res, err := r.data.db.CronZK.Get(ctx, id)
	if err != nil {
		if model.IsNotFound(err) {
			return nil, errors.NotFound("CRONJOB_NOT_FOUND", fmt.Sprintf("cron job with id %d not found", id))
		}
		return nil, errors.InternalServer("GET_CRONJOB_FAILED", fmt.Sprintf("failed to get cron job: %v", err))
	}

	return &biz.CronZK{
		ID:        res.ID,
		Name:      res.Name,
		Spec:      res.Spec,
		Command:   res.Command,
		Enabled:   res.Enabled,
		LastRunAt: res.LastRunAt,
		NextRunAt: res.NextRunAt,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}

func (r *cronZKRepo) UpdateCronZK(ctx context.Context, cj *biz.CronZK) (*biz.CronZK, error) {
	updater := r.data.db.CronZK.UpdateOneID(cj.ID)
	if cj.Name != "" {
		updater.SetName(cj.Name)
	}
	if cj.Spec != "" {
		updater.SetSpec(cj.Spec)
	}
	if cj.Command != "" {
		updater.SetCommand(cj.Command)
	}
	updater.SetEnabled(cj.Enabled) // Always update enabled status

	if cj.LastRunAt != nil {
		updater.SetLastRunAt(*cj.LastRunAt)
	}
	if cj.NextRunAt != nil {
		updater.SetNextRunAt(*cj.NextRunAt)
	}

	res, err := updater.Save(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return nil, errors.NotFound("CRONJOB_NOT_FOUND", fmt.Sprintf("cron job with id %d not found", cj.ID))
		}
		if model.IsConstraintError(err) {
			return nil, errors.BadRequest("CRONJOB_EXISTS", fmt.Sprintf("cron job with name %s already exists", cj.Name))
		}
		return nil, errors.InternalServer("UPDATE_CRONJOB_FAILED", fmt.Sprintf("failed to update cron job: %v", err))
	}

	return &biz.CronZK{
		ID:        res.ID,
		Name:      res.Name,
		Spec:      res.Spec,
		Command:   res.Command,
		Enabled:   res.Enabled,
		LastRunAt: res.LastRunAt,
		NextRunAt: res.NextRunAt,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}

func (r *cronZKRepo) DeleteCronZK(ctx context.Context, id int64) error {
	err := r.data.db.CronZK.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return errors.NotFound("CRONJOB_NOT_FOUND", fmt.Sprintf("cron job with id %d not found", id))
		}
		return errors.InternalServer("DELETE_CRONJOB_FAILED", fmt.Sprintf("failed to delete cron job: %v", err))
	}
	return nil
}

func (r *cronZKRepo) ListCronZK(ctx context.Context, page, pageSize int64) ([]*biz.CronZK, int64, error) {
	offset := (page - 1) * pageSize
	total, err := r.data.db.CronZK.Query().Count(ctx)
	if err != nil {
		return nil, 0, errors.InternalServer("LIST_CRONJOB_FAILED", fmt.Sprintf("failed to count cron jobs: %v", err))
	}

	res, err := r.data.db.CronZK.Query().
		Offset(int(offset)).
		Limit(int(pageSize)).
		All(ctx)
	if err != nil {
		return nil, 0, errors.InternalServer("LIST_CRONJOB_FAILED", fmt.Sprintf("failed to list cron jobs: %v", err))
	}

	cronJobs := make([]*biz.CronZK, 0, len(res))
	for _, cj := range res {
		cronJobs = append(cronJobs, &biz.CronZK{
			ID:        cj.ID,
			Name:      cj.Name,
			Spec:      cj.Spec,
			Command:   cj.Command,
			Enabled:   cj.Enabled,
			LastRunAt: cj.LastRunAt,
			NextRunAt: cj.NextRunAt,
			CreatedAt: cj.CreatedAt,
			UpdatedAt: cj.UpdatedAt,
		})
	}
	return cronJobs, int64(total), nil
}

func (r *cronZKRepo) ListEnabledCronZKs(ctx context.Context) ([]*biz.CronZK, error) {
	res, err := r.data.db.CronZK.Query().
		Where(cronzk.Enabled(true)).
		All(ctx)
	if err != nil {
		return nil, errors.InternalServer("LIST_ENABLED_CRONJOB_FAILED", fmt.Sprintf("failed to list enabled cron jobs: %v", err))
	}

	cronJobs := make([]*biz.CronZK, 0, len(res))
	for _, cj := range res {
		cronJobs = append(cronJobs, &biz.CronZK{
			ID:        cj.ID,
			Name:      cj.Name,
			Spec:      cj.Spec,
			Command:   cj.Command,
			Enabled:   cj.Enabled,
			LastRunAt: cj.LastRunAt,
			NextRunAt: cj.NextRunAt,
			CreatedAt: cj.CreatedAt,
			UpdatedAt: cj.UpdatedAt,
		})
	}
	return cronJobs, nil
}
