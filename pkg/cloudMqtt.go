package pkg

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

const (
	MQTTHost           = "MQTTHost"
	MQTTUser           = "MQTTUser"
	MQTTPass           = "MQTTPass"
	MQTTPort           = "MQTTPort"
	Qos                = "Qos"
	KeepAlive          = "KeepAlive"
	CertFilename       = "MQTTCert"
	PrivateKeyFilename = "MQTTKey"
	SkipCertVerify     = "SkipCertVerify"
	PersistOnError     = "PersistOnError"
)

var log logger.LoggingClient

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

func getAppSetting(setting map[string]string, name string) string {
	value, ok := setting[name]

	if ok {
		log.Debug(value)
		return value
	}
	log.Error(fmt.Sprintf("ApplicationName application setting %s not found", name))
	return ""
}

func LoadMqttConfig(sdk *appsdk.AppFunctionsSDK) (*MqttConfig, error) {
	if sdk == nil {
		return nil, errors.New("Invalid AppFunctionsSDK")
	}

	log = sdk.LoggingClient

	var MqttHost, MqttPort, MqttUser, MqttPass, MqttCert, MqttKey, MqttQos, MqttKeepAlive string
	var skipCertVerify bool
	// var persistOnError bool
	var errSkip, errPersist error

	appSettings := sdk.ApplicationSettings()
	if appSettings != nil {
		MqttUser = getAppSetting(appSettings, MQTTUser)
		MqttHost = getAppSetting(appSettings, MQTTHost)
		MqttPort = getAppSetting(appSettings, MQTTPort)
		MqttPass = getAppSetting(appSettings, MQTTPass)
		MqttCert = getAppSetting(appSettings, CertFilename)
		MqttKey = getAppSetting(appSettings, PrivateKeyFilename)
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
	//config.PersistOnError = persistOnError

	if isSkipCertVerify(skipCertVerify) {
		config.MqttScheme = "tcp"
		config.MqttCertFile = ""
		config.MqttKeyFile = ""
		return config, nil
	}
	config.MqttScheme = "tls"
	config.MqttCertFile = MqttCert
	config.MqttKeyFile = MqttKey
	return config, nil
}

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
