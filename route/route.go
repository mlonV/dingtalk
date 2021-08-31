package route

import (
	"github.com/gin-gonic/gin"
	"github.com/mlonV/dingtalk/controller"
)

func RegisterRoutes() *gin.Engine {
	router := gin.Default()

	alertController := &controller.AlterController{}
	router.GET("/", alertController.GetIndex)
	router.POST("/sendmsg", alertController.SendMsg)
	return router
}
