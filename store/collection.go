package store

import (
	"fmt"
	"strings"

	"github.com/orangeseeds/blitzbase/utils/schema"
)

const (
	BaseType = "BASE"
	AuthType = "AUTH"
)

type Collection interface {
	TableName() string
	TableSchema() map[string]string
	RawSchema() []schema.FieldSchema
	IsAuth() bool
}

type BaseColletion struct {
	Id     int
	Name   string
	Type   string
	System bool
	Schema []schema.FieldSchema
}

func NewBaseCollection(name string, ctype string, system bool, schema []schema.FieldSchema) *BaseColletion {
	c := BaseColletion{
		Id:     0,
		Name:   name,
		Type:   ctype,
		System: system,
		Schema: schema,
	}
	c.addBasicFields()
	return &c
}

func (bc BaseColletion) TableName() string {
	if bc.System {
		return bc.Name
	}

	switch strings.ToUpper(bc.Type) {
	case BaseType:
		return fmt.Sprintf("_base_collection_%s", bc.Name)
	case AuthType:
		return fmt.Sprintf("_auth_collection_%s", bc.Name)
	default:
		return fmt.Sprintf("_%s", bc.Name)
	}
}

func (bc *BaseColletion) addBasicFields() {

	//     switch strings.ToUpper(bc.Type) {
	//     case BaseType:
	//
	// }
	bc.Schema = append(bc.Schema, schema.FieldSchema{
		Name:     "id",
		Type:     schema.INTEGER_FIELD,
		Options:  schema.IntFieldOptions{},
		Required: true,
	})

	bc.Schema = append(bc.Schema, schema.FieldSchema{
		Name:     "created",
		Type:     schema.TIMESTAMP_FIELD,
		Options:  schema.TimeStampFieldOptions{},
		Required: true,
	})

	bc.Schema = append(bc.Schema, schema.FieldSchema{
		Name:     "updated",
		Type:     schema.TIMESTAMP_FIELD,
		Options:  schema.TimeStampFieldOptions{},
		Required: true,
	})

}

func (bc BaseColletion) TableSchema() map[string]string {
	tableSchema := map[string]string{}

	for _, s := range bc.Schema {
		if s.Name == "id" {
			tableSchema[s.Name] = "INTEGER PRIMARY KEY"
			continue
		}

		tableSchema[s.Name] = s.SQL()
	}
	return tableSchema
}

func (bc BaseColletion) RawSchema() []schema.FieldSchema {
	return bc.Schema
}

func (bc BaseColletion) IsAuth() bool {
	return strings.ToUpper(bc.Type) == AuthType
}
