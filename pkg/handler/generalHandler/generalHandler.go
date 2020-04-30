package generalHandler

import "github.com/tuanldchainos/app-functions-sdk-go/appsdk"

type gHandler struct {
	sdk *appsdk.AppFunctionsSDK
}

func (f *gHandler) NewGeneralHandler(sdk *appsdk.AppFunctionsSDK) *gHandler {
	return &gHandler{
		sdk: sdk,
	}
}

func (f *gHandler) StartListeningGeneralIncoming() {
	f.onHandleMqttIncomming()
}

func (f *gHandler) onHandleMqttIncomming() {
	var success bool
	if success {
		f.gResponseSuccess()
	} else {
		f.gResponseFail()
	}
}

func (f *gHandler) gResponseSuccess() {
	return
}

func (f *gHandler) gResponseFail() {
	return
}
