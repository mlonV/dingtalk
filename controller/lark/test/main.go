package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

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

// LarkMessage holds Lark message structure
type LarkMessage struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

// LarkInteractiveMessage represents an interactive message
type LarkInteractiveMessage struct {
	MsgType string `json:"msg_type"`
	Card    struct {
		Elements []struct {
			Tag  string `json:"tag"`
			Text struct {
				Tag     string `json:"tag"`
				Content string `json:"content"`
			} `json:"text"`
			URL string `json:"url"`
		} `json:"elements"`
	} `json:"card"`
}

// var jenkinsConfig JenkinsConfig
var jenkinsConfig = JenkinsConfig{
	URL:      "http://192.168.254.100:18080",
	Username: "admin",
	APIToken: "abcd1234!",
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
		if strings.Contains(v.Name, job) {
			return true
		}
	}
	return false
}

// Trigger Jenkins build
func TriggerJenkinsBuild(jobName string) error {
	client := resty.New()
	_, err := client.R().
		SetBasicAuth(jenkinsConfig.Username, jenkinsConfig.APIToken).
		Post(fmt.Sprintf("%s/job/%s/build", jenkinsConfig.URL, jobName))

	return err
}

// Send interactive message to Lark
func SendMessageToLarkWebhook(jobs []Job, larkWebhookURL string) error {
	message := LarkInteractiveMessage{
		MsgType: "interactive",
	}

	for _, job := range jobs {
		element := struct {
			Tag  string `json:"tag"`
			Text struct {
				Tag     string `json:"tag"`
				Content string `json:"content"`
			} `json:"text"`
			URL string `json:"url"`
		}{
			Tag: "a",
			Text: struct {
				Tag     string `json:"tag"`
				Content string `json:"content"`
			}{
				Tag:     "lark_md",
				Content: fmt.Sprintf("**%s**", job.Name),
			},
			URL: job.URL,
		}
		message.Card.Elements = append(message.Card.Elements, element)
	}

	client := resty.New()
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(message).
		Post(larkWebhookURL)

	return err
}

type FreeStyleProject struct {
	XMLName xml.Name `xml:"freeStyleProject"`
	Actions []Action `xml:"action"`
}

type Action struct {
	XMLName              xml.Name              `xml:"action"`
	ParameterDefinitions []ParameterDefinition `xml:"parameterDefinition"`
}

type ParameterDefinition struct {
	XMLName          xml.Name              `xml:"parameterDefinition"`
	Name             string                `xml:"name"`
	Description      string                `xml:"description,omitempty"`
	DefaultParameter DefaultParameterValue `xml:"defaultParameterValue"`
	Type             string                `xml:"type"`
	AllValueItems    []ValueItem           `xml:"allValueItems>value,omitempty"`
	Choices          []string              `xml:"choice,omitempty"`
}

type DefaultParameterValue struct {
	XMLName xml.Name `xml:"defaultParameterValue"`
	Name    string   `xml:"name"`
	Value   string   `xml:"value"`
}

type ValueItem struct {
	XMLName xml.Name `xml:"value"`
	Name    string   `xml:"name"`
	Value   string   `xml:"value"`
}

func main() {
	// jobname := "test-jenkins"
	jobname := "in-user-recommend.ohlaapp.com"
	client := resty.New()
	resp, err := client.R().
		SetBasicAuth(jenkinsConfig.Username, jenkinsConfig.APIToken).
		SetHeader("Content-Type", "application/json").
		Get(fmt.Sprintf("%s/job/%s/api/xml?depth=1", jenkinsConfig.URL, jobname))
	if err != nil {
		fmt.Println(err)
	}
	// 创建一个 FreeStyleProject 实例
	var project FreeStyleProject

	// 使用 XML 解码器解析数据
	decoder := xml.NewDecoder(strings.NewReader(string(resp.Body())))
	err = decoder.Decode(&project)
	if err != nil {
		fmt.Sprintln("Failed to decode XML:", err)
	}
	// 打印解析结果
	for _, action := range project.Actions {
		for _, param := range action.ParameterDefinitions {
			fmt.Printf("Parameter Name: %s\nDescription: %s\nDefault Value: %s\nType: %s\n", param.Name, param.Description, param.DefaultParameter.Value, param.Type)
			if len(param.AllValueItems) > 0 {
				fmt.Println("All Value Items:")
				for _, valueItem := range param.AllValueItems {
					fmt.Printf("\tName: %s, Value: %s\n", valueItem.Name, valueItem.Value)
				}
			}
			if len(param.Choices) > 0 {
				fmt.Println("Choices:")
				for _, choice := range param.Choices {
					fmt.Printf("\tChoice: %s\n", choice)
				}
			}
			fmt.Println()
		}
	}
}
