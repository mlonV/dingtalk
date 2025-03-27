package super

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/mlonV/dingtalk/config"
	"github.com/mlonV/dingtalk/types/supervisor"
	"gorm.io/gorm"
)

var Supervisors = &config.Conf.Supervisors

func StartProcess(c *gin.Context) {
	client := resty.New()
	processname, ok := c.Params.Get("processname")
	if !ok {
		c.String(http.StatusOK, "processname 参数异常")
		return
	}
	// reqBody := SupervisorRequest{Method: "supervisor.getAllProcessInfo"}
	reqBody := supervisor.SupervisorRequest{Method: "supervisor.startProcess", Params: []supervisor.ReqParam{{Value: processname}}}
	reqXML, _ := xml.Marshal(reqBody)
	reqBody.Params = []supervisor.ReqParam{{Value: c.Query("process")}}
	resp, err := client.R().
		SetHeader("Content-Type", "text/xml").
		SetBasicAuth("super", "Super@123").
		SetBody(reqXML).
		Post("http://172.31.1.243:9001/RPC2")

	if err != nil {
		fmt.Println(resp.Body(), err)
	}
	c.String(200, resp.String())
}

func GetSupervisorList(c *gin.Context) {
	c.JSON(http.StatusOK, Supervisors)
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
