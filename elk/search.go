package elk

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mlonV/dingtalk/config"
	"github.com/mlonV/dingtalk/utils"
	"github.com/olivere/elastic/v7"
)

var esClient *elastic.Client
var eslog *log.Logger
var err error
var esalarm *config.ESconfig
var ticker *time.Ticker

// 写死的message
type EsSource struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func init() {
	// 解析配置文件到 结构体变量
	fmt.Printf("%#v", esalarm)
	config.Viper.Unmarshal(&esalarm)

	eslog = log.New(os.Stdout, "[ElasticLog] ", log.LstdFlags)
	esClient, err = elastic.NewClient(
		elastic.SetErrorLog(eslog),
		elastic.SetURL(esalarm.Hosts...),
		elastic.SetBasicAuth(esalarm.User, esalarm.Pass),
	)
	if err != nil {
		panic(err)
	}
	// 查询每个版本和连接信息
	for _, host := range esalarm.Hosts {

		info, code, err := esClient.Ping(host).Do(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Printf("es return code %d ,clusterName : %s ,es Name : %s \n", code, info.ClusterName, info.Name)
		res, err := esClient.ElasticsearchVersion(host)
		if err != nil {
			panic(err)
		}
		fmt.Printf("esversion is : %s \n", res)
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

	// 查询数量超过设置num值，则触发告警
	if as.IsAlarm && as.StartTime != "" {
		as.AlarmCount = as.AlarmCount + 1
		trigger(q, res, as)
		// 重复告警的间隔
		time.Sleep(time.Second * time.Duration(q.RepeatInterval))
	}

	// 查询数量没超过设置num值，但是之前有设置过startTime了，则恢复告警
	if !as.IsAlarm && as.StartTime != "" {
		as.EndTime = time.Now().Format("2006-01-02 15:04:05.00")
		if q.IsResolved {
			resolved(q, res, as)
		}
		// 重置startTime
		as.StartTime = ""
		as.AlarmCount = 0
	}
}

// Query 查询
func query(ql config.Query, esC *elastic.Client) (*elastic.SearchResult, error) {
	fmt.Printf("%#v \n", ql)

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
		fmt.Sprintf("开始时间 : %s", as.StartTime),
		fmt.Sprintf("告警内容: %s 超过告警阈值%d条", ql.LogKey, ql.Num),
		fmt.Sprintf("%d分钟内%s数量: %d", ql.TimeRange, ql.LogKey, res.Hits.TotalHits.Value),
		fmt.Sprintf("信息如下[%d]条 : ", ql.SendMsgNum),
	}
	// 遍历钉钉组发送
	for _, dg := range ql.DingGroup {
		alarm(dg.DingURL, dg.Dingsecret, strings.Join(msgList, "\n"))
	}
}

// 恢复通知
func resolved(ql config.Query, res *elastic.SearchResult, as *config.AlarmStatus) {
	// 需要发送的消息
	msgList := []string{
		fmt.Sprintf("告警名 : ELK log %s", strings.Join(ql.LogKey, " ")),
		fmt.Sprintf("告警状态 : %s", string("resolved")),
		fmt.Sprintf("开始时间 : %s", as.StartTime),
		fmt.Sprintf("结束时间 : %s", as.EndTime),
		fmt.Sprintf("告警内容: %s 超过告警阈值%d条", ql.LogKey, ql.Num),
		fmt.Sprintf("%d分钟内%s数量: %d", ql.TimeRange, ql.LogKey, res.Hits.TotalHits.Value),
		fmt.Sprintln("信息 : 告警已恢复"),
	}
	// 遍历钉钉组发送
	for _, dg := range ql.DingGroup {
		alarm(dg.DingURL, dg.Dingsecret, strings.Join(msgList, "\n"))
	}
}

func alarm(url, secret, msg string) {

	sendurl := utils.GetSendUrl(url, secret)
	// 发送钉钉数据的结构
	DingData := &config.Msg{}
	DingData.Msgtype = "text"
	DingData.Text.Content = msg

	data, _ := json.Marshal(DingData)
	// 真正发送消息的地方
	fmt.Println("开始告警！！！！！！： ", sendurl, string(data))
	// body, err := utils.SendMsg(sendurl, data)
	// if err != nil {
	// 	log.Fatal("send dingtalk err : ", err)
	// }
	// fmt.Println(string(body))
}

// 定时器
func Ticker() {
	// interval := time.Duration(esalarm.Time) * time.Minute
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

// func main() {
// 	result, err := esClient.ElasticsearchVersion("http://192.168.103.113:9200")
// 	eslog.Println(result, err)
// 	Ticker()
// }
