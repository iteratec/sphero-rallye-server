package mqtt

import (
	"github.com/spf13/viper"
	"fmt"
)

type TopicType int

func (topicType TopicType) GetName() string {
	names := [4]string{
		"possibleActionTypes",
		"plannedActions",
		"errors",
		"roundEnd",
	}
	if topicType < POSSIBLE_ACTION_TYPES || topicType > ROUND_END {
		return "unknown"
	}
	return names[topicType]
}

const (
	POSSIBLE_ACTION_TYPES TopicType = iota
	PLANNED_ACTIONS
	ERRORS
	ROUND_END
)

func GetRoundEndTopicPath() string {
	topicPrefix := viper.GetString("mqtt.topicPrefix")
	return fmt.Sprintf("%s/%s", topicPrefix, ROUND_END.GetName())
}
