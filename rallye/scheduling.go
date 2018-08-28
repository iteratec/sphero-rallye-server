package rallye

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"github.com/iteratec/sphero-rallye-server/log"
	"github.com/iteratec/sphero-rallye-server/rallye/actions"
	"github.com/iteratec/sphero-rallye-server/sphero"
	"github.com/iteratec/sphero-rallye-server/conf"
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

	for _, p := range conf.Players {
		actionsOfRound := actions.GetActions(p.Name)
		for _, action := range actionsOfRound {
			sphero.RunAction(p.Name,action)
		}
		actions.SetActions(nil, p.Name)
	}

	actions.ProvideNextActionTypes()
	SendNextRoundEnd()

}
