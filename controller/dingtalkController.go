package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mlonV/dingtalk/config"
	"github.com/mlonV/dingtalk/utils"
)

type AlterController struct {
}

func (a AlterController) GetIndex(ctx *gin.Context) {
	ctx.String(http.StatusOK, "瞎访问个什么玩意儿呢？\n 嗯？")

}

func (a AlterController) SendMsg(ctx *gin.Context) {

	alertmanagerMsg := &config.AlertmanagerMsg{}
	ctx.ShouldBindJSON(alertmanagerMsg)
	// fmt.Println("Alerts  :  ", alertmanagerMsg.Alerts)
	// fmt.Println("GroupKey  :  ", alertmanagerMsg.GroupKey)
	// fmt.Println("GroupLabels  :  ", alertmanagerMsg.GroupLabels)
	// fmt.Println("Receiver  :  ", alertmanagerMsg.Receiver)
	// fmt.Println("Status  :  ", alertmanagerMsg.Status)
	// fmt.Println("TruncatedAlerts  :  ", alertmanagerMsg.TruncatedAlerts)
	// fmt.Println("Version  :  ", alertmanagerMsg.Version)
	// fmt.Println("Annotations  :  ", alertmanagerMsg.Alerts[0].Annotations)
	// fmt.Println("---------------------------------------------------------")
	for _, alert := range alertmanagerMsg.Alerts {
		var msgList []string
		// 定义发送消息的内容
		msgList = append(msgList,
			[]string{fmt.Sprintf("告警名 : %s", alert.Labels["alertname"]),
				fmt.Sprintf("主机名 : %s", alert.Labels["host"]),
				fmt.Sprintf("实例名 : %s", alert.Labels["instance"]),
				fmt.Sprintf("job名 : %s", alert.Labels["job"]),
				fmt.Sprintf("状态 : %s", alert.Status),
				fmt.Sprintf("开始时间 : %s", utils.GetParseTime(alert.StartsAt)),
				fmt.Sprintf("结束时间 : %s", utils.GetParseTime(alert.EndsAt)),
				fmt.Sprintf("概览 : %s", alert.Annotations["summary"]),
				fmt.Sprintf("详情 : %s", alert.Annotations["description"])}...,
		)
		// 拼接字符串 换行
		sendAlertMsg := strings.Join(msgList, "\n")

		// 发送at相关人设置,还未实现
		// dingtalkAt := make(map[string]interface{})
		// dingtalkAt["atMobiles"] = []string{"15515787648", "13262685201"}
		// dingtalkAt["isAtAll"] = false
		// DingData := &config.Msg{Msgtype: "text", At: dingtalkAt, Text: config.Text{Content: sendAlertMsg}}

		DingData := &config.Msg{Msgtype: "text", Text: config.Text{Content: sendAlertMsg}}
		// 遍历配置文件匹配 job发送给对应的群组
		// fmt.Println(sendAlertMsg, alert.Labels["job"])
		dingtalk := config.Viper.Get("dingtalk").([]interface{})
		for _, v := range dingtalk {
			dingtalkConfigMap := v.(map[interface{}]interface{})
			dingtalkJob := dingtalkConfigMap["job"].(string)
			dingtalkSecret := dingtalkConfigMap["secret"].(string)
			dingtalkUrl := dingtalkConfigMap["url"].(string)

			if dingtalkJob == alert.Labels["job"] {
				sendurl := utils.GetSendUrl(dingtalkUrl, dingtalkSecret)
				// 真正发送消息的地方
				data, _ := json.Marshal(DingData)
				body, err := utils.SendMsg(sendurl, data)
				if err != nil {
					log.Fatal("send dingtalk err : ", err)
				}
				fmt.Println(string(body))

			}
		}

	}

	ctx.String(http.StatusOK, "sendmsg ok")
}
