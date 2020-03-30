package pkg

import (
	"fmt"
	"sync"

	"App-Mqtt-Go/constant/topic"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var wg sync.WaitGroup

func StartListeningMqttIncoming(client MQTT.Client, topiclist func(int) string, config *MqttConfig) {
	numOfTopic := topic.GetNumOfTopic(topiclist)
	wg.Add(numOfTopic)
	for i := 1; i < numOfTopic+1; i++ {
		go func(num int) {
			token := client.Subscribe(topiclist(num), byte(config.MqttQos), OnHandleMqttIncomming)
			if token.Wait() && token.Error() != nil {
				log.Info(fmt.Sprintf("[Incoming listener] Stop incoming data listening. Cause:%v", token.Error()))
				return
			}
			log.Info("[Incoming listener] Start incoming data listening.")
		}(i)
	}
	wg.Wait()
	select {}
}

func OnHandleMqttIncomming(client MQTT.Client, message MQTT.Message) {

	wg.Done()
}
