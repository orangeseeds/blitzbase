package store

import (
	// "database/sql"
	"log"

	"github.com/go-ozzo/ozzo-dbx"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	DB        *dbx.DB
	Driver    string
	Path      string
	Publisher *Publisher
}

func connectDB(driver, path string) *dbx.DB {
	db, err := dbx.Open(driver, path)
	// db, err := sql.Open(driver, path)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.DB().Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

func NewStorage(dbPath string) *Storage {
	log.Println("connecting to ", dbPath)
	publisher := NewPublisher()
	return &Storage{
		Path:      dbPath,
		Publisher: publisher,
		// DB: connectDB(driver, dbPath),
	}
}

func (s *Storage) Connect() {
	driver := initDriverWithUpdateHook(s.Publisher)
	s.Driver = driver
	s.DB = connectDB(s.Driver, s.Path)
}
