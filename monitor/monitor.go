package monitor

import (
	"github.com/mlonV/dingtalk/config"
	der "github.com/mlonV/dingtalk/monitor/docker"
	"github.com/mlonV/dingtalk/monitor/elk"
	"github.com/mlonV/dingtalk/monitor/rediskey"
)

// 运行所有的监控进程
func init() {

	//es 的日志告警
	if config.Conf.ESAlarm.IsOpen {
		config.Log.Info("开启elk 功能 , isopen: %t", config.Conf.ESAlarm.IsOpen)
		go elk.Ticker()
	}

	// 启动监控docker 内进程的监控
	if config.Conf.MonitorDocker.IsOpen {
		config.Log.Info("开启docker Monitor 功能 , isopen: %t", config.Conf.MonitorDocker.IsOpen)
		go der.Ticker()
	}

	// 启动监控redis Key(list llen) 内进程的监控
	if config.Conf.RedisKey.IsOpen {
		config.Log.Info("开启Redis QueueKey 功能 , isopen: %t", config.Conf.RedisKey.IsOpen)
		go rediskey.Ticker()
	}
}
