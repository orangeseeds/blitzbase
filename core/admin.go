package core

import (
	"log"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

func (a *App) CreateNewAdmin(email string, password string) error {
	// query := fmt.Sprintf("Insert into _admins (email, password) values ('%s', '%s')", email, password)
	res, err := a.Store.DB.Insert("_admins", dbx.Params{
		"email":    email,
		"password": password,
	}).Execute()
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	log.Println("new admin id: ", id)
	return nil
}

func (a *App) ListAdmins() error {
	// query := "select email from _admins"
	q := a.Store.DB.Select("email").From("_admins")
	var admins []string

	err := q.All(&admins)
	if err != nil {
		return err
	}
	return nil
}
