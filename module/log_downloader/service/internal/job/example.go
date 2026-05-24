package job111

import (
	"context"
	"fmt"
	"mall-go/module/log_downloader/service/internal/biz"
	"time"
)

type ExampleJob struct {
	uc *biz.GreeterUsecase
}

func NewExampleJob(uc *biz.GreeterUsecase) *ExampleJob {
	job := &ExampleJob{
		uc: uc,
	}
	return job
}

func (s *ExampleJob) Init() {
	DefaultJobs = map[string]JobFunc{
		"one": s.DoMyWork,
		"two": s.DoOtherWork,
	}
}

func (s *ExampleJob) DoMyWork() {
	s.uc.CreateGreeter(context.Background(), &biz.Greeter{})
	fmt.Printf("waktu saat ini JOB 1 %v \n", time.Now().Unix())
}

func (s *ExampleJob) DoOtherWork() {
	fmt.Printf("Waktu saat ini dengan JOB 2 %v \n", time.Now().Unix())
}
