package server

import (
	"context"
	"mall-go/module/log_downloader/service/internal/conf"
	"mall-go/module/log_downloader/service/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
)

type CronWorker struct {
	c    *conf.Job
	sche *cron.Cron
}

// func NewCronWorker(c *conf.Job, jobService *service.JobService, downloadService *service.DownloadJobService) (s *CronWorker) {
func NewCronWorker(c *conf.Job, downloadService *service.DownloadJobService) (s *CronWorker) {
	// jobService.Init()
	//Download Log
	// // downloadService.Init()
	// s = &CronWorker{
	// 	c:    c,
	// 	sche: cron.New(),
	// }
	// for _, j := range c.Jobs {
	// 	job, ok := service.DefaultJobs[j.Name]
	// 	if !ok {
	// 		log.Warnf("can not find job: %s", j.Name)

	// 		continue
	// 	}
	// 	s.sche.AddFunc(j.Schedule, job)
	// }
	// log.Info("Jumlah pekerjaan yang dimuat :", len(c.Jobs))
	//Download Log
	// downloadService.Init()
	// jobService.Init()
	downloadService.Init()
	s = &CronWorker{
		c:    c,
		sche: cron.New(),
	}
	for _, j := range c.Jobs {

		job, ok := service.DefaultJobs[j.Name]
		if !ok {
			log.Warnf("can not find job: %s", j.Name)

			continue
		}
		s.sche.AddFunc(j.Schedule, job)
	}
	log.Infof("Registered Jobs: %+v", service.DefaultJobs)

	for k := range service.DefaultJobs {
		log.Infof("registered job: %s", k)
	}
	log.Info("Jumlah pekerjaan yang dimuat :", len(c.Jobs))
	return s
}

func (s *CronWorker) Start(c context.Context) error {
	s.sche.Start()
	return nil
}

func (s *CronWorker) Stop(c context.Context) error {
	s.sche.Stop()
	return nil
}

func (s *CronWorker) RunSrv(name string) {
	// log.Info("run job{%s}", name)
	log.Infof("run job{%s}", name)
	//switch name {
	//case s.c.ExampleJob.JonName:
	//	s.job.DoMyWork()
	//default:
	//	s.HeartBeat()
	//}
}

// HeartBeat .
func (s *CronWorker) HeartBeat() {
	log.Info("alive...")
}
