package store

import (
	"fmt"
	"log"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	model "github.com/orangeseeds/blitzbase/models"
)

type SQliteStore struct {
	db             *dbx.DB
	path           string
	migrationsPath string
}

func NewSQliteStore(db *dbx.DB) *SQliteStore {
	return &SQliteStore{
		db:             db,
		path:           "",
		migrationsPath: "",
	}
}
func (s *SQliteStore) DB() *dbx.DB            { return s.db }
func (s *SQliteStore) Path() string           { return s.path }
func (s *SQliteStore) MigrationsPath() string { return s.migrationsPath }

func (s *SQliteStore) TableExists(db DBExector, tableName string) bool {
	var exists bool
	err := db.Select("count(*)").
		From("sqlite_schema").
		AndWhere(dbx.HashExp{"type": []any{"table", "view"}}).
		AndWhere(dbx.NewExp("LOWER([[name]])=LOWER({:tableName})", dbx.Params{"tableName": tableName})).
		Limit(1).
		Row(&exists)

	return exists || (err != nil)
}

func (s *SQliteStore) createTable(db DBExector, name string, defn map[string]string) error {
	query := db.CreateTable(name, defn)
	if _, err := query.Execute(); err != nil {
		return err
	}

	// Possible place to add migrations
	return nil
}

// Make the default fields fit in better
func (s *SQliteStore) CreateCollectionTable(db DBExector, c *model.Collection) error {
	if s.TableExists(db, c.GetName()) {
		return fmt.Errorf("Table with name %s already exists.", c.GetName())
	}

	return s.createTable(db, c.GetName(), c.DataDefn())
}

// make sure if you are running this check if previous admin tables are present
func (s *SQliteStore) CreateAdminTable(db DBExector) error {
	var a *model.Admin

	if s.TableExists(db, a.TableName()) {
		return fmt.Errorf("Table with name %s already exists.", a.TableName())
	}
	return s.createTable(db, a.TableName(), a.ExportForTable())
}

func (s *SQliteStore) CreateCollectionMetaTable(db DBExector) error {
	var col model.Collection

	if s.TableExists(db, col.TableName()) {
		return fmt.Errorf("Table with name %s already exists.", col.GetName())
	}

	err := s.createTable(db, col.TableName(), col.MetaDataDefn())
	if err != nil {
		return err
	}
	return nil
}

func (s *SQliteStore) FindCollectionByNameorId(db DBExector, query string) (*model.Collection, error) {
	var col model.Collection
	err := db.Select().From(
		col.TableName(),
	).Where(
		dbx.HashExp{model.FieldName: query},
	).OrWhere(
		dbx.HashExp{model.FieldId: query},
	).One(&col)
	return &col, err
}

func (s *SQliteStore) SaveCollection(db DBExector, col *model.Collection) error {
	json, err := col.Schema.MarshalJSON()
	if err != nil {
		return err
	}
	col.SetName(col.Name)
	params := dbx.Params{
		"Id":     col.GetID(),
		"Name":   col.GetName(),
		"Type":   col.Type,
		"Schema": string(json),
		// "Rule":       col.Rule,
		"IndexRule":  col.IndexRule,
		"CreateRule": col.CreateRule,
		"DetailRule": col.DetailRule,
		"UpdateRule": col.UpdateRule,
		"DeleteRule": col.DeleteRule,

		"CreatedAt": time.Now().String(),
		"UpdatedAt": time.Now().String(),
	}
	_, err = db.Insert(col.TableName(), params).Execute()
	return err
}

func (s *SQliteStore) DeleteCollection(db DBExector, col *model.Collection) error {
	err := db.Model(col).Delete()
	if err != nil {
		return err
	}
	_, err = db.DropTable(col.GetName()).Execute()
	return err
}

func (s *SQliteStore) FindAdminById(db DBExector, id string) (*model.Admin, error) {
	var admin model.Admin
	err := db.Select().
		From(admin.TableName()).
		Where(dbx.HashExp{
			"Id": id,
		}).One(&admin)
	return &admin, err
}

func (s *SQliteStore) FindAdminByEmail(db DBExector, email string) (*model.Admin, error) {
	var admin model.Admin
	err := db.(*dbx.DB).Select().
		From(admin.TableName()).
		Where(dbx.HashExp{
			"Email": email,
		}).One(&admin)
	return &admin, err
}

func (s *SQliteStore) FindAdminByToken(db DBExector, token string) (*model.Admin, error) {
	var admin model.Admin
	err := db.Select().
		From(admin.TableName()).
		Where(dbx.HashExp{
			"Token": token,
		}).One(&admin)
	return &admin, err
}

func (s *SQliteStore) CheckAdminEmailIsUnique(db DBExector, email string) bool {
	if email == "" {
		return false
	}
	exists := false
	var admin model.Admin

	err := db.Select("count(*)").
		From(admin.TableName()).
		Where(dbx.HashExp{"email": email}).
		Limit(1).Row(&exists)
	if err != nil {
		return false
	}
	return !exists
}

func (s *SQliteStore) SaveAdmin(db DBExector, a *model.Admin) error {
	if !s.CheckAdminEmailIsUnique(db, a.Email) {
		return fmt.Errorf("admin email %s not unique", a.Email)
	}
	return db.(*dbx.DB).Model(a).Insert()
}

