package types

// 查询rediskey[type list] 的长度,监控redis队列使用
type DnsQuery struct {
	IsOpen          bool     `json:"isopen"`
	Domain          string   `json:"domain"`
	ARecords        []string `json:"arecords"`
	CnameRecords    []string `json:"cnamerecords"`
	Servers         []string `json:"servers"`
	Interval        int64    `json:"interval"`
	QueryTimeOut    int64    `json:"querytimeout"`
	DelayUnregister int64    `json:"delayunregister"`
}
