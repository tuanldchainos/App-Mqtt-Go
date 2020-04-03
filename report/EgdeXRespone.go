package report

type EdgeXResponse struct {
	Service     string      `json:"Service"`
	Id          string      `json:"Id"`
	Method      string      `json:"Method"`
	HttpRequest string      `json:"HttpRequest"`
	Body        interface{} `json:"Body"`
}
