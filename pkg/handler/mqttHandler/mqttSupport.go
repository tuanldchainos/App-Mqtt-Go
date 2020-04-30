package mqttHandler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tuanldchainos/app-functions-sdk-go/appsdk"
)

const (
	RequestTopic  = "RequestTopic"
	ResponseTopic = "ResponseTopic"
	Qos           = "Qos"
)

func GetRequestTopic(sdk *appsdk.AppFunctionsSDK) string {
	log := sdk.LoggingClient
	appSettings := sdk.ApplicationSettings()

	var topic string
	if appSettings != nil {
		topic = getAppSetting(sdk, appSettings, RequestTopic)
	} else {
		log.Error("No request topic found!")
	}

	return topic
}

// GetResponseTopicList return list of response topic
func GetResponseTopicList(sdk *appsdk.AppFunctionsSDK) []string {
	log := sdk.LoggingClient
	appSettings := sdk.ApplicationSettings()

	var topics string
	if appSettings != nil {
		topics = getAppSetting(sdk, appSettings, ResponseTopic)
	} else {
		log.Error("No response topic found!")
	}

	topicLists := strings.Split(topics, ", ")

	return topicLists
}

// GetMqttQos return qos of mqtt report
func GetMqttQos(sdk *appsdk.AppFunctionsSDK) int {
	log := sdk.LoggingClient
	appSettings := sdk.ApplicationSettings()

	var MqttQos int
	if appSettings != nil {
		MqttQos, _ = strconv.Atoi(getAppSetting(sdk, appSettings, Qos))
	} else {
		log.Error("No application-specific settings found")
	}

	return MqttQos
}

func getAppSetting(sdk *appsdk.AppFunctionsSDK, setting map[string]string, name string) string {
	log := sdk.LoggingClient
	value, ok := setting[name]

	if ok {
		log.Debug(value)
		return value
	}
	log.Error(fmt.Sprintf("ApplicationName application setting %s not found", name))
	return ""
}
