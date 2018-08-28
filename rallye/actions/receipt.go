package actions

import (
	"fmt"
	"sync"
	"encoding/json"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/iteratec/sphero-rallye-server/log"
	"github.com/iteratec/sphero-rallye-server/mqtt"
	"github.com/iteratec/sphero-rallye-server/rallye/player"
	"errors"
	"github.com/iteratec/sphero-rallye-server/sphero"
	"github.com/iteratec/sphero-rallye-server/conf"
)

var (
	plannedActions map[string][]mqtt.SpheroAction
	actionMutex    *sync.Mutex
)

func init() {
	plannedActions = make(map[string][]mqtt.SpheroAction)
	for _, p := range conf.Players {
		plannedActions[p.Name] = []mqtt.SpheroAction{}
	}
	actionMutex = &sync.Mutex{}
}

func HandleIncomingSpheroActions() {

	client := mqtt.GetClient()
	players := conf.Players

	for _, p := range players {
		go handleIncomingSpheroActions(p, client)
	}

}

func handleIncomingSpheroActions(p player.RallyePlayer, client MQTT.Client) {
	topicsToSubscribe := map[string]byte{
		player.GetTopicName(p, mqtt.PLANNED_ACTIONS): byte(mqtt.AT_LEAST_ONCE),
		player.GetTopicName(p, mqtt.AD_HOC_ACTION):   byte(mqtt.AT_LEAST_ONCE),
	}

	msgQueue := make(chan MQTT.Message)
	client.SubscribeMultiple(topicsToSubscribe, func(i MQTT.Client, msg MQTT.Message) {
		msgQueue <- msg
	})
	for {
		incomingMsg := <-msgQueue
		topicName := string(incomingMsg.Topic())
		log.Info.Printf("TOPIC '%s' -> received MESSAGE: %s\n", topicName, string(incomingMsg.Payload()))

		switch topicName {
		case player.GetOtherTopicOfSameplayer(topicName, mqtt.PLANNED_ACTIONS):
			addPlannedActions(client, incomingMsg)
		case player.GetOtherTopicOfSameplayer(topicName, mqtt.AD_HOC_ACTION):
			runAdHocAction(client, incomingMsg)
		}
	}
}

func runAdHocAction(client MQTT.Client, incomingMsg MQTT.Message) {

	action := mqtt.SpheroAction{}
	err := json.Unmarshal(incomingMsg.Payload(), &action)

	if err != nil {
		publishError(incomingMsg, err, fmt.Sprintf("Die folgende Aktion hat das falsche Format: %s", string(incomingMsg.Payload())))
	} else {
		playerName := player.GetPlayerName(string(incomingMsg.Topic()))
		log.Info.Printf("Run the following action for player '%s': %v\n", playerName, action)
		go sphero.RunAction(playerName, action)
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
		if valid, validationMsg := ValidateActions(actions, playerName); !valid {
			publishError(incomingMsg, errors.New(validationMsg), validationMsg)
		} else {
			SetActions(actions, playerName)
		}
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

func SetActions(actions []mqtt.SpheroAction, playerName string) {
	actionMutex.Lock()
	plannedActions[playerName] = actions
	actionMutex.Unlock()
}
