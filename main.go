package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"imaotai/config"
	"imaotai/db"
	"imaotai/models"
	"imaotai/reqfunc"
	"imaotai/service"
	"imaotai/task"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var configpath string
	var dbpath string
	var address string
	var codephone string
	var phone string
	var code string

	flag.StringVar(&configpath, "path", "config.yml", "配置文件路径")
	flag.StringVar(&dbpath, "db", "imaotai.db", "数据库路径")

	flag.StringVar(&address, "address", "", "查询的地址")

	flag.StringVar(&codephone, "codephone", "", "要登录的手机号")

	flag.StringVar(&phone, "phone", "", "登录手机号")
	flag.StringVar(&code, "code", "", "登录验证码")

	// 从arguments中解析注册的flag。必须在所有flag都注册好而未访问其值时执行。未注册却使用flag -help时，会返回ErrHelp。
	flag.Parse()

	// 获取地址的经纬度
	if len(address) > 0 {
		rt, err := reqfunc.GetLocationByAddress(address)
		if err != nil {
			panic(err)
		}
		data, err := json.Marshal(rt)
		if err != nil {
			panic(err)
		}
		println(string(data))
		os.Exit(0)
	}

	// 获取手机验证码
	if len(codephone) > 0 {
		rt, err := reqfunc.SendCode(codephone)
		if err != nil {
			panic(err)
		}
		if rt {
			println("登录验证码发送成功")
		} else {
			println("登录验证码发送失败")
		}
		os.Exit(0)
	}

	// 登录
	if (len(phone) == 0 && len(code) > 0) || (len(phone) > 0 && len(code) == 0) {
		panic("登录时，phone和code不能为空")
	}
	if len(phone) > 0 && len(code) > 0 {
		rt, err := reqfunc.Login(phone, code)
		if err != nil {
			panic(err)
		}
		data, err := json.Marshal(rt)
		if err != nil {
			panic(err)
		}
		println(string(data))
		os.Exit(0)
	}

	// 加载系统，启动
	err := config.GetConfig(configpath)
	if err != nil {
		panic(err)
	}

	// 数据库
	err = db.Init(dbpath)
	if err != nil {
		panic(err)
	}
	models.Init()

	// 每次启动 重新初始化数据
	err = service.RefreshData(config.Configs)
	if err != nil {
		panic(err)
	}

	// 启动定时任务
	crontask := task.Init()
	crontask.AddTask()
	crontask.Task.Start()
	defer crontask.Task.Stop()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	// 创建一个通道来接收信号
	signalChan := make(chan os.Signal, 1)

	// 监听中断信号和终止信号
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动一个 goroutine 在后台等待信号
	go func() {
		// 阻塞等待信号
		sig := <-signalChan
		fmt.Printf("接收到信号: %v\n", sig)

		// 执行清理和关闭操作

		// 退出程序
		os.Exit(0)
	}()

	fmt.Println("程序正在运行...")

	// 阻塞主 goroutine，使程序持续运行
	select {}
}
