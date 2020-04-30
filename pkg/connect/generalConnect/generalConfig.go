package generalConfig

import (
	"github.com/tuanldchainos/app-functions-sdk-go/appsdk"
)

type gConfig struct{}

type gConnect struct {
	config *gConfig
	sdk    *appsdk.AppFunctionsSDK
}

func (f *gConnect) NewGeneralConnect(sdk *appsdk.AppFunctionsSDK) *gConnect {
	gConfig := new(gConfig)
	return &gConnect{
		config: gConfig,
		sdk:    sdk,
	}
}

func (f *gConnect) LoadGeneralConfig() error {
	return nil
}

func (f *gConnect) CreateGeneralClient() (interface{}, error) {
	var gClient interface{}
	return gClient, nil
}
