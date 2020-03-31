package topic

import (
	"fmt"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

var log logger.LoggingClient

func SetTopicList() (func(string) string, error) {

	topicLists := map[string]string{
		"Request":  "RequestToic",
		"Response": "ResponseTopic",
	}

	return func(key string) string {
		return topicLists[key]
	}, nil
}

func GetTopicList() func(string) string {
	TopicLists, err := SetTopicList()
	if err != nil {
		log.Error(fmt.Sprintln("Can not get list of topic!"))
	}
	return TopicLists
}
