package pkg

import (
	"fmt"

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

func setServiceUrlList(sdk *appsdk.AppFunctionsSDK) (func(string) string, error) {

	urlList := map[string]string{
		CoreDataClientName: sdk.GetServiceUrl(CoreDataClientName),
	}

	return func(key string) string {
		return urlList[key]
	}, nil
}

func GetServiceUrlList(sdk *appsdk.AppFunctionsSDK) func(string) string {
	urlList, err := setServiceUrlList(sdk)
	if err != nil {
		log.Error(fmt.Sprintln("Can not get list of topic!"))
	}
	return urlList
}
