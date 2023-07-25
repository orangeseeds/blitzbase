package core

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	DB       *sql.DB
}

func connectDB(driver, path string) *sql.DB {
	db, err := sql.Open(driver, path)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

func NewStorage(publisher *Publisher) *Storage {
	path := "./test.db"
	log.Println("connecting to ", path)
	driver := initDriverWithUpdateHook(publisher)
	return &Storage{
		DB:            connectDB(driver, path),
	}
}
