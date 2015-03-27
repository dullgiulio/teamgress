package teamgress

import (
	"time"
)

type listener struct {
	accept Filter
	ch     chan<- Event
}

func newListener(ch chan<- Event, accept Filter) *listener {
	return &listener{
		ch:     ch,
		accept: accept,
	}
}

// Send an event to a listener withing timeout time.
func (l *listener) emitEvent(event Event, timeout time.Duration) {
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
