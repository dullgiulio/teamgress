package libteamgress

import (
	"sort"
	"sync"
)

type buckets struct {
	buckets map[int64][]Event
	mux     *sync.Mutex
	secs    int64
	max     int
}

func newBuckets(secs int64, max int) *buckets {
	return &buckets{
		buckets: make(map[int64][]Event, 0),
		mux:     &sync.Mutex{},
		secs:    secs,
		max:     max,
	}
}

func (b *buckets) add(e Event) {
	b.mux.Lock()
	defer b.mux.Unlock()

	eTime := e.Time.Unix()
	key := eTime - (eTime % b.secs)

	if _, found := b.buckets[key]; !found {
		b._trim()
		b.buckets[key] = make([]Event, 0)
	}

	b.buckets[key] = append(b.buckets[key], e)
}

func (b *buckets) _trim() {
	keys := b._keys()

	if len(keys) > b.max {
		key := keys[0]
		b.buckets[key] = nil
		delete(b.buckets, key)
	}
}

func (b *buckets) get(key int64) (t []Event) {
	b.mux.Lock()
	defer b.mux.Unlock()

	if events, found := b.buckets[key]; found {
		t = make([]Event, len(events))
		copy(t, events)
	}

	return t
}

func (b *buckets) _keys() (keys []int64) {
	keys = make([]int64, len(b.buckets))
	i := 0

	for k, _ := range b.buckets {
		keys[i] = k
		i += 1
	}

	sort.Sort(Int64Slice(keys))

	return keys
}

func (b *buckets) keys() []int64 {
	b.mux.Lock()
	defer b.mux.Unlock()

	return b._keys()
}
