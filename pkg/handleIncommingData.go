package pkg

import (
	"context"
	"encoding/json"
	"fmt"

	"App-Mqtt-Go/helper/topic"
	"App-Mqtt-Go/report"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
)

//var wg sync.WaitGroup

func StartListeningMqttIncoming(client MQTT.Client, config *MqttConfig) {
	//wg.Add(1)
	go func() {
		topiclist := topic.GetTopicList()
		token := client.Subscribe(topiclist("Request"), byte(config.MqttQos), OnHandleMqttIncomming)
		if token.Wait() && token.Error() != nil {
			log.Info(fmt.Sprintf("[Incoming listener] Stop incoming data listening. Cause:%v", token.Error()))
		}
		log.Info("[Incoming listener] Start incoming data listening.")
		return
	}()
	//wg.Wait()
	select {}
}

func OnHandleMqttIncomming(client MQTT.Client, message MQTT.Message) {
	var MqttRequest = report.MqttRequest{}
	err := json.Unmarshal(message.Payload(), &MqttRequest)
	if err != nil {
		log.Error(fmt.Sprintln("Exception when parse json report: %s", err))
	}

	var dataCheck map[string]interface{}
	json.Unmarshal(message.Payload(), &dataCheck)
	if !checkDataWithKey(dataCheck, "Method") || !checkDataWithKey(dataCheck, "Service") || !checkDataWithKey(dataCheck, "Path") || !checkDataWithKey(dataCheck, "Body") {
		return
	}

	fmt.Println(MqttRequest)
	//wg.Done()
}

func getUrlEndpoint(MqttRequest *report.MqttRequest) string {

}

func sendingHttpRequest(request *report.MqttRequest) (string, error) {
	ctx := context.Background()

	dataHttpResponse, err := clients.GetRequest(ctx)
	if err != nil {

	}
	select {}
}

//func createMqttResponse(EdgeX)

func checkDataWithKey(data map[string]interface{}, key string) bool {
	val, ok := data[key]
	if !ok {
		log.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No %v found : msg=%v", key, data))
		return false
	}

	switch val.(type) {
	case string:
		return true
	default:
		log.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. %v should be string : msg=%v", key, data))
		return false
	}
}
