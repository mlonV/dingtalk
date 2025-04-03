package types

// 查询rediskey[type list] 的长度,监控redis队列使用
type DnsQuery struct {
	IsOpen   bool     `json:"isopen"`
	Domains  []string `json:"domains"`
	Servers  []string `json:"servers"`
	Interval int64    `json:"interval"`
}
