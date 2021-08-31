package controller

import (
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
	ctx.String(http.StatusOK, "return ")

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
			[]string{fmt.Sprintf("alertname : %s", alert.Labels["alertname"]),
				fmt.Sprintf("instance : %s", alert.Labels["instance"]),
				fmt.Sprintf("job : %s", alert.Labels["job"]),
				fmt.Sprintf("status : %s", alert.Status),
				fmt.Sprintf("StartsAt : %s", utils.GetParseTime(alert.StartsAt)),
				fmt.Sprintf("EndsAt : %s", utils.GetParseTime(alert.EndsAt)),
				fmt.Sprintf("summary : %s", alert.Annotations["summary"]),
				fmt.Sprintf("description : %s", alert.Annotations["description"])}...,
		)

		// 拼接字符串 换行
		sendAlertMsg := strings.Join(msgList, "\n")
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
				body, err := utils.PostMsg(sendurl, DingData)
				if err != nil {
					log.Fatal("send dingtalk err : ", err)
				}
				fmt.Println(string(body))

			}
		}

	}

	ctx.String(http.StatusOK, "return ")
}
