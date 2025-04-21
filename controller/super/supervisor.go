package super

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mlonV/dingtalk/types/supervisor"
	"gorm.io/gorm"
)

// 直接全部加载的
func GetetAllProcessInfo(c *gin.Context) {
	hosts, err := GetHosts()
	if err != nil {
		fmt.Println(err)
	}
	var processList supervisor.ProcessInfoList

	for index, host := range hosts {
		resp, err := host.SuperReq("supervisor.getAllProcessInfo", []supervisor.ReqParam{})
		if err != nil {
			// host.Status = err.Error()
			// hosts = append(hosts, host)
		}

		var methodResponse supervisor.MethodResponse
		if err := xml.Unmarshal(resp.Body(), &methodResponse); err != nil {
			fmt.Println("Error unmarshalling XML:", err)
		}
		var process supervisor.ProcessInfo
		process.ID = index + 1
		process.Host = host.Name
		process.Name = "-"
		process.Statename = "-"
		process.Description = fmt.Sprintf("主机: %s ,进程数量: %d", host.Name, len(methodResponse.Params.Param.Value.Array.Data.Value))
		for k, v := range methodResponse.Params.Param.Value.Array.Data.Value {
			var processChildren supervisor.ProcessInfo
			processChildren.ID = ((index + 1) * 1000) + k
			processChildren.Host = host.Name
			for _, member := range v.Struct.Members {
				switch member.Name {
				case "name":
					processChildren.Name = member.Value
				case "group":
					processChildren.Group = member.Value
				case "statename":
					processChildren.Statename = member.Value
				case "spawnerr":
					processChildren.Spawnerr = member.Value
				case "exitstatus":
					processChildren.Exitstatus = member.Code
				case "pid":
					processChildren.Pid = member.Code
				case "logfile":
					processChildren.Logfile = member.Value
				case "stdout_logfile":
					processChildren.Stdout_logfile = member.Value
				case "stderr_logfile":
					processChildren.Stderr_logfile = member.Value
				case "state":
					processChildren.State = member.Code
				case "now":
					processChildren.Now = member.Value
				case "start":
					processChildren.Start = member.Code
				case "stop":
					processChildren.Stop = member.Code
				case "description":
					processChildren.Description = member.Value
				}
			}
			process.Children = append(process.Children, processChildren)
		}
		processList.Items = append(processList.Items, process)
		// fmt.Println(string(resp.Body()), err)
		// fmt.Println(processList, err)
	}
	c.JSON(http.StatusOK, supervisor.Response{
		Code:    0,
		Data:    processList,
		Message: "主机process列表已解析 0.0",
	})
}

// hasChildren 设置为true,请求另一个API来获取数据，单个主机的数据
func GetAllProcessInfo_LazyLoad(c *gin.Context) {
	hosts, err := GetHosts()
	if err != nil {
		fmt.Println(err)
	}
	var processList supervisor.ProcessInfoList

	for index, host := range hosts {
		resp, err := host.SuperReq("supervisor.getAllProcessInfo", []supervisor.ReqParam{})
		if err != nil {
			// host.Status = err.Error()
			// hosts = append(hosts, host)
		}

		var methodResponse supervisor.MethodResponse
		if err := xml.Unmarshal(resp.Body(), &methodResponse); err != nil {
			fmt.Println("Error unmarshalling XML:", err)
		}
		var process supervisor.ProcessInfo
		process.ID = index + 1
		process.Host = host.Name
		process.Name = "-"
		process.Statename = "-"
		process.Description = fmt.Sprintf("主机: %s ,进程数量: %d", host.Name, len(methodResponse.Params.Param.Value.Array.Data.Value))
		if len(methodResponse.Params.Param.Value.Array.Data.Value) > 0 {
			process.HasChildren = true
		}
		processList.Items = append(processList.Items, process)

	}
	c.JSON(http.StatusOK, supervisor.Response{
		Code:    0,
		Data:    processList,
		Message: "主机process列表已解析 0.0",
	})
}
func GetProcessInfo(c *gin.Context) {
	hostname, ok1 := c.Params.Get("host")
	if !ok1 {
		c.String(http.StatusOK, "processname 参数异常")
		return
	}
	host, err := GetHostByHostname(hostname)
	if err != nil {
		fmt.Println(err)
	}
	processInfoList, err := host.GetAllProcessInfo()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    0,
			Data:    nil,
			Message: "主机GetProcessInfo 有错误 0.0",
		})
		return
	}
	c.JSON(http.StatusOK, supervisor.Response{
		Code:    0,
		Data:    processInfoList,
		Message: "主机列表加载成功",
	})

}

