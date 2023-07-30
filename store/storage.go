package store

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	DB        *sql.DB
	Driver    string
	Path      string
	Publisher *Publisher
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

func NewStorage(dbPath string) (*Storage) {
	log.Println("connecting to ", dbPath)
	publisher := NewPublisher()
	return &Storage{
		Path:      dbPath,
		Publisher: publisher,
		// DB: connectDB(driver, dbPath),
	}
}

func (s *Storage) Connect(){
	driver := initDriverWithUpdateHook(s.Publisher)
	s.Driver = driver
	s.DB = connectDB(s.Driver, s.Path)
}
