package config

type YearningMsg struct {
	Msgtype  string                 `json:"msgtype,omitempty"`
	At       map[string]interface{} `json:"at,omitempty"` // 为了实现发送消息at相关人员准备的
	Markdown `json:"markdown,omitempty"`
}

type Markdown struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}
