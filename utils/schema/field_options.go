package schema

import "fmt"

type Options interface {
	Tags() string
}

type TextFieldOptions struct {
	MinLen int
	MaxLen int
	Regex  int
}

func (o TextFieldOptions) Tags() string {
	return ""
}

type IntFieldOptions struct {
	Max int
	Min int
}

func (o IntFieldOptions) Tags() string {
	return ""
}

type RelationFieldOptions struct {
	RelatedTable      string
	RelatedTableField string
	CascadingDel      bool
}

func (o RelationFieldOptions) Tags() string {
	if o.CascadingDel {
		return fmt.Sprintf("REFERENCES %s(%s) ON DELETE CASCADE", o.RelatedTable, o.RelatedTableField)
	}
	return fmt.Sprintf("REFERENCES %s(%s)", o.RelatedTable, o.RelatedTableField)
}
