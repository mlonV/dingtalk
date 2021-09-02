package route

import (
	"github.com/gin-gonic/gin"
	"github.com/mlonV/dingtalk/controller"
)

func RegisterRoutes() *gin.Engine {
	router := gin.Default()

	alertController := &controller.AlterController{}
	router.GET("/", alertController.GetIndex)
	router.POST("/", alertController.GetIndex)

	router.POST("/sendmsg", alertController.SendMsg)

	yearningController := &controller.YearningController{}
	router.POST("/yearning", yearningController.SendYearning)

	return router
}
