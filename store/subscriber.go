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
	topics   map[string]TopicInfo
	active   bool
	mut      sync.RWMutex
}

func NewSubscriber(buffer int) *Subscriber {

	b := make([]byte, 12)

	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	id := fmt.Sprintf("%X-%X-%X", b[0:4], b[4:8], b[8:12])

	return &Subscriber{
		id:       id,
		buffer:   buffer,
		messages: make(chan Message, buffer),
		topics:   map[string]TopicInfo{},
		active:   true,
	}
}

func (s *Subscriber) GetTopics() []string {
	topics := []string{}

	for topic := range s.topics {
		topics = append(topics, topic)
	}
	return topics
}

func (s *Subscriber) addTopic(topic string, rule TopicInfo) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.topics[topic] = rule
}

func (s *Subscriber) removeTopic(topic string) {
	s.mut.Lock()
	defer s.mut.Unlock()
	delete(s.topics, topic)
}

func (s *Subscriber) IsActive() bool {
	return s.active
}

func (s *Subscriber) Notify(m *Message) {
	s.mut.Lock()
	defer s.mut.Unlock()
	log.Println(*m)
	if s.active {
		s.messages <- *m
	}
}

func (s *Subscriber) Listen() <-chan Message {
	return s.messages
}
