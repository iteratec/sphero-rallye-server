package main

import (
	"os"
	"sync"
	"os/signal"
	"syscall"
	"github.com/iteratec/sphero-rallye-server/conf"
	"github.com/iteratec/sphero-rallye-server/log"
	"github.com/iteratec/sphero-rallye-server/mqtt"
	"github.com/iteratec/sphero-rallye-server/rallye"
)

var wg sync.WaitGroup

func init() {
	conf.Init()
	mqtt.InitClient()
	rallye.InitSchedules()
}

func main() {

	go rallye.InitGame()

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
