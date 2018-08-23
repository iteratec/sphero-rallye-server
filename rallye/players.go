package rallye

import (
	"github.com/iteratec/sphero-rallye-server/conf"
)

var players []conf.RallyePlayer

func InitPlayers() {
	players = conf.GetPlayers()
}
