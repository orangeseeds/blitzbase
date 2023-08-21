package core

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
)

func (a *App) CreateNewAdmin(email string, password string) (int64, error) {
	res, err := a.Store.DB.Insert("_admins", dbx.Params{
		"email":    email,
		"password": password,
	}).Execute()
	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()
	return id, nil
}

type adminRecord struct {
	ID       int64
	Email    string
	Password string
}

func (a *App) ListAdmins() ([]adminRecord, error) {
	var admins []adminRecord

	q := a.Store.DB.Select("email").From("_admins")

	err := q.All(&admins)
	if err != nil {
		return nil, err
	}

	return admins, nil
}
