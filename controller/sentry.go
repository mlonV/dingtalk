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

type SentryController struct {
}

// [Sentry] event.title
// Type: event.type
// Project: project
// URL: url

func (sc *SentryController) WebHookForText(ctx *gin.Context) {
	var sentrystuct types.SentryAlert
	ctx.ShouldBindJSON(&sentrystuct)
	config.Log.Debug("%#v", &sentrystuct)
	// 筛选需要的数据
	var msgList []string
	msgList = append(msgList,
		[]string{fmt.Sprintf("[Sentry] %s ", sentrystuct.Event.Title),
			fmt.Sprintf("type : %s", sentrystuct.Event.Type),
			fmt.Sprintf("Project : %s", sentrystuct.Project),
			fmt.Sprintf("Url : %s", sentrystuct.URL)}...,
	)
	DingData := &types.Msg{}
	DingData.Msgtype = "text"
	DingData.Text.Content = strings.Join(msgList, "\n")
	data, _ := json.Marshal(DingData)

	// 获取Url后面参数（钉钉的access_token和secret  因为发送URL的a
	access_token := ctx.Query("access_token")
	DingURL := "https://oapi.dingtalk.com/robot/send?access_token=" + access_token
	secret := ctx.Query("secret")
	sendUrl := utils.GetSendUrl(DingURL, secret)

	body, err := utils.SendMsg(sendUrl, data)
	if err != nil {
		config.Log.Error("Sentry SendMsg Error: %s", err)
	}
	config.Log.Debug("发送的dingUrl: %s \n发送的数据: %s\n返回的body: %s\n", sendUrl, string(data), string(body))

	// fmt.Println(data)
	ctx.String(http.StatusOK, "sentry ok")
}

// Markdown
func (sc *SentryController) WebHookForMarkdown(ctx *gin.Context) {
	var sentrystuct types.SentryAlert
	ctx.ShouldBindJSON(&sentrystuct)
	config.Log.Debug("%#v", &sentrystuct)
	// 筛选需要的数据
	var text []string
	text = append(text,
		// 	markdown 两个空格加上回车换行
		[]string{fmt.Sprintf("### [Sentry] %s ", sentrystuct.Event.Title),
			fmt.Sprintf("type : %s  ", sentrystuct.Event.Type),
			fmt.Sprintf("Project : %s  ", sentrystuct.Project),
			fmt.Sprintf("Url : [详情请点击](%s)  ", sentrystuct.URL)}...,
	)
	DingData := &types.Msg{}
	DingData.Msgtype = "markdown"
	DingData.Markdown.Title = "Sentry Markdown Notify"
	DingData.Markdown.Text = strings.Join(text, "\n")
	data, _ := json.Marshal(DingData)

	// 获取Url后面参数（钉钉的access_token和secret  因为发送URL的a
	access_token := ctx.Query("access_token")
	DingURL := "https://oapi.dingtalk.com/robot/send?access_token=" + access_token
	secret := ctx.Query("secret")
	sendUrl := utils.GetSendUrl(DingURL, secret)

	body, err := utils.SendMsg(sendUrl, data)
	if err != nil {
		config.Log.Error("Sentry SendMsg Error: %s", err)
	}
	config.Log.Debug("发送的dingUrl: %s \n发送的数据: %s\n返回的body: %s\n", sendUrl, string(data), string(body))

	// fmt.Println(data)
	ctx.String(http.StatusOK, "sentry ok")
}
