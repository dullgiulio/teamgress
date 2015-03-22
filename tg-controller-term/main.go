package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
    "time"

	tg "github.com/dullgiulio/teamgress/libteamgress"
)

func printEvents(s *tg.Store) {
	eventsCh := make(chan tg.Event)
	filter := tg.GetFromTime(time.Now())
	listener := s.Subscribe(eventsCh, filter)

    interruptCh := make(chan os.Signal)
    signal.Notify(interruptCh, os.Interrupt)

	go func() {
		<-interruptCh
        s.Cancel(listener)
	    
        // Exit successfully.
        os.Exit(1)
    }()

    // Just print the received events back to stdout.
	for e := range eventsCh {
		fmt.Printf("%s\n", e.ToJSON())
	}
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

	// Read events from stdin.
    go tg.ReadJSONEvents(os.Stdin, evs)

	// Listen to events. Never returns.
	go store.Listen(evs)

    printEvents(store)
	os.Exit(1)
}
