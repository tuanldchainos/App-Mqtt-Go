package mqttHandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"App-Mqtt-Go/pkg"
	"App-Mqtt-Go/pkg/connect/mqttConnect"
	"App-Mqtt-Go/report"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/tuanldchainos/app-functions-sdk-go/appsdk"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// MqttHandle struct is a object, that have responsibility for handling a mqtt request
type MqttHandle struct {
	sdk         *appsdk.AppFunctionsSDK
	requestData chan *report.MqttRequest
	wg          sync.WaitGroup
	mux         sync.Mutex
}

// NewMqttHandle return a handle mqtt struct
func NewMqttHandle(sdk *appsdk.AppFunctionsSDK) *MqttHandle {
	var requestData = make(chan *report.MqttRequest)
	return &MqttHandle{
		requestData: requestData,
		sdk:         sdk,
	}
}

// StartListeningMqttIncoming listen mqtt report incomming and handle request
func (f *MqttHandle) StartListeningMqttIncoming(client MQTT.Client) {
	log := f.sdk.LoggingClient
	responseTopicLists := mqttConnect.GetResponseTopicList(f.sdk)
	Qos := mqttConnect.GetMqttQos(f.sdk)
	requestTopic := mqttConnect.GetRequestTopic(f.sdk)
	f.wg.Add(2)
	go func() {
		token := client.Subscribe(requestTopic, byte(Qos), f.onHandleMqttIncomming)
		if token.Wait() && token.Error() != nil {
			log.Info(fmt.Sprintf("[Incoming listener] Stop incoming data listening. Cause:%v", token.Error()))
		}
		log.Info("[Incoming listener] Start incoming data listening.")
	}()
	go func() {
		for {
			select {
			case MqttRequest := <-f.requestData:
				f.mux.Lock()
				msg, err := f.createHTTPRequest(MqttRequest)
				if err != nil {
					f.mqttResponseFail(client, MqttRequest, err.Error(), responseTopicLists, Qos)
					log.Error(fmt.Sprintf("Exception when handling http method: %s", err))
				} else {
					f.mqttResponseSuccess(client, MqttRequest, msg, responseTopicLists, Qos)
					log.Info(fmt.Sprintln("Successfully handle request!"))
				}
				f.mux.Unlock()
			}
		}
	}()
	f.wg.Wait()
}

func (f *MqttHandle) onHandleMqttIncomming(client MQTT.Client, message MQTT.Message) {
	log := f.sdk.LoggingClient
	fmt.Println("Request receving: ", string(message.Payload()))
	log.Info(fmt.Sprintf("Request receving: %s", string(message.Payload())))
	var MqttRequest = report.MqttRequest{}
	err := json.Unmarshal(message.Payload(), &MqttRequest)
	if err != nil {
		log.Error(fmt.Sprintf("Exception when parse json report: %s", err))
	}
	f.requestData <- &MqttRequest
}

func mqttCommandFilter(req *report.MqttRequest) error {
	switch req.Key {
	case "":
		if req.Service == pkg.CoreCommandClientName || req.Service == pkg.CoreDataClientName || req.Service == pkg.MetadataClientName || req.Service == pkg.LoggingClientName || req.Service == pkg.SystemAgentClienName || req.Service == pkg.SchedulerClientName || req.Service == pkg.NotificationsClientName {
			return errors.New("Not have permission")
		}
		return nil
	case pkg.SecretDevKey:
		return nil
	default:
		return errors.New("Error secret key")
	}
}

func (f *MqttHandle) createHTTPRequest(req *report.MqttRequest) (string, error) {
	err := mqttCommandFilter(req)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	httpURL := pkg.UrlList[req.Service] + req.Path

	switch req.Method {
	case "Get":
		res, err := clients.GetRequestWithURL(ctx, httpURL)
		return string(res), err
	case "Post":
		res, err := clients.PostJSONRequestWithURL(ctx, httpURL, req.Body)
		return res, err
	case "Put":
		urlPre := local.New(pkg.UrlList[req.Service])
		putReqBody, _ := json.Marshal(req.Body)
		res, err := clients.PutRequest(ctx, req.Path, putReqBody, urlPre)
		return res, err
	case "Delete":
		urlPre := local.New(pkg.UrlList[req.Service])
		err := clients.DeleteRequest(ctx, req.Path, urlPre)
		return "", err
	default:
		return "", errors.New("Unknown http method")
	}
}

func (f *MqttHandle) mqttResponseSuccess(client MQTT.Client, MqttRequest *report.MqttRequest, msg string, responseTopicLists []string, qos int) {
	log := f.sdk.LoggingClient
	var EdgeXResponse report.EdgeXResponse
	EdgeXResponse.Service = (*MqttRequest).Service
	EdgeXResponse.Method = (*MqttRequest).Service
	EdgeXResponse.Id = (*MqttRequest).Id
	EdgeXResponse.HttpRequest = (*MqttRequest).Path
	EdgeXResponse.Body = msg

	Response, _ := json.Marshal(EdgeXResponse)
	for i := 0; i < len(responseTopicLists); i++ {
		client.Publish(responseTopicLists[i], byte(qos), false, string(Response))
	}
	fmt.Println("Response sending: ", string(Response))
	log.Info(fmt.Sprintf("response sending: %s", string(Response)))
}

func (f *MqttHandle) mqttResponseFail(client MQTT.Client, MqttRequest *report.MqttRequest, err string, responseTopicLists []string, qos int) {
	log := f.sdk.LoggingClient
	var EdgeXResponse report.EdgeXResponse
	EdgeXResponse.Service = (*MqttRequest).Service
	EdgeXResponse.Method = (*MqttRequest).Method
	EdgeXResponse.Id = (*MqttRequest).Id
	EdgeXResponse.HttpRequest = (*MqttRequest).Path
	EdgeXResponse.Body = err

	Response, _ := json.Marshal(EdgeXResponse)
	for i := 0; i < len(responseTopicLists); i++ {
		client.Publish(responseTopicLists[i], byte(qos), false, string(Response))
	}
	fmt.Println("Response sending: ", string(Response))
	log.Info(fmt.Sprintf("response sending: %s", string(Response)))
}
