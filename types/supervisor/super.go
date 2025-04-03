package supervisor

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
)

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
	// var hosts []HostData
	for index, host := range sl.Items {

		client := resty.New().SetTimeout(1 * time.Second)
		reqBody := SupervisorRequest{Method: "supervisor.getState", Params: []ReqParam{}}
		reqXML, _ := xml.Marshal(reqBody)
		resp, err := client.R().
			SetHeader("Content-Type", "text/xml").
			SetBasicAuth(host.Username, host.Password).
			SetBody(reqXML).
			Post(host.URL)

		if err != nil {
			sl.Items[index].Status = err.Error()
			// host.Status = err.Error()
			// hosts = append(hosts, host)
			// fmt.Println(resp.Body(), err)
			continue
		}

		var methodResponse MethodResponse
		if err := xml.Unmarshal(resp.Body(), &methodResponse); err != nil {
			// fmt.Println("Error unmarshalling XML:", err)
			sl.Items[index].Status = err.Error()
			// host.Status = err.Error()
			// hosts = append(hosts, host)
			continue
		}
		// 查找 statecode 的值
		for _, member := range methodResponse.Params.Param.Value.Struct.Members {
			if member.Name == "statename" {
				sl.Items[index].Status = member.Value
				// hosts = append(hosts, host)
				break
			}
		}
		if sl.Items[index].Status == "" {
			sl.Items[index].Status = "未查询到状态"
		}
		// hosts = append(hosts, host)
	}
	// sl.Items = hosts
}

// hostdata请求process列表
func (hd *HostData) SuperReq(method string, params []ReqParam) (*resty.Response, error) {
	client := resty.New()
	client.SetTimeout(3 * time.Second)
	reqBody := SupervisorRequest{Method: method, Params: params}
	reqXML, _ := xml.Marshal(reqBody)

	return client.R().
		SetHeader("Content-Type", "text/xml").
		SetBasicAuth(hd.Username, hd.Password).
		SetBody(reqXML).
		Post(hd.URL)
}

func (hd *HostData) StartProcess(processName string) (bool, error) {
	resp, err := hd.SuperReq("supervisor.startProcess", []ReqParam{{Value: ReqValue{StringValue: &processName}}})
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

	resp, err := hd.SuperReq("supervisor.stopProcess", []ReqParam{{Value: ReqValue{StringValue: &processName}}})
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
		{Value: ReqValue{StringValue: &processName}},
		{Value: ReqValue{IntValue: &offset}},
		{Value: ReqValue{IntValue: &length}},
	})
	if err != nil {
		// return false, err
	}

	fmt.Println(resp.RawBody())
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func cleanXML(data []byte) []byte {
	re := regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F-\x9F]`)
	return re.ReplaceAll(data, []byte{})
}

func (hd *HostData) StreamLogsWS(c *gin.Context, processName, method string) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	offset := 0
	length := 0

	var logResponse LogResponse
	resp, err := hd.SuperReq(method, []ReqParam{
		{Value: ReqValue{StringValue: &processName}},
		{Value: ReqValue{IntValue: &offset}},
		{Value: ReqValue{IntValue: &length}},
	})
	if err := xml.Unmarshal(resp.Body(), &logResponse); err != nil {
		fmt.Println("Error unmarshalling XML:", err)
		return
	}

	offset = logResponse.Params.Param.Value.Array.Data.Values[1].IntValue - 1024 // 默认偏移量
	length = 1024

	// print(resp)
	done := make(chan struct{})

	// 处理 WebSocket 关闭事件
	go func() {
		defer close(done)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("WebSocket closed:", err)
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			fmt.Println("Stopping log stream due to WebSocket closure")
			return
		default:
			resp, err := hd.SuperReq(method, []ReqParam{
				{Value: ReqValue{StringValue: &processName}},
				{Value: ReqValue{IntValue: &offset}},
				{Value: ReqValue{IntValue: &length}},
			})
			if err != nil {
				break
			}

			cleanedBody := cleanXML(resp.Body())

			var logResponse LogResponse
			if err := xml.Unmarshal(cleanedBody, &logResponse); err != nil {
				// fmt.Println(string(resp.Body()))
				fmt.Println("xml 解析 err : ", err)
			}
			if offset >= logResponse.Params.Param.Value.Array.Data.Values[1].IntValue || len(logResponse.Params.Param.Value.Array.Data.Values) < 2 {
				time.Sleep(1000 * time.Millisecond)
				continue
			}

			logs := logResponse.Params.Param.Value.Array.Data.Values[0].StringValue
			length = logResponse.Params.Param.Value.Array.Data.Values[1].IntValue - offset
			offset += length // 更新偏移量

			if err := conn.WriteMessage(websocket.TextMessage, []byte(logs)); err != nil {
				break
			}

			time.Sleep(1 * time.Second)
		}
	}
}
