package store

import (
	"fmt"
	"sync"
	// "github.com/orangeseeds/blitzbase/core"
)

type Publisher struct {
	subs map[string]*Subscriber // map of id->Subscriber
	mut  sync.RWMutex
}

func NewPublisher() *Publisher {
	return &Publisher{
		subs: make(map[string]*Subscriber),
	}
}

func (p *Publisher) Subscribe(s *Subscriber) {
	p.mut.Lock()
	defer p.mut.Unlock()
	p.subs[s.ID()] = s
}

func (p *Publisher) Unsubscribe(id string) {
	p.mut.Lock()
	defer p.mut.Unlock()
	if sub, ok := p.subs[id]; ok {
		sub.Deactivate()
		delete(p.subs, id)
	}
}

func (p *Publisher) SubByID(id string) (*Subscriber, error) {
	p.mut.RLock()
	defer p.mut.RUnlock()

	// fmt.Println(p.subs)
	// fmt.Println(id)
	if sub, ok := p.subs[id]; ok {
		return sub, nil
	}

	return nil, fmt.Errorf("[PUB]no subscriber with the ID %s exists", id)
}

func (p *Publisher) SubCount() int {
	p.mut.Lock()
	defer p.mut.Unlock()
	return len(p.subs)
}

// here topic means collection
func (p *Publisher) Dispatch(e DBHookEvent) {
	// log.Println(e)
	for _, sub := range p.subs {

		if !sub.HasTopic(e.Message.Record.Collection) || !p.IsAllowed(sub) {
			continue
		}
		go func(s *Subscriber, e DBHookEvent) {
			s.Notify(&e)

		}(sub, e)
	}
}

func (p *Publisher) IsAllowed(s *Subscriber) bool {
	return true
}
