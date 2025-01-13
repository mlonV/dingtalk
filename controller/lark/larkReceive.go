package conlark

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	// larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

var APP_ID = "cli_a6739d811c78d010"
var APP_SECRET = "1CLkLeCJEaFH6SYeOoPgec5EmHSOkFBd"
var VerificationToken = "3ZRPgFqRNwJslEkhArpNecN8jTt6B4fx"

type Content struct {
	Text string `json:"text"`
}

// type cardChannel map[type]type

func LarkEventHandler() *dispatcher.EventDispatcher {
	// 注册消息处理器
	// 这里替换成自己应用后台里的verificationToken\eventEncryptKey，没有配置就填空字符串
	handler := dispatcher.NewEventDispatcher(VerificationToken, "")
	handler.OnP2MessageReceiveV1(larkOnP2MessageReceiveV1)
	handler.OnP2MessageReadV1(larkOnP2MessageReadV1)
	handler.OnP2CardActionTrigger(larkOnP2CardActionTrigger)
	return handler
}

// 收到卡片消息并且 更新卡片
func larkOnP2CardActionTrigger(ctx context.Context, event *callback.CardActionTriggerEvent) (*callback.CardActionTriggerResponse, error) {
	// 处理消息 event，这里简单打印消息的内容

	fmt.Println(" larkOnP2CardActionTrigger:  ", larkcore.Prettify(event))
	fmt.Println(" larkOnP2CardActionTrigger:  ", event.RequestId())

	if event.Event.Action.Value["action"] == "addparam" {
		job, _ := event.Event.Action.Value["job"].(string)

		bd, _ := BuildMap.LoadOrStore(job, &BuildData{})

		bd.(*BuildData).AddParamData(event)
		fmt.Println(bd)
	}
	if event.Event.Action.Value["action"] == "build" {
		job, _ := event.Event.Action.Value["job"].(string)

		bd, _ := BuildMap.LoadOrStore(job, &BuildData{})

		bd.(*BuildData).AddJobAndUrlData(event)
		fmt.Println(bd)
		// err := TriggerJenkinsBuild(bd.(*BuildData))
		// fmt.Println("err := TriggerJenkinsBuild(bd.(*BuildData))", err)
		// catr := &callback.CardActionTriggerResponse{}
		// catr.Card = &callback.Card{Type: "raw", Data: GetCard("test-jenkins")}
		// toast := &callback.Toast{Content: "testttt", Type: "info"}
		// catr.Toast = toast

		req := larkim.NewGetMessageReqBuilder().MessageId(event.Event.Context.OpenMessageID).UserIdType("open_id").Build()
		client := lark.NewClient(APP_ID, APP_SECRET)
		resp, err := client.Im.Message.Get(context.Background(), req)
		if err != nil {
			fmt.Println("client.Im.Message.Get(context.Background(), req) : ", err)
		}
		fmt.Println(larkcore.Prettify(resp.Data.Items))
		return nil, nil
	}

	return nil, nil
}

func larkOnP2MessageReceiveV1(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	// 处理消息 event，这里简单打印消息的内容
	fmt.Println(" larkOnP2MessageReceiveV1:  ", larkcore.Prettify(event))
	fmt.Println(" larkOnP2MessageReceiveV1:  ", event.RequestId())
	if *event.Event.Message.ChatType == "group" && *event.Event.Message.MessageType == "text" {
		if *event.Event.Message.Mentions[0].Name == "黑熊精" {
			content := &Content{}
			if err := json.Unmarshal([]byte(*event.Event.Message.Content), &content); err != nil {
				return err
			}
			msg := removeAtUsernames(content.Text)
			parts := strings.Split(msg, " ")
			switch parts[0] {
			case "jenkins":
				SendJenkinsList(event, parts[1])
			case "build":
				if len(parts) < 2 {
					msg := string("无效输入，以下示例 \n    请输入: \"jenkins <job name>\" 来获取过滤后的job任务名和jenkins链接\n    请输入: \"build <job name>\"来请求触发构建任务")
					SendMessage(event, msg)
					return fmt.Errorf("无效输入: %s", parts[0])
				}
				if !JobIsExists(parts[1]) {
					msg := "job任务不存在: " + parts[1]
					SendMessage(event, msg)
					return fmt.Errorf("job任务不存在: %s", parts[1])
				}
				SendCard(event, parts[1])
			default:
				// 收到无效指令，返回Help
				msg := string("无效输入，以下示例 \n    请输入: \"jenkins <job name>\" 来获取过滤后的job任务名和jenkins链接\n    请输入: \"build <job name>\"来请求触发构建任务")
				SendMessage(event, msg)
			}
		}
	}
	return nil
}

func larkOnP2MessageReadV1(ctx context.Context, event *larkim.P2MessageReadV1) error {
	// 处理消息 event，这里简单打印消息的内容
	fmt.Println(" larkOnP2MessageReadV1:  ", larkcore.Prettify(event))
	fmt.Println(" larkOnP2MessageReadV1:  ", event.RequestId())
	return nil
}
