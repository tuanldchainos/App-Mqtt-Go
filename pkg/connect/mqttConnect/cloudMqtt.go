package mqttConnect

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/tuanldchainos/app-functions-sdk-go/appsdk"
)

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

type MqttConnect struct {
	config *MqttConfig
	sdk    *appsdk.AppFunctionsSDK
}

func NewMqttConnect(sdk *appsdk.AppFunctionsSDK) *MqttConnect {
	config := new(MqttConfig)
	return &MqttConnect{
		config: config,
		sdk:    sdk,
	}
}

// LoadMqttConfig create config to mqtt connect
func (f *MqttConnect) LoadMqttConfig() error {
	log := f.sdk.LoggingClient
	if f.sdk == nil {
		log.Error("Invalid AppFunctionsSDK")
		return errors.New("Invalid AppFunctionsSDK")
	}

	var MqttHost, MqttPort, MqttUser, MqttPass, MqttCertData, MqttKeyData, MqttQos, MqttKeepAlive string
	var skipCertVerify bool
	// var persistOnError bool
	var errSkip, errPersist error

	appSettings := f.sdk.ApplicationSettings()
	if appSettings != nil {
		MqttUser = f.getAppSetting(appSettings, MQTTUser)
		MqttHost = f.getAppSetting(appSettings, MQTTHost)
		MqttPort = f.getAppSetting(appSettings, MQTTPort)
		MqttPass = f.getAppSetting(appSettings, MQTTPass)
		MqttCertData = f.getAppSetting(appSettings, MQTTCertData)
		MqttKeyData = f.getAppSetting(appSettings, MQTTKeyData)
		MqttQos = f.getAppSetting(appSettings, Qos)
		MqttKeepAlive = f.getAppSetting(appSettings, KeepAlive)
		skipCertVerify, errSkip = strconv.ParseBool(f.getAppSetting(appSettings, SkipCertVerify))
		// persistOnError, errPersist = strconv.ParseBool(f.getAppSetting(appSettings, PersistOnError))
		_, errPersist = strconv.ParseBool(f.getAppSetting(appSettings, PersistOnError))

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

		err = writeDataToKeyFile(MqttKeyData)
		if err != nil {
			log.Error(fmt.Sprintln("error while writing data to key file ", err))
		}
	} else {
		log.Error("No application-specific settings found")
	}

	f.config.MqttUser = MqttUser
	f.config.MqttPassword = MqttPass
	f.config.MqttQos, _ = strconv.Atoi(MqttQos)
	f.config.MqttKeepAlive, _ = strconv.Atoi(MqttKeepAlive)
	f.config.MqttHost = MqttHost
	f.config.MqttPort = MqttPort
	f.config.MqttCertFile = MQTTCertDir
	f.config.MqttKeyFile = MQTTKeyDir
	//config.PersistOnError = persistOnError

	if isSkipCertVerify(skipCertVerify) {
		f.config.MqttScheme = "tcp"
	} else {
		f.config.MqttScheme = "tls"
	}
	return nil
}

// CreateClient return a client, that connect to mqtt cloud successfully
func (f *MqttConnect) CreateClient() (MQTT.Client, error) {
	log := f.sdk.LoggingClient
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s:%s", f.config.MqttScheme, f.config.MqttHost, f.config.MqttPort))
	opts.SetUsername(f.config.MqttUser)
	opts.SetPassword(f.config.MqttPassword)
	opts.SetKeepAlive(time.Second * time.Duration(f.config.MqttKeepAlive))
	opts.SetAutoReconnect(true)

	if f.config.MqttScheme == "tls" {
		cert, err := tls.LoadX509KeyPair(f.config.MqttCertFile, f.config.MqttKeyFile)
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
				if f.config.PersistOnError {
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

func (f *MqttConnect) getAppSetting(setting map[string]string, name string) string {
	log := f.sdk.LoggingClient
	value, ok := setting[name]

	if ok {
		log.Debug(value)
		return value
	}
	log.Error(fmt.Sprintf("ApplicationName application setting %s not found", name))
	return ""
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
