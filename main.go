package main

import (
	"fmt"
	"os"

	"App-Mqtt-Go/pkg"

	"github.com/tuanldchainos/app-functions-sdk-go/appsdk"
)

const (
	serviceKey = "MqttExport"
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

	config, err := pkg.LoadMqttConfig(edgexSdk)
	if err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("Failed to load MQTT configurations: %v\n", err))
		os.Exit(-1)
	}

	client, err := pkg.CreateClient(config)
	if err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("Failed to create MQTT client: %v\n", err))
		os.Exit(-1)
	}

	go pkg.NewMqttHandle().StartListeningMqttIncoming(client, edgexSdk)
	err = edgexSdk.MakeItRun()
	if err != nil {
		edgexSdk.LoggingClient.Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	defer func() {
		if client.IsConnected() {
			client.Disconnect(5000)
		}
	}()

	os.Exit(0)
}
