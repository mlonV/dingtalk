package elk

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
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

// Query ERROR
func query(ql config.Query, esC *elastic.Client) {
	fmt.Printf("%#v \n", ql)
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

	// 默认只查询两条
	if ql.SendMsgNum == 0 {
		ql.SendMsgNum = 2
	}
	res, err := esC.Search(ql.Index).Size(int(ql.SendMsgNum)).Query(bq).Do(context.Background())
	if err != nil {
		panic(err)
	}
	// 超过数量就告警
	if res.Hits.TotalHits.Value >= ql.Num {
		alarm(res, ql)
	}
	fmt.Printf("时间范围内: [%v,%v] 查询内容: %s ,查询条数%d \n",
		timebefore,
		timenow,
		ql.LogKey,
		res.Hits.TotalHits.Value,
	)

	// var ess EsSource
	// for _, item := range res.Each(reflect.TypeOf(ess)) {
	// 	t := item.(EsSource)
	// 	fmt.Printf("\n --  --  -- \n%#v \n\n\n", t)
	// }
}

// 告警
func alarm(res *elastic.SearchResult, ql config.Query) {
	msgList := []string{
		fmt.Sprintf("告警名 : ELK log %s", ql.LogKey),
		fmt.Sprintf("告警内容: %s 超过告警阈值%d条", ql.LogKey, ql.Num),
		fmt.Sprintf("%d分钟内%s数量: %d", ql.TimeRange, ql.LogKey, res.Hits.TotalHits.Value),
		fmt.Sprintf("信息如下[%d]条 : ", ql.SendMsgNum),
	}

	// 发送告警
	for _, dinggroup := range ql.DingGroup {
		dingurl := dinggroup.DingURL
		dingsecret := dinggroup.Dingsecret

		sendurl := utils.GetSendUrl(dingurl, dingsecret)
		sendAlertMsg := strings.Join(msgList, "\n")

		// 在添加两条message索引里的消息
		// 遍历每个值
		var ess EsSource
		for k, item := range res.Each(reflect.TypeOf(ess)) {
			if k < int(ql.SendMsgNum) {
				sendAlertMsg += fmt.Sprintf("\n%s", item.(EsSource))
			}
			// t := item.(EsSource)
			// fmt.Printf("%#v \n", t)
		}

		// DingData := &config.Msg{Msgtype: "text", At: config.At{AtMobiles: confList.AtMobiles}, Text: config.Text{Content: sendAlertMsg}}

		// 发送钉钉数据的结构
		DingData := &config.Msg{}
		DingData.Msgtype = "text"
		DingData.Text.Content = sendAlertMsg

		data, _ := json.Marshal(DingData)
		// 真正发送消息的地方
		body, err := utils.SendMsg(sendurl, data)
		if err != nil {
			log.Fatal("send dingtalk err : ", err)
		}
		fmt.Println(string(body))
	}
	// 静默5分钟
	time.Sleep(time.Duration(ql.RepeatInterval) * time.Second)
}

// 定时器
func Ticker() {
	// interval := time.Duration(esalarm.Time) * time.Minute
	for _, ql := range esalarm.QueryList {

		go func(ql config.Query) {
			if esalarm.IsOpen {
				interval := time.Second * time.Duration(ql.Interval)
				ticker = time.NewTicker(interval)

				// 定时任务处理逻辑
				for {
					// 调用Reset方法对timer对象进行定时器重置
					// 	ticker.Reset(interval)
					<-ticker.C
					query(ql, esClient)
				}
			}
		}(ql)
	}

	// ticker.Stop()
}

// func main() {
// 	result, err := esClient.ElasticsearchVersion("http://192.168.103.113:9200")
// 	eslog.Println(result, err)
// 	Ticker()
// }
