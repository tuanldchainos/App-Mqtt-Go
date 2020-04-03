package report

/*
--------Here is MqttRequest report example-------
{
	"Method" : "Get",
	"Service": "CoreData",
	"Id" : "1",
	"Path" : "/api/v1/ping",
	"Body" : ""
}

*/

type MqttRequest struct {
	Service string      `json:"Service"`
	Id      string      `json:"Id"`
	Method  string      `json:"Method"`
	Path    string      `json:"Path"`
	Body    interface{} `json:"Body"`
}
