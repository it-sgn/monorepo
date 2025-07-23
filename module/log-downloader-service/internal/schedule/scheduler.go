package schedule

import (
	"log"
	"time"
)

func getInterval() time.Duration {
	now := time.Now()
	hour := now.Hour()
	if (hour == 7) || (hour == 14) || (hour == 22) {
		return 2 * time.Minute
	}
	return 5 * time.Hour
}

func Start(devices []string, downloader func(string) error) {
	go func() {
		for {
			interval := getInterval()
			log.Printf("Next interval: %v", interval)
			for _, ip := range devices {
				go func(ip string) {
					if err := downloader(ip); err != nil {
						log.Printf("Download error from %s: %v", ip, err)
					}
				}(ip)
			}
			time.Sleep(interval)
		}
	}()
	select {} // keep running
}
