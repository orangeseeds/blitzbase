package model

import (
	"encoding/json"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/spf13/cast"
)

type Record struct {
	BaseModel

	collection *Collection
	data       map[string]any
	// remaining expanded relations
}

func NewRecord(c *Collection) *Record {
	return &Record{
		collection: c,
		data:       make(map[string]any),
	}
}

func (r *Record) LoadNullStringMap(data dbx.NullStringMap) {
	for key, val := range data {
		if val.Valid {
			r.Set(key, val.String)
		}
	}
}

func (r *Record) Load(data map[string]any) {
	for key, val := range data {
		r.Set(key, val)
	}
}

func (r *Record) TableName() string {
	return r.Collection().GetName()
}

func (r *Record) Collection() *Collection {
	return r.collection
}

func (r *Record) Set(key string, val any) {
	switch key {
	case Id:
		r.SetID(cast.ToString(val))
	default:
		if r.Collection().Schema.HasField(key) {
			r.data[key] = val
		}
	}
}

func (r *Record) Export() map[string]any {
	toExport := make(map[string]any)
	if r.data != nil {
		toExport = r.data
	}
	toExport[Id] = r.Id
	return toExport
}

func (r Record) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Export())
}

func (r *Record) UnmarshalJSON(data []byte) error {
	res := map[string]any{}

	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	r.Load(res)

	return nil
}
