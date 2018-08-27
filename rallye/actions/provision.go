package actions

import (
	"gobot.io/x/gobot"
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/iteratec/sphero-rallye-server/log"
	"github.com/iteratec/sphero-rallye-server/mqtt"
	"github.com/iteratec/sphero-rallye-server/rallye/player"
	"sync"
)

var (
	actionsTypes    map[string][]mqtt.ActionType
	actionTypeMutex *sync.Mutex
)

func init() {
	actionsTypes = make(map[string][]mqtt.ActionType)
	for _, p := range player.GetPlayers() {
		actionsTypes[p.Name] = []mqtt.ActionType{}
	}
	actionTypeMutex = &sync.Mutex{}
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
	providedActionTypesPerRound := viper.GetInt("rallye.providedActionTypesPerRound")
	actionTypes := make([]mqtt.ActionType, providedActionTypesPerRound)
	actionTypesStrings := make([]string, providedActionTypesPerRound)
	for i := 0; i < providedActionTypesPerRound; i++ {
		actionType := getRandomActionType()
		actionTypes[i] = actionType
		actionTypesStrings[i] = actionType.String()
	}
	var jsonActionTypes, _ = json.Marshal(actionTypesStrings)
	strActionTypes := string(jsonActionTypes)
	topicName := player.GetTopicName(p, mqtt.POSSIBLE_ACTION_TYPES)
	log.Debug.Printf("RallyePlayer '%v': nextActionTypes=%v (publishing to topic '%s')", p.Name, strActionTypes, topicName)
	client.Publish(topicName, byte(mqtt.AT_MOST_ONCE), false, strActionTypes)
	SetActionTypes(actionTypes, p.Name)
}

func GetActionTypes(playerName string) []mqtt.ActionType {
	actionTypeMutex.Lock()
	actionTypesOfPlayer :=  actionsTypes[playerName]
	actionTypeMutex.Unlock()
	return actionTypesOfPlayer
}
func SetActionTypes(actionTypes []mqtt.ActionType, playerName string)  {
	actionTypeMutex.Lock()
	actionsTypes[playerName] = actionTypes
	actionTypeMutex.Unlock()
}
