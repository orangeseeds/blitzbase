package model

import (
	"testing"
)

func TestCollectionMarshal(t *testing.T) {
	col := NewCollection("id", "test_collection", BASE)
	col.Schema.AddField(&Field{"1", "FieldName", FieldTypeText, nil})
	col.Schema.AddField(&Field{"2", "FieldName2", FieldTypeText, nil})

    expected := `{"Id":"id","Name":"test_collection","Type":"base","Schema":[{"Id":"1","Name":"FieldName","Type":"text","Options":null},{"Id":"2","Name":"FieldName2","Type":"text","Options":null}],"Rule":""}`

	json, err := col.MarshalJSON()
	if err != nil {
		t.Error("Marshal Error", err)
	}
	got := string(json)

	if expected != got {
		t.Errorf("expected %s got %s", expected, got)
	}
}
