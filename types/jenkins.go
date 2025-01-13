package types

import "encoding/xml"

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