func StartProcess(c *gin.Context) {
	hostname, ok1 := c.Params.Get("host")
	processname, ok2 := c.Params.Get("name")
	if !ok1 || !ok2 {
		c.String(http.StatusOK, "processname 参数异常")
		return
	}

	host, err := GetHostByHostname(hostname)
	if err != nil {
		fmt.Println(err)
	}
	ok, err := host.StartProcess(processname)

	if err != nil || !ok {
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    0,
			Data:    "",
			Type:    "error",
			Message: fmt.Sprintf("请求启动失败" + err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, supervisor.Response{
		Code:    0,
		Data:    "",
		Type:    "success",
		Message: fmt.Sprintf("启动成功 :  host: %s , process: %s  ", hostname, processname),
	})
}

func StopProcess(c *gin.Context) {
	hostname, ok1 := c.Params.Get("host")
	processname, ok2 := c.Params.Get("name")
	if !ok1 || !ok2 {
		c.String(http.StatusOK, "processname 参数异常")
		return
	}

	host, err := GetHostByHostname(hostname)
	if err != nil {
		fmt.Println(err)
	}
	ok, err := host.StopProcess(processname)
	if err != nil || !ok {
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    0,
			Data:    "",
			Type:    "error",
			Message: fmt.Sprintf("请求停止失败" + err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, supervisor.Response{
		Code:    0,
		Data:    "",
		Type:    "success",
		Message: fmt.Sprintf("停止成功 :  host: %s , process: %s  ", hostname, processname),
	})
}

// todo --
func TailStdoutLog(c *gin.Context) {
	method := c.Query("method")
	hostname, ok1 := c.Params.Get("host")
	processName, ok2 := c.Params.Get("name")
	if !ok1 || !ok2 {
		c.String(http.StatusOK, "processname 参数异常")
		return
	}
	// offsetStr := c.DefaultQuery("offset", "0")
	// lengthStr := c.DefaultQuery("length", "1024") // 每次最多读取 1024 字节

	host, err := GetHostByHostname(hostname)
	if err != nil {
		fmt.Println(err)
	}

	// host.TailProcessStdoutLog(processName, offset, length)
	host.StreamLogsWS(c, processName, method)
}

func GetHosts() ([]supervisor.HostData, error) {
	var hosts []supervisor.HostData
	result := db.Find(&hosts)
	if result.Error != nil {
		return nil, result.Error
	}
	return hosts, nil
}

func GetHostByHostname(hostname string) (*supervisor.HostData, error) {
	var host *supervisor.HostData
	result := db.Where("name = ?", hostname).First(&host)
	if result.Error != nil {
		return nil, result.Error
	}
	return host, nil
}

func AddHost(c *gin.Context) {
	// 初始化hostData变量，用于存储请求中的主机信息
	var hostData supervisor.HostData
	if err := c.ShouldBindJSON(&hostData); err != nil {
		// 如果绑定JSON数据失败，返回400 Bad Request响应，并附带错误信息
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    10020,
			Data:    "",
			Message: "解析请求数据失败" + err.Error(),
		})
		return
	}

	// 检查hostData是否已经存在于数据库中，如果存在则返回错误
	// 这里假设有一个函数CheckHostExists用于检查主机是否存在
	if CheckHostExists(hostData) {
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    10021,
			Data:    "",
			Message: "URL已存在 : " + hostData.URL,
		})
		return
	}

	// 将新的hostData保存到数据库
	// 这里假设有一个函数SaveHost用于保存主机信息到数据库
	if err := SaveHost(hostData); err != nil {
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    10022,
			Data:    "",
			Message: "添加主机失败 : " + err.Error(),
		})
		return
	}

	// // 打印hostData，用于调试（生产环境中应移除）
	// fmt.Println(hostData)

	// 返回成功响应
	c.JSON(http.StatusOK, supervisor.Response{
		Code:    0,
		Data:    "",
		Message: "添加主机成功 : " + hostData.Name,
	})
}

func CheckHostExists(hostData supervisor.HostData) bool {
	var result supervisor.HostData
	// 使用主机名查询数据库中的url记录
	result = supervisor.HostData{}
	err := db.Where("url = ?", hostData.URL).First(&result).Error

	// 如果查询出错或者没有找到匹配的记录，则认为主机不存在
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			fmt.Printf("Host with URL %s not found: %v", hostData.URL, err)
			return false
		}
		fmt.Printf("Error checking URL existence: %v", err)
		return false
	}

	// 如果找到匹配的记录，则认为主机存在
	fmt.Printf("Host with Url %s found", hostData.URL)
	return true
}

func SaveHost(hostData supervisor.HostData) error {
	result := db.Create(&hostData)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func ListHost(c *gin.Context) {
	var resp supervisor.SuperList
	var hosts []supervisor.HostData

	result := db.Find(&hosts)
	if result.Error != nil {
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    0,
			Data:    nil,
			Message: result.Error.Error(),
		})
		return
	}
	resp.Total = len(hosts)
	resp.Items = hosts
	resp.AddStatus()

	c.JSON(http.StatusOK, supervisor.Response{
		Code:    0,
		Data:    resp,
		Message: "响应数据成功",
	})
}

func DelHost(c *gin.Context) {
	type Req struct {
		ID string `json:"id" uri:"id" binding:"required"`
	}
	var req Req
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    10023,
			Data:    nil,
			Message: "DELETE id 参数异常",
		})
		return
	}

	err := db.Where("id = ?", req.ID).Delete(&supervisor.HostData{}).Error
	if err != nil {
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    10024,
			Data:    nil,
			Message: "删除主机失败 : " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, supervisor.Response{
		Code:    0,
		Data:    nil,
		Message: "删除主机成功",
	})

}

func UpdateHost(c *gin.Context) {

	var hd supervisor.HostData
	if err := c.ShouldBind(&hd); err != nil {
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    10024,
			Data:    nil,
			Message: "Update 数据异常",
		})
		return
	}

	err := db.Where("id = ?", hd.ID).Updates(hd).Error
	if err != nil {
		c.JSON(http.StatusOK, supervisor.Response{
			Code:    10024,
			Data:    nil,
			Message: "修改失败 : " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, supervisor.Response{
		Code:    0,
		Data:    nil,
		Message: "修改成功成功",
	})

}
