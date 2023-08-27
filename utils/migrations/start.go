package migrations

import (
	"fmt"

	"github.com/orangeseeds/blitzbase/store"
	"github.com/orangeseeds/blitzbase/utils/schema"
)

func CreateInitTable(s *store.Storage) error {
	admins := s.DB.CreateTable("_admins", map[string]string{
		"id":       "INTEGER PRIMARY KEY",
		"email":    "varchar(255) UNIQUE NOT NULL",
		"password": "varchar(255) NOT NULL",
	})

	// if _, err := admins.Execute(); err != nil {
	// 	return err
	// }

	collections := s.DB.CreateTable("_collections", map[string]string{
		"id":     "INTEGER PRIMARY KEY",
		"name":   "varchar(255) NOT NULL",
		"type":   "integer NOT NULL",
		"schema": "json NOT NULL",
	})

	// if _, err := collections.Execute(); err != nil {
	// 	return err
	// }

	sc := []schema.FieldSchema{
		{
			Name:     "username",
			Type:     schema.TEXT_FIELD,
			Options:  schema.TextFieldOptions{},
			Required: true,
		},
		{
			Name:     "email",
			Type:     schema.TEXT_FIELD,
			Options:  schema.TextFieldOptions{},
			Required: true,
		},
		{
			Name:     "password",
			Type:     schema.TEXT_FIELD,
			Options:  schema.TextFieldOptions{},
			Required: true,
		},
	}

	c := *store.NewBaseCollection("users", store.AuthType, true, sc)

	user := s.DB.CreateTable(c.Name, c.TableSchema())

	up := fmt.Sprintf("%s;\n%s;\n%s;", admins.SQL(), collections.SQL(), user.SQL())
	down := "DROP TABLE IF EXISTS _admins;\nDROP TABLE IF EXISTS _collections;\nDROP TABLE IF EXISTS users;"
	if err := CreateNewMigration(up, down, s.MigrationsPath); err != nil {
		return err
	}

	err := MigrateUp(s)
	if err != nil {
		return err
	}
	if err := AddCollectionRecord(s, c); err != nil {
		return err
	}

	return nil
}
