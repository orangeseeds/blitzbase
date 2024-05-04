package pubsub

// package pubsub
//
// import (
// 	"fmt"
// 	"sync"
// )
//
// type Producer struct {
// 	subs map[string]*Consumer // map of id->Subscriber
// 	mut  sync.RWMutex
// }
//
// func NewProducer() *Producer {
// 	return &Producer{
// 		subs: make(map[string]*Consumer),
// 	}
// }
//
// func (p *Producer) Subscribe(s *Consumer) {
// 	p.mut.Lock()
// 	defer p.mut.Unlock()
// 	p.subs[s.ID()] = s
// }
//
// func (p *Producer) Unsubscribe(id string) {
// 	p.mut.Lock()
// 	defer p.mut.Unlock()
// 	if sub, ok := p.subs[id]; ok {
// 		sub.Deactivate()
// 		delete(p.subs, id)
// 	}
// }
//
// func (p *Producer) SubByID(id string) (*Consumer, error) {
// 	p.mut.RLock()
// 	defer p.mut.RUnlock()
//
// 	// fmt.Println(p.subs)
// 	// fmt.Println(id)
// 	if sub, ok := p.subs[id]; ok {
// 		return sub, nil
// 	}
//
// 	return nil, fmt.Errorf("[PUB]no subscriber with the ID %s exists", id)
// }
//
// func (p *Producer) SubCount() int {
// 	p.mut.Lock()
// 	defer p.mut.Unlock()
// 	return len(p.subs)
// }
//
// // here topic means collection
// func (p *Producer) Dispatch(e struct{}) {
// 	// log.Println(e)
// 	// for _, sub := range p.subs {
//
// 	// 	if !sub.HasTopic(e.Message.Record.Collection) || !p.IsAllowed(sub) {
// 	// 		continue
// 	// 	}
// 	// 	go func(s *Subscriber, e DBHookEvent) {
// 	// 		s.Notify(&e)
// 	//
// 	// 	}(sub, e)
// 	// }
// }
