package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

/**
Topic
    collection: "",
    field_name: "",
    field_value: "",

**/

type Message struct {
	topic string
	data  HookData
}

func NewMessage(topic string, meta HookData) *Message {
	return &Message{
		topic: topic,
		data:  meta,
	}
}

func (m *Message) GetTopic() string {
	return m.topic
}

func (m *Message) GetData() HookData {
	return m.data
}

func (m *Message) FormatSSE() (string, error) {
	buff := bytes.NewBuffer([]byte{})

	encoder := json.NewEncoder(buff)

	err := encoder.Encode(m.data)
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("event: %s\n", m.topic))
	sb.WriteString(fmt.Sprintf("data: %v\n\n", buff.String()))
	// log.Println(sb.String())

	return sb.String(), nil
}
