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
