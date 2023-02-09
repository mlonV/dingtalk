package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mlonV/dingtalk/monitor/docker"
	"github.com/mlonV/dingtalk/prome"
)

type Prome struct {
}

// 取消单个注册prometheus
func (p Prome) Unregister(ctx *gin.Context) {

	containerName, ok := ctx.Params.Get("containername")
	if !ok {
		ctx.String(http.StatusOK, "请使用正确的Url,且输入正确的容器名")
		return
	}
	dmap := docker.GetdMap()
	value, ok := dmap.Load(containerName)
	if !ok {
		ctx.String(http.StatusOK, "容器名不存在,请输入正确的容器名")
		return
	}
	ps, _ := (value).(docker.ProcessStatus)
	dmap.Delete(containerName)
	Gauge := prome.PromeRegister.Unregister(ps.PromeGauge)
	PIDGauge := prome.PromeRegister.Unregister(ps.PromePIDGauge)
	result := fmt.Sprintf("进程存活注销: %t ,进程PID注销: %t ", Gauge, PIDGauge)
	ctx.String(http.StatusOK, result)
}

// 取消全部注册prometheus
func (p Prome) UnregisterAll(ctx *gin.Context) {
	keyList := []string{}
	dmap := docker.GetdMap()
	dmap.Range(func(containerName, value interface{}) bool {
		ps, _ := (value).(docker.ProcessStatus)
		dmap.Delete(containerName)
		prome.PromeRegister.Unregister(ps.PromeGauge)
		prome.PromeRegister.Unregister(ps.PromePIDGauge)
		dmap.Delete(containerName)
		keyList = append(keyList, containerName.(string))
		return true
	})

	result := fmt.Sprintf("已删除KeyList: %v", keyList)
	ctx.String(http.StatusOK, result)
}
