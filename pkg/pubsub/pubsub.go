package pubsub

import (
	"errors"
	"log"
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
	mu         sync.RWMutex
}

type TTLData struct {
	TTL  time.Time
	Data any
}

func NewPubSub(ttl time.Duration) *PubSub {
	p := &PubSub{
		mu:         sync.RWMutex{},
		data:       treemap.NewWithStringComparator(),
		DefaultTtl: ttl,
	}

	go p.gc()
	return p
}

func (p *PubSub) Get(id string) (any, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

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

	log.Println("data got", id)

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
		expired := 0

		p.mu.RLock()
		tree := p.data.Select(func(key, value interface{}) bool {
			data := value.(TTLData)
			expire := !time.Now().After(data.TTL)
			if !expire {
				log.Printf("expired %#v\n", data)
				expired++
			}

			return expire
		})
		p.mu.RUnlock()

		if expired != 0 {
			p.mu.Lock()
			p.data = tree
			p.mu.Unlock()
		}
	}
}
