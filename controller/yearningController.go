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

	// yearningData := &config.Msg{"msgtype": "markdown", "markdown": {"title": "Yearning sql审计平台", "text": "%s"}}
	// 从yearning的webhook Post过来的数据
	yearningData := &config.YearningMsg{}
	ctx.ShouldBindJSON(&yearningData)
	fmt.Println(yearningData.Text)
	// 读取配置文件的yearning配置
	yearningFileConf := &config.YearningFileConf{}
	err := config.Viper.Unmarshal(&yearningFileConf)
	test := config.Viper.Get("yearning")
	if err != nil {
		fmt.Println("ViperUnmarsh err : ", err)
	}
	fmt.Println(yearningData, "1", yearningFileConf.Yearning, test)
	text := yearningData.Text
	// 遍历发送所有群组
	for _, yconf := range yearningFileConf.Yearning {
		yearningData.Text = text
		yearningData.At.AtMobiles = yconf.AtMobiles
		yearningData.At.IsAtAll = yconf.IsAtAll

		// 如果有@全体，则不单独@个人
		if !yearningData.At.IsAtAll {
			for _, v := range yearningData.At.AtMobiles {
				yearningData.Text += fmt.Sprintf("@%s ", v)
			}
		}

		// 获取Post 到钉钉接口的url
		sendurl := utils.GetSendUrl(yconf.Url, yconf.Secret)
		yData, _ := json.Marshal(yearningData)
		// 真正发送消息的地方
		body, err := utils.SendMsg(sendurl, yData)
		if err != nil {
			log.Fatal("send dingtalk err : ", err)
		}
		fmt.Println(string(body))
	}
}
