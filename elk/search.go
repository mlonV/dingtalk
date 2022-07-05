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

type EsSource struct {
	Message string `json:"message"`
}

func init() {
	// 解析配置文件到 结构体变量
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

func count() {
	// res, err := esClient.Search("phplog_6").Do(context.Background())

	timenow := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	timebefore := time.Now().UTC().Add(time.Minute * -time.Duration(esalarm.Time)).Format("2006-01-02T15:04:05.000Z")
	fmt.Println(timenow, timebefore)
	bq := elastic.NewBoolQuery()
	bq.Must(
		elastic.NewMatchQuery(esalarm.IndexField, esalarm.LogKey),
		elastic.NewRangeQuery(esalarm.TimeField).Gt(timebefore).Lt(timenow).Format("strict_date_optional_time||epoch_millis").TimeZone("+08:00"),
	)
	// bq.Must(elastic.NewRangeQuery("timestamp").Gt(timebefore).Lt(timenow).Format("strict_date_optional_time||epoch_millis"))

	// 默认只查询两条
	res, err := esClient.Search(esalarm.Index).Size(2).Query(bq).Do(context.Background())
	// res, err := esClient.Search("phplog_deflector").Query(bq).Do(context.Background())
	if err != nil {
		panic(err)
	}

	// 超过数量就告警
	if res.Hits.TotalHits.Value >= esalarm.Num {
		alarm(res)
	}
	fmt.Printf("%v", res.Hits.TotalHits)

	var ess EsSource
	for _, item := range res.Each(reflect.TypeOf(ess)) {
		t := item.(EsSource)
		fmt.Printf("%#v \n", t)
	}
}

// 告警
func alarm(res *elastic.SearchResult) {
	msgList := []string{
		fmt.Sprintf("告警名 : ELK log %s", esalarm.LogKey),
		fmt.Sprintf("告警内容: %s 超过告警阈值%d条", esalarm.LogKey, esalarm.Num),
		fmt.Sprintf("%d分钟内%s数量: %d", esalarm.Time, esalarm.LogKey, res.Hits.TotalHits.Value),
		fmt.Sprintf("信息如下[%d]条 : ", esalarm.SendMsgNum),
	}

	// 发送告警
	AlertmanagerFileConf := &config.AlertmanagerFileConf{}
	config.Viper.Unmarshal(&AlertmanagerFileConf)
	for _, confList := range AlertmanagerFileConf.Alertmanager {
		confSecret := confList.Secret
		confUrl := confList.Url

		sendurl := utils.GetSendUrl(confUrl, confSecret)
		sendAlertMsg := strings.Join(msgList, "\n")

		// 在添加两条message索引里的消息
		// 遍历每个值
		var ess EsSource
		for k, item := range res.Each(reflect.TypeOf(ess)) {
			if k < int(esalarm.SendMsgNum) {
				sendAlertMsg += fmt.Sprintf("\n%s", item.(EsSource))
			}
			// t := item.(EsSource)
			// fmt.Printf("%#v \n", t)
		}

		// 如果有@全体，则不单独@个人(需要吧@+手机号写入到发送的信息里，有at全体的话则不生效)
		if !confList.IsAtAll {
			for _, v := range confList.AtMobiles {
				sendAlertMsg += fmt.Sprintf(", @%s ", v)
			}
		}

		// DingData := &config.Msg{Msgtype: "text", At: config.At{AtMobiles: confList.AtMobiles}, Text: config.Text{Content: sendAlertMsg}}

		// 发送钉钉数据的结构
		DingData := &config.Msg{}
		DingData.Msgtype = "text"
		DingData.At.AtMobiles = confList.AtMobiles
		DingData.At.IsAtAll = confList.IsAtAll
		DingData.Text.Content = sendAlertMsg

		data, _ := json.Marshal(DingData)
		// 真正发送消息的地方
		body, err := utils.SendMsg(sendurl, data)
		// 静默5分钟
		time.Sleep(time.Duration(esalarm.RepeatInterval) * time.Second)

		if err != nil {
			log.Fatal("send dingtalk err : ", err)
		}
		fmt.Println(string(body))
	}
}

// 定时器
func Ticker() {
	// interval := time.Duration(esalarm.Time) * time.Minute
	if esalarm.IsOpen {
		interval := time.Second * time.Duration(esalarm.Interval)
		ticker = time.NewTicker(interval)

		// 定时任务处理逻辑
		for {
			// 调用Reset方法对timer对象进行定时器重置
			// 	ticker.Reset(interval)
			<-ticker.C
			count()
		}
	}
	// ticker.Stop()
}

// func main() {
// 	result, err := esClient.ElasticsearchVersion("http://192.168.103.113:9200")
// 	eslog.Println(result, err)
// 	Timer()
// }
