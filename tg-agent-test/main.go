package main

import (
	"log"
	"os"
	"time"

	tg "github.com/dullgiulio/teamgress/libteamgress"
)

func makeEvent(conf *tg.Conf) *tg.Event {
	e := tg.NewEvent()
	e.ID = "deadbeef"
	e.Time = time.Now()
	e.User = *conf.Users[0]
	e.EnvName = "test-service"
	e.EnvStage = "staging"
	e.Type = "testing"
	e.Data = "Just a simple\ntest"
	e.Level = "info"

	return e
}

func emitEvents(conf *tg.Conf) {
	for {
		e := makeEvent(conf)
		e.Emit()

		<-time.After(1 * time.Second)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Must specify a configuration file.")
	}

	conf := tg.NewConf()
	conf.LoadFile(os.Args[1])

	emitEvents(conf)

	os.Exit(1)
}
