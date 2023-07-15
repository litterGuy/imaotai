package task

import (
	"github.com/robfig/cron/v3"
	"imaotai/config"
	"imaotai/service"
	"log"
	"os"
	"time"
)

func StartTask() {
	// 启动定时任务
	c := cron.New(cron.WithSeconds(), cron.WithChain(
		cron.SkipIfStillRunning(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))),
		cron.Recover(cron.DefaultLogger),
	))
	// 每日7点刷新一次数据
	_, _ = c.AddFunc("0 0 7 * * ?", func() {
		log.Printf("start task, and time = %d\n", time.Now().Unix())
		err := service.RefreshData(config.Configs)
		log.Println(err.Error())
	})
	// 每天9点20预约
	_, _ = c.AddFunc("0 20 9 * * ?", func() {
		log.Printf("start task, and time = %d\n", time.Now().Unix())
		err := service.Reservation(config.Configs)
		log.Println(err.Error())
	})
	c.Start()
	defer c.Stop()
}
