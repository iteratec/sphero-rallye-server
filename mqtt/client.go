package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"github.com/iteratec/sphero-rallye-server/log"
)

type QualityOfService byte
const (
	AT_MOST_ONCE QualityOfService = iota
	AT_LEAST_ONCE
	EXACTLY_ONCE
	id = "SPHERO_RALLYE_SERVER"
)
var(
	broker string
	client MQTT.Client
	store string
)

func InitClient() {
	broker = viper.GetString("mqtt.broker")
	log.Info.Printf("Connecting to mqtt broker =%v\n", broker)
	store = ":memory:"
	client = MQTT.NewClient(getMqttOptions())
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func GetClient() MQTT.Client {
	return client
}

func getMqttOptions() *MQTT.ClientOptions {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(id)
	opts.SetKeepAlive(5)
	opts.SetCleanSession(false)
	if store != ":memory:" {
		opts.SetStore(MQTT.NewFileStore(store))
	}
	log.Info.Printf("Initialize mqtt client with: broker=%v\n", broker)
	return opts
}
