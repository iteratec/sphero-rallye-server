package sphero

import (
	"github.com/iteratec/sphero-rallye-server/rallye/player"
	"gobot.io/x/gobot/platforms/ble"
	"os"
	"gobot.io/x/gobot/platforms/sphero/sprkplus"
	"gobot.io/x/gobot"
	"time"
	"fmt"
	"github.com/iteratec/sphero-rallye-server/mqtt"
	"github.com/iteratec/sphero-rallye-server/log"
)

var (
	actionsToExecuteByPlayername map[string]chan mqtt.SpheroAction
)

func init() {
	actionsToExecuteByPlayername = make(map[string]chan mqtt.SpheroAction)
}

func InitSpheros(){

	for _,p := range player.GetPlayers(){
		actionsToExecuteByPlayername[p.Name] = make(chan mqtt.SpheroAction)
		bluetoothId := p.BluetoothId
		bleAdaptor := ble.NewClientAdaptor(bluetoothId)
		sprk := sprkplus.NewDriver(bleAdaptor)
		work := func() {
			for {
				nextMove := <-actionsToExecuteByPlayername[p.Name]
				log.Info.Printf("Player %s will execute action %v now.", p.Name, nextMove)
			}
			//gobot.Every(1*time.Second, func() {
			//	sprk.Roll(50, 0)
			//})
			//gobot.After(3*time.Second, func() {
			//	sprk.Stop()
			//})
		}
		robot := gobot.NewRobot(fmt.Sprintf("sprk_robot_%s", p.Name),
			[]gobot.Connection{bleAdaptor},
			[]gobot.Device{sprk},
			work,
		)
		robot.Start()
	}
}

func MovePlayerSpheros()  {
	//for _,p := range player.GetPlayers(){
	//}
}
