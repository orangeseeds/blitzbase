package model

import (
	"testing"

	"github.com/orangeseeds/blitzbase/utils"
)

func TestRecordLoad(t *testing.T) {
	recordId := utils.RandStr(10)
	collectionId := utils.RandStr(10)
	col := NewCollection(collectionId, "test_collection", BASE)
	col.Schema.AddField(&Field{recordId, "field_1", FieldTypeNumber, nil})

	rec := NewRecord(col)
	rec.Load(map[string]any{
		"Id":      recordId,
		"field_1": 1000,
	})

	json, err := rec.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(json))

}
