package conlark

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/mlonV/dingtalk/types"
)

type JobParam struct {
	Name  string
	Value string
}

type BuildData struct {
	Job   Job
	Param []JobParam
}

// var jenkinsConfig JenkinsConfig
var jenkinsConfig = JenkinsConfig{
	URL:      "http://192.168.254.100:18080",
	Username: "admin",
	APIToken: "11d4e409fb4ac779b4033658ae0a11bb78",
}

// JenkinsConfig holds Jenkins connection details
type JenkinsConfig struct {
	URL      string
	Username string
	APIToken string
}

// Job holds Jenkins job details
type Job struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Fetch Jenkins jobs
func GetJenkinsJobs() ([]Job, error) {
	client := resty.New()
	resp, err := client.R().
		SetBasicAuth(jenkinsConfig.Username, jenkinsConfig.APIToken).
		SetHeader("Content-Type", "application/json").
		Get(fmt.Sprintf("%s/api/json", jenkinsConfig.URL))

	if err != nil {
		return nil, err
	}

	var jobs struct {
		Jobs []Job `json:"jobs"`
	}
	err = json.Unmarshal(resp.Body(), &jobs)
	if err != nil {
		return nil, err
	}

	return jobs.Jobs, nil
}

// Fetch Jenkins jobs
func GetJenkinsUrlByJob(jobName string) string {

	jobs, _ := GetJenkinsJobs()
	// 过滤出需要的job
	for _, v := range jobs {
		if v.Name == jobName {
			return v.URL
		}
	}
	return ""
}

// 根据job名过滤获取对应的job
func GetJenkinsJobsFilter(job string) string {
	jobs, _ := GetJenkinsJobs()
	resJobList := []Job{}
	// 过滤出需要的job
	for _, v := range jobs {
		if strings.Contains(v.Name, job) {
			resJobList = append(resJobList, v)
		}
	}

	var jobList string
	jobList += fmt.Sprintf("jenkins job 匹配数量: %v \n", len(resJobList))
	// 返回给Lark要字符串的形式
	for k, v := range resJobList {
		if k > 4 {
			break
		}
		jobList += v.Name + " : " + v.URL + "\n"
	}
	return jobList
}

// 查看job是否存在
func JobIsExists(job string) bool {
	jobs, _ := GetJenkinsJobs()
	// 过滤出需要的job
	for _, v := range jobs {
		if v.Name == job {
			return true
		}
	}
	return false
}

// Trigger Jenkins build
func TriggerJenkinsBuild(bd *BuildData) error {
	client := resty.New().R().SetBasicAuth(jenkinsConfig.Username, jenkinsConfig.APIToken)

	postStr := ""
	if bd.Param == nil {
		postStr = fmt.Sprintf("%s/job/%s/build", jenkinsConfig.URL, bd.Job.Name)
		resp, err := client.Post(postStr)
		fmt.Println(string(resp.Body()))
		return err

	} else {
		postStr = fmt.Sprintf("%s/job/%s/buildWithParameters", jenkinsConfig.URL, bd.Job.Name)
		paraMap := map[string]string{}
		for _, v := range bd.Param {
			paraMap[v.Name] = v.Value
		}
		resp, err := client.SetPathParams(paraMap).Post(postStr)
		fmt.Println(string(resp.Body()))
		return err

	}

}

// 获取jenkins的参数配置
func GetJobConfig(jobname string) (*types.FreeStyleProject, error) {
	client := resty.New()
	resp, err := client.R().
		SetBasicAuth(jenkinsConfig.Username, jenkinsConfig.APIToken).
		SetHeader("Content-Type", "application/json").
		Get(fmt.Sprintf("%s/job/%s/api/xml?depth=1", jenkinsConfig.URL, jobname))
	if err != nil {
		return nil, err
	}
	// 创建一个 FreeStyleProject 实例
	var project *types.FreeStyleProject

	// 使用 XML 解码器解析数据
	decoder := xml.NewDecoder(strings.NewReader(string(resp.Body())))
	err = decoder.Decode(&project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func BuildJob(jobname string) error {
	client := resty.New()
	resp, err := client.R().
		SetBasicAuth(jenkinsConfig.Username, jenkinsConfig.APIToken).
		SetHeader("Content-Type", "application/json").
		// SetPathParams().
		Post(fmt.Sprintf("%s/job/%s/buildWithParameters", jenkinsConfig.URL, jobname))
	if err != nil {
		return err
	}
	fmt.Println(resp.Body())
	return nil
}
