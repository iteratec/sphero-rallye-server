package player

import (
	"github.com/spf13/viper"
	"fmt"
	"strings"
	"github.com/iteratec/sphero-rallye-server/log"
	"github.com/iteratec/sphero-rallye-server/mqtt"
)

var playersByTopicPrefix map[string]RallyePlayer

func init() {
	playersByTopicPrefix = make(map[string]RallyePlayer)
}

type RallyePlayer struct {
	Name        string `json:"name"`
	BluetoothId string `json:"bluetoothId"`
}

func GetPlayers() []RallyePlayer {
	var players []RallyePlayer
	err := viper.UnmarshalKey("rallye.players", &players)
	if err != nil {
		log.Error.Printf("unable to decode into struct, %v", err)
	}
	return players
}

func GetTopicName(player RallyePlayer, topicType mqtt.TopicType) string {
	return fmt.Sprintf("%s/%s/%s", viper.GetString("mqtt.topicPrefix"), player.Name, topicType.GetName())
}
func GetOtherTopicOfSameplayer(srcTopicName string, targetTopicType mqtt.TopicType) string {
	splitted := strings.Split(srcTopicName, "/")
	return fmt.Sprintf("%s/%s/%s",
		splitted[0], splitted[1], targetTopicType.GetName())
}
func GetPlayerName(topicName string) string {
	return strings.Split(topicName, "/")[1]
}
