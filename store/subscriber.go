package store

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
	id       string
	buffer   int // size of messages channel
	messages chan Message
	topics   map[string]struct{}
	active   bool
	mut      sync.RWMutex
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
		id:       genRandomID(),
		buffer:   buffer,
		messages: make(chan Message, buffer),
		topics:   make(map[string]struct{}),
		active:   true,
	}
}

func (s *Subscriber) ID() string {
	s.mut.RLock()
	defer s.mut.RUnlock()
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
	if _, ok := s.topics[topic]; ok {
		return true
	}
	return false
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

func (s *Subscriber) Deactivate() {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.active = false
}

func (s *Subscriber) IsActive() bool {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return s.active
}

func (s *Subscriber) Notify(m *Message) {
	s.mut.Lock()
	defer s.mut.Unlock()
	// log.Println(*m)
	if s.active {
		s.messages <- *m
	}
}

func (s *Subscriber) Listen() chan Message {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return s.messages
}
