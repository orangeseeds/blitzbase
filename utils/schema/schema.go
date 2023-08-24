package schema

import "fmt"

const (
	INTEGER_FIELD   = "INTEGER"
	TEXT_FIELD      = "TEXT"
	TIMESTAMP_FIELD = "TIMESTAMP"
)

type FieldSchema struct {
	Name     string
	Type     string
	Options  Options
	Required bool
}

func (fs *FieldSchema) nullOpt() string {
	if fs.Required {
		return "NOT NULL"
	}
	return "NULL"
}

func (fs *FieldSchema) SQL() string {
	return fmt.Sprintf(
		"%s %s %s", fs.Type, fs.Options.Tags(), fs.nullOpt(),
	)
}
