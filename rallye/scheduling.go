package rallye

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"github.com/iteratec/sphero-rallye-server/log"
	"github.com/iteratec/sphero-rallye-server/rallye/actions"
	"github.com/iteratec/sphero-rallye-server/sphero"
	"github.com/iteratec/sphero-rallye-server/conf"
	"github.com/iteratec/sphero-rallye-server/rallye/player"
)

func InitSchedules() {
	cron := cron.New()
	scheduleRoundEnd(cron)
	scheduleSpheroWakeups(cron)
	cron.Start()
}

func scheduleRoundEnd(cron *cron.Cron) {
	roundLengthCron := fmt.Sprintf("@every %ds", viper.GetInt("rallye.roundLengthInSeconds"))
	log.Info.Printf("Schedule provision of next actions now: %s", roundLengthCron)
	cron.AddFunc(roundLengthCron, handleRoundEnd)
}
func scheduleSpheroWakeups(cron *cron.Cron) {
	if !viper.GetBool("rallye.mutePlayerControl") {
		cron.AddFunc("@every 60s", wakeUpSpheros)
	}
}
func wakeUpSpheros() {
	for _, p := range conf.Players {
		log.Debug.Printf("Waking up bot of player %s", p.Name)
		sphero.Wakeup(p.Name)
	}
}

func handleRoundEnd() {

	if !viper.GetBool("rallye.mutePlayerControl") {
		runPlayerActions()
	}

	StartNextRound()

}

func StartNextRound() {
	if !viper.GetBool("rallye.muteRoundLogic") {
		actions.ProvideNextActionTypes()
		SendNextRoundEnd()
	}
}

func runPlayerActions() {
	for _, p := range conf.Players {
		actionsOfRound := actions.GetActions(p.Name)
		for _, action := range actionsOfRound {
			sphero.RunAction(p.Name, action)
		}
		actions.SetActions(nil, p.Name)
	}
}
