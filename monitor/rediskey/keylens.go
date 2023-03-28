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
	RedisKeyList []string // keys queue:*的结果
}

type Metrics struct {
	RedisCli           *redis.Client
	RedisRegexpKeyName string // queue:*
	RedisKeyName       string
	RedisGauge         prometheus.Gauge
}

var (
	keymap     map[string]*Metrics            = make(map[string]*Metrics)
	metricsMap map[string]map[string]*Metrics = make(map[string]map[string]*Metrics, 200)
	redisCli   *redis.Client
)

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
func Handler() {
	for _, mmap := range metricsMap {
		for _, metrics := range mmap {
			ic := metrics.RedisCli.LLen(context.Background(), metrics.RedisKeyName)
			if ic.Err() != nil {
				config.Log.Error("Get RedisKey Llen Err: %s ", ic.Err())
				// 如果查询报错，直接设置0
				metrics.RedisGauge.Set(float64(0))
				continue
			}
			metrics.RedisGauge.Set(float64(ic.Val()))
			config.Log.Debug("metrics.RedisGauge.Set Key %s: %d ", metrics.RedisKeyName, ic.Val())
		}
	}
}

func worker() {
	config.Log.Info("Start Redis QueueKeyLens ,Interval: %ds ", config.Conf.RedisKey.Interval)
	for _, regexpKey := range config.Conf.RedisKey.Keys {
		ssc := redisCli.Keys(context.Background(), regexpKey)
		s, err := ssc.Result()
		if err != nil {
			config.Log.Error("Get RedisKey by regexpKey Err: %s ", err)
		}
		for _, key := range s {
			// 如果是空的话就新设置一个 map[string]*Metrics
			if metricsMap[regexpKey][key] == nil {
				config.Log.Debug("metricsMap[%s][%s] : %v \n", regexpKey, key, metricsMap[regexpKey][key])
				metrics := &Metrics{
					RedisCli:           redisCli,
					RedisRegexpKeyName: regexpKey,
					RedisKeyName:       key,
					RedisGauge:         SetRedisMetricsGauge(key),
				}
				// 注册到prometheus
				prome.PromeRegister.Register(metrics.RedisGauge)
				keymap[key] = metrics
				metricsMap[regexpKey] = keymap
			}
		}
	}
	config.Log.Debug("metricsMap ALL : %#v \n", metricsMap)

	// 处理map所有的数据
	Handler()
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
