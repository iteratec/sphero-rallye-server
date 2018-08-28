package conf

import (
	"github.com/spf13/viper"
	"fmt"
	"github.com/iteratec/sphero-rallye-server/log"
	"github.com/iteratec/sphero-rallye-server/rallye/player"
)

var Players []player.RallyePlayer

func Init() {
	viper.AddConfigPath("/etc/sphero-rallye-server")
	viper.AddConfigPath("./conf/")
	viper.SetConfigName("conf")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error.Printf("Fatal error config file: %s", err)
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	readInPlayers()
}

func readInPlayers() {
	err := viper.UnmarshalKey("rallye.Players", &Players)
	if err != nil {
		log.Error.Printf("unable to decode into struct, %v", err)
	}
}