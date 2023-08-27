package schema

import (
	"fmt"
)

type Options interface {
	Tags() string
}

type TextFieldOptions struct {
	MinLen int
	MaxLen int
	Regex  int
}

func (o TextFieldOptions) Tags() string {
	return `DEFAULT ""`
}

type IntFieldOptions struct {
	Max int
	Min int
}

func (o IntFieldOptions) Tags() string {
	return `DEFAULT 0`
}

type RelationFieldOptions struct {
	RelatedTable      string `validate:"required"`
	RelatedTableField string
	CascadingDel      bool `validate:"boolean,required"`
}

func (o RelationFieldOptions) Tags() string {
	if o.CascadingDel {
		return fmt.Sprintf("REFERENCES %s(%s) ON DELETE CASCADE", o.RelatedTable, o.RelatedTableField)
	}
	return fmt.Sprintf("REFERENCES %s(%s)", o.RelatedTable, o.RelatedTableField)
}

type TimeStampFieldOptions struct {
	Before string
	After  string
}

func (o TimeStampFieldOptions) Tags() string {
	return `DEFAULT CURRENT_TIMESTAMP`
}
