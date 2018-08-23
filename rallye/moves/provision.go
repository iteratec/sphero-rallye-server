package moves

import (
	"github.com/iteratec/sphero-rallye-server/log"
	"github.com/iteratec/sphero-rallye-server/mqtt"
	"gobot.io/x/gobot"
	"encoding/json"
	"github.com/iteratec/sphero-rallye-server/conf"
	"github.com/spf13/viper"
)

// Provides one random ActionType
func getRandomActionType() mqtt.ActionType {
	return mqtt.ActionType(gobot.Rand(3))
}

// Provides a list of Sphero ActionTypes for the next round and
// sends them to the players mqtt topics
func ProvideNextActionTypes() {

	for _, p := range conf.GetPlayers() {
		go provideActionTypes(p)
	}

}

func provideActionTypes(player conf.RallyePlayer) {
	client := mqtt.GetClient()
	var actionTypes [9]string
	for i := 0; i < viper.GetInt("rallye.numberOfMovesPerRound"); i++ {
		actionTypes[i] = getRandomActionType().String()
	}
	var jsonActionTypes, _ = json.Marshal(actionTypes)
	strActionTypes := string(jsonActionTypes)
	log.Info.Printf("RallyePlayer '%v': nextActionTypes=%v (publishing to topic '%s')", player.Name, strActionTypes, player.TopicPossibleActionTypes)
	client.Publish(player.TopicPossibleActionTypes, byte(mqtt.AT_MOST_ONCE), false, strActionTypes)
}
