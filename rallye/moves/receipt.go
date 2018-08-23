package moves

import (
	"encoding/json"
	"github.com/iteratec/sphero-rallye-server/log"
	"github.com/iteratec/sphero-rallye-server/conf"
	"github.com/iteratec/sphero-rallye-server/mqtt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func HandleIncomingSpheroActions() {

	client := mqtt.GetClient()
	defer client.Disconnect(250)

	players := conf.GetPlayers()
	for _, p := range players {
		client.Subscribe(p.TopicPlannedActions, byte(mqtt.AT_MOST_ONCE), func(client MQTT.Client, msg MQTT.Message) {
			actions := []mqtt.SpheroAction{}
			err := json.Unmarshal(msg.Payload(), &actions)
			if err != nil {
				log.Error.Printf("An error occurred unmarshalling json from mqtt topic:\n\tjson=%v\n\terr=%v\n", string(msg.Payload()), err)
			} else {
				p.PlannedActions = append(p.PlannedActions, actions...)
			}
		})
	}

}