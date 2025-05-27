package route

import (
	"github.com/gin-gonic/gin"
	"github.com/mlonV/dingtalk/controller"
	"github.com/mlonV/dingtalk/controller/search"
	"github.com/mlonV/dingtalk/controller/super"
	"github.com/mlonV/dingtalk/prome"
	"github.com/mlonV/dingtalk/utils"
)

func RegisterRoutes() *gin.Engine {

	router := gin.Default()

	alertController := &controller.AlterController{}

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

	// DBSearch
	router.GET("/", search.Index)
	router.GET("/search/schema/:schema", search.SearchSchemaHandler)
	router.GET("/search/table/:table", search.SearchTableHandler)
	router.GET("/search/column/:column", search.SearchColumnHandler)

	// supervisor
	router.Use(utils.Cors())

	g := router.Group("api", utils.AuthMiddleware())
	{
		// denglu
		g.POST("/login", super.Login)
		g.GET("/info", super.UserInfo)
		g.POST("/logout", super.LogOut)

		g.GET("/supervisor/list", super.ListHost)
		g.POST("/supervisor/addhost", super.AddHost)
		g.DELETE("/supervisor/host/:id", super.DelHost)
		g.PATCH("/supervisor/host", super.UpdateHost)

		// g.GET("/supervisor/status", super.GetSupervisorStatus)
		// g.GET("/supervisor/process/all", super.GetetAllProcessInfo)  // 暂时不用了这个
		g.GET("/supervisor/process/all", super.GetAllProcessInfo_LazyLoad)
		g.GET("/supervisor/process/host/:host", super.GetProcessInfo)

		// control process
		g.POST("/supervisor/process/start/host/:host/name/:name", super.StartProcess)
		g.POST("/supervisor/process/stop/host/:host/name/:name", super.StopProcess)
		// log
		g.GET("/supervisor/process/log/host/:host/name/:name", super.TailStdoutLog)
	}

	return router
}
