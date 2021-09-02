package config

type Text struct {
	Content string `json:"content,omitempty"`
}

// 发送钉钉消息格式 Msg  josn格式模板
// message = """{
//     "at": {
//         "atMobiles":[
//             "15555555555"
//         ],
//         "isAtAll": false
//     },
//     "text": {
//         "content":"test"
//     },
//     "msgtype":"text"
// }"""
// 发送钉钉消息格式 Msg
type Msg struct {
	Msgtype string                 `json:"msgtype,omitempty"`
	At      map[string]interface{} `json:"at,omitempty"`
	Text    `json:"text,omitempty"`
}

// alertmanager 官方接口发送数据的类型
// {
// 	"version": "4",
// 	"groupKey": <string>,              // key identifying the group of alerts (e.g. to deduplicate)
// 	"truncatedAlerts": <int>,          // how many alerts have been truncated due to "max_alerts"
// 	"status": "<resolved|firing>",
// 	"receiver": <string>,
// 	"groupLabels": <object>,
// 	"commonLabels": <object>,
// 	"commonAnnotations": <object>,
// 	"externalURL": <string>,           // backlink to the Alertmanager.
// 	"alerts": [
// 	  {
// 		"status": "<resolved|firing>",
// 		"labels": <object>,
// 		"annotations": <object>,
// 		"startsAt": "<rfc3339>",
// 		"endsAt": "<rfc3339>",
// 		"generatorURL": <string>,      // identifies the entity that caused the alert
// 		"fingerprint": <string>        // fingerprint to identify the alert
// 	  },
// 	  ...
// 	]
//   }
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

type DingtalkConfig struct {
	Url     string `json:"url,omitempty"`
	Secret  string `json:"secret,omitempty"`
	Jobname string `json:"jobname,omitempty"`
}
