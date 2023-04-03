package domain

import (
	"github.com/gerdong/auto-github-hosts/config"
	"time"
)

// 定时任务：按照config.C.Updater.Interval时间间隔，执行domain.UpdateHosts()
func Cron() {
	interval := config.UpdateInterval
	if interval <= 0 {
		interval = 12
	}

	if config.ImmediatelyAtStartup {
		doWork()
	}

	ticker := time.NewTicker(time.Duration(time.Duration(interval)*time.Hour) * time.Second)
	for {
		select {
		case <-ticker.C:
			doWork()
		}
	}
}

func doWork() {
	result := UpdateHosts()
	if len(result) > 0 {

	}
}
