package store

import (
	"fmt"

	"github.com/orangeseeds/blitzbase/utils/schema"
)

// type CollectionType int
//
// const (
// 	Base CollectionType = 0
// 	Auth CollectionType = 1
// )
//
// type Collection struct {
// 	ID   int
// 	Name string
// 	Type CollectionType
// }

type CollectionType string

const (
	Base CollectionType = "base"
	Auth CollectionType = "auth"
)

type Collection interface {
	TableName() string
	TableSchema() map[string]string
}

type BaseColletion struct {
	ID     int
	Name   string
	Type   string
	Schema []schema.FieldSchema
}

func (bc BaseColletion) TableName() string {
	return fmt.Sprintf("_base_collection_%s", bc.Name)
}

func (bc BaseColletion) TableSchema() map[string]string {
	tableSchema := map[string]string{}
	for _, s := range bc.Schema {
		tableSchema[s.Name] = s.SQL()
	}
	return tableSchema
}
