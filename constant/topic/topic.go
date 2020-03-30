package topic

import (
	"fmt"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

var log logger.LoggingClient

func SetTopicList() (func(int) string, error) {

	topicLists := map[int]string{
		1: "topictest1",
		2: "topictest2",
	}

	return func(key int) string {
		return topicLists[key]
	}, nil
}

func GetNumOfTopic(topicList func(int) string) int {
	for i := 1; ; i++ {
		if topicList(i) == "" {
			return i - 1
		}
	}
}

func GetTopicList() func(int) string {
	TopicLists, err := SetTopicList()
	if err != nil {
		log.Error(fmt.Sprintln("Can not get list of topic!"))
	}
	return TopicLists
}
