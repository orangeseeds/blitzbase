package model

import "testing"

func TestSchemaMarshal(t *testing.T) {
	schema := Schema{
		Fields: []*Field{
			{"1", "Name", FieldTypeText, nil},
			{"2", "Addr", FieldTypeNumber, nil},
		},
	}
	expected := `[{"Id":"1","Name":"Name","Type":"text","Options":null},{"Id":"2","Name":"Addr","Type":"number","Options":null}]`

	json, err := schema.MarshalJSON()
	if err != nil {
		t.Error(err)
	}

	got := string(json)
	if got != expected {
		t.Errorf("expected %s got %s", expected, got)
	}
}

func TestSchemaUnmarshal(t *testing.T) {
	expected := `[{"Id":"1","Name":"Name","Type":"text","Options":null},{"Id":"2","Name":"Addr","Type":"number","Options":null}]`
    var schema Schema

    err := schema.UnmarshalJSON([]byte(expected))
    if err != nil {
        t.Error(err)
    }
    t.Logf("%#v", schema)
}

