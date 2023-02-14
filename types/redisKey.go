package types

// 查询rediskey[type list] 的长度,监控redis队列使用
type RedisKey struct {
	IsOpen   bool     `json:"isopen"`
	Hostname string   `json:"hostname"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Port     int      `json:"port"`
	Keys     []string `json:"keys"`
	Interval int64    `json:"interval"`
}
