package conf

import (
	"github.com/spf13/viper"
	"fmt"
	"github.com/iteratec/sphero-rallye-server/log"
	"github.com/iteratec/sphero-rallye-server/mqtt"
)

func Init() {
	viper.AddConfigPath("/etc/rallye-rallye-server")
	viper.AddConfigPath("./conf/")
	viper.SetConfigName("conf")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error.Printf("Fatal error config file: %s", err)
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}

type RallyePlayer struct {
	Name                     string              `json:"name"`
	TopicPossibleActionTypes string              `json:"topicPossibleActionTypes"`
	PossibleActionTypes      []mqtt.ActionType   `json:"possibleActionTypes"`
	TopicPlannedActions      string              `json:"topicPlannedActions"`
	PlannedActions           []mqtt.SpheroAction `json:"plannedActions"`
}

func GetPlayers() []RallyePlayer {
	var players []RallyePlayer
	err := viper.UnmarshalKey("rallye.players", &players)
	if err != nil {
		log.Error.Printf("unable to decode into struct, %v", err)
	}
	return players
}
