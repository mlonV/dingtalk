package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mlonV/dingtalk/config"
)

func ReloadConfig(ctx *gin.Context) {
	if err := config.LaodConfig(); err != nil {
		ctx.String(http.StatusOK, " config reload err", err.Error())
	}
	ctx.String(http.StatusOK, " config reload ok")
}

// 返回help帮助信息
func Help(ctx *gin.Context) {
	var result string = `	
	ALL API Usage:

	GET	    "/help"     "帮助信息"
	GET     "/"         None
	GET	    "/metrics"  "获取prometheus收集的监控信息,自定义监控"

	POST    "/"         None
	POST    "/sendmsg"  "Alertmanager用于发送告警信息"
	POST    "/-/reload" "不可用🚫"

	DELETE  "/prome/delete/[containername]" "取消注册到prometheus的单个指标[中是容器名]"
	DELETE  "/prome/all"                    "取消注册所有的容器进程监控"
	
	POST    "/sentry"         "sentry 发送webhook"`
	ctx.String(http.StatusOK, result)
}
