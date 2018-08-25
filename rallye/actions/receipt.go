package actions

import (
	"fmt"
	"sync"
	"encoding/json"
	"github.com/spf13/viper"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/iteratec/sphero-rallye-server/log"
	"github.com/iteratec/sphero-rallye-server/mqtt"
	"github.com/iteratec/sphero-rallye-server/rallye/player"
	"errors"
)

var (
	plannedActions map[string][]mqtt.SpheroAction
	actionMutex    *sync.Mutex
)

func init() {
	plannedActions = make(map[string][]mqtt.SpheroAction)
	for _, p := range player.GetPlayers() {
		plannedActions[p.Name] = []mqtt.SpheroAction{}
	}
	actionMutex = &sync.Mutex{}
}

func HandleIncomingSpheroActions() {

	client := mqtt.GetClient()

	players := player.GetPlayers()
	log.Debug.Printf("Subscribing to topics for incoming planned actions of players '%v'\n", players)

	for _, p := range players {

		log.Debug.Printf("Subscribing to topic for incoming planned actions of player '%s'\n", p.Name)
		client.Subscribe(player.GetTopicName(p, mqtt.PLANNED_ACTIONS), byte(mqtt.AT_MOST_ONCE), addPlannedActions)

	}

}

func addPlannedActions(client MQTT.Client, incomingMsg MQTT.Message) {

	actions := []mqtt.SpheroAction{}
	err := json.Unmarshal(incomingMsg.Payload(), &actions)

	if err != nil {
		publishError(incomingMsg, err, fmt.Sprintf("Die folgende Aktion hat das falsche Format: %s", string(incomingMsg.Payload())))
	} else {
		playerName := player.GetPlayerName(string(incomingMsg.Topic()))
		log.Info.Printf("Set the following actions for player '%s': %v\n", playerName, actions)
		if valid, validationMsg := validate(actions); !valid {
			publishError(incomingMsg, errors.New(validationMsg), validationMsg)
		}
		setActions(actions, playerName)
	}
}

func publishError(incomingMsg MQTT.Message, err error, msgToPublish string) {
	errorTopicName := player.GetOtherTopicOfSameplayer(incomingMsg.Topic(), mqtt.ERRORS)
	log.Error.Printf("An error occurred unmarshalling SpheroAction:\nmqtt topic: %s\npayload: %s\nerr=%v\n-> Sending an error message to the following topic: %s",
		string(incomingMsg.Topic()), string(incomingMsg.Payload()), err, errorTopicName)
	mqtt.GetClient().Publish(
		errorTopicName,
		byte(mqtt.AT_LEAST_ONCE),
		false,
		msgToPublish)
}

func GetActions(playerName string) []mqtt.SpheroAction {
	actionMutex.Lock()
	actions := plannedActions[playerName]
	actionMutex.Unlock()
	return actions
}

func setActions(actions []mqtt.SpheroAction, playerName string) {
	actionMutex.Lock()
	plannedActions[playerName] = actions
	actionMutex.Unlock()
}

// Checks whether list of SpheroActions
// 	* contains correct amount of actions
func validate(actions []mqtt.SpheroAction) (bool, string) {
	validationMsg := ""
	numberOfMovesPerRound := viper.GetInt("rallye.numberOfMovesPerRound")
	valid := len(actions) == numberOfMovesPerRound
	if !valid {
		validationMsg += fmt.Sprintf("Die Liste muss genau %v Aktionen enthalten. Hier war(en) es: %v.", numberOfMovesPerRound, len(actions))
	}
	return valid, validationMsg
}
