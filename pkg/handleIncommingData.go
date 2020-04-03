package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"App-Mqtt-Go/report"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var requestData = make(chan *report.MqttRequest)

var wg sync.WaitGroup

var mux sync.Mutex

func StartListeningMqttIncoming(client MQTT.Client, config *MqttConfig, urlList func(string) string) {
	wg.Add(2)
	go func() {
		token := client.Subscribe(Resquest_Topic, byte(config.MqttQos), OnHandleMqttIncomming)
		if token.Wait() && token.Error() != nil {
			log.Info(fmt.Sprintf("[Incoming listener] Stop incoming data listening. Cause:%v", token.Error()))
		}
		log.Info("[Incoming listener] Start incoming data listening.")
	}()

	go func() {
		for {
			select {
			case MqttRequest := <-requestData:
				mux.Lock()
				msg, err := createHttpRequest(MqttRequest, urlList)
				if err != nil {
					mqttResponseFail(client, MqttRequest, err.Error())
				} else {
					mqttResponseSuccess(client, MqttRequest, msg)
				}
				log.Info(fmt.Sprintln("Successfully handle request!"))
				mux.Unlock()
			}
		}
	}()
	wg.Wait()
}

func OnHandleMqttIncomming(client MQTT.Client, message MQTT.Message) {
	fmt.Println("Request receving: ", string(message.Payload()))
	var MqttRequest = report.MqttRequest{}
	err := json.Unmarshal(message.Payload(), &MqttRequest)
	if err != nil {
		log.Error(fmt.Sprintln("Exception when parse json report: %s", err))
	}

	// var dataCheck map[string]interface{}
	// json.Unmarshal(message.Payload(), &dataCheck)
	// if !checkDataWithKey(dataCheck, "Method") || !checkDataWithKey(dataCheck, "Service") || !checkDataWithKey(dataCheck, "Path") || !checkDataWithKey(dataCheck, "Body") || !checkDataWithKey(dataCheck, "Id") {
	// 	log.Info(fmt.Sprintln("Unrecognize mqtt request"))
	// 	mqttResponseFail(client, &MqttRequest, "Unrecognize mqtt request")
	// }
	requestData <- &MqttRequest
}

func createHttpRequest(req *report.MqttRequest, urlList func(string) string) (string, error) {
	ctx := context.TODO()
	httpUrl := urlList(req.Service) + req.Path

	switch req.Method {
	case "Get":
		res, err := clients.GetRequestWithURL(ctx, httpUrl)
		return string(res), err
	case "Post":
		res, err := clients.PostJSONRequestWithURL(ctx, httpUrl, req.Body)
		return res, err
	//case "Put":
	default:
		return "", errors.New("Unknow http method")
	}
}

// func checkDataWithKey(data map[string]interface{}, key string) bool {
// 	val, ok := data[key]
// 	if !ok {
// 		log.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. No %v found : msg=%v", key, data))
// 		return false
// 	}

// 	switch val.(type) {
// 	case string:
// 		return true
// 	default:
// 		log.Warn(fmt.Sprintf("[Incoming listener] Incoming reading ignored. %v should be string : msg=%v", key, data))
// 		return false
// 	}
// }

func mqttResponseSuccess(client MQTT.Client, MqttRequest *report.MqttRequest, msg string) {
	var EdgeXResponse report.EdgeXResponse
	EdgeXResponse.Service = (*MqttRequest).Service
	EdgeXResponse.Method = (*MqttRequest).Service
	EdgeXResponse.Id = (*MqttRequest).Id
	EdgeXResponse.HttpRequest = (*MqttRequest).Path
	EdgeXResponse.Body = msg

	Response, _ := json.Marshal(EdgeXResponse)
	client.Publish(Response_Topic, 0, false, string(Response))
	fmt.Println("Response sending: ", string(Response))
	wg.Done()
}

func mqttResponseFail(client MQTT.Client, MqttRequest *report.MqttRequest, err string) {
	var EdgeXResponse report.EdgeXResponse
	EdgeXResponse.Service = (*MqttRequest).Service
	EdgeXResponse.Method = (*MqttRequest).Service
	EdgeXResponse.Id = (*MqttRequest).Id
	EdgeXResponse.HttpRequest = (*MqttRequest).Path
	EdgeXResponse.Body = err

	Response, _ := json.Marshal(EdgeXResponse)
	client.Publish(Response_Topic, 0, false, string(Response))
	fmt.Println("Response sending: ", string(Response))
	wg.Done()
}
