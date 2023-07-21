package task

import (
	"github.com/robfig/cron/v3"
	"imaotai/config"
	"imaotai/msg"
	"imaotai/service"
	"log"
	"os"
	"time"
)

type CronTask struct {
	Task *cron.Cron
}

func Init() *CronTask {
	// 启动定时任务
	c := cron.New(cron.WithSeconds(), cron.WithChain(
		cron.SkipIfStillRunning(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))),
		cron.Recover(cron.DefaultLogger),
	))
	return &CronTask{Task: c}
}

func (c *CronTask) AddTask() {

	// 每日7点刷新一次数据
	_, _ = c.Task.AddFunc("0 10 7,8 * * ?", func() {
		log.Printf("start task, and time = %d\n", time.Now().Unix())
		err := service.RefreshData(config.Configs)
		if err != nil {
			log.Println(err.Error())
			msg.SendPushPlus(err.Error())
			return
		}
		msg.SendPushPlus("刷新列表完成")
	})
	// 每天9点20预约
	_, _ = c.Task.AddFunc("0 20 9 * * ?", func() {
		log.Printf("start task, and time = %d\n", time.Now().Unix())
		rt, err := service.Reservation(config.Configs)
		if err != nil {
			log.Println(err.Error())
			msg.SendPushPlus(err.Error())
		} else {
			msg.SendPushPlus(rt)
		}
	})
}
