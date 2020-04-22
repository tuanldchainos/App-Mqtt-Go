package pkg

// Mqtt params constant
const (
	MQTTHost       = "MQTTHost"
	MQTTUser       = "MQTTUser"
	MQTTPass       = "MQTTPass"
	MQTTPort       = "MQTTPort"
	Qos            = "Qos"
	KeepAlive      = "KeepAlive"
	MQTTCertData   = "MQTTCertData"
	MQTTKeyData    = "MQTTKeyData"
	SkipCertVerify = "SkipCertVerify"
	PersistOnError = "PersistOnError"
	RequestTopic   = "RequestTopic"
	ResponseTopic  = "ResponseTopic"
	MQTTCertDir    = "./certs/user.client.crt"
	MQTTKeyDir     = "./certs/user.client.key"
)

// Edgex constant
const (
	LoggingClientName       = "Logging"
	CoreCommandClientName   = "Command"
	CoreDataClientName      = "CoreData"
	NotificationsClientName = "Notifications"
	MetadataClientName      = "Metadata"
	SchedulerClientName     = "Scheduler"
)

// Edgex constant
const (
	CoreCommandServiceKey          = "edgex-core-command"
	CoreDataServiceKey             = "edgex-core-data"
	CoreMetaDataServiceKey         = "edgex-core-metadata"
	SupportLoggingServiceKey       = "edgex-support-logging"
	SupportNotificationsServiceKey = "edgex-support-notifications"
	SupportSchedulerServiceKey     = "edgex-support-scheduler"
)
