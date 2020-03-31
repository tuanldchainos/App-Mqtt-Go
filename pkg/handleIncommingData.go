package pkg

import (
	"fmt"
	"sync"

	"App-Mqtt-Go/constant/topic"
	"App-Mqtt-Go/report"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var wg sync.WaitGroup

func StartListeningMqttIncoming(client MQTT.Client, config *MqttConfig) {
	wg.Add(1)
	go func() {
		topiclist := topic.GetTopicList()
		token := client.Subscribe(topiclist("Request"), byte(config.MqttQos), OnHandleMqttIncomming)
		if token.Wait() && token.Error() != nil {
			log.Info(fmt.Sprintf("[Incoming listener] Stop incoming data listening. Cause:%v", token.Error()))
		}
		log.Info("[Incoming listener] Start incoming data listening.")
		return
	}()
	wg.Wait()
	select {}
}

func OnHandleMqttIncomming(client MQTT.Client, message MQTT.Message) {
	dataInComming := string(message.Payload())

	wg.Done()
}

func sendingHttpRequest(request *report.MqttRequest) (string, error) {

}

func createMqttResponse(EdgeX)

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
