package libteamgress

import (
	"sync"
	"time"
)

type Listener struct {
	accept Filter
	ch     chan<- Event
}

type Store struct {
	buckets     map[int64][]Event
	conf        *Conf
	mux         *sync.Mutex
	Listeners   map[*Listener]struct{}
	ListenersCh chan *Listener
	eventsCh    chan Event
	timeout     time.Duration
	bucketSecs  int64
}

func NewStore(conf *Conf) *Store {
	s := &Store{
		mux:         &sync.Mutex{},
		conf:        conf,
		buckets:     make(map[int64][]Event, 0),
		Listeners:   make(map[*Listener]struct{}),
		ListenersCh: make(chan *Listener),
		eventsCh:    make(chan Event, 5),    // Can use some buffering here.
		timeout:     time.Millisecond * 500, // TODO: From config
		bucketSecs:  10,                     // TODO: From config
	}

	// Remove listeners when the are cancelled.
	go s.handleCancelled()
	// Broadcast events to all listeners.
	go s.broadcast()

	return s
}

func (s *Store) Cancel(l *Listener) {
	s.ListenersCh <- l
}

func (s *Store) broadcast() {
	for e := range s.eventsCh {
		s.mux.Lock()

		for l, _ := range s.Listeners {
			// The storage will be locked for s.timeout * len(s.Listeners) at max.
			l.emitEvent(e, s.timeout)
		}

		s.mux.Unlock()
	}
}

func (s *Store) handleCancelled() {
	for l := range s.ListenersCh {
		s.mux.Lock()
		delete(s.Listeners, l)
		s.mux.Unlock()

		close(l.ch)
	}
}

func (s *Store) addToBucket(e Event) {
	s.mux.Lock()
	defer s.mux.Unlock()

	eTime := e.Time.Unix()
	key := eTime - (eTime % s.bucketSecs)

	if _, found := s.buckets[key]; !found {
		s.buckets[key] = make([]Event, 1)
	}

	s.buckets[key] = append(s.buckets[key], e)
}

func (s *Store) getBucket(key int64) (t []Event) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if events, found := s.buckets[key]; found {
		t = make([]Event, len(events))
		copy(t, events)
	}

	return t
}

func (s *Store) Listen(evs <-chan Event) {
	for e := range evs {
		// Copy the event in the Store
		s.addToBucket(e)

		s.eventsCh <- e
	}
}

func (s *Store) addListener(l *Listener) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.Listeners[l] = struct{}{}
}

func (s *Store) bucketsKeys() (keys []int64) {
	s.mux.Lock()
	defer s.mux.Unlock()

	keys = make([]int64, len(s.buckets))
	for k, _ := range s.buckets {
		keys = append(keys, k)
	}

	return keys
}

func (l *Listener) emitEvent(event Event, timeout time.Duration) {
	if l.accept(event) {
		// A client can only listen to a s.timeout periond of time
		// or it will be skipped. The storage will be locked for
		// s.timeout * len(s.Listeners) at max.
		select {
		case l.ch <- event:
		case <-time.After(timeout):
		}
	}
}

func (s *Store) startListener(l *Listener) {
	bucketsKeys := s.bucketsKeys()

	for _, key := range bucketsKeys {
		events := s.getBucket(key)

		if events == nil {
			continue
		}

		for _, e := range events {
			l.emitEvent(e, s.timeout)
		}
	}

	s.addListener(l)
}

func (s *Store) Subscribe(evs chan<- Event, accept Filter) *Listener {
	l := &Listener{
		ch:     evs,
		accept: accept,
	}

	go s.startListener(l)

	return l
}

type Filter func(Event) bool

func GetByUser(user string) Filter {
	return func(e Event) bool {
		return e.User.UnixName == user
	}
}

func GetFromTime(time time.Time) Filter {
	unixTime := time.Unix()

	return func(e Event) bool {
		return e.Time.Unix() >= unixTime
	}
}
