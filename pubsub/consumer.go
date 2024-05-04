package pubsub
// package pubsub
//
// import (
// 	"sync"
//
// 	"github.com/google/uuid"
// 	model "github.com/orangeseeds/blitzbase/models"
// )
//
// type TopicType string
//
// const (
// 	TopicCreate TopicType = "Create"
// 	TopicView   TopicType = "View"
// 	TopicUpdate TopicType = "Update"
// 	TopicDelete TopicType = "Delete"
// 	TopicAuth   TopicType = "Auth"
// )
//
// type Topic struct {
// 	Type         TopicType
// 	CollectionId string
// 	RecordId     string
// 	Data         *model.Record
// }
//
// type Consumer struct {
// 	id     string
// 	bus    chan struct{}
// 	mut    sync.RWMutex
// 	topics []Topic
// }
//
// func NewSubscriber(topics ...Topic) *Consumer {
// 	return &Consumer{
// 		id:     uuid.NewString(),
// 		bus:    make(chan struct{}, 1),
// 		topics: append([]Topic{}, topics...),
// 	}
// }
//
// func (s *Consumer) ID() string {
// 	return s.id
// }
//
// func (s *Consumer) Topics() []Topic {
// 	s.mut.RLock()
// 	defer s.mut.RUnlock()
// 	return s.topics
// }
//
// func (s *Consumer) HasTopic(topic string) bool {
// 	s.mut.RLock()
// 	defer s.mut.RUnlock()
// 	_, exists := s.topics[topic]
// 	return exists
// }
//
// func (s *Consumer) AddTopics(topics ...string) {
// 	s.mut.Lock()
// 	defer s.mut.Unlock()
//
// 	for _, t := range topics {
// 		s.topics[t] = struct{}{}
// 	}
// }
//
// func (s *Consumer) RemoveTopic(topics ...string) {
// 	s.mut.Lock()
// 	defer s.mut.Unlock()
//
// 	for _, t := range topics {
// 		if _, ok := s.topics[t]; ok {
// 			delete(s.topics, t)
// 		}
// 	}
// }
//
// func (s *Consumer) Activate() {
// 	s.mut.Lock()
// 	defer s.mut.Unlock()
// 	s.active = true
// }
//
// func (s *Consumer) Deactivate() {
// 	s.mut.Lock()
// 	defer s.mut.Unlock()
// 	s.active = false
// }
//
// func (s *Consumer) Close() {
// 	s.mut.Lock()
// 	defer s.mut.Unlock()
// 	close(s.bus)
// }
//
// func (s *Consumer) IsActive() bool {
// 	s.mut.RLock()
// 	defer s.mut.RUnlock()
// 	return s.active
// }
//
// func (s *Consumer) Notify(m *struct{}) {
// 	s.mut.Lock()
// 	defer s.mut.Unlock()
// 	// log.Println(*m)
// 	if s.active {
// 		s.bus <- *m
// 	}
// }
//
// func (s *Consumer) Listen() chan struct{} {
// 	s.mut.RLock()
// 	defer s.mut.RUnlock()
// 	return s.bus
// }
