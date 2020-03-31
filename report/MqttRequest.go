package report

type MqttRequest struct {
	Method  string `json:"Method"`
	Service string `json:"Service"`
	Path    string `json:"Path"`
	Body    string `json:"Body"`
}
