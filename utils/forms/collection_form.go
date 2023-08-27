package forms

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/orangeseeds/blitzbase/store"
	"github.com/orangeseeds/blitzbase/utils/schema"
)

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()
}

type CreateCollectionForm struct {
	Name   string `validate:"required"`
	Type   string `validate:"required,oneof=base auth"`
	Schema []struct {
		Name     string `validate:"required"`
		Type     string `validate:"required,oneof=TEXT RELATION INTEGER"`
		Required bool   `validate:"boolean"`
		Unique   bool   `validate:"boolean"`
		Options  map[string]any
	}
	CreateRule string
	ListRule   string
}

func (s *CreateCollectionForm) IsValid() error {
	err := validate.Struct(s)
	if err != nil {
		return err
	}

	for _, field := range s.Schema {
		switch strings.ToUpper(field.Type) {
		case schema.RELATION_FIELD:
			var f schema.RelationFieldOptions
			if err = mapstructure.Decode(field.Options, &f); err != nil {
				return err
			}
			if err := validate.Struct(f); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *CreateCollectionForm) ToCollection() store.BaseColletion {
	sch := []schema.FieldSchema{}
	for _, field := range f.Schema {
		fs := schema.FieldSchema{
			Name:     field.Name,
			Type:     strings.ToUpper(field.Type),
			Required: field.Required,
		}

		switch strings.ToUpper(field.Type) {
		case schema.TEXT_FIELD:
			fs.Options = schema.TextFieldOptions{}
		case schema.RELATION_FIELD:
			fs.Options = schema.RelationFieldOptions{
				RelatedTable:      field.Options["relatedTable"].(string),
				RelatedTableField: "id",
				CascadingDel:      false,
			}
		case schema.INTEGER_FIELD:
			fs.Options = schema.IntFieldOptions{}
		default:
			continue
		}
		sch = append(sch, fs)
	}

	return *store.NewBaseCollection(f.Name, store.BaseType, false, sch)
}
