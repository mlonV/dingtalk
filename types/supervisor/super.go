package supervisor

import (
	"encoding/xml"
	"fmt"
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

// 机器上的进程信息
// State 状态标识符
// STOPPED (0)
// STARTING (10)
// RUNNING (20)
// BACKOFF (30)
// STOPPING (40)
// EXITED (100)
// FATAL (200)
// UNKNOWN (1000)
type ProcessInfoList struct {
	Total int           `json:"total"`
	Items []ProcessInfo `json:"items"`
}

type ProcessInfo struct {
	ID             int    `json:"id"`
	Host           string `json:"host"`
	Name           string `json:"name"`
	Group          string `json:"group"`
	Statename      string `json:"statename"`
	Spawnerr       string `json:"spawnerr"`
	Exitstatus     int    `json:"exitstatus"`
	Pid            int    `json:"pid"`
	Logfile        string `json:"logfile"`
	Stdout_logfile string `json:"stdout_logfile"`
	Stderr_logfile string `json:"stderr_logfile"`
	State          int    `json:"state"`
	Now            string `json:"now"`
	Start          int    `json:"start"`
	Stop           int    `json:"stop"`
	Description    string `json:"description"`

	// 给vue用的树状表格
	Children []ProcessInfo `json:"children,omitempty"`
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

// hostdata请求process列表
func (hd *HostData) SuperReq(method string, params []ReqParam) (*resty.Response, error) {
	client := resty.New()
	reqBody := SupervisorRequest{Method: method, Params: params}
	reqXML, _ := xml.Marshal(reqBody)

	return client.R().
		SetHeader("Content-Type", "text/xml").
		SetBasicAuth(hd.Username, hd.Password).
		SetBody(reqXML).
		Post(hd.URL)
}

func (hd *HostData) StartProcess(processName string) (bool, error) {
	resp, err := hd.SuperReq("supervisor.startProcess", []ReqParam{{Value: processName}})
	if err != nil {
		return false, err
	}
	var methodResponseSuccess MethodResponseSuccess
	var methodResponseFault MethodResponseFault

	if err := xml.Unmarshal(resp.Body(), &methodResponseFault); err == nil {
		// 错误响应
		for _, member := range methodResponseFault.Fault.Value.Struct.Members {
			if member.Name == "faultString" {
				return false, fmt.Errorf("Failed to start process: %s - %s", processName, member.Value)
			}

		}
	}

	if err := xml.Unmarshal(resp.Body(), &methodResponseSuccess); err == nil {
		// 成功响应
		if methodResponseSuccess.Params.Param.Value.Boolean == 1 {
			return true, nil
		}
	}

	return false, fmt.Errorf("Unexpected response format")
}

func (hd *HostData) StopProcess(processName string) (bool, error) {

	resp, err := hd.SuperReq("supervisor.stopProcess", []ReqParam{{Value: processName}})
	if err != nil {
		return false, err
	}
	var methodResponseSuccess MethodResponseSuccess
	var methodResponseFault MethodResponseFault

	if err := xml.Unmarshal(resp.Body(), &methodResponseFault); err == nil {
		// 错误响应
		for _, member := range methodResponseFault.Fault.Value.Struct.Members {
			if member.Name == "faultString" {
				return false, fmt.Errorf("Failed to Stop process: %s - %s", processName, member.Value)
			}

		}
	}

	if err := xml.Unmarshal(resp.Body(), &methodResponseSuccess); err == nil {
		// 成功响应
		if methodResponseSuccess.Params.Param.Value.Boolean == 1 {
			return true, nil
		}
	}

	return false, fmt.Errorf("Unexpected response format")
}

func (hd *HostData) TailProcessStdoutLog(processName string, offset, length int) {
	resp, err := hd.SuperReq("supervisor.tailProcessStdoutLog", []ReqParam{
		{Value: processName},
		{IntValue: offset},
		{IntValue: length},
	})
	if err != nil {
		// return false, err
	}

	fmt.Println(resp)
}
