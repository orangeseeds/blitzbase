package pubsub

import (
	"crypto/rand"
	"fmt"
	"log"
	"sync"
)

type TopicInfo struct {
	Collection string `json:"collection"`
	FieldName  string `json:"field_name"`
	FieldValue any    `json:"field_value"`
}

type Subscriber struct {
	id     string
	buffer int // size of messages channel
	events chan struct{}
	topics map[string]struct{}
	active bool
	mut    sync.RWMutex
}

func genRandomID() string {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%X-%X-%X", b[0:4], b[4:8], b[8:12])
}

func NewSubscriber(buffer int) *Subscriber {
	return &Subscriber{
		id:     genRandomID(),
		buffer: buffer,
		events: make(chan struct{}, buffer),
		topics: make(map[string]struct{}),
		active: true,
	}
}

func (s *Subscriber) ID() string {
	return s.id
}

func (s *Subscriber) Topics() map[string]struct{} {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return s.topics
}

func (s *Subscriber) HasTopic(topic string) bool {
	s.mut.RLock()
	defer s.mut.RUnlock()
	_, exists := s.topics[topic]
	return exists
}

func (s *Subscriber) AddTopics(topics ...string) {
	s.mut.Lock()
	defer s.mut.Unlock()

	for _, t := range topics {
		s.topics[t] = struct{}{}
	}
}

func (s *Subscriber) RemoveTopic(topics ...string) {
	s.mut.Lock()
	defer s.mut.Unlock()

	for _, t := range topics {
		if _, ok := s.topics[t]; ok {
			delete(s.topics, t)
		}
	}
}

func (s *Subscriber) Activate() {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.active = true
}

func (s *Subscriber) Deactivate() {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.active = false
}

func (s *Subscriber) Close() {
	s.mut.Lock()
	defer s.mut.Unlock()
	close(s.events)
}

func (s *Subscriber) IsActive() bool {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return s.active
}

func (s *Subscriber) Notify(m *struct{}) {
	s.mut.Lock()
	defer s.mut.Unlock()
	// log.Println(*m)
	if s.active {
		s.events <- *m
	}
}

func (s *Subscriber) Listen() chan struct{} {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return s.events
}
