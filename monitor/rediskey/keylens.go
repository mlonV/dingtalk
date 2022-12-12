package rediskey

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mlonV/dingtalk/config"
	"github.com/mlonV/dingtalk/prome"
	"github.com/prometheus/client_golang/prometheus"
)

type RedisMetrics struct {
	RedisCli           *redis.Client
	MetricsList        []Metrics
	RedisRegexpKeyName string   // queue:*
	RedisKeyList       []string // keys queue:*的结果
}

type Metrics struct {
	RedisKeyName string
	RedisGauge   prometheus.Gauge
}

var redisCli *redis.Client

func init() {
	redisCli = redis.NewClient(&redis.Options{
		Addr:     config.Conf.RedisKey.Hostname + ":" + fmt.Sprintf("%d", config.Conf.RedisKey.Port),
		Username: config.Conf.RedisKey.Username,
		Password: config.Conf.RedisKey.Password,
		PoolSize: 100,
	})
}

func NewRedisMetrics() *RedisMetrics {
	return &RedisMetrics{}
}

func SetRedisMetricsGauge(RedisKeyName string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "victory",
		Subsystem: "Redis",
		Name:      "QueueLens",
		Help:      "redis queue 剩余未消费的数量",
		ConstLabels: map[string]string{
			"app":           "app-victory",
			"rediskey_name": RedisKeyName,
		},
	})
}

// 处理每一个根据Keys queues* 查出来的队列
func Handler(rm *RedisMetrics) error {
	for _, key := range rm.RedisKeyList {

		// 遍历每一个key 查询长度，设置Gauge
		ic := rm.RedisCli.LLen(context.Background(), key)
		if ic.Err() != nil {
			config.Log.Error("Get RedisKey Llen Err: %s ", ic.Err())
			continue
		}
		metrics := Metrics{
			RedisKeyName: key,
			RedisGauge:   SetRedisMetricsGauge(key),
		}
		metrics.RedisGauge.Set(float64(ic.Val()))
		prome.PromeRegister.Register(metrics.RedisGauge)
	}
	return nil
}

func worker() {
	for _, regexpKey := range config.Conf.RedisKey.Keys {
		config.Log.Info("Start Redis QueueKeyLens ,Interval: %ds ", config.Conf.RedisKey.Interval)
		rm := NewRedisMetrics()
		rm.RedisCli = redisCli
		rm.RedisRegexpKeyName = regexpKey
		ssc := rm.RedisCli.Keys(context.Background(), rm.RedisRegexpKeyName)
		s, err := ssc.Result()
		if err != nil {
			config.Log.Error("Get RedisKey by regexpKey Err: %s ", err)
		}
		rm.RedisKeyList = s
		config.Log.Debug("Get Key By %s : %s", rm.RedisRegexpKeyName, s)
		if err := Handler(rm); err != nil {
			config.Log.Error("Handler RedisKey err : ", err)
		}

	}
}

// 定时器
func Ticker() {

	// interval := time.Duration(esalarm.Time) * time.Minute
	interval := time.Second * time.Duration(config.Conf.RedisKey.Interval)
	ticker := time.NewTicker(interval)
	config.Log.Info("Start Redis QueueKeyLens ,Interval: %ds ", config.Conf.RedisKey.Interval)
	for {
		// 调用Reset方法对timer对象进行定时器重置
		// 	ticker.Reset(interval)
		<-ticker.C
		worker()
	}
	// ticker.Stop()
}
