package core

import (
	"fmt"
	"log"
)

func (a *App) CreateNewAdmin(email string, password string) error {
	query := fmt.Sprintf("Insert into _admins (email, password) values ('%s', '%s')", email, password)
	res, err := a.Store.DB.Exec(query)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	log.Println("new admin id: ", id)
	return nil
}

func (a *App) ListAdmins() error {
	query := "select email from _admins"
	rows, err := a.Store.DB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			return err
		}
		log.Println(email)
	}

	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}
