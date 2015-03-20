package main

import (
	"sync"
	"time"

	tg "github.com/dullgiulio/teamgress/libteamgress"
)

type listener struct {
	accept filter
	ch     chan<- tg.Event
}

type store struct {
	events      []tg.Event
	conf        *tg.Conf
	mux         *sync.Mutex
	listeners   map[*listener]struct{}
	listenersCh chan *listener
	eventsCh    chan tg.Event
	timeout     time.Duration
}

func newStore(conf *tg.Conf) *store {
	s := &store{
		mux:         &sync.Mutex{},
		conf:        conf,
		events:      make([]tg.Event, 0),
		listeners:   make(map[*listener]struct{}),
		listenersCh: make(chan *listener),
		eventsCh:    make(chan tg.Event, 5), // Can use some buffering here.
		timeout:     time.Millisecond * 500,
	}

	return s
}

func (s *store) cancel(l *listener) {
	s.listenersCh <- l
}

func (s *store) broadcast() {
	for e := range s.eventsCh {
		s.mux.Lock()

		for l, _ := range s.listeners {
			if l.accept(e) {
				// A client can only listen to a s.timeout periond of time
				// or it will be skipped. The storage will be locked for
				// s.timeout * len(s.listeners) at max.
				select {
				case l.ch <- e:
				case <-time.After(s.timeout):
				}
			}
		}

		s.mux.Unlock()
	}
}

func (s *store) handleCancelled() {
	for l := range s.listenersCh {
		s.mux.Lock()
		delete(s.listeners, l)
		s.mux.Unlock()

		close(l.ch)
	}
}

func (s *store) listen(evs <-chan tg.Event) {
	for e := range evs {
		// Copy the event in the store
		s.mux.Lock()
		s.events = append(s.events, e)
		s.mux.Unlock()

		s.eventsCh <- e
	}
}

func (s *store) subscribe(evs chan<- tg.Event, accept filter) *listener {
	l := &listener{
		ch:     evs,
		accept: accept,
	}

	go func() {
		s.mux.Lock()
		events := make([]tg.Event, len(s.events))
		copy(events, s.events)
		s.mux.Unlock()

		for _, e := range events {
			if accept(e) {
				// A client can only listen to a s.timeout periond of time
				// or it will be skipped.
				select {
				case evs <- e:
				case <-time.After(s.timeout):
				}
			}
		}

		s.mux.Lock()
		s.listeners[l] = struct{}{}
		s.mux.Unlock()

	}()

	return l
}

type filter func(tg.Event) bool

func getByUser(user string) filter {
	return func(e tg.Event) bool {
		return e.User.UnixName == user
	}
}

func getFromTime(time time.Time) filter {
	unixTime := time.Unix()

	return func(e tg.Event) bool {
		return e.Time.Unix() >= unixTime
	}
}
