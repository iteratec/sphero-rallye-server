package rallye

import (
	"github.com/iteratec/sphero-rallye-server/rallye/actions"
	"github.com/iteratec/sphero-rallye-server/mqtt"
	"time"
	"github.com/spf13/viper"
	"github.com/iteratec/sphero-rallye-server/log"
)

func InitGame() {
	actions.HandleIncomingSpheroActions()
	StartNextRound()
}
func SendNextRoundEnd() {
	client := mqtt.GetClient()
	topicPath := mqtt.GetRoundEndTopicPath()
	roundEndTime := getNextRoundEndTime().Format(time.RFC1123Z)
	log.Debug.Printf("Sending round end '%v' to topic '%s'", roundEndTime, topicPath)
	client.Publish(topicPath, byte(mqtt.AT_MOST_ONCE), false, roundEndTime)
}
func getNextRoundEndTime() time.Time {
	roundLengthInSeconds := viper.GetInt("rallye.roundLengthInSeconds")
	return time.Now().Add(time.Duration(roundLengthInSeconds) * time.Second)
}

