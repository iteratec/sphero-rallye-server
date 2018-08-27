package actions

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/iteratec/sphero-rallye-server/mqtt"
	"github.com/iteratec/sphero-rallye-server/sphero"
)

// Checks whether given list of SpheroActions
// 	* Contains correct amount of allowed actions
//	* Contains just allowed amount of ActionTypes
//	* All actions in the list contain correct configuration attributes
func ValidateActions(actions []mqtt.SpheroAction, playerName string) (bool, string) {
	valid := true
	validationMsg := ""
	allowedActionsPerRound := viper.GetInt("rallye.allowedActionsPerRound")
	correctNumberOfActions := len(actions) == allowedActionsPerRound
	if !correctNumberOfActions {
		valid = false
		validationMsg += fmt.Sprintf("Die Liste muss genau %v Aktionen enthalten. Hier war(en) es: %v.\n", allowedActionsPerRound, len(actions))
	}
	if validActionTypes, validActionTypesMsg := validateActionTypes(actions, playerName); !validActionTypes {
		valid = false
		validationMsg += validActionTypesMsg
	}
	for _, action := range actions {
		if singleValid, singleValidMsg := validateSingleAction(action); !singleValid {
			valid = false
			validationMsg += singleValidMsg
		}

	}
	return valid, validationMsg
}

func validateActionTypes(actions []mqtt.SpheroAction, playerName string) (bool, string) {
	valid := true
	validationMsg := ""

	testableActions := make([]mqtt.ActionTypeTestable, len(actions))
	for i, action := range actions {
		testableActions[i] = action
	}
	countGivenSetRgb := countActionTypes(mqtt.SET_RGB, testableActions)
	countGivenRoll := countActionTypes(mqtt.ROLL, testableActions)
	countGivenRotate := countActionTypes(mqtt.ROTATE, testableActions)

	allowedActionTypes := GetActionTypes(playerName)
	testableAllowedActionTypes := make([]mqtt.ActionTypeTestable, len(allowedActionTypes))
	for i, allowedType := range allowedActionTypes {
		testableAllowedActionTypes[i] = allowedType
	}
	countAllowedSetRgb := countActionTypes(mqtt.SET_RGB, testableAllowedActionTypes)
	countAllowedRoll := countActionTypes(mqtt.ROLL, testableAllowedActionTypes)
	countAllowedRotate := countActionTypes(mqtt.ROTATE, testableAllowedActionTypes)

	if countGivenRoll > countAllowedRoll {
		valid = false
		validationMsg += fmt.Sprintf("In dieser Runde sind nur %d ROLL-Aktionen erlaubt. %d sind zu viele!\n", countAllowedRoll, countGivenRoll)
	}
	if countGivenRotate > countAllowedRotate {
		valid = false
		validationMsg += fmt.Sprintf("In dieser Runde sind nur %d ROTATE-Aktionen erlaubt. %d sind zu viele!\n", countAllowedRotate, countGivenRotate)
	}
	if countGivenSetRgb > countAllowedSetRgb {
		valid = false
		validationMsg += fmt.Sprintf("In dieser Runde sind nur %d SET_RGB-Aktionen erlaubt. %d sind zu viele!\n", countAllowedSetRgb, countGivenSetRgb)
	}

	return valid, validationMsg
}

func countActionTypes(actionType mqtt.ActionType, listToCountIn []mqtt.ActionTypeTestable) int {
	counter := 0
	for _, testable := range listToCountIn {
		if testable.IsType(actionType) {
			counter++
		}
	}
	return counter
}

func validateSingleAction(action mqtt.SpheroAction) (bool, string) {
	valid := true
	validationMsg := ""
	switch action.ActionType {
	case mqtt.SET_RGB:
		if red, ok := action.Config[sphero.ActionConfKey_Red]; !ok || red < 0 || red > 255 {
			valid = false
			validationMsg += fmt.Sprintf("Jede SET_RGB Action muss das Attribut %s mit einem Wert zwischen 0 und 255 in der Config haben.\n", sphero.ActionConfKey_Red)
		}
		if green, ok := action.Config[sphero.ActionConfKey_Green]; !ok || green < 0 || green > 255 {
			valid = false
			validationMsg += fmt.Sprintf("Jede SET_RGB Action muss das Attribut %s mit einem Wert zwischen 0 und 255 in der Config haben.\n", sphero.ActionConfKey_Green)
		}
		if blue, ok := action.Config[sphero.ActionConfKey_Blue]; !ok || blue < 0 || blue > 255 {
			valid = false
			validationMsg += fmt.Sprintf("Jede SET_RGB Action muss das Attribut %s mit einem Wert zwischen 0 und 255 in der Config haben.\n", sphero.ActionConfKey_Blue)
		}
	case mqtt.ROLL:
		if speed, ok := action.Config[sphero.ActionConfKey_Speed]; !ok || speed < 0 || speed > 255 {
			valid = false
			validationMsg += fmt.Sprintf("Jede SET_RGB Action muss das Attribut %s mit einem Wert zwischen 0 und 255 in der Config haben.\n", sphero.ActionConfKey_Speed)
		}
		if duration, ok := action.Config[sphero.ActionConfKey_DurationInSecs]; !ok || duration < 0 || duration > 3 {
			valid = false
			validationMsg += fmt.Sprintf("Jede SET_RGB Action muss das Attribut %s mit einem Wert zwischen 0 und 255 in der Config haben.\n", sphero.ActionConfKey_DurationInSecs)
		}
	case mqtt.ROTATE:
		if heading, ok := action.Config[sphero.ActionConfKey_Heading]; !ok || heading < 0 || heading > 360 {
			valid = false
			validationMsg += fmt.Sprintf("Jede SET_RGB Action muss das Attribut %s mit einem Wert zwischen 0 und 255 in der Config haben.\n", sphero.ActionConfKey_Heading)
		}
	}
	return valid, validationMsg
}
