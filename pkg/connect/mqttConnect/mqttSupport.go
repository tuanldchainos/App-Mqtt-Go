package mqttConnect

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/prometheus/common/log"
	"github.com/tuanldchainos/app-functions-sdk-go/appsdk"
)

func GetRequestTopic(sdk *appsdk.AppFunctionsSDK) string {
	log := sdk.LoggingClient
	appSettings := sdk.ApplicationSettings()

	var topic string
	if appSettings != nil {
		topic = getAppSetting(appSettings, RequestTopic)
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
		topics = getAppSetting(appSettings, ResponseTopic)
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
		MqttQos, _ = strconv.Atoi(getAppSetting(appSettings, Qos))
	} else {
		log.Error("No application-specific settings found")
	}

	return MqttQos
}

func getAppSetting(setting map[string]string, name string) string {
	value, ok := setting[name]

	if ok {
		fmt.Println(name + ":" + value)
		return value
	}
	log.Error(fmt.Sprintf("ApplicationName application setting %s not found", name))
	return ""
}
