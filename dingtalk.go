package main

import (
	"github.com/mlonV/dingtalk/config"
	"github.com/mlonV/dingtalk/elk"
	"github.com/mlonV/dingtalk/route"
)

// "github.com/mlonV/dingtalk/config"

func main() {

	r := route.RegisterRoutes()

	//es 的日志告警
	elk.Ticker()

	r.Run(":" + config.Port)

}
