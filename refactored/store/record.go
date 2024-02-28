package store

import (
	"fmt"

	dbx "github.com/go-ozzo/ozzo-dbx"
	model "github.com/orangeseeds/blitzbase/refactored/models"
)

func (s *BaseStore) FindRecordById(db any, id string, colName string, filters ...FilterFunc) (*model.Record, error) {
	var selectQuery *dbx.SelectQuery
	switch db.(type) {
	case *dbx.Tx:
		selectQuery = db.(*dbx.Tx).Select().From(colName)
	case *dbx.DB:
		selectQuery = db.(*dbx.DB).Select().From(colName)
	default:
		return nil, fmt.Errorf("Type didnot fit in FindCollection!")
	}

	col, err := s.FindCollectionByNameorId(db, colName)
	if err != nil {
		return nil, err
	}
	rec := model.NewRecord(col)

	q := selectQuery.From(col.Name).AndWhere(dbx.HashExp{
		"Id": id,
	})
	for _, filter := range filters {
		if filter == nil {
			continue
		}
		err := filter(q)
		if err != nil {
			return nil, err
		}
	}

	resp := dbx.NullStringMap{}
	err = q.Limit(1).One(resp)
	if err != nil {
		return nil, err
	}
	rec.LoadNullStringMap(resp)
	return rec, nil
}

// Make sure to get the collection from the table before sending the record
func (s *BaseStore) SaveRecord(db any, r *model.Record, filters ...FilterFunc) error {
	// do a dry run to save the record
	// check whether you can view the record or not
	// if you can view the record save, else rollback the transaction
	err := s.DB().Transactional(func(tx *dbx.Tx) error {
		_, err := tx.Insert(r.TableName(), r.Export()).Execute()
		if err != nil {
			return err
		}
		rec, err := s.FindRecordById(tx, r.Id, r.Collection().Name, filters...)
		if err != nil {
			return err
		}
		if rec == nil {
			return fmt.Errorf("record save failed %v.", r)
		}
		return nil
	})
	return err
}

func (s *BaseStore) DeleteRecord(r *model.Record) error {
	return s.DB().Model(r).Delete()
}
