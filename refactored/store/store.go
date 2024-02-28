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

	FindCollectionByNameorId(any, string) (*model.Collection, error)
	SaveCollection(any, *model.Collection) error
	DeleteCollection(any, *model.Collection) error

	FindAdminById(any, string) (*model.Admin, error)
	FindAdminByEmail(any, string) (*model.Admin, error)
	FindAdminByToken(any, string) (*model.Admin, error)
	CheckAdminEmailIsUnique(any, string) bool
	SaveAdmin(any, *model.Admin) error
	UpdateAdmin(any, *model.Admin) error
	DeleteAdmin(any, *model.Admin) error

	FindRecordsAll(any, string,...FilterFunc) ([]*model.Record, error)
	FindRecordById(any, string, string, ...FilterFunc) (*model.Record, error)
	// FindRecordsByExpr(...dbx.Expression) ([]*model.Record, error)
	SaveRecord(any, *model.Record, ...FilterFunc) error
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
