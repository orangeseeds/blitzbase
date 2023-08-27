package store

import (
	// "database/sql"
	"log"

	"github.com/go-ozzo/ozzo-dbx"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	DB             *dbx.DB
	Driver         string
	Path           string
	Publisher      *Publisher
	MigrationsPath string
	connected      bool
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

// [Important!] Make sure to run s.Connect() after creating a new storage
func NewStorage(dbPath string, migrationsPath string) *Storage {
	// log.Println("connecting to ", dbPath)
	publisher := NewPublisher()
	return &Storage{
		Path:           dbPath,
		Publisher:      publisher,
		MigrationsPath: migrationsPath,
		// DB: connectDB(driver, dbPath),
	}
}

func (s *Storage) Connect() {
	driver := initDriverWithUpdateHook(s.Publisher)
	s.Driver = driver
	s.DB = connectDB(s.Driver, s.Path)
	s.connected = true
}
