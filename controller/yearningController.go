package controller

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mlonV/dingtalk/config"
	"github.com/mlonV/dingtalk/utils"
)

type YearningController struct {
}

func (yc YearningController) SendYearning(ctx *gin.Context) {

	// DingData := &config.Msg{Msgtype: "text", At: dingtalkAt, Text: config.Text{Content: sendAlertMsg}}
	// yearningData := &config.Msg{"msgtype": "markdown", "markdown": {"title": "Yearning sql审计平台", "text": "%s"}}
	yearningData := &config.YearningMsg{}
	ctx.ShouldBindJSON(&yearningData)
	yData, _ := json.Marshal(yearningData)
	yearning := config.Viper.Get("yearning").([]interface{})
	for _, v := range yearning {
		dingtalkConfigMap := v.(map[interface{}]interface{})
		dingtalkSecret := dingtalkConfigMap["secret"].(string)
		dingtalkUrl := dingtalkConfigMap["url"].(string)

		sendurl := utils.GetSendUrl(dingtalkUrl, dingtalkSecret)
		// 真正发送消息的地方
		body, err := utils.SendMsg(sendurl, yData)
		if err != nil {
			log.Fatal("send dingtalk err : ", err)
		}
		fmt.Println(string(body))
	}
}
