package types

import "github.com/mlonV/tools/loger"

// 加载配置文件用
type DingtalkConfig struct {
	ESAlarm       ESAlarm        `yaml:"esalarm"`
	AlarmStatus                  //告警状态
	Yearning      []Yearning     `json:"yearning,omitempty"`
	MonitorDocker MonitorDocker  `json:"monitordocker,omitempty"`
	Alertmanager  []Alertmanager `json:"alertmanager,omitempty"`
	RedisKey      RedisKey       `json:"rediskey,omitempty"`
	LogSet        LogSet         `json:"logset,omitempty"`
	Sentry        Sentry         `json:"sentry,omitempty"`
}

// Sentry dingding告警的地址
type Sentry struct {
}

// 日志设置 struct
type LogSet struct {
	ToFile          bool           `json:"tofile,omitempty"`
	Level           loger.LogLevel `json:"level,omitempty"`
	FileName        string         `json:"filename,omitempty"`
	FilePath        string         `json:"filepath,omitempty"`
	FileSize        int64          `json:"filesize,omitempty"`
	WithFuncAndFile bool           `json:"withfunc,omitempty"`
	FileSaveNum     uint64         `json:"filesavenum,omitempty"`
}

// ES连接设置 struct
type ESAlarm struct {
	IsOpen    bool     `json:"isopen"` //是否使用
	User      string   `json:"user"`
	Pass      string   `json:"pass"`
	Hosts     []string `json:"hosts"` //["http://192.168.103.113:9200"]
	QueryList []Query  `json:"querylist"`
}

// 叮叮告警url和secret
type DingNotify struct {
	DingURL    string `json:"dingurl"`
	Dingsecret string `json:"dingsecret"`
}

// 告警的发送和恢复通知处理
type AlarmStatus struct {
	AlarmCount int64
	StartTime  string
	EndTime    string
	IsAlarm    bool // 当前是否在告警
}

// yearning消息struct
type YearningMsg struct {
	Msgtype  string                `json:"msgtype,omitempty"`
	At       `json:"at,omitempty"` // 为了实现发送消息at相关人员准备的
	Markdown `json:"markdown,omitempty"`
}

// 发送叮叮的markdown
type Markdown struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}

// 发送钉钉根据手机号@组员
type At struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}

type Yearning struct {
	Url       string   `json:"url,omitempty"`
	Secret    string   `json:"secret,omitempty"`
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}

// 钉钉消息text
type Text struct {
	Content string `json:"content,omitempty"`
}

// 发送钉钉消息格式 Msg  josn格式模板
//
//	message = """{
//	    "at": {
//	        "atMobiles":[
//	            "15555555555"
//	        ],
//	        "isAtAll": false
//	    },
//	    "text": {
//	        "content":"test"
//	    },
//	    "msgtype":"text"
//	}"""
//
// 发送钉钉消息格式 Text类型
type Msg struct {
	Msgtype  string      `json:"msgtype,omitempty"`
	At       `json:"at"` // 为了实现发送消息at相关人员准备的
	Text     `json:"text"`
	Markdown `json:"markdown"`
}

// alertmanager 官方接口发送数据的类型
//
//	{
//		"version": "4",
//		"groupKey": <string>,              // key identifying the group of alerts (e.g. to deduplicate)
//		"truncatedAlerts": <int>,          // how many alerts have been truncated due to "max_alerts"
//		"status": "<resolved|firing>",
//		"receiver": <string>,
//		"groupLabels": <object>,
//		"commonLabels": <object>,
//		"commonAnnotations": <object>,
//		"externalURL": <string>,           // backlink to the Alertmanager.
//		"alerts": [
//		  {
//			"status": "<resolved|firing>",
//			"labels": <object>,
//			"annotations": <object>,
//			"startsAt": "<rfc3339>",
//			"endsAt": "<rfc3339>",
//			"generatorURL": <string>,      // identifies the entity that caused the alert
//			"fingerprint": <string>        // fingerprint to identify the alert
//		  },
//		  ...
//		]
//	  }
type AlertmanagerMsg struct {
	Version         string            `json:"version,omitempty"`
	GroupKey        string            `json:"groupKey,omitempty"`
	TruncatedAlerts int               `json:"truncatedAlerts,omitempty"`
	Status          string            `json:"status,omitempty"`
	Receiver        string            `json:"receiver,omitempty"`
	GroupLabels     map[string]string `json:"groupLabels,omitempty"`
	Alerts          []AlertMsg        `json:"alerts,omitempty"`
}

type AlertMsg struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Status      string            `json:"status,omitempty"`
	StartsAt    string            `json:"startsAt,omitempty"`
	EndsAt      string            `json:"endsAt,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

type DingtalkUrlInfo struct {
	Url     string `json:"url,omitempty"`
	Secret  string `json:"secret,omitempty"`
	Jobname string `json:"jobname,omitempty"`
}

type Alertmanager struct {
	Job       string   `json:"job,omitempty"`
	Url       string   `json:"url,omitempty"`
	Secret    string   `json:"secret,omitempty"`
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
	IsDefault bool     `json:"isDefault,omitempty"`
}
