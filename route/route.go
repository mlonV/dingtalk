package route

import (
	"github.com/gin-gonic/gin"
	"github.com/mlonV/dingtalk/controller"
	"github.com/mlonV/dingtalk/prome"
)

func RegisterRoutes() *gin.Engine {
	router := gin.Default()

	alertController := &controller.AlterController{}
	router.GET("/", alertController.GetIndex)
	router.POST("/", alertController.GetIndex)

	router.POST("/sendmsg", alertController.SendMsg)

	yearningController := &controller.YearningController{}
	router.POST("/yearning", yearningController.SendYearning)

	// reload 接口
	router.POST("/-/reload", controller.ReloadConfig)
	// help 接口
	router.GET("/help", controller.Help)

	// 注册prometheus的监控指标
	router.GET("/metrics", prome.PromeHTTPFunc())

	// 取消注册prometheus.Register
	p := &controller.Prome{}
	router.DELETE("/prome/delete/:containername", p.Unregister)
	router.DELETE("/prome/all", p.UnregisterAll)

	// 发送sentry告警
	sc := &controller.SentryController{}
	router.POST("/sentry", sc.WebHook)
	return router
}
