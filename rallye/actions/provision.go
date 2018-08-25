package actions

import (
	"gobot.io/x/gobot"
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/iteratec/sphero-rallye-server/log"
	"github.com/iteratec/sphero-rallye-server/mqtt"
	"github.com/iteratec/sphero-rallye-server/rallye/player"
)

var topicSuffixPossibleActionTypes string

func init() {
	topicSuffixPossibleActionTypes = "possibleActionTypes"
}

// Provides one random ActionType
func getRandomActionType() mqtt.ActionType {
	return mqtt.ActionType(gobot.Rand(3))
}

// Provides a list of Sphero ActionTypes for the next round and
// sends them to the players mqtt topics
func ProvideNextActionTypes() {

	for _, p := range player.GetPlayers() {
		go provideActionTypes(p)
	}

}

func provideActionTypes(p player.RallyePlayer) {
	client := mqtt.GetClient()
	var actionTypes [9]string
	for i := 0; i < viper.GetInt("rallye.numberOfMovesPerRound"); i++ {
		actionTypes[i] = getRandomActionType().String()
	}
	var jsonActionTypes, _ = json.Marshal(actionTypes)
	strActionTypes := string(jsonActionTypes)
	topicName := player.GetTopicName(p, mqtt.POSSIBLE_ACTION_TYPES)
	log.Debug.Printf("RallyePlayer '%v': nextActionTypes=%v (publishing to topic '%s')", p.Name, strActionTypes, topicName)
	client.Publish(topicName, byte(mqtt.AT_MOST_ONCE), false, strActionTypes)
}
