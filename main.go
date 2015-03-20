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

func importEvents(files []string, evs chan<- event, conf *conf) {
	for _, file := range files {
		go tailFile(file, evs, conf)
	}
}

// TODO: Implement this as a WebSocket handler.
func _test(s *store) {
	evs := make(chan event)
	filter := getFromTime(time.Now().AddDate(0, 0, -1))
	listener := s.stream(evs, filter)

	go func() {
		<-time.After(5 * time.Second)
		s.cancel(listener)
	}()

	for e := range evs {
		fmt.Printf("%s\n", e.JSON())
	}

	fmt.Printf("_test has exited\n")
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Must specify which files to follow.")
	}

	conf := newConf()
	conf.loadFile(os.Args[1])

	files := os.Args[2:]

	evs := make(chan event)

	// Will start a routine for each file.
	importEvents(files, evs, conf)

	// Start store handler.
	store := newStore(conf)
	// Remove listeners when the are cancelled.
	go store.handleCancelled()
	// Broadcast events to all listeners.
	go store.broadcast()

	// Example of handler.
	go _test(store)

	// Listen to events. Never returns.
	// XXX: This can go in the background; the main thread will handle the connections.
	store.listen(evs)

	os.Exit(1)
}
