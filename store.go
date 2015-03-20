package main

import (
	"sync"
)

type store struct {
	events  []*event
	indices map[string]map[string][]*event
	mux     *sync.Mutex
}

func newStore() *store {
	s := &store{}

	s.mux = &sync.Mutex{}
	s.events = make([]*event, 0)
	s.indices = make(map[string]map[string][]*event)
	s.indices["user"] = make(map[string][]*event)
	s.indices["envName"] = make(map[string][]*event)

	return s
}

func (s *store) _addToIndex(index, key string, e *event) {
	if _, found := s.indices[index][key]; !found {
		s.indices[index][key] = make([]*event, 0)
	}

	s.indices[index][key] = append(s.indices[index][key], e)
}

func (s *store) listen(evs <-chan event) {
	for e := range evs {
		s.mux.Lock()

		s.events = append(s.events, &e)
		s._addToIndex("user", e.user, &e)
		s._addToIndex("envName", e.envName, &e)

		s.mux.Unlock()
	}
}

func (s *store) getByUser(user string, ch chan<- event) {
	s.mux.Lock()
	defer s.mux.Unlock()
	defer close(ch)

	events, found := s.indices["user"][user]
	if !found {
		return
	}

	for _, e := range events {
		ch <- *e
	}
}
