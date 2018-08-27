package mqtt

import (
	"fmt"
	"bytes"
	"encoding/json"
)

type ActionType int

const (
	SET_RGB ActionType = iota
	ROLL
	ROTATE
)

var (
	actionType2Name = map[ActionType]string{
		SET_RGB: "SET_RGB",
		ROLL:    "ROLL",
		ROTATE:  "ROTATE",
	}
	name2actionType = map[string]ActionType{
		"SET_RGB": SET_RGB,
		"ROLL":    ROLL,
		"ROTATE":  ROTATE,
	}
)

func (at ActionType) String() string {
	return actionType2Name[at]
}
func (at *ActionType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(actionType2Name[*at])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (at *ActionType) UnmarshalJSON(b []byte) error {
	// unmarshal as string
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	// lookup value
	*at = name2actionType[s]
	return nil
}

func (at ActionType) IsType(actionType ActionType) bool {
	return at == actionType
}

type SpheroAction struct {
	ActionType ActionType        `json:"ActionType"`
	Config     map[string]uint16 `json:"Config,omitempty"`
}

func (action SpheroAction) IsType(actionType ActionType) bool {
	return action.ActionType == actionType
}
func (action SpheroAction) String() string {
	return fmt.Sprintf("type=%s|Config=", action.ActionType.String(), action.Config)
}

type ActionTypeTestable interface {
	IsType(actionType ActionType) bool
}
