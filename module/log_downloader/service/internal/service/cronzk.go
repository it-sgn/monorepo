package service

import (
	"context"
	v1 "mall-go/api/log_downloader/service/v1"
	"mall-go/module/log_downloader/service/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type CronZKService struct {
	v1.UnimplementedCronJobServiceServer

	uc  *biz.CronZKUseCase
	log *log.Helper
}

func NewCronZKService(uc *biz.CronZKUseCase, logger log.Logger) *CronZKService {
	return &CronZKService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "module", "service/cronjob")),
	}
}

func (s *CronZKService) CreateCronZK(ctx context.Context, req *v1.CreateCronJobRequest) (*v1.CreateCronJobReply, error) {
	cj, err := s.uc.Create(ctx, &biz.CronZK{
		Name:    req.Name,
		Spec:    req.Spec,
		Command: req.Command,
		Enabled: req.Enabled,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateCronJobReply{
		CronJob: &v1.CronJob{
			Id:        cj.ID,
			Name:      cj.Name,
			Spec:      cj.Spec,
			Command:   cj.Command,
			Enabled:   cj.Enabled,
			CreatedAt: cj.CreatedAt.Format(time.RFC3339),
			UpdatedAt: cj.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *CronZKService) UpdateCronZK(ctx context.Context, req *v1.UpdateCronJobRequest) (*v1.UpdateCronJobReply, error) {
	cj, err := s.uc.Update(ctx, &biz.CronZK{
		ID:      req.Id,
		Name:    req.Name,
		Spec:    req.Spec,
		Command: req.Command,
		Enabled: req.Enabled,
	})
	if err != nil {
		return nil, err
	}

	lastRunAt := ""
	if cj.LastRunAt != nil {
		lastRunAt = cj.LastRunAt.Format(time.RFC3339)
	}
	nextRunAt := ""
	if cj.NextRunAt != nil {
		nextRunAt = cj.NextRunAt.Format(time.RFC3339)
	}

	return &v1.UpdateCronJobReply{
		CronJob: &v1.CronJob{
			Id:        cj.ID,
			Name:      cj.Name,
			Spec:      cj.Spec,
			Command:   cj.Command,
			Enabled:   cj.Enabled,
			LastRunAt: lastRunAt,
			NextRunAt: nextRunAt,
			CreatedAt: cj.CreatedAt.Format(time.RFC3339),
			UpdatedAt: cj.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *CronZKService) DeleteCronZK(ctx context.Context, req *v1.DeleteCronJobRequest) (*v1.DeleteCronJobReply, error) {
	err := s.uc.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteCronJobReply{}, nil
}

func (s *CronZKService) GetCronZK(ctx context.Context, req *v1.GetCronJobRequest) (*v1.GetCronJobReply, error) {
	cj, err := s.uc.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	lastRunAt := ""
	if cj.LastRunAt != nil {
		lastRunAt = cj.LastRunAt.Format(time.RFC3339)
	}
	nextRunAt := ""
	if cj.NextRunAt != nil {
		nextRunAt = cj.NextRunAt.Format(time.RFC3339)
	}

	return &v1.GetCronJobReply{
		CronJob: &v1.CronJob{
			Id:        cj.ID,
			Name:      cj.Name,
			Spec:      cj.Spec,
			Command:   cj.Command,
			Enabled:   cj.Enabled,
			LastRunAt: lastRunAt,
			NextRunAt: nextRunAt,
			CreatedAt: cj.CreatedAt.Format(time.RFC3339),
			UpdatedAt: cj.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *CronZKService) ListCronZK(ctx context.Context, req *v1.ListCronJobRequest) (*v1.ListCronJobReply, error) {
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.PageToken == 0 {
		req.PageToken = 1
	}

	cronJobs, total, err := s.uc.List(ctx, req.PageToken, req.PageSize)
	if err != nil {
		return nil, err
	}

	v1CronZKs := make([]*v1.CronJob, 0, len(cronJobs))
	for _, cj := range cronJobs {
		lastRunAt := ""
		if cj.LastRunAt != nil {
			lastRunAt = cj.LastRunAt.Format(time.RFC3339)
		}
		nextRunAt := ""
		if cj.NextRunAt != nil {
			nextRunAt = cj.NextRunAt.Format(time.RFC3339)
		}

		v1CronZKs = append(v1CronZKs, &v1.CronJob{
			Id:        cj.ID,
			Name:      cj.Name,
			Spec:      cj.Spec,
			Command:   cj.Command,
			Enabled:   cj.Enabled,
			LastRunAt: lastRunAt,
			NextRunAt: nextRunAt,
			CreatedAt: cj.CreatedAt.Format(time.RFC3339),
			UpdatedAt: cj.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &v1.ListCronJobReply{
		CronJobs:   v1CronZKs,
		TotalCount: total,
	}, nil
}
