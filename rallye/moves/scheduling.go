package moves

import (
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"fmt"
	"github.com/iteratec/sphero-rallye-server/log"
)

func InitSchedules() {
	cron := cron.New()
	scheduleRoundEnd(cron)
	cron.Start()
}

func scheduleRoundEnd(cron *cron.Cron) {
	roundLengthCron := fmt.Sprintf("@every %s", viper.GetString("rallye.roundLength"))
	log.Info.Printf("Schedule provision of next moves now: %s", roundLengthCron)
	cron.AddFunc(roundLengthCron, handleRoundEnd)
}

func handleRoundEnd(){

	//TODO: Run this round player moves

	ProvideNextActionTypes()

}
