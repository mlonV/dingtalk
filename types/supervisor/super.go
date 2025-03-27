package supervisor

import (
	"encoding/xml"
	"time"

	"github.com/go-resty/resty/v2"
)

type Supervisor struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

const (
	FATAL      string = "FATAL"
	RUNNING    string = "RUNNING"
	RESTARTING string = "RESTARTING"
	SHUTDOWN   string = "SHUTDOWN"
)

type SuperList struct {
	Total int        `json:"total"`
	Items []HostData `json:"items"`
}

type HostData struct {
	ID       int    `json:"id" gorm:"primaryKey;table:host"`
	Name     string `json:"name" gorm:"column:name"`
	URL      string `json:"url" gorm:"column:url"`
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
	Env      string `json:"env" gorm:"column:env"`
	Desc     string `json:"desc" gorm:"column:desc"`
	Status   string `json:"status" gorm:"-"`
}

// 明确指定表名为 user
func (HostData) TableName() string {
	return "host"
}

func (sl *SuperList) AddStatus() {
	var hosts []HostData
	for _, host := range sl.Items {

		client := resty.New().SetTimeout(2 * time.Second)
		reqBody := SupervisorRequest{Method: "supervisor.getState", Params: []ReqParam{}}
		reqXML, _ := xml.Marshal(reqBody)
		resp, err := client.R().
			SetHeader("Content-Type", "text/xml").
			SetBasicAuth(host.Username, host.Password).
			SetBody(reqXML).
			Post(host.URL)

		if err != nil {
			host.Status = err.Error()
			hosts = append(hosts, host)
			// fmt.Println(resp.Body(), err)
			continue
		}

		var methodResponse MethodResponse
		if err := xml.Unmarshal(resp.Body(), &methodResponse); err != nil {
			// fmt.Println("Error unmarshalling XML:", err)
			host.Status = err.Error()
			hosts = append(hosts, host)
			continue
		}
		// 查找 statecode 的值
		for _, member := range methodResponse.Params.Param.Value.Struct.Members {
			if member.Name == "statename" {
				host.Status = member.Value
				hosts = append(hosts, host)
			}
		}
	}
	sl.Items = hosts
}
