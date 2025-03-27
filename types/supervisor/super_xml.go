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

// // UnmarshalXML 解析 member 中的 value
// func (m *Member) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	type T Member
// 	var t T
// 	if err := d.DecodeElement(&t, &start); err != nil {
// 		return err
// 	}

// 	*m = Member{
// 		Name: t.Name,
// 	}

// 	// 解析 value 的类型
// 	var value struct {
// 		Int    *XmlIntValue    `xml:"int"`
// 		String *XmlStringValue `xml:"string"`
// 	}

// 	if err := d.DecodeElement(&value, &start); err != nil {
// 		return err
// 	}

// 	if value.Int != nil {
// 		m.Value = value.Int
// 	} else if value.String != nil {
// 		m.Value = value.String
// 	}

// 	return nil
// }
