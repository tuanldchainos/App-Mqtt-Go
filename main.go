package main

import (
	"fmt"
	"os"

	"App-Mqtt-Go/pkg"
	"App-Mqtt-Go/pkg/connect/mqttConnect"
	"App-Mqtt-Go/pkg/handler/mqttHandler"

	"github.com/tuanldchainos/app-functions-sdk-go/appsdk"
)

const (
	serviceKey = "AppService-mqtt-export"
)

func main() {
	os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")
	edgexSdk := &appsdk.AppFunctionsSDK{ServiceKey: serviceKey}
	if err := edgexSdk.Initialize(); err != nil {
		message := fmt.Sprintf("SDK initialization failed: %v\n", err)
		if edgexSdk.LoggingClient != nil {
			edgexSdk.LoggingClient.Error(message)
		} else {
			fmt.Println(message)
		}
		os.Exit(-1)
	}

	go pkg.SetServiceURIList(edgexSdk)
	go pkg.UpdateServiceURI(edgexSdk)

	MqttConnect := mqttConnect.NewMqttConnect(edgexSdk)

	err := MqttConnect.LoadMqttConfig()
	if err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("Failed to load MQTT configurations: %v\n", err))
		os.Exit(-1)
	}

	MqttClient, err := MqttConnect.CreateClient()
	if err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("Failed to create MQTT client: %v\n", err))
		os.Exit(-1)
	}

	go mqttHandler.NewMqttHandle(edgexSdk).StartListeningMqttIncoming(MqttClient)
	err = edgexSdk.MakeItRun()
	if err != nil {
		edgexSdk.LoggingClient.Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	defer func() {
		if MqttClient.IsConnected() {
			MqttClient.Disconnect(5000)
		}
	}()

	os.Exit(0)
}
