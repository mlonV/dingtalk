package main

import (
	// "net/http"
	// _ "net/http/pprof"

	"github.com/mlonV/dingtalk/config"
	_ "github.com/mlonV/dingtalk/monitor"
	"github.com/mlonV/dingtalk/route"
)

// "github.com/mlonV/dingtalk/config"

func main() {

	r := route.RegisterRoutes()
	// go func() {
	// 	http.ListenAndServe("0.0.0.0:9090", nil)
	// }()

	// 启动gin监控告警
	if err := r.Run(":" + config.Port); err != nil {
		panic(err.Error())
	}
}
