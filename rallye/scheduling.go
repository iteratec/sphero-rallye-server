package rallye

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"github.com/iteratec/sphero-rallye-server/log"
	"github.com/iteratec/sphero-rallye-server/rallye/player"
	"github.com/iteratec/sphero-rallye-server/rallye/actions"
)

func InitSchedules() {
	cron := cron.New()
	scheduleRoundEnd(cron)
	cron.Start()
}

func scheduleRoundEnd(cron *cron.Cron) {
	roundLengthCron := fmt.Sprintf("@every %s", viper.GetString("rallye.roundLength"))
	log.Info.Printf("Schedule provision of next actions now: %s", roundLengthCron)
	cron.AddFunc(roundLengthCron, handleRoundEnd)
}

func handleRoundEnd() {

	//TODO: Run this round player actions
	for _, player := range player.GetPlayers() {
		log.Info.Printf("Starting actions for player %s: %v", player.Name, actions.GetActions(player.Name))
	}
	actions.ProvideNextActionTypes()

}
