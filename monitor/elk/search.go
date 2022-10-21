package elk

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/mlonV/dingtalk/config"
	"github.com/mlonV/dingtalk/utils"
	"github.com/mlonV/tools/loger"
	"github.com/olivere/elastic/v7"
)

var (
	esClient *elastic.Client
	eslog    *loger.Loger
	err      error
	ticker   *time.Ticker
	esalarm  *config.ESAlarm
)

// 写死的message
type EsSource struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func init() {
	elkInit()
}

func elkInit() {
	// 解析配置文件到 结构体变量
	esalarm = &config.Conf.ESAlarm
	eslog = config.Log
	esClient, err = elastic.NewClient(
		elastic.SetErrorLog(eslog),
		elastic.SetURL(esalarm.Hosts...),
		elastic.SetBasicAuth(esalarm.User, esalarm.Pass),
	)
	if err != nil {
		config.Log.Fatal(err.Error())
	}
	// 查询每个版本和连接信息
	for _, host := range esalarm.Hosts {

		info, code, err := esClient.Ping(host).Do(context.Background())
		if err != nil {
			panic(err)
		}
		eslog.Info("es return code %d ,clusterName : %s ,es Name : %s", code, info.ClusterName, info.Name)
		res, err := esClient.ElasticsearchVersion(host)
		if err != nil {
			panic(err)
		}
		eslog.Info("esversion is : %s", res)
	}
}

func NewAlarmStatus() config.AlarmStatus {
	return config.AlarmStatus{
		AlarmCount: 0,
		IsAlarm:    false,
	}
}

// 负责处理查询，告警恢复等
func worker(q config.Query, as *config.AlarmStatus) {
	// 1.先查询
	if q.SendMsgNum == 0 { // 默认查询两条
		q.SendMsgNum = 2
	}
	res, err := query(q, esClient)
	if err != nil {
		panic(err)
	}
	// 2.判断是否需要告警
	// 超过数量就告警
	as.IsAlarm = false
	if res.Hits.TotalHits.Value >= q.Num {
		as.IsAlarm = true
		if as.StartTime == "" {
			as.StartTime = time.Now().Format("2006-01-02 15:04:05.00")
		}
	}
	eslog.Info("Index : %s,LogKey : %s , Hit : %d, ", q.Index, q.LogKey, res.Hits.TotalHits.Value)

	// 查询数量超过设置num值，则触发告警
	if as.IsAlarm && as.StartTime != "" {
		as.AlarmCount = as.AlarmCount + 1
		trigger(q, res, as)
		// 重复告警的间隔
		time.Sleep(time.Second * time.Duration(q.RepeatInterval))
	}
}

// Query 查询
func query(ql config.Query, esC *elastic.Client) (*elastic.SearchResult, error) {

	// 时间范围默认一分钟
	if ql.TimeRange == 0 {
		ql.TimeRange = 1
	}
	timenow := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	timebefore := time.Now().UTC().Add(time.Minute * -time.Duration(ql.TimeRange)).Format("2006-01-02T15:04:05.000Z")
	bq := elastic.NewBoolQuery()

	// matchQuery := elastic.NewMatchQuery(ql.IndexField, ql.LogKey)
	rangeQuery := elastic.NewRangeQuery(ql.TimeField).Gt(timebefore).Lt(timenow).Format("strict_date_optional_time||epoch_millis").TimeZone("+08:00")

	mustBQ := elastic.NewBoolQuery()

	for _, v := range ql.LogKey {

		mustBQ.Must(elastic.NewMatchPhraseQuery(ql.IndexField, v))
	}
	// termsQuery := elastic.NewTermQuery(ql.IndexField, ql.LogKey)
	bq.Must(rangeQuery, mustBQ)

	return esC.Search(ql.Index).Size(int(ql.SendMsgNum)).Query(bq).Do(context.Background())
}

// 告警发送
func trigger(ql config.Query, res *elastic.SearchResult, as *config.AlarmStatus) {

	// 需要发送的消息
	msgList := []string{
		fmt.Sprintf("告警名 : ELK log %s", strings.Join(ql.LogKey, " ")),
		fmt.Sprintf("告警状态 : %s", string("firing")),
		fmt.Sprintf("开始时间 : %s", as.StartTime),
		fmt.Sprintf("告警内容: %s 超过告警阈值%d条", ql.LogKey, ql.Num),
		fmt.Sprintf("%d分钟内%s数量: %d", ql.TimeRange, ql.LogKey, res.Hits.TotalHits.Value),
		fmt.Sprintf("信息如下[%d]条 : ", ql.SendMsgNum),
	}
	sendAlertMsg := strings.Join(msgList, "\n")
	var ess EsSource
	for k, item := range res.Each(reflect.TypeOf(ess)) {
		if k < int(ql.SendMsgNum) {
			sendAlertMsg += fmt.Sprintf("\n%s", item.(EsSource))
		}
		// t := item.(EsSource)
	}

	// 遍历钉钉组发送
	for _, dg := range ql.DingGroup {
		alarm(dg.DingURL, dg.Dingsecret, sendAlertMsg)
	}
}

func alarm(url, secret, msg string) {

	sendurl := utils.GetSendUrl(url, secret)
	// 发送钉钉数据的结构
	DingData := &config.Msg{}
	DingData.Msgtype = "text"
	DingData.Text.Content = msg

	data, _ := json.Marshal(DingData)
	body, err := utils.SendMsg(sendurl, data)
	if err != nil {
		config.Log.Error("send dingtalk err : ", err)
	}
	config.Log.Warning("发送的告警信息为: ", string(data))
	config.Log.Warning("发送告警接口返回: ", string(body))

}

// 定时器
func Ticker() {
	for _, q := range esalarm.QueryList {
		go func(q config.Query) {
			if esalarm.IsOpen {
				interval := time.Second * time.Duration(q.Interval)
				ticker = time.NewTicker(interval)

				// 定时任务处理逻辑
				//运行每个goroutine的query都先初始化下告警状态
				alarmStatus := NewAlarmStatus()
				for {
					// 调用Reset方法对timer对象进行定时器重置
					// 	ticker.Reset(interval)
					<-ticker.C
					worker(q, &alarmStatus)
				}
			}
		}(q)
	}

	// ticker.Stop()
}
