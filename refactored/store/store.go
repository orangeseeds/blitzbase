package store

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
	model "github.com/orangeseeds/blitzbase/refactored/models"
)

type Store interface {
	DB() *dbx.DB
	Path() string
	MigrationsPath() string

	CreateAdminTable() error
	CreateCollectionMetaTable() error
	CreateCollectionTable(*model.Collection) error
	// CreateMigrationsTable() error

	FindCollectionByNameorId(string) (*model.Collection, error)
	SaveCollection(*model.Collection) error
	DeleteCollection(*model.Collection) error

	FindAdminByEmail(string) (*model.Admin, error)
	FindAdminByToken(string) (*model.Admin, error)
	CheckAdminEmailIsUnique(string) bool
	SaveAdmin(*model.Admin) error
	DeleteAdmin(*model.Admin) error

	FindRecordById(string, string, ...FilterFunc) (*model.Record, error)
	// FindRecordsByExpr(...dbx.Expression) ([]*model.Record, error)
	SaveRecord(*model.Record, ...FilterFunc) error
	DeleteRecord(*model.Record) error

	// ExpandRecord()
}

type BaseStore struct {
	db             *dbx.DB
	tx             *dbx.Tx
	isTransaction  bool
	path           string
	migrationsPath string
}

func NewBaseStore(db *dbx.DB) *BaseStore {
	return &BaseStore{
		db:             db,
		path:           "",
		migrationsPath: "",
	}
}

func (s *BaseStore) Tx() *dbx.Tx {
	return s.tx
}
func (s *BaseStore) DB() *dbx.DB {
	return s.db
}

func (s *BaseStore) Path() string {
	return s.path
}

func (s *BaseStore) MigrationsPath() string {
	return s.migrationsPath
}
