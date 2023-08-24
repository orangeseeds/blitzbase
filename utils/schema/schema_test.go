package schema

import (
	"testing"
)

func TestSchemaSQL(t *testing.T) {

	fsRel := FieldSchema{
		Name: "test_col",
		Type: INTEGER_FIELD,
		Options: RelationFieldOptions{
			RelatedTable:      "test_relation",
			RelatedTableField: "ID",
			CascadingDel:      true,
		},
		Required: true,
	}

	fs := FieldSchema{
		Name:     "test_col",
		Type:     INTEGER_FIELD,
		Options:  IntFieldOptions{},
		Required: true,
	}

	cases := []struct {
		Expected string
		Got      string
	}{
		{"INTEGER  NOT NULL", fs.SQL()},
		{"INTEGER REFERENCES test_relation(ID) ON DELETE CASCADE NOT NULL", fsRel.SQL()},
	}

	for _, c := range cases {
		if c.Expected != c.Got {
			t.Errorf(`Expected "%s" GOT "%s"`, c.Expected, c.Got)
		}
	}

}
