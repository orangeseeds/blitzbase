package model

import (
	"encoding/json"
	"strings"
)

type CollectionType string

const (
	BASE CollectionType = "base"
	AUTH CollectionType = "auth"
)

// check if names contains whitespace
type Collection struct {
	BaseModel

	Name   string         `db:"Name"`
	Type   CollectionType `db:"Type"` // base,auth
	Schema Schema         `db:"Schema"`
	Rule   string         `db:"Rule"` // for now all rules, this needs to later be expanded to fit list, view, update, create and delete
}

func NewCollection(id string, name string, colType CollectionType) *Collection {
	col := &Collection{
		Type: colType,
	}
	col.SetID(id)
	col.SetName(name)
	return col
}

func (c *Collection) SetName(name string) {
	c.Name = strings.ReplaceAll(name, " ", "_")
}

func (c *Collection) GetName() string {
	return c.Name
}

func (c *Collection) TableName() string {
	return "_collection"
}

func (c *Collection) IsAuth() bool {
	return c.Type == AUTH
}

// Provides an key,val pair of col name and datatype to build a new collection table
func (c *Collection) DataDefn() map[string]string {
	toExport := make(map[string]string)
	toExport["Id"] = string(Text)
	for _, f := range c.Schema.GetFields() {
		toExport[f.Name] = string(f.Type)
	}
	return toExport
}

func (c *Collection) MetaDataDefn() map[string]string {
	return map[string]string{
		"Id":     string(Text) + " primary key",
		"Name":   string(Text),
		"Type":   string(Text),
		"Schema": string(Json),
		"Rule":   string(Text),
	}
}

func (c Collection) MarshalJSON() ([]byte, error) {
	type alias Collection // prevent recursion
	return json.Marshal(alias(c))
}
