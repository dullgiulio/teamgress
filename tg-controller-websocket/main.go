package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tg "github.com/dullgiulio/teamgress/libteamgress"
)

// TODO: Implement this as a WebSocket handler.
func _test(s *store) {
	evs := make(chan tg.Event)
	filter := getFromTime(time.Now().AddDate(0, 0, -1))
	listener := s.stream(evs, filter)

	go func() {
		<-time.After(5 * time.Second)
		s.cancel(listener)
	}()

	for e := range evs {
		fmt.Printf("%s\n", e.ToJSON())
	}

	fmt.Printf("_test has exited\n")
}

// TODO: Read input and generate Events

func main() {
	if len(os.Args) < 1 {
		log.Fatal("Must specify a configuration file.")
	}

	conf := tg.NewConf()
	conf.LoadFile(os.Args[1])

	evs := make(chan tg.Event)

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
