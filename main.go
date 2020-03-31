package main

import (
	"fmt"
	"os"

	"App-Mqtt-Go/pkg"

	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
)

const (
	serviceKey = "MqttExport"
)

func main() {
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

	pkg.StartListeningMqttIncoming(client, config)
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

func printDataToConsole(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {

	if len(params) < 1 {
		// We didn't receive a result
		return false, nil
	}

	fmt.Println(params[0].(string))

	// Leverage the built in logging service in EdgeX
	edgexcontext.LoggingClient.Debug("Printed to console")

	edgexcontext.Complete([]byte(params[0].(string)))
	return false, nil

}
