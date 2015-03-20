package main

import (
	"sync"
	"time"
)

type store struct {
	events   []event
	conf     *conf
	indicesS map[string]map[string][]*event
	indicesI map[string]map[int64][]*event
	mux      *sync.Mutex
}

func newStore(conf *conf) *store {
	s := &store{
		mux:      &sync.Mutex{},
		conf:     conf,
		events:   make([]event, 0),
		indicesS: make(map[string]map[string][]*event),
		indicesI: make(map[string]map[int64][]*event),
	}

	s.indicesS["user"] = make(map[string][]*event)
	s.indicesS["envName"] = make(map[string][]*event)
	s.indicesI["time"] = make(map[int64][]*event)

	return s
}

func (s *store) _addToIndexString(index, key string, e *event) {
	if _, found := s.indicesS[index][key]; !found {
		s.indicesS[index][key] = make([]*event, 0)
	}

	s.indicesS[index][key] = append(s.indicesS[index][key], e)
}

func (s *store) _addToIndexInt64(index string, key int64, e *event) {
	if _, found := s.indicesI[index][key]; !found {
		s.indicesI[index][key] = make([]*event, 0)
	}

	s.indicesI[index][key] = append(s.indicesI[index][key], e)
}

func (s *store) listen(evs <-chan event) {
	for e := range evs {
		s.mux.Lock()

		// Copy the event in the store.
		s.events = append(s.events, e)
		// Point to the copy in the indices.
		ep := &s.events[len(s.events)-1]

		s._addToIndexString("user", e.User.UnixName, ep)
		s._addToIndexString("envName", e.EnvName, ep)
		s._addToIndexInt64("time", e.Time.Unix(), ep)

		s.mux.Unlock()
	}
}

func (s *store) getByUser(user string, ch chan<- event) {
	s.mux.Lock()
	defer s.mux.Unlock()
	defer close(ch)

	events, found := s.indicesS["user"][user]
	if !found {
		return
	}

	for _, e := range events {
		ch <- *e
	}
}

func (s *store) getFromTime(time time.Time, ch chan<- event) {
	s.mux.Lock()
	defer s.mux.Unlock()
	defer close(ch)

	unixTime := time.Unix()

	for k, events := range s.indicesI["time"] {
		if k >= unixTime {
			for _, e := range events {
				ch <- *e
			}
		}
	}
}
