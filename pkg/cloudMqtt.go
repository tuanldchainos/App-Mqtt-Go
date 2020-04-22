package pkg

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/tuanldchainos/app-functions-sdk-go/appsdk"
)

var log logger.LoggingClient

// MqttConfig struct is struct of params, required to connecting to mqtt cloud
type MqttConfig struct {
	MqttScheme    string
	MqttHost      string
	MqttPort      string
	MqttQos       int
	MqttUser      string
	MqttPassword  string
	MqttKeepAlive int
	MqttCertFile  string
	MqttKeyFile   string
	// MqttPersistOnError string
}

// LoadMqttConfig create config to mqtt connect
func LoadMqttConfig(sdk *appsdk.AppFunctionsSDK) (*MqttConfig, error) {
	if sdk == nil {
		return nil, errors.New("Invalid AppFunctionsSDK")
	}

	log = sdk.LoggingClient

	var MqttHost, MqttPort, MqttUser, MqttPass, MqttCertData, MqttKeyData, MqttQos, MqttKeepAlive string
	var skipCertVerify bool
	// var persistOnError bool
	var errSkip, errPersist error

	appSettings := sdk.ApplicationSettings()
	if appSettings != nil {
		MqttUser = getAppSetting(appSettings, MQTTUser)
		MqttHost = getAppSetting(appSettings, MQTTHost)
		MqttPort = getAppSetting(appSettings, MQTTPort)
		MqttPass = getAppSetting(appSettings, MQTTPass)
		MqttCertData = getAppSetting(appSettings, MQTTCertData)
		MqttKeyData = getAppSetting(appSettings, MQTTKeyData)
		fmt.Println(MqttKeyData)
		MqttQos = getAppSetting(appSettings, Qos)
		MqttKeepAlive = getAppSetting(appSettings, KeepAlive)
		skipCertVerify, errSkip = strconv.ParseBool(getAppSetting(appSettings, SkipCertVerify))
		// persistOnError, errPersist = strconv.ParseBool(getAppSetting(appSettings, PersistOnError))
		_, errPersist = strconv.ParseBool(getAppSetting(appSettings, PersistOnError))

		if errSkip != nil {
			log.Error("Unable to parse " + SkipCertVerify + " value")
		}
		if errPersist != nil {
			log.Error("Unable to parse " + PersistOnError + " value")
		}

		err := writeDataToCerFile(MqttCertData)
		if err != nil {
			log.Error(fmt.Sprintln("error while writing data to cer file ", err))
		}

		err = writeDataToKeyFile(MqttCertData)
		if err != nil {
			log.Error(fmt.Sprintln("error while writing data to key file ", err))
		}
	} else {
		log.Error("No application-specific settings found")
	}

	config := new(MqttConfig)

	config.MqttUser = MqttUser
	config.MqttPassword = MqttPass
	config.MqttQos, _ = strconv.Atoi(MqttQos)
	config.MqttKeepAlive, _ = strconv.Atoi(MqttKeepAlive)
	config.MqttHost = MqttHost
	config.MqttPort = MqttPort
	config.MqttCertFile = MQTTCertDir
	config.MqttKeyFile = MQTTKeyDir
	//config.PersistOnError = persistOnError

	if isSkipCertVerify(skipCertVerify) {
		config.MqttScheme = "tcp"
		return config, nil
	}
	config.MqttScheme = "tls"
	return config, nil
}

// CreateClient return a client, that connect to mqtt cloud successfully
func CreateClient(config *MqttConfig) (MQTT.Client, error) {
	log.Info(fmt.Sprintf("Create MQTT client and connection"))
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s:%s", config.MqttScheme, config.MqttHost, config.MqttPort))
	opts.SetUsername(config.MqttUser)
	opts.SetPassword(config.MqttPassword)
	opts.SetKeepAlive(time.Second * time.Duration(config.MqttKeepAlive))
	opts.SetAutoReconnect(true)

	if config.MqttScheme == "tls" {
		cert, err := tls.LoadX509KeyPair(config.MqttCertFile, config.MqttKeyFile)
		if err != nil {
			log.Error("Failed loading x509 data")
			return nil, errors.New("Can not create a mqtt client")
		}

		tlsConfig := &tls.Config{
			ClientCAs:    nil,
			Certificates: []tls.Certificate{cert},
		}
		opts.SetTLSConfig(tlsConfig)
	}

	client := MQTT.NewClient(opts)
	if !client.IsConnected() {
		token := client.Connect()
		if token.Wait() && token.Error() != nil {
			/*
				if config.PersistOnError {
					subMessage = "persisting Event for later retry"
				}
			*/
			log.Error(fmt.Sprintln(token.Error()))
			log.Info(fmt.Sprintf("Create MQTT client and connection succesful"))
			return client, token.Error()
		}
	}

	return client, nil
}

func isSkipCertVerify(SkipCertVerify bool) bool {
	return SkipCertVerify
}

func writeDataToCerFile(data string) error {
	dataWrite := "-----BEGIN CERTIFICATE-----\n" + data + "\n-----END CERTIFICATE-----"
	err := ioutil.WriteFile(MQTTCertDir, []byte(dataWrite), 0644)
	return err
}

func writeDataToKeyFile(data string) error {
	dataWrite := "-----BEGIN PRIVATE KEY-----\n" + data + "\n-----END PRIVATE KEY-----"
	err := ioutil.WriteFile(MQTTCertDir, []byte(dataWrite), 0644)
	return err
}
