package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ActiveState/tail"
)

func tailFile(file string, evs chan<- event) {
	t, err := tail.TailFile(file, tail.Config{Follow: true})
	if err != nil {
		log.Fatal(err)
	}

	for line := range t.Lines {
		// TODO: Maybe only filter the events in the last N days?
		evs <- makeEvent(line.Text)
	}
}

func _test(s *store) {
	<-time.After(1 * time.Second)

	fmt.Printf("Do query\n")

	evs := make(chan event)
	go s.getByUser("giotti", evs)

	for e := range evs {
		fmt.Printf("%s\n", &e)
	}
}

func main() {
	files := os.Args[1:]

	if len(files) == 0 {
		log.Fatal("Must specify which files to follow.")
	}

	store := newStore()
	evs := make(chan event)

	for _, file := range files {
		go tailFile(file, evs)
	}

	go _test(store)

	store.listen(evs)
	os.Exit(1)
}
