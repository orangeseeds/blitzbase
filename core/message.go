package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Message struct {
	topic string
	body  string
}

func (m *Message) NewMessage(topic, body string) *Message {
	return &Message{
		topic: topic,
		body:  body,
	}
}

func (m *Message) GetTopic() string {
	return m.topic
}

func (m *Message) GetBody() string {
	return m.body
}

func (m *Message) FormatSSE() (string, error) {
	data := map[string]any{
		"data": m.body,
	}

	buff := bytes.NewBuffer([]byte{})

	encoder := json.NewEncoder(buff)

	err := encoder.Encode(data)
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("event: %s\n", m.topic))
	sb.WriteString(fmt.Sprintf("data: %v\n\n", buff.String()))

	return sb.String(), nil
}
