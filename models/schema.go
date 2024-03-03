package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/orangeseeds/blitzbase/utils"
)

const (
	FieldTypeText     = "text"
	FieldTypeNumber   = "number"
	FieldTypeBool     = "bool"
	FieldTypeDate     = "date"
	FieldTypeJson     = "json"
	FieldTypeFile     = "file"
	FieldTypeRelation = "relation"
)

type FieldName string

// basic fields
const (
	FieldId       = "Id"
	FieldEmail    = "Email"
	FieldToken    = "Token"
	FieldPassword = "Password"
	// FieldRule     = "Rule"

	FieldIndexRule  = "IndexRule"
	FieldDetailRule = "DetailRule"
	FieldUpdateRule = "UpdateRule"
	FieldCreateRule = "CreateRule"
	FieldDeleteRule = "DeleteRule"

	FieldCreatedAt = "CreatedAt"
	FieldUpdatedAt = "UpdatedAt"
)

func AuthRecordFields() []string {
	return []string{
		FieldPassword,
		FieldEmail,
		FieldToken,
		FieldPassword,
		// FieldCreatedAt,
		// FieldUpdatedAt,
	}
}

// contains Id,CreatedAt,UpdatedAt
func BaseFields() []string {
	return []string{
		FieldId,
		FieldCreatedAt,
		FieldUpdatedAt,
	}
}

func (f FieldName) String() string {
	return string(f)
}

type FieldOption interface {
	Validate() error
}

type Field struct {
	Id   string
	Name string
	Type string

	// js code string
	Options FieldOption
}

type Schema struct {
	Fields []*Field `json:"Fields"`
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

// into the current Schema instance.
func (s *Schema) Scan(value any) error {
	var data []byte
	switch v := value.(type) {
	case nil:
		// no cast needed
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("Failed to unmarshal Schema value %q.", value)
	}

	if len(data) == 0 {
		data = []byte("[]")
	}

	return s.UnmarshalJSON(data)
}

func (s *Schema) HasField(name string) bool {
	for _, f := range s.GetFields() {
		if f.Name == name {
			return true
		}
	}
	return false
}
