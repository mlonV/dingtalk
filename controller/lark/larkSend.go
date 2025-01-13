package conlark

import (
	"context"
	"encoding/json"
	"fmt"

	"regexp"
	"strings"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func SendJenkinsList(event *larkim.P2MessageReceiveV1, job string) error {

	joblist := GetJenkinsJobsFilter(job)
	content := map[string]string{
		"text": joblist,
	}
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		return err
	}
	client := lark.NewClient(APP_ID, APP_SECRET)
	// 创建请求对象
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`chat_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(*event.Event.Message.ChatId).
			MsgType(`text`).
			Content(string(jsonBytes)).
			Uuid(event.EventReq.RequestId()).
			Build()).Build()
	// 发起请求
	resp, err := client.Im.Message.Create(context.Background(), req)
	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}

	return nil
}

// 需要可解析的json  content {text: " "}
func SendMessage(event *larkim.P2MessageReceiveV1, msg string) error {
	content := map[string]string{
		"text": msg,
	}
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		return err
	}
	client := lark.NewClient(APP_ID, APP_SECRET)
	// 创建请求对象
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`chat_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(*event.Event.Message.ChatId).
			MsgType(`text`).
			Content(string(jsonBytes)).
			Uuid(event.EventReq.RequestId()).
			Build()).Build()
	// 发起请求
	resp, err := client.Im.Message.Create(context.Background(), req)
	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}
	return nil
}

func SendCard(event *larkim.P2MessageReceiveV1, jobname string) error {
	// 发送card消息
	content, err := GetJenkinsCard(jobname)
	if err != nil {
		return err
	}

	CardMap.Store(jobname, content)

	jsonBytes, err := json.Marshal(content)
	if err != nil {
		return err
	}
	client := lark.NewClient(APP_ID, APP_SECRET)
	// 创建请求对象
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`chat_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(*event.Event.Message.ChatId).
			MsgType(`interactive`).
			Content(string(jsonBytes)).
			Uuid(event.EventReq.RequestId()).
			Build()).Build()
	// 发起请求
	resp, err := client.Im.Message.Create(context.Background(), req)
	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}

	return nil
}

// 去掉@user的消息
func removeAtUsernames(input string) string {
	// 正则表达式匹配模式：以 @ 开头，后面是字母、数字、下划线的连续字符串
	re := regexp.MustCompile(`@\w+`)
	// 使用 ReplaceAllString 将匹配到的 @username 替换为空字符串
	result := re.ReplaceAllString(input, "")
	// 去除前后空格
	return strings.TrimSpace(result)

}
