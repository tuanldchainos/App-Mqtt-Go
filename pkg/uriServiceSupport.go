package pkg

import (
	"time"

	"github.com/tuanldchainos/app-functions-sdk-go/appsdk"
)

var (
	coredataURL      string
	commandURL       string
	loggerURL        string
	notificationsURL string
	metadataURL      string
	shedulerURL      string
	agentURL         string
	UrlList          map[string]string
)

func SetServiceURIList(sdk *appsdk.AppFunctionsSDK) {
	coredataURL = sdk.GetServiceURLViaConfigFile(CoreDataClientName)

	commandURL = sdk.GetServiceURLViaConfigFile(CoreCommandClientName)

	loggerURL = sdk.GetServiceURLViaConfigFile(LoggingClientName)

	notificationsURL = sdk.GetServiceURLViaConfigFile(NotificationsClientName)

	metadataURL = sdk.GetServiceURLViaConfigFile(MetadataClientName)

	shedulerURL = sdk.GetServiceURLViaConfigFile(SchedulerClientName)

	agentURL = sdk.GetServiceURLViaConfigFile(SystemAgentClienName)

	UrlList = map[string]string{
		CoreDataServiceKey:              coredataURL,
		SupportSchedulerServiceKey:      shedulerURL,
		CoreCommandServiceKey:           commandURL,
		CoreMetaDataServiceKey:          metadataURL,
		SupportNotificationsServiceKey:  notificationsURL,
		SupportLoggingServiceKey:        loggerURL,
		SystemManagementAgentServiceKey: agentURL,
	}
}

func UpdateServiceURI(sdk *appsdk.AppFunctionsSDK) {
	for i := 0; ; i++ {
		time.Sleep(60 * time.Second)

		coredataURL, _ = sdk.GetServiceURLViaRegistry(CoreDataServiceKey)
		commandURL, _ = sdk.GetServiceURLViaRegistry(CoreCommandServiceKey)
		loggerURL, _ = sdk.GetServiceURLViaRegistry(SupportLoggingServiceKey)
		notificationsURL, _ = sdk.GetServiceURLViaRegistry(SupportNotificationsServiceKey)
		metadataURL, _ = sdk.GetServiceURLViaRegistry(CoreMetaDataServiceKey)
		shedulerURL, _ = sdk.GetServiceURLViaRegistry(SupportSchedulerServiceKey)
		agentURL, _ = sdk.GetServiceURLViaRegistry(SystemManagementAgentServiceKey)

		UrlList = map[string]string{
			CoreDataServiceKey:              coredataURL,
			SupportSchedulerServiceKey:      shedulerURL,
			CoreCommandServiceKey:           commandURL,
			CoreMetaDataServiceKey:          metadataURL,
			SupportNotificationsServiceKey:  notificationsURL,
			SupportLoggingServiceKey:        loggerURL,
			SystemManagementAgentServiceKey: agentURL,
		}
	}
}
