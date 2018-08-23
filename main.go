package main

import (
	"github.com/iteratec/sphero-rallye-server/conf"
	"github.com/iteratec/sphero-rallye-server/mqtt"
	"gobot.io/x/gobot/platforms/ble"
	"os"
	"gobot.io/x/gobot/platforms/sphero/sprkplus"
	"gobot.io/x/gobot"
	"time"
	"github.com/iteratec/sphero-rallye-server/rallye/moves"
	"sync"
	"os/signal"
	"syscall"
	"github.com/iteratec/sphero-rallye-server/log"
)

var wg sync.WaitGroup

func init() {
	conf.Init()
	mqtt.InitClient()
	moves.InitSchedules()
}

func main() {

	runUntilInterrupt()

}

func runUntilInterrupt() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		signal := <-sigs
		log.Error.Printf("Got os signal '%v' -> Exiting.\n", signal)
		mqtt.GetClient().Disconnect(250)
		os.Exit(1)
	}()

	wg.Add(1)
	wg.Wait()
}

func blink() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	sprk := sprkplus.NewDriver(bleAdaptor)

	work := func() {
		gobot.Every(1*time.Second, func() {
			r := uint8(gobot.Rand(255))
			g := uint8(gobot.Rand(255))
			b := uint8(gobot.Rand(255))
			sprk.SetRGB(r, g, b)
		})
	}

	robot := gobot.NewRobot("sprk",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{sprk},
		work,
	)

	robot.Start()
}
