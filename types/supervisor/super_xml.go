package supervisor

import "encoding/xml"

type MethodResponse struct {
	Params Params `xml:"params"`
}

type Params struct {
	Param Param `xml:"param"`
}

type Param struct {
	Value Value `xml:"value"`
}

type Value struct {
	Struct  Struct `xml:"struct"`
	Array   Array  `xml:"array"`
	Boolean int    `xml:"boolean"`
}

type Array struct {
	Data Data `xml:"data"`
}

// GetetAllProcessInfo use
type Data struct {
	Value []AllProcessValue `xml:"value"`
}

type AllProcessValue struct {
	Struct Struct `xml:"struct"`
}

type Struct struct {
	Members []Member `xml:"member"`
}

type Member struct {
	Name  string `xml:"name"`
	Value string `xml:"value>string"`
	Code  int    `xml:"value>int"`
}

type SupervisorRequest struct {
	XMLName xml.Name   `xml:"methodCall"`
	Method  string     `xml:"methodName"`
	Params  []ReqParam `xml:"params>param"`
}

type ReqParam struct {
	Value string `xml:"value>string"`
}

// 响应结构体

// MethodResponseSuccess 定义成功响应结构体
type MethodResponseSuccess struct {
	Params Params `xml:"params"`
}

// MethodResponseFault 定义错误响应结构体
type MethodResponseFault struct {
	Fault Fault `xml:"fault"`
}

type Fault struct {
	Value Value `xml:"value"`
}
