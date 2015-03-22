package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tg "github.com/dullgiulio/teamgress/libteamgress"
)

// TODO: Implement this as a WebSocket handler.
func _test(s *tg.Store) {
	fmt.Printf("Start receiving events\n")

	evs := make(chan tg.Event)
	filter := tg.GetFromTime(time.Now().AddDate(0, 0, -1))
	listener := s.Subscribe(evs, filter)

	go func() {
		<-time.After(5 * time.Second)
		s.Cancel(listener)
	}()

	for e := range evs {
		fmt.Printf("%s\n", e.ToJSON())
	}

	fmt.Printf("Finished receiving events\n")
}

func main() {
	if len(os.Args) < 1 {
		log.Fatal("Must specify a configuration file.")
	}

	conf := tg.NewConf()
	conf.LoadFile(os.Args[1])

	evs := make(chan tg.Event)

	// Start store handler.
	store := tg.NewStore(conf)

	go tg.ReadJSONEvents(os.Stdin, evs)

	// Example of handler.
	go func() {
		for {
			_test(store)
			<-time.After(2 * time.Second)
		}
	}()

	// Listen to events. Never returns.
	// XXX: This can go in the background; the main thread will handle the connections.
	store.Listen(evs)

	os.Exit(1)
}
