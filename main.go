package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ActiveState/tail"
)

func tailFile(file string, evs chan<- event, conf *conf) {
	t, err := tail.TailFile(file, tail.Config{Follow: true})
	if err != nil {
		log.Fatal(err)
	}

	for line := range t.Lines {
		// TODO: Maybe only filter the events in the last N days?
		e := newEvent(line.Text, conf)
		evs <- *e
	}
}

func _test(s *store) {
	<-time.After(1 * time.Second)

	evs := make(chan event)
	//go s.getByUser("giotti", evs)
	go s.getFromTime(time.Now().AddDate(0, 0, -1), evs)

	for e := range evs {
		fmt.Printf("%s\n", e.JSON())
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Must specify which files to follow.")
	}

	conf := newConf()
	conf.loadFile(os.Args[1])

	files := os.Args[2:]

	store := newStore(conf)
	evs := make(chan event)

	for _, file := range files {
		go tailFile(file, evs, conf)
	}

	go _test(store)

	store.listen(evs)
	os.Exit(1)
}
