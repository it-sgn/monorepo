package job111

import (
	"context"
	"fmt"
	"mall-go/module/log_downloader/service/internal/biz"
	"time"
)

func Test(s *ExampleJob) {
	s.uc.CreateGreeter(context.Background(), &biz.Greeter{})
	fmt.Printf("waktu saat ini %v \n", time.Now().Unix())
}
