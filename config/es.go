package config

// es告警struct
type ESconfig struct {
	ESAlarm `yaml:"esalarm"`
}

type ESAlarm struct {
	User       string   `json:"user"`
	Pass       string   `json:"pass"`
	Hosts      []string `json:"hosts"`  //["http://192.168.103.113:9200"]
	Index      string   `json:"index"`  //"xaas*" #根据kibana前缀匹配的索引 xaax*
	LogKey     string   `json:"logKey"` //"ERROR"  #
	Num        int64    `json:"num"`    //#数量
	Time       int64    `json:"time"`
	IndexField string   `json:"indexfield"`
	SendMsgNum int64    `json:"sendMsgNum"`
	IsOpen     bool     `json:"isOpen"` //是否使用
}
