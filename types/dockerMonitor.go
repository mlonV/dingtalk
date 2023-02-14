package types

// docker 相关
type MonitorDocker struct {
	IsOpen    bool     `json:"isopen"` // 是否开启这个功能
	Username  string   `json:"username"`
	Port      int64    `json:"port"`
	Interval  int64    `json:"interval"` // 时间间隔
	Num       int64    `json:"num"`      // 时间间隔
	Hosts     []string `json:"hosts"`    //主机列表
	Process   string   `json:"process"`
	GameXPath string   `json:"gamexpath"`
}
