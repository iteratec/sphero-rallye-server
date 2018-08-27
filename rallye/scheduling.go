package rallye

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"github.com/iteratec/sphero-rallye-server/log"
	"github.com/iteratec/sphero-rallye-server/rallye/player"
	"github.com/iteratec/sphero-rallye-server/rallye/actions"
	"github.com/iteratec/sphero-rallye-server/sphero"
)

func InitSchedules() {
	cron := cron.New()
	scheduleRoundEnd(cron)
	cron.Start()
}

func scheduleRoundEnd(cron *cron.Cron) {
	roundLengthCron := fmt.Sprintf("@every %ds", viper.GetInt("rallye.roundLengthInSeconds"))
	log.Info.Printf("Schedule provision of next actions now: %s", roundLengthCron)
	cron.AddFunc(roundLengthCron, handleRoundEnd)
}

func handleRoundEnd() {

	for _, p := range player.GetPlayers() {
		actionsOfRound := actions.GetActions(p.Name)
		for _, action := range actionsOfRound {
			sphero.RunAction(p.Name,action)
		}
		actions.SetActions(nil, p.Name)
	}

	actions.ProvideNextActionTypes()
	SendNextRoundEnd()

}
