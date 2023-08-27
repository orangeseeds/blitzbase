package schema

import "fmt"

const (
	INTEGER_FIELD   = "INTEGER"
	TEXT_FIELD      = "TEXT"
	TIMESTAMP_FIELD = "TIMESTAMP"
	RELATION_FIELD  = "RELATION"
)

type FieldSchema struct {
	Name     string
	Type     string
	Options  Options
	Required bool
	Unique   bool
}

func (fs *FieldSchema) uniqueOpt() string {
	if fs.Unique {
		return "UNIQUE"
	}
	return ""
}

func (fs *FieldSchema) nullOpt() string {
	if fs.Required {
		return "NOT NULL"
	}
	return "NULL"
}

func (fs *FieldSchema) SQL() string {
	if fs.Type == RELATION_FIELD {
		return fmt.Sprintf(
			"INTEGER %s %s %s", fs.nullOpt(), fs.uniqueOpt(), fs.Options.Tags(),
		)
	}
	return fmt.Sprintf(
		"%s %s %s %s", fs.Type, fs.Options.Tags(), fs.uniqueOpt(), fs.nullOpt(),
	)
}
