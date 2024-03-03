package store

import (
	"fmt"
	"log"

	dbx "github.com/go-ozzo/ozzo-dbx"
	model "github.com/orangeseeds/blitzbase/refactored/models"
)

func (s *BaseStore) FindRecordsAll(db any, colName string, filters ...FilterFunc) ([]*model.Record, error) {
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

	q := selectQuery.From(col.Name)
	for _, filter := range filters {
		if filter == nil {
			continue
		}
		err := filter(q)
		if err != nil {
			return nil, err
		}
	}
	records := []*model.Record{}
	resp := []dbx.NullStringMap{}
	err = q.All(&resp)
	if err != nil {
		return nil, err
	}

	for i := range resp {
		rec := model.NewRecord(col)
		rec.LoadNullStringMap(resp[i])
		records = append(records, rec)
	}
	return records, nil

}

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
		log.Println("collection not found")
		return nil, err
	}
	rec := model.NewRecord(col)

	q := selectQuery.AndWhere(dbx.HashExp{
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
	err = q.One(&resp)
	if err != nil {
		log.Println("record not found")
		return nil, err
	}
	rec.LoadNullStringMap(resp)
	return rec, nil
}

func (s *BaseStore) AuthRecordEmailIsUnique(db any, collection string, email string) bool {
	var selectQuery *dbx.SelectQuery
	switch db.(type) {
	case *dbx.Tx:
		selectQuery = db.(*dbx.Tx).Select("count(*)")
	case *dbx.DB:
		selectQuery = db.(*dbx.DB).Select("count(*)")
	default:
		return false
	}
	if email == "" {
		return false
	}
	query := selectQuery.From(collection).
		Where(dbx.HashExp{model.FieldEmail: email}).
		Limit(1)

	var exists bool

	return query.Row(&exists) == nil && !exists
}

// Make sure to get the collection from the table before sending the record
func (s *BaseStore) SaveRecord(db any, r *model.Record, filters ...FilterFunc) error {
	// do a dry run to save the record
	// check whether you can view the record or not
	// if you can view the record save, else rollback the transaction

	if r.Collection().IsAuth() {
		for _, v := range model.AuthRecordFields() {
			if r.GetString(v) == "" {
				return fmt.Errorf("auth record needs the field %s", v)
			}
		}
		if !s.AuthRecordEmailIsUnique(db, r.TableName(), r.GetString(model.FieldEmail)) {
			return fmt.Errorf("record in %s collection already exists with email %s", r.TableName(), r.GetString(model.FieldEmail))
		}
	}

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

func (s *BaseStore) DeleteRecord(db any, r *model.Record) error {
	switch db.(type) {
	case *dbx.Tx:
		return db.(*dbx.Tx).Model(r).Delete()
	case *dbx.DB:
		return db.(*dbx.DB).Model(r).Delete()
	default:
		return fmt.Errorf("Type didnot fit in FindCollection!")
	}
}

func (s *BaseStore) FindAuthRecordByEmail(db any, collectionName string, email string) (*model.Record, error) {
	coll, err := s.FindCollectionByNameorId(db, collectionName)
	if err != nil {
		return nil, err
	}

	if !coll.IsAuth() {
		return nil, fmt.Errorf("Collection %s is not a auth collection.", collectionName)
	}

	var selectQuery *dbx.SelectQuery
	switch db.(type) {
	case *dbx.Tx:
		selectQuery = db.(*dbx.Tx).Select().From(coll.Name)
	case *dbx.DB:
		selectQuery = db.(*dbx.DB).Select().From(coll.Name)
	default:
		return nil, fmt.Errorf("Type didnot fit in FindCollection!")
	}

	record := model.NewRecord(coll)
	q := selectQuery.AndWhere(dbx.HashExp{
		model.FieldEmail: email,
	}).Limit(1)
	if err != nil {
		return nil, err
	}

	resp := dbx.NullStringMap{}
	err = q.One(&resp)
	if err != nil {
		log.Println("record not found")
		return nil, err
	}
	record.LoadNullStringMap(resp)

	return record, nil
}

func (s *BaseStore) FindAuthRecordByToken(db any, collectionName string, token string) (*model.Record, error) {
	coll, err := s.FindCollectionByNameorId(db, collectionName)
	if err != nil {
		return nil, err
	}

	if !coll.IsAuth() {
		return nil, fmt.Errorf("Collection %s is not a auth collection.", collectionName)
	}

	var selectQuery *dbx.SelectQuery
	switch db.(type) {
	case *dbx.Tx:
		selectQuery = db.(*dbx.Tx).Select().From(coll.Name)
	case *dbx.DB:
		selectQuery = db.(*dbx.DB).Select().From(coll.Name)
	default:
		return nil, fmt.Errorf("Type didnot fit in FindCollection!")
	}

	record := model.NewRecord(coll)
	q := selectQuery.AndWhere(dbx.HashExp{
		model.FieldToken: token,
	}).Limit(1)
	if err != nil {
		return nil, err
	}

	resp := dbx.NullStringMap{}
	err = q.One(&resp)
	if err != nil {
		log.Println("record not found")
		return nil, err
	}
	record.LoadNullStringMap(resp)

	return record, nil
}

// Make sure to get the collection from the table before sending the record
func (s *BaseStore) UpdateRecord(db any, collection string, r *model.Record, filters ...FilterFunc) error {
	// do a dry run to save the record
	// check whether you can view the record or not
	// if you can view the record save, else rollback the transaction

	col, err := s.FindCollectionByNameorId(db, collection)
	if err != nil {
		return err
	}
	params := dbx.Params{}
	for _, f := range col.Schema.GetFields() {
		params[f.Name] = r.Get(f.Name)
	}

	switch db.(type) {
	case *dbx.Tx:
		_, err := db.(*dbx.Tx).Update(col.Name, params, dbx.HashExp{
			model.FieldId: r.GetID(),
		}).Execute()
		return err
	case *dbx.DB:
		res, err := db.(*dbx.DB).Update(col.Name, params, dbx.HashExp{
			model.FieldId: r.GetID(),
		}).Execute()
		log.Println(res, err)
		return err
	default:
		return fmt.Errorf("Type didnot fit in FindCollection!")
	}
}
