package model

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
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

type CommonFieldName string

// basic fields
const (
	FieldId       = "id"
	FieldEmail    = "email"
	FieldName     = "name"
	FieldToken    = "token"
	FieldPassword = "password"
	FieldType     = "type"
	// FieldRule     = "Rule"

	FieldSchema  = "schema"
	FieldOptions = "options"

	FieldIndexRule  = "index_rule"
	FieldDetailRule = "detail_rule"
	FieldUpdateRule = "update_rule"
	FieldCreateRule = "create_rule"
	FieldDeleteRule = "delete_rule"

	FieldCreatedAt = "created_at"
	FieldUpdatedAt = "updated_at"
)

func CollectionFieldsFields() []string {
	return []string{
		FieldId,
		FieldName,
		FieldType,
		FieldSchema,
		FieldOptions,
		FieldIndexRule,
		FieldDetailRule,
		FieldUpdateRule,
		FieldDeleteRule,
	}
}

func AuthRecordFields() []string {
	return []string{
		FieldEmail,
		FieldPassword,
		FieldToken,
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

func (f CommonFieldName) String() string {
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
    Fields []*Field `json:"fields" validate:"unique=Name"`
}

func (s *Schema) GetFields() []*Field {
	return s.Fields
}

func (s *Schema) AddField(newField *Field) {
	if newField.Id == "" {
		newField.Id = uuid.NewString()
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
