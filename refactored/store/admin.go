package store

import (
	"fmt"

	dbx "github.com/go-ozzo/ozzo-dbx"
	model "github.com/orangeseeds/blitzbase/refactored/models"
)

func (s *BaseStore) FindAdminById(db any, id string) (*model.Admin, error) {
	var admin model.Admin
	var selectQuery *dbx.SelectQuery
	switch db.(type) {
	case *dbx.Tx:
		selectQuery = db.(*dbx.Tx).Select()
	case *dbx.DB:
		selectQuery = db.(*dbx.DB).Select()
	default:
		return nil, fmt.Errorf("Type didnot fit in FindAdminByEmail!")
	}

	err := selectQuery.From(admin.TableName()).Where(dbx.HashExp{
		"Id": id,
	}).One(&admin)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (s *BaseStore) FindAdminByEmail(db any, email string) (*model.Admin, error) {
	var admin model.Admin
	var selectQuery *dbx.SelectQuery
	switch db.(type) {
	case *dbx.Tx:
		selectQuery = db.(*dbx.Tx).Select()
	case *dbx.DB:
		selectQuery = db.(*dbx.DB).Select()
	default:
		return nil, fmt.Errorf("Type didnot fit in FindAdminByEmail!")
	}

	err := selectQuery.From(admin.TableName()).Where(dbx.HashExp{
		"Email": email,
	}).One(&admin)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (s *BaseStore) FindAdminByToken(db any, token string) (*model.Admin, error) {
	var admin model.Admin

	var selectQuery *dbx.SelectQuery
	switch db.(type) {
	case *dbx.Tx:
		selectQuery = db.(*dbx.Tx).Select()
	case *dbx.DB:
		selectQuery = db.(*dbx.DB).Select()
	default:
		return nil, fmt.Errorf("Type didnot fit in FindAdminByEmail!")
	}

	err := selectQuery.From(admin.TableName()).Where(dbx.HashExp{
		"Token": token,
	}).One(&admin)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (s *BaseStore) CheckAdminEmailIsUnique(db any, email string) bool {
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
	var admin model.Admin

	query := selectQuery.From(admin.TableName()).
		Where(dbx.HashExp{"email": email}).
		Limit(1)

	var exists bool

	return query.Row(&exists) == nil && !exists
}

func (s *BaseStore) SaveAdmin(db any, a *model.Admin) error {
	if !s.CheckAdminEmailIsUnique(db, a.Email) {
		return fmt.Errorf("admin email %s not unique", a.Email)
	}

	switch db.(type) {
	case *dbx.Tx:
		return db.(*dbx.Tx).Model(a).Insert()
	case *dbx.DB:
		return db.(*dbx.DB).Model(a).Insert()
	default:
		return fmt.Errorf("Type didnot fit in FindCollection!")
	}
}

func (s *BaseStore) UpdateAdmin(db any, a *model.Admin) error {
	switch db.(type) {
	case *dbx.Tx:
		return db.(*dbx.Tx).Model(a).Update()
	case *dbx.DB:
		return db.(*dbx.DB).Model(a).Update()
	default:
		return fmt.Errorf("Type didnot fit in FindCollection!")
	}
}

func (s *BaseStore) DeleteAdmin(db any, a *model.Admin) error {
	switch db.(type) {
	case *dbx.Tx:
		return db.(*dbx.Tx).Model(a).Delete()
	case *dbx.DB:
		return db.(*dbx.DB).Model(a).Delete()
	default:
		return fmt.Errorf("Type didnot fit in FindCollection!")
	}
}
