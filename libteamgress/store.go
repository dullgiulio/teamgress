package libteamgress

import (
	"time"
)

type Store struct {
	buckets          *buckets
	conf             *Conf
	listeners        map[*listener]struct{}
	cancelListenerCh chan *listener
	addListenerCh    chan *listener
	eventsCh         chan Event
	timeout          time.Duration
	bucketSecs       int64
	bucketMax        int
}

func NewStore(conf *Conf) *Store {
	s := &Store{
		conf:             conf,
		listeners:        make(map[*listener]struct{}),
		cancelListenerCh: make(chan *listener),
		addListenerCh:    make(chan *listener),
		eventsCh:         make(chan Event, 5),    // Can use some buffering here.
		timeout:          time.Millisecond * 500, // TODO: From config
	}

	// Take seconds to rotate buckets and buckets to keep from config.
	s.buckets = newBuckets(3600, 24)

	go s.handlerLoop()

	return s
}

// Get a stream of all events that match the accept filter.
func (s *Store) Subscribe(evs chan<- Event, accept Filter) *listener {
	l := newListener(evs, accept)

	go s.startListener(l)

	return l
}

// Import events in the store
func (s *Store) Listen(evs <-chan Event) {
	for e := range evs {
		// Copy the event in the Store
		s.buckets.add(e)

		s.eventsCh <- e
	}
}

// Cancel a listener (will close its channel)
func (s *Store) Cancel(l *listener) {
	s.cancelListenerCh <- l
}

// Main loop to handle all events
func (s *Store) handlerLoop() {
	for {
		select {
		case e := <-s.eventsCh:
			s.broadcast(e)
		case l := <-s.cancelListenerCh:
			delete(s.listeners, l)
			close(l.ch)
		case l := <-s.addListenerCh:
			s.listeners[l] = struct{}{}
		}
	}
}

// Broadcast events to all listeners.
func (s *Store) broadcast(e Event) {
	for l, _ := range s.listeners {
		// The storage will be locked for s.timeout * len(s.listeners) at max.
		l.emitEvent(e, s.timeout)
	}
}

// Add a listener to all incoming events.
func (s *Store) addListener(l *listener) {
	s.addListenerCh <- l
}

// A new listener will receive all old messages (filtered) and new ones.
func (s *Store) startListener(l *listener) {
	bucketsKeys := s.buckets.keys()

	// Copy one bucket at a time
	for _, key := range bucketsKeys {
		events := s.buckets.get(key)

		// This bucket might have been garbage collected
		// while we are sending.
		if events == nil {
			continue
		}

		// Emit all events in this bucket
		for _, e := range events {
			l.emitEvent(e, s.timeout)
		}
	}

	// Receive future events by listening
	s.addListener(l)
}
