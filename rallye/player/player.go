package player

import (
	"github.com/spf13/viper"
	"fmt"
	"strings"
	"github.com/iteratec/sphero-rallye-server/mqtt"
)

type RallyePlayer struct {
	Name        string `json:"name"`
	BluetoothId string `json:"bluetoothId"`
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
