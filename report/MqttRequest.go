package report

/*
--------Here is MqttRequest report example-------
{
	"Method" : "Get",
	"Service": "edgex-core-data",
	"Path" : "/api/v1/ping",
	"Body" : {
		"OnOff" : "true"
	}
}

*/

type MqttRequest struct {
	Method  string            `json:"Method"`
	Service string            `json:"Service"`
	Path    string            `json:"Path"`
	Body    map[string]string `json:"Body"`
}
