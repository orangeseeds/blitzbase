package store

import (
	"sync"
)

type Subscribers map[string]*Subscriber

type Publisher struct {
	topics map[string]Subscribers // map of topic->Subscriber
	mut    sync.RWMutex
}

func NewPublisher() *Publisher {
	return &Publisher{
		topics: map[string]Subscribers{},
	}
}

func (p *Publisher) Subscribe(s *Subscriber, topic string, rule TopicInfo) {
	p.mut.Lock()
	defer p.mut.Unlock()
	// for _, topic := range topics {
	if _, ok := p.topics[topic]; !ok {
		p.topics[topic] = Subscribers{}
	}
	s.addTopic(topic, rule)
	p.topics[topic][s.id] = s
	// }
}

func (p *Publisher) Unsubscribe(s *Subscriber, topics ...string) {
	p.mut.Lock()
	defer p.mut.Unlock()
	for _, topic := range topics {
		delete(p.topics[topic], s.id)
		s.removeTopic(topic)
	}
}

func (p *Publisher) SubCount(topic string) int {
	if subs, ok := p.topics[topic]; ok {
		return len(subs)
	}
	return 0
}

func (p *Publisher) Broadcast(data HookData, topic string) {
	// p.mut.Lock()
	// defer p.mut.Unlock()

	// for _, topic := range topics {
	if subs, ok := p.topics[topic]; ok {
		for _, sub := range subs {
			go func(sub *Subscriber, topic string) {
				msg := NewMessage(topic, data)
				collection := sub.topics[topic].Collection
				if collection == "*" || data.CollectionName == collection {
					sub.Notify(msg)
				}
			}(sub, topic)
		}
	}
	// }
}
