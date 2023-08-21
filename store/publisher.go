package store

import (
	"sync"
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

func (p *Publisher) SubByID(id string) *Subscriber {
	p.mut.RLock()
	defer p.mut.RUnlock()

	if sub, ok := p.subs[id]; ok {
		return sub
	}
	return nil
}

func (p *Publisher) SubCount() int {
	p.mut.Lock()
	defer p.mut.Unlock()
	return len(p.subs)
}

func (p *Publisher) Broadcast(data HookData, topic string) {
	for _, sub := range p.subs {

		if !sub.HasTopic(topic) {
			continue
		}
		go publish(sub, topic, data)
	}
}

func publish(s *Subscriber, topic string, data HookData) {
	msg := NewMessage(topic, data)
	// collection := s.topics[topic].Collection
	// if collection == "*" || data.CollectionName == collection {
	s.Notify(msg)
	// }
}
