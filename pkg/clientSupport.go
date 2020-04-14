package pkg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tuanldchainos/app-functions-sdk-go/appsdk"
)

const (
	LoggingClientName       = "Logging"
	CoreCommandClientName   = "Command"
	CoreDataClientName      = "CoreData"
	NotificationsClientName = "Notifications"
	MetadataClientName      = "Metadata"
	SchedulerClientName     = "Scheduler"
)

const (
	CoreCommandServiceKey          = "edgex-core-command"
	CoreDataServiceKey             = "edgex-core-data"
	CoreMetaDataServiceKey         = "edgex-core-metadata"
	SupportLoggingServiceKey       = "edgex-support-logging"
	SupportNotificationsServiceKey = "edgex-support-notifications"
	SupportSchedulerServiceKey     = "edgex-support-scheduler"
)

// GetServiceURLList return list url of egdex service
func GetServiceURLList(sdk *appsdk.AppFunctionsSDK) func(string) string {
	urlList, err := setServiceURLList(sdk)
	if err != nil {
		log.Error(fmt.Sprintln("Can not get list of service url!"))
	}
	return urlList
}

// GetRequestTopic return request topic
func GetRequestTopic(sdk *appsdk.AppFunctionsSDK) string {
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
		log.Debug(value)
		return value
	}
	log.Error(fmt.Sprintf("ApplicationName application setting %s not found", name))
	return ""
}

func setServiceURLList(sdk *appsdk.AppFunctionsSDK) (func(string) string, error) {
	coredataURL, err := sdk.GetServiceURLViaRegistry(CoreDataServiceKey)
	if err != nil {
		coredataURL = sdk.GetServiceURLViaConfigFile(CoreDataClientName)
		log.Info(fmt.Sprintf("Can not take service url via registry, take url in config file: %s", err))
	}
	fmt.Println(coredataURL)
	commandURL, err := sdk.GetServiceURLViaRegistry(CoreCommandServiceKey)
	if err != nil {
		commandURL = sdk.GetServiceURLViaConfigFile(CoreCommandClientName)
		log.Info(fmt.Sprintf("Can not take service url via registry, take url in config file: %s", err))
	}

	loggerURL, err := sdk.GetServiceURLViaRegistry(SupportLoggingServiceKey)
	if err != nil {
		loggerURL = sdk.GetServiceURLViaConfigFile(LoggingClientName)
		log.Info(fmt.Sprintf("Can not take service url via registry, take url in config file: %s", err))
	}

	notificationsURL, err := sdk.GetServiceURLViaRegistry(SupportNotificationsServiceKey)
	if err != nil {
		notificationsURL = sdk.GetServiceURLViaConfigFile(NotificationsClientName)
		log.Info(fmt.Sprintf("Can not take service url via registry, take url in config file: %s", err))
	}

	metadataURL, err := sdk.GetServiceURLViaRegistry(CoreMetaDataServiceKey)
	if err != nil {
		metadataURL = sdk.GetServiceURLViaConfigFile(MetadataClientName)
		log.Info(fmt.Sprintf("Can not take service url via registry, take url in config file: %s", err))
	}

	shedulerURL, err := sdk.GetServiceURLViaRegistry(SupportSchedulerServiceKey)
	if err != nil {
		shedulerURL = sdk.GetServiceURLViaConfigFile(SchedulerClientName)
		log.Info(fmt.Sprintf("Can not take service url via registry, take url in config file: %s", err))
	}

	urlList := map[string]string{
		CoreDataClientName:      coredataURL,
		SchedulerClientName:     shedulerURL,
		CoreCommandClientName:   commandURL,
		MetadataClientName:      metadataURL,
		NotificationsClientName: notificationsURL,
		LoggingClientName:       loggerURL,
	}

	return func(key string) string {
		return urlList[key]
	}, nil
}
