package conlark

import (
	"fmt"
	"sync"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"
)

// 用于存储用户的状态和数据
var BuildMap = sync.Map{} // 全局 map 存储所有用户的状态
// 用于存储卡片信息
var CardMap = sync.Map{} // 全局 map 存储所有用户的状态

// 自定义测试用的构建卡片
func GetCard(job string) *larkcard.MessageCard {
	// config
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(true).
		UpdateMulti(true).
		Build()

	// header
	header := larkcard.NewMessageCardHeader().
		Template("blue").
		Title(larkcard.NewMessageCardPlainText().
			Content("jenkins构建").
			Build()).
		Build()

	// Elements
	titleElement := larkcard.NewMessageCardMarkdown().Content("<font color='green'> Jenkins构建 </font>")

	param1 := &larkcard.MessageCardPlainText{}
	param1.Content("参数化构建a")
	select1 := &larkcard.MessageCardEmbedSelectMenuStatic{}
	select1.MessageCardEmbedSelectMenuStatic(
		larkcard.NewMessageCardEmbedSelectMenuBase().
			Options(
				[]*larkcard.MessageCardEmbedSelectOption{
					larkcard.NewMessageCardEmbedSelectOption().Text(GetCardText("选项111")).Value("value111").Build(),
					larkcard.NewMessageCardEmbedSelectOption().Text(GetCardText("选项222")).Value("value222").Build(),
				}).
			Placeholder(GetCardText("branch")).Build()).
		Value(map[string]interface{}{"key_a": "value_a"}).
		Build()

	param2 := &larkcard.MessageCardPlainText{}
	param2.Content("参数话构建b")
	select2 := &larkcard.MessageCardEmbedSelectMenuStatic{}
	select2.MessageCardEmbedSelectMenuStatic(
		larkcard.NewMessageCardEmbedSelectMenuBase().
			Options(
				[]*larkcard.MessageCardEmbedSelectOption{
					larkcard.NewMessageCardEmbedSelectOption().Text(GetCardText("选项333")).Value("value333").Build(),
					larkcard.NewMessageCardEmbedSelectOption().Text(GetCardText("选项444")).Value("value444").Build(),
				}).
			Placeholder(GetCardText("branch")).Build()).
		Value(map[string]interface{}{"key_b": "value_b"}).
		Build()

	button := larkcard.NewMessageCardEmbedButton().
		Value(map[string]interface{}{
			"action":  "submit_form",
			"value_a": "${key_a}",
			"value_b": "${key_b}",
		}).Text(GetCardText("提交构建")).
		Confirm(
			larkcard.NewMessageCardActionConfirm().Text(GetCardText("是否要发起构建")).Title(GetCardText("确认")).Build(),
		)
	divElement1 := larkcard.NewMessageCardDiv().Text(param1).Extra(select1)
	divElement2 := larkcard.NewMessageCardDiv().Text(param2).Extra(select2)
	divElement3 := larkcard.NewMessageCardDiv().Extra(button)
	// 卡片消息体
	messageCard := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements([]larkcard.MessageCardElement{titleElement, divElement1, divElement2, divElement3}).
		// CardLink(cardLink).
		Build()

	return messageCard
}

func GetCardText(str string) *larkcard.MessageCardPlainText {
	t := &larkcard.MessageCardPlainText{}
	t.Content(str)
	return t
}

// 根据job信息来获取job参数，构建卡片结构返回给lark
func GetJenkinsCard(job string) (*larkcard.MessageCard, error) {

	paramList, err := GetJobConfig(job)
	if err != nil {
		return nil, err
	}
	// config
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(true).
		UpdateMulti(true).
		Build()

	// header
	header := larkcard.NewMessageCardHeader().
		Template("blue").
		Title(larkcard.NewMessageCardPlainText().
			Content("jenkins构建").
			Build()).
		Build()

	emeList := []larkcard.MessageCardElement{}
	// Elements
	titleElement := larkcard.NewMessageCardMarkdown().Content(fmt.Sprintf("<font color='green'> Jenkins构建 : %s </font>", job))
	emeList = append(emeList, titleElement)
	for _, action := range paramList.Actions {
		for _, param := range action.ParameterDefinitions {
			// 参数名称
			mcpt := &larkcard.MessageCardPlainText{}
			mcpt.Content(param.Name)

			// 遍历所有的可选参数
			selectOptionList := []*larkcard.MessageCardEmbedSelectOption{}
			if len(param.AllValueItems) > 0 {
				for _, valueItem := range param.AllValueItems {
					selectOptionList = append(
						selectOptionList,
						larkcard.NewMessageCardEmbedSelectOption().Text(GetCardText(valueItem.Name)).Value(valueItem.Value).Build(),
					)
				}
			}
			if len(param.Choices) > 0 {
				for _, choice := range param.Choices {
					selectOptionList = append(
						selectOptionList,
						larkcard.NewMessageCardEmbedSelectOption().Text(GetCardText(choice)).Value(choice).Build(),
					)
				}
			}
			selectMenu := &larkcard.MessageCardEmbedSelectMenuStatic{}
			selectMenu.MessageCardEmbedSelectMenuStatic(
				larkcard.NewMessageCardEmbedSelectMenuBase().
					Options(selectOptionList).
					Value(map[string]interface{}{"action": "addparam", "job": job, "tag": param.Name}).
					Placeholder(GetCardText(param.Description)).Build()).
				Build()
			divElement := larkcard.NewMessageCardDiv().Text(mcpt).Extra(selectMenu)
			emeList = append(emeList, divElement)
		}
	}

	button := larkcard.NewMessageCardEmbedButton().
		Value(map[string]interface{}{
			"action": "build",
			"job":    job,
			"url":    GetJenkinsUrlByJob(job),
		}).Text(GetCardText("提交构建")).
		Confirm(
			larkcard.NewMessageCardActionConfirm().
				Title(GetCardText("Jenkins")).
				Text(GetCardText(fmt.Sprintf("是否要发起构建:%s", job))).
				Build(),
		)

	divElementButton := larkcard.NewMessageCardDiv().Extra(button)
	emeList = append(emeList, divElementButton)
	// 卡片消息体
	messageCard := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements(emeList).
		// CardLink(cardLink).
		Build()

	return messageCard, nil
}

func (bd *BuildData) AddParamData(event *callback.CardActionTriggerEvent) *BuildData {
	p := &JobParam{}
	p.Name = event.Event.Action.Value["tag"].(string)
	p.Value = event.Event.Action.Option
	for k, v := range bd.Param {
		if v.Name == event.Event.Action.Value["tag"].(string) {
			bd.Param[k] = *p
			return bd
		}
	}
	bd.Param = append(bd.Param, *p)
	return bd
}

func (bd *BuildData) AddJobAndUrlData(event *callback.CardActionTriggerEvent) *BuildData {
	bd.Job.Name = event.Event.Action.Value["job"].(string)
	bd.Job.URL = event.Event.Action.Value["url"].(string)
	return bd
}
