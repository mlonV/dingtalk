package route

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mlonV/dingtalk/controller"
	"github.com/mlonV/dingtalk/controller/super"
	"github.com/mlonV/dingtalk/prome"
	// "github.com/mlonV/dingtalk/utils"
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
	router.POST("/sentry/text", sc.WebHookForText)
	router.POST("/sentry/markdown", sc.WebHookForMarkdown)

	// supervisor

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:9528"}, // vue-element-admin 默认端口
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	g := router.Group("api")
	g.POST("/login", super.Login)

	// 暂时关闭验证
	// g.Use(utils.AuthMiddleware())

	g.GET("/info", super.UserInfo)
	g.POST("/logout", super.LogOut)

	g.GET("/supervisor/list", super.ListHost)
	g.POST("/supervisor/addhost", super.AddHost)
	g.DELETE("/supervisor/host/:id", super.DelHost)
	g.PATCH("/supervisor/host", super.UpdateHost)

	// g.GET("/supervisor/status", super.GetSupervisorStatus)
	g.GET("/supervisor/process/all", super.GetetAllProcessInfo)

	// control process
	g.POST("/supervisor/process/start/host/:host/name/:name", super.StartProcess)
	g.POST("/supervisor/process/stop/host/:host/name/:name", super.StopProcess)
	// log
	g.GET("/supervisor/process/log/host/:host/name/:name", super.TailStdoutLog)
	return router
}
