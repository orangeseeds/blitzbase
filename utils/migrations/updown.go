package migrations

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/orangeseeds/blitzbase/store"
)

// migrations name need to be like 1_install.up.sql
func MigrateUp(s *store.Storage) error {

	driver, err := sqlite3.WithInstance(s.DB.DB(), &sqlite3.Config{})
	if err != nil {
		log.Println("db instance err: ", err)
		return err
	}

	fs, err := (&file.File{}).Open(s.MigrationsPath)
	if err != nil {
		log.Println("opening file err: ", err)
		return err
	}

	mg, err := migrate.NewWithInstance("file", fs, "myDB", driver)
	if err != nil {
		log.Println("migrate err: ", err)
		return err
	}

	if err = mg.Up(); err != nil {
		log.Println("migration up err: ", err)
		return err
	}

	return nil
}

func MigrateDown(s *store.Storage) error {

	driver, err := sqlite3.WithInstance(s.DB.DB(), &sqlite3.Config{})
	if err != nil {
		log.Println("db instance err: ", err)
		return err
	}

	fs, err := (&file.File{}).Open(s.MigrationsPath)
	if err != nil {
		log.Println("opening file err: ", err)
		return err
	}

	mg, err := migrate.NewWithInstance("file", fs, "myDB", driver)
	if err != nil {
		log.Println("migrate err: ", err)
		return err
	}

	if err = mg.Down(); err != nil {
		log.Println("migration up err: ", err)
		return err
	}

	return nil
}
