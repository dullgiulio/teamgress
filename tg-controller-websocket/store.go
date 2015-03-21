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
	buckets     map[int64][]tg.Event
	conf        *tg.Conf
	mux         *sync.Mutex
	listeners   map[*listener]struct{}
	listenersCh chan *listener
	eventsCh    chan tg.Event
	timeout     time.Duration
	bucketSecs  int64
}

func newStore(conf *tg.Conf) *store {
	s := &store{
		mux:         &sync.Mutex{},
		conf:        conf,
		buckets:     make(map[int64][]tg.Event, 0),
		listeners:   make(map[*listener]struct{}),
		listenersCh: make(chan *listener),
		eventsCh:    make(chan tg.Event, 5), // Can use some buffering here.
		timeout:     time.Millisecond * 500, // TODO: From config
		bucketSecs:  10,                     // TODO: From config
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
			// The storage will be locked for s.timeout * len(s.listeners) at max.
			l.emitEvent(e, s.timeout)
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

func (s *store) _addToBucket(e tg.Event) {
	s.mux.Lock()
	defer s.mux.Unlock()

	eTime := e.Time.Unix()
	key := eTime - (eTime % s.bucketSecs)

	if _, found := s.buckets[key]; !found {
		s.buckets[key] = make([]tg.Event, 1)
	}

	s.buckets[key] = append(s.buckets[key], e)
}

func (s *store) _getBucket(key int64) (t []tg.Event) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if events, found := s.buckets[key]; found {
		t = make([]tg.Event, len(events))
		copy(t, events)
	}

	return t
}

func (s *store) listen(evs <-chan tg.Event) {
	for e := range evs {
		// Copy the event in the store
		s._addToBucket(e)

		s.eventsCh <- e
	}
}

func (s *store) _addListener(l *listener) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.listeners[l] = struct{}{}
}

func (s *store) _bucketsKeys() (keys []int64) {
	s.mux.Lock()
	defer s.mux.Unlock()

	keys = make([]int64, len(s.buckets))
	for k, _ := range s.buckets {
		keys = append(keys, k)
	}

	return keys
}

func (l *listener) emitEvent(event tg.Event, timeout time.Duration) {
	if l.accept(event) {
		// A client can only listen to a s.timeout periond of time
		// or it will be skipped. The storage will be locked for
		// s.timeout * len(s.listeners) at max.
		select {
		case l.ch <- event:
		case <-time.After(timeout):
		}
	}
}

func (s *store) startListener(l *listener) {
	bucketsKeys := s._bucketsKeys()

	for _, key := range bucketsKeys {
		events := s._getBucket(key)

		if events == nil {
			continue
		}

		for _, e := range events {
			l.emitEvent(e, s.timeout)
		}
	}

	s._addListener(l)
}

func (s *store) subscribe(evs chan<- tg.Event, accept filter) *listener {
	l := &listener{
		ch:     evs,
		accept: accept,
	}

	go s.startListener(l)

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
