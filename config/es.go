package config

// es告警struct
type ESconfig struct {
	ESAlarm `yaml:"esalarm"`
}

type ESAlarm struct {
	IsOpen    bool     `json:"isOpen"` //是否使用
	User      string   `json:"user"`
	Pass      string   `json:"pass"`
	Hosts     []string `json:"hosts"` //["http://192.168.103.113:9200"]
	QueryList []Query  `json:"querylist"`
}

type Query struct {
	Index          string       `json:"index"`          //"xaas*" #根据kibana前缀匹配的索引 xaax*
	IndexField     string       `json:"indexfield"`     //索引中的字段
	LogKey         []string     `json:"logKey"`         //"ERROR"  #
	TimeField      string       `json:"timefield"`      //索引时间段的字段名
	Num            int64        `json:"num"`            //#数量
	TimeRange      int64        `json:"timerange"`      // 时间范围
	SendMsgNum     int64        `json:"sendMsgNum"`     //发送消息数量
	Interval       int64        `json:"interval"`       //查询es的间隔
	RepeatInterval int64        `json:"repeatinterval"` //重复告警间隔
	DingGroup      []DingNotify `json:"dinggroup"`      //钉钉告警组
}

type DingNotify struct {
	DingURL    string `json:"dingurl"`
	Dingsecret string `json:"dingsecret"`
}
