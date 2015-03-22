package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ActiveState/tail"
	tg "github.com/dullgiulio/teamgress/libteamgress"
)

func makeEvent(str string, conf *tg.Conf) *tg.Event {
	e := tg.NewEvent()
	split := strings.SplitN(str, " ", 6)

	e.ID = split[0]

	if unixTime, err := strconv.ParseInt(split[1], 10, 0); err == nil {
		e.Time = time.Unix(int64(unixTime), 0)
	}

	e.User = conf.GetByUsername(split[2])
	env := strings.SplitN(split[3], ".", 2)
	e.Project = tg.MakeProject(env[0], "deploy-log", env[1])
	e.Type = split[4]
	e.Data = split[5]
	e.Level = "info"

	return e
}

func tailFile(file string, evs chan<- tg.Event, conf *tg.Conf) {
	t, err := tail.TailFile(file, tail.Config{Follow: true})
	if err != nil {
		log.Fatal(err)
	}

	for line := range t.Lines {
		// TODO: Maybe only filter the events in the last N days?
		e := makeEvent(line.Text, conf)
		evs <- *e
	}
}

func importEvents(files []string, evs chan tg.Event, conf *tg.Conf) {
	for _, file := range files {
		go tailFile(file, evs, conf)
	}

	for e := range evs {
		e.Emit()
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Must specify a configuration file and which files to follow.")
	}

	conf := tg.NewConf()
	conf.LoadFile(os.Args[1])

	files := os.Args[2:]

	evs := make(chan tg.Event)

	// Will start a routine for each file.
	importEvents(files, evs, conf)

	os.Exit(1)
}
