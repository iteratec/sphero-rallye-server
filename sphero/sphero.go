package sphero

import (
	"fmt"
	"time"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/sphero/sprkplus"
	"github.com/iteratec/sphero-rallye-server/rallye/player"
	"github.com/iteratec/sphero-rallye-server/mqtt"
	"github.com/iteratec/sphero-rallye-server/log"
	"github.com/iteratec/sphero-rallye-server/conf"
)

const (
	ActionConfKey_Red            = "red"
	ActionConfKey_Green          = "green"
	ActionConfKey_Blue           = "blue"
	ActionConfKey_Speed          = "speed"
	ActionConfKey_Heading        = "heading"
	ActionConfKey_DurationInSecs = "durationInSecs"
)

var (
	botsByPlayername map[string]*SpheroBot
)

func init() {
	botsByPlayername = make(map[string]*SpheroBot)
}

func InitSpheros() {

	players := conf.Players
	log.Debug.Printf("Players for which spheros get initialized now: %v", players)
	for _, p := range players {
		adaptor := ble.NewClientAdaptor(p.BluetoothId)
		bot := &SpheroBot{
			RallyePlayer: p,
			BleAdaptor:   adaptor,
			SprkDriver:   sprkplus.NewDriver(adaptor),
			ActionChan:   make(chan mqtt.SpheroAction),
			Heading:      0,
		}
		go bot.awaitActions()
		botsByPlayername[p.Name] = bot
	}
	log.Debug.Printf("botsByPlayername=%v", botsByPlayername)

}

type SpheroBot struct {
	RallyePlayer player.RallyePlayer
	BleAdaptor   *ble.ClientAdaptor
	SprkDriver   *sprkplus.SPRKPlusDriver
	ActionChan   chan mqtt.SpheroAction
	Heading      uint16
}

func (sb *SpheroBot) awaitActions() {
	log.Debug.Printf("awaitActions: start for player %s", sb.RallyePlayer.Name)
	work := func() {
		for {
			log.Info.Printf("RallyePlayer %s waits for the next action now.", sb.RallyePlayer.Name)
			nextAction := <-sb.ActionChan
			log.Info.Printf("RallyePlayer %s will execute action %v now.", sb.RallyePlayer.Name, nextAction)
			switch nextAction.ActionType {
			case mqtt.SET_RGB:
				sb.setColor(nextAction.Config)
			case mqtt.ROTATE:
				sb.rotate(nextAction.Config)
			case mqtt.ROLL:
				sb.roll(nextAction.Config)
			}
		}
	}
	log.Debug.Printf("awaitActions: work defined for player %s", sb.RallyePlayer.Name)
	robot := gobot.NewRobot(fmt.Sprintf("sprk_robot_%s", sb.RallyePlayer.Name),
		[]gobot.Connection{sb.BleAdaptor},
		[]gobot.Device{sb.SprkDriver},
		work,
	)
	log.Debug.Printf("awaitActions: starting robot now for player %s", sb.RallyePlayer.Name)
	robot.Start()
}

func (sb *SpheroBot) setColor(config map[string]uint16) {
	sb.SprkDriver.SetRGB(uint8(config[ActionConfKey_Red]), uint8(config[ActionConfKey_Green]), uint8(config[ActionConfKey_Blue]))
}
func (sb *SpheroBot) rotate(config map[string]uint16) {
	heading := config[ActionConfKey_Heading]
	sb.Heading = heading
	ticker := gobot.Every(1*time.Second, func() {
		sb.SprkDriver.Roll(0, heading)
	})
	gobot.After(3*time.Second, func() {
		ticker.Stop()
	})
}
func (sb *SpheroBot) roll(config map[string]uint16) {
	ticker := gobot.Every(1*time.Second, func() {
		sb.SprkDriver.Roll(uint8(config[ActionConfKey_Speed]), sb.Heading)
	})
	gobot.After(time.Duration(config[ActionConfKey_DurationInSecs])*time.Second, func() {
		ticker.Stop()
	})
func (sb *SpheroBot) Wakeup() {
	sb.SprkDriver.Wake()
}

func RunAction(playerName string, action mqtt.SpheroAction) {
	log.Debug.Printf("Adding the following action to channel of player %s: %v", playerName, action.String())
	botsByPlayername[playerName].ActionChan <- action
}
func Wakeup(playerName string) {
	botsByPlayername[playerName].Wakeup()
}
