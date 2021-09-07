package config

type YearningMsg struct {
	Msgtype  string                `json:"msgtype,omitempty"`
	At       `json:"at,omitempty"` // 为了实现发送消息at相关人员准备的
	Markdown `json:"markdown,omitempty"`
}

type Markdown struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}

// 发送钉钉根据手机号@组员
type At struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}

// 配置文件序列化成结构体
type YearningFileConf struct {
	Yearning []Yearning `json:"yearning,omitempty"`
}

type Yearning struct {
	Url       string   `json:"url,omitempty"`
	Secret    string   `json:"secret,omitempty"`
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}
