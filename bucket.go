package teamgress

import (
	"sort"
	"sync"
)

type buckets struct {
	buckets     map[int64][]Event
	mux         *sync.Mutex
	key         int64
	size        int64
	maxSize     int64
	maxNumber   int
	currentSize int64
}

func newBuckets(maxBucketSize int64, nBuckets int) *buckets {
	return &buckets{
		buckets:   make(map[int64][]Event, 0),
		mux:       &sync.Mutex{},
		maxSize:   maxBucketSize,
		maxNumber: nBuckets,
	}
}

func (b *buckets) add(e Event) {
	b.mux.Lock()
	defer b.mux.Unlock()

	esize := int64(e.Size())
	b.currentSize += esize

	if b.currentSize > b.maxSize {
		b.key += 1

		b._trim(b.key - int64(b.maxNumber))
		b.buckets[b.key] = make([]Event, 0)
	}

	b.buckets[b.key] = append(b.buckets[b.key], e)
}

func (b *buckets) _trim(lowKey int64) {
	keys := b._keys()

	for _, k := range keys {
		if k < lowKey {
			key := keys[0]
			b.buckets[key] = nil
			delete(b.buckets, key)
		}
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
