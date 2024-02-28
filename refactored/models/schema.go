package model

import (
	"encoding/json"
	"strings"

	"github.com/orangeseeds/blitzbase/utils"
)

type FieldType string

const (
	Text     FieldType = "text"
	Number   FieldType = "number"
	Bool     FieldType = "bool"
	Email    FieldType = "email"
	Url      FieldType = "url"
	Editor   FieldType = "editor"
	Date     FieldType = "date"
	Select   FieldType = "select"
	Json     FieldType = "json"
	File     FieldType = "file"
	Relation FieldType = "relation"
)

// basic fields
const (
	Id = "Id"
)

type FieldOption interface {
	Validate() error
}

type Field struct {
	Id   string
	Name string
	Type FieldType

	// js code string
	Options FieldOption
}

type Schema struct {
    Fields []*Field 
}

func (s *Schema) GetFields() []*Field {
	return s.Fields
}

func (s *Schema) AddField(newField *Field) {
	if newField.Id == "" {
		newField.Id = strings.ToLower(utils.RandStr(10))
	}

	for i, field := range s.Fields {
		if field.Id == newField.Id {
			s.Fields[i] = newField
			return
		}
	}

	s.Fields = append(s.Fields, newField)
}

func (s Schema) MarshalJSON() ([]byte, error) {
	if s.Fields == nil {
		s.Fields = []*Field{}
	}

	return json.Marshal(s.Fields)
}

func (s *Schema) UnmarshalJSON(data []byte) error {
	var fields []Field

	s.Fields = []*Field{}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	for _, f := range fields {
		s.AddField(&f)
	}
	return nil
}

func (s *Schema) HasField(name string) bool {
	for _, f := range s.GetFields() {
		if f.Name == name {
			return true
		}
	}
	return false
}
