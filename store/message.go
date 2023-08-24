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
	Action string
	Record struct {
		ID         string
		Collection string
	}
}

func NewMessage(action, topic, data string) *Message {
	return &Message{
		Action: action,
		Record: struct {
			ID         string
			Collection string
		}{
			ID:         data,
			Collection: topic,
		},
	}
}

func (m *Message) FormatSSE() (string, error) {

	buff := bytes.NewBuffer([]byte{})

	encoder := json.NewEncoder(buff)

	err := encoder.Encode(m)
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}

	// sb.WriteString(fmt.Sprintf("subID: %s\n", subID))
	sb.WriteString(fmt.Sprintf("event: %s\n", m.Record.Collection))
	sb.WriteString(fmt.Sprintf("data: %v\n\n", buff.String()))
	// log.Println(sb.String())

	return sb.String(), nil
}
