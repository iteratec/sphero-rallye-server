package mqtt

type ActionType int

func (at ActionType) String() string {
	names := [3]string{
		"SET_COLOR",
		"MOVE",
		"TURN_AROUND",
	}
	if at < SET_COLOR || at > TURN_AROUND {
		return "unknown"
	}
	return names[at]
}

const (
	SET_COLOR   ActionType = iota
	MOVE
	TURN_AROUND
)

type (
	SpheroAction struct {
		actionType ActionType        `json:"actionType"`
		config     map[string]uint16 `json:"config,omitempty"`
	}
)
