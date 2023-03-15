package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mlonV/dingtalk/config"
	"github.com/mlonV/dingtalk/types"
	"github.com/mlonV/dingtalk/utils"
)

type AlterController struct {
}

func (a AlterController) GetIndex(ctx *gin.Context) {
	ctx.String(http.StatusOK, "瞎访问个什么玩意儿呢？\n 嗯？")

}

func (a AlterController) SendMsg(ctx *gin.Context) {

	alertmanagerMsg := &types.AlertmanagerMsg{}
	ctx.ShouldBindJSON(alertmanagerMsg)
	config.Log.Info("绑定altermanager结构体的内容: %#v", alertmanagerMsg)
	for _, alertMsg := range alertmanagerMsg.Alerts {
		var msgList []string
		// 定义发送消息的内容
		msgList = append(msgList,
			[]string{fmt.Sprintf("告警名 : %s", alertMsg.Labels["alertname"]),
				fmt.Sprintf("主机名 : %s", alertMsg.Labels["host"]),
				fmt.Sprintf("实例名 : %s", alertMsg.Labels["instance"]),
				fmt.Sprintf("job名 : %s", alertMsg.Labels["job"]),
				fmt.Sprintf("状态 : %s", alertMsg.Status),
				fmt.Sprintf("开始时间 : %s", utils.GetParseTime(alertMsg.StartsAt)),
				fmt.Sprintf("结束时间 : %s", utils.GetParseTime(alertMsg.EndsAt)),
				fmt.Sprintf("概览 : %s", alertMsg.Annotations["summary"]),
				fmt.Sprintf("详情 : %s", alertMsg.Annotations["description"])}...,
		)

		for _, confList := range config.Conf.Alertmanager {
			// confJob := confList.Job
			confSecret := confList.Secret
			confUrl := confList.Url

			sendurl := utils.GetSendUrl(confUrl, confSecret)
			// 根据prometheus里面设置的job名和配置文件里面的job名来分别发送不同群聊
			// if confJob == alertMsg.Labels["job"] {
			sendAlertMsg := strings.Join(msgList, "\n")
			// 如果有@全体，则不单独@个人(需要吧@+手机号写入到发送的信息里，有at全体的话则不生效)
			if !confList.IsAtAll {
				for _, v := range confList.AtMobiles {
					sendAlertMsg += fmt.Sprintf(", @%s ", v)
				}
			}
			// DingData := &types.Msg{Msgtype: "text", At: types.At{AtMobiles: confList.AtMobiles}, Text: types.Text{Content: sendAlertMsg}}

			// 发送钉钉数据的结构
			DingData := &types.Msg{}
			DingData.Msgtype = "text"
			DingData.At.AtMobiles = confList.AtMobiles
			DingData.At.IsAtAll = confList.IsAtAll
			DingData.Text.Content = sendAlertMsg

			data, _ := json.Marshal(DingData)
			// 真正发送消息的地方
			body, err := utils.SendMsg(sendurl, data)
			config.Log.Info("发送的数据URL: %s ,发送的数据: %s", sendurl, data)

			if err != nil {
				config.Log.Fatal("send dingtalk err : ", err)
			}
			config.Log.Info(string(body))
			// }
		}

	}

	ctx.String(http.StatusOK, "sendmsg ok")
}