func (s *SQliteStore) UpdateAdmin(db DBExector, a *model.Admin) error {
	return db.Model(a).Update()
}

func (s *SQliteStore) DeleteAdmin(db DBExector, a *model.Admin) error {
	return db.Model(a).Delete()
}

func (s *SQliteStore) FindRecordsAll(db DBExector, collectionId string, filters ...FilterFunc) ([]*model.Record, error) {
	col, err := s.FindCollectionByNameorId(db, collectionId)
	if err != nil {
		return nil, err
	}

	query := db.Select().From(col.GetName())
	for _, filter := range filters {
		if filter == nil {
			continue
		}
		err := filter(query)
		if err != nil {
			return nil, err
		}
	}

	records := []*model.Record{}
	resp := []dbx.NullStringMap{}
	err = query.All(&resp)
	if err != nil {
		return nil, err
	}

	for i := range resp {
		record := model.NewRecord(col).LoadNullStringMap(resp[i])
		records = append(records, record)
	}
	return records, nil

}

func (s *SQliteStore) FindRecordById(db DBExector, id string, collectionId string, filters ...FilterFunc) (*model.Record, error) {
	col, err := s.FindCollectionByNameorId(db, collectionId)
	if err != nil {
		return nil, err
	}

	query := db.Select().From(col.GetName()).
		AndWhere(dbx.HashExp{
			"Id": id,
		})
	for _, filter := range filters {
		if filter == nil {
			continue
		}

		err := filter(query)
		if err != nil {
			return nil, err
		}
	}

	record := model.NewRecord(col)
	resp := dbx.NullStringMap{}
	err = query.One(&resp)
	if err != nil {
		return nil, err
	}
	record.LoadNullStringMap(resp)
	return record, nil
}

func (s *SQliteStore) AuthRecordEmailIsUnique(db DBExector, collection string, email string) bool {
	if email == "" {
		return false
	}
	exists := false
	err := db.Select("count(*)").
		From(collection).
		Where(dbx.HashExp{model.FieldEmail: email}).
		Limit(1).Row(&exists)

	return !exists || err != nil
}

// Make sure to get the collection from the table before sending the record
func (s *SQliteStore) SaveRecord(db DBExector, r *model.Record, filters ...FilterFunc) error {
	// do a dry run to save the record
	// check whether you can view the record or not
	// if you can view the record save, else rollback the transaction

	if r.Collection().IsAuth() {
		for _, v := range model.AuthRecordFields() {
			if r.GetString(v) == "" {
				return fmt.Errorf("auth record needs the field %s", v)
			}
		}
		if s.AuthRecordEmailIsUnique(db, r.TableName(), r.GetString(model.FieldEmail)) {
			return fmt.Errorf("record in %s collection already exists with email %s", r.TableName(), r.GetString(model.FieldEmail))
		}

	}

	err := s.DB().Transactional(func(tx *dbx.Tx) error {
		r.CreatedAt = time.Now().String()
		r.UpdatedAt = time.Now().String()
		_, err := tx.Insert(r.TableName(), r.Export()).Execute()
		if err != nil {
			return err
		}
		record, err := s.FindRecordById(&TxWrapper{tx}, r.Id, r.Collection().Name, filters...)
		if err != nil {
			return err
		}

		if record == nil {
			return fmt.Errorf("record save failed %v.", r)
		}
		return nil
	})
	return err
}

func (s *SQliteStore) DeleteRecord(db DBExector, r *model.Record) error {
	return db.Model(r).Delete()
}

func (s *SQliteStore) FindAuthRecordByEmail(db DBExector, collectionId string, email string) (*model.Record, error) {
	coll, err := s.FindCollectionByNameorId(db, collectionId)
	if err != nil {
		return nil, err
	}

	if !coll.IsAuth() {
		return nil, fmt.Errorf("Collection %s.%s is not a auth collection.", coll.GetName(), coll.GetID())
	}

	resp := dbx.NullStringMap{}

	err = db.Select().From(coll.GetName()).AndWhere(dbx.HashExp{
		model.FieldEmail: email,
	}).Limit(1).One(&resp)
	if err != nil {
		log.Println("record not found")
		return nil, err
	}

	record := model.NewRecord(coll).LoadNullStringMap(resp)
	return record, nil
}

func (s *SQliteStore) FindAuthRecordByToken(db DBExector, collectionId string, token string) (*model.Record, error) {
	coll, err := s.FindCollectionByNameorId(db, collectionId)
	if err != nil {
		return nil, err
	}

	if !coll.IsAuth() {
		return nil, fmt.Errorf("Collection %s.%s is not a auth collection.", coll.GetName(), coll.GetID())
	}

	resp := dbx.NullStringMap{}

	err = db.Select().From(coll.Name).
		AndWhere(dbx.HashExp{
			model.FieldToken: token,
		}).
		Limit(1).
		One(&resp)
	if err != nil {
		return nil, err
	}

	record := model.NewRecord(coll).LoadNullStringMap(resp)
	return record, nil
}

// Make sure to get the collection from the table before sending the record
func (s *SQliteStore) UpdateRecord(db DBExector, collection string, r *model.Record, filters ...FilterFunc) error {
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

	_, err = db.Update(col.Name, params, dbx.HashExp{
		model.FieldId: r.GetID(),
	}).Execute()
	return err
}
