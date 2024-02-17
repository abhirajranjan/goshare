package pubsub

import (
	"errors"
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
)

var (
	ErrNotFound   = errors.New("value not found")
	ErrTTLExpired = errors.New("ttl expired")
)

type PubSub struct {
	DefaultTtl time.Duration
	data       *treemap.Map
	mu         sync.Mutex
}

type TTLData struct {
	TTL  time.Time
	Data any
}

func NewPubSub(ttl time.Duration) *PubSub {
	p := &PubSub{
		mu:         sync.Mutex{},
		data:       treemap.NewWithStringComparator(),
		DefaultTtl: ttl,
	}

	go p.gc()
	return p
}

func (p *PubSub) Get(id string) (any, error) {
	ival, ok := p.data.Get(id)
	if !ok {
		return nil, ErrNotFound
	}

	data := ival.(TTLData)
	if ok := time.Now().After(data.TTL); ok {
		return data.Data, ErrTTLExpired
	}

	return data.Data, nil
}

func (p *PubSub) Set(id string, data any) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.data.Put(id, TTLData{
		TTL:  time.Now().Add(p.DefaultTtl),
		Data: data,
	})
}

func (p *PubSub) gc() {
	ticker := time.NewTicker(30 * time.Millisecond)
	defer ticker.Stop()

	for {
		<-ticker.C
		p.mu.Lock()

		p.data = p.data.Select(func(key, value interface{}) bool {
			data := value.(TTLData)
			return !time.Now().After(data.TTL)
		})

		p.mu.Unlock()
		ticker.Reset(5 * time.Second)
	}
}
