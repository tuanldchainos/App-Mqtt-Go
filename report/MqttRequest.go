package report

/*
--------Here is MqttRequest report example-------
{
	"Method" : "Get",
	"Service": "CoreData",
	"Id" : "1",
	"Key" :"",
	"Path" : "/api/v1/ping",
	"Body" : ""
}

*/

type MqttRequest struct {
	Service string      `json:"Service"`
	Id      string      `json:"Id"`
	Key     string      `json:"Key"`
	Method  string      `json:"Method"`
	Path    string      `json:"Path"`
	Body    interface{} `json:"Body"`
}
