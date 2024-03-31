package store

import (
	"fmt"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	model "github.com/orangeseeds/blitzbase/models"
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

	err := selectQuery.From(
		col.TableName(),
	).Where(
		dbx.HashExp{"Name": query},
	).OrWhere(
		dbx.HashExp{"Id": query},
	).One(&col)
	if err != nil {
		return nil, err
	}

	return &col, nil
}

func (s *BaseStore) SaveCollection(db any, col *model.Collection) error {
	err := db.(*dbx.DB).Transactional(func(tx *dbx.Tx) error {
		json, err := col.Schema.MarshalJSON()
		if err != nil {
			return err
		}
		col.SetName(col.Name)

		params := dbx.Params{
			model.FieldId:         col.GetID(),
			model.FieldName:       col.GetName(),
			model.FieldType:       col.Type,
			model.FieldSchema:     string(json),
			model.FieldIndexRule:  col.IndexRule,
			model.FieldCreateRule: col.CreateRule,
			model.FieldDetailRule: col.DetailRule,
			model.FieldUpdateRule: col.UpdateRule,
			model.FieldDeleteRule: col.DeleteRule,
			model.FieldCreatedAt:  time.Now().String(),
			model.FieldUpdatedAt:  time.Now().String(),
		}

		switch db.(type) {
		case *dbx.Tx:
			_, err = db.(*dbx.Tx).Insert(col.TableName(), params).Execute()
			if err != nil {
				return err
			}
		case *dbx.DB:
			_, err = db.(*dbx.DB).Insert(col.TableName(), params).Execute()
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Type didnot fit in saveCollection!")
		}
		return nil
	})
	return err
}

func (s *BaseStore) DeleteCollection(db any, col *model.Collection) error {
	err := s.DB().Model(col).Delete()
	if err != nil {
		return err
	}

	_, err = s.DB().DropTable(col.GetName()).Execute()
	return err
}
