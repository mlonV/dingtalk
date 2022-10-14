package main

import (
	"log"

	"github.com/mlonV/dingtalk/config"
	der "github.com/mlonV/dingtalk/monitor/docker"
	"github.com/mlonV/dingtalk/monitor/elk"
	"github.com/mlonV/dingtalk/route"
)

// "github.com/mlonV/dingtalk/config"

func main() {

	r := route.RegisterRoutes()

	//es 的日志告警
	if config.Conf.ESAlarm.IsOpen {
		log.Println("开启elk 功能 ，isopen: ", config.Conf.ESAlarm.IsOpen)
		go elk.Ticker()
	}

	// 启动监控docker 内进程的监控
	if config.Conf.MonitorDocker.IsOpen {
		log.Println("开启docker Monitor 功能 ，isopen: ", config.Conf.ESAlarm.IsOpen)
		go der.Ticker()
	}

	// 启动gin监控告警
	if err := r.Run(":" + config.Port); err != nil {
		panic(err.Error())
	}

}
