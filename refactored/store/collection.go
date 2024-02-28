package store

import (
	"encoding/json"
	"fmt"

	dbx "github.com/go-ozzo/ozzo-dbx"
	model "github.com/orangeseeds/blitzbase/refactored/models"
)

type FilterFunc func(q *dbx.SelectQuery) error

// TODO: If possible try to implement the Unmarshal directly inside schema
func (s *BaseStore) FindCollectionByNameorId(db any, query string) (*model.Collection, error) {
	var col model.Collection
	var selectQuery *dbx.SelectQuery
	switch db.(type) {
	case *dbx.Tx:
		selectQuery = db.(*dbx.Tx).Select().From(
			col.TableName(),
		)
	case *dbx.DB:
		selectQuery = db.(*dbx.DB).Select().From(
			col.TableName(),
		)
	default:
		return nil, fmt.Errorf("Type didnot fit in FindCollection!")
	}

	res, err := selectQuery.From(
		col.TableName(),
	).Where(
		dbx.HashExp{"Name": query},
	).OrWhere(
		dbx.HashExp{"Id": query},
	).Rows()
	if err != nil {
		return nil, err
	}

	var id, name, tp, schema, rule string
	for res.Next() {
		res.Scan(&id, &name, &tp, &schema, &rule)
	}

	col.Id = id
	col.Name = name
	col.Type = model.CollectionType(tp)
	col.Rule = rule

	err = json.Unmarshal([]byte(schema), &col.Schema)
	if err != nil {
		return nil, err
	}

	return &col, nil
}

func (s *BaseStore) SaveCollection(db any, col *model.Collection) error {
	json, err := col.Schema.MarshalJSON()
	if err != nil {
		return err
	}

	switch db.(type) {
	case *dbx.Tx:
		_, err = db.(*dbx.Tx).Insert(col.TableName(), dbx.Params{
			"Id":     col.GetID(),
			"Name":   col.GetName(),
			"Type":   col.Type,
			"Schema": string(json),
			"Rule":   col.Rule,
		}).Execute()
		if err != nil {
			return err
		}
	case *dbx.DB:
		_, err = db.(*dbx.DB).Insert(col.TableName(), dbx.Params{
			"Id":     col.GetID(),
			"Name":   col.GetName(),
			"Type":   col.Type,
			"Schema": string(json),
			"Rule":   col.Rule,
		}).Execute()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Type didnot fit in saveCollection!")
	}

	return nil
}

func (s *BaseStore) DeleteCollection(db any, col *model.Collection) error {
	err := s.DB().Model(col).Delete()
	if err != nil {
		return err
	}

	_, err = s.DB().DropTable(col.GetName()).Execute()
	return err
}
