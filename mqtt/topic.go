package mqtt

type TopicType int

func (topicType TopicType) GetName() string {
	names := [3]string{
		"possibleActionTypes",
		"plannedActions",
		"errors",
	}
	if topicType < POSSIBLE_ACTION_TYPES || topicType > ERRORS {
		return "unknown"
	}
	return names[topicType]
}

const (
	POSSIBLE_ACTION_TYPES TopicType = iota
	PLANNED_ACTIONS
	ERRORS
)