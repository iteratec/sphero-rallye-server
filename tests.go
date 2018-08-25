package main

import (
	"os"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/sphero/sprkplus"
	"time"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	sprk := sprkplus.NewDriver(bleAdaptor)

	work := func() {
		gobot.Every(1*time.Second, func() {
			sprk.Roll(50, 0)
		})
		gobot.After(3*time.Second, func() {
			sprk.Stop()
		})

	}

	robot := gobot.NewRobot("sprk",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{sprk},
		work,
	)

	robot.Start()
	robot.Stop()
}