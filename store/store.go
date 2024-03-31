package store

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
	model "github.com/orangeseeds/blitzbase/models"
)

type Store interface {
	DB() *dbx.DB
	Path() string
	MigrationsPath() string

	TableExists(db DBExector, tableName string) bool
	CreateCollectionTable(db DBExector, c *model.Collection) error
	CreateAdminTable(db DBExector) error
	CreateCollectionMetaTable(db DBExector) error

	// CreateMigrationsTable() error

	FindCollectionByNameorId(db DBExector, query string) (*model.Collection, error)
	SaveCollection(db DBExector, col *model.Collection) error
	DeleteCollection(db DBExector, col *model.Collection) error

	FindAdminById(db DBExector, id string) (*model.Admin, error)
	FindAdminByEmail(db DBExector, email string) (*model.Admin, error)
	FindAdminByToken(db DBExector, token string) (*model.Admin, error)
	CheckAdminEmailIsUnique(db DBExector, email string) bool
	SaveAdmin(db DBExector, a *model.Admin) error
	UpdateAdmin(db DBExector, a *model.Admin) error
	DeleteAdmin(db DBExector, a *model.Admin) error

	FindRecordsAll(db DBExector, collectionId string, filters ...FilterFunc) ([]*model.Record, error)
	FindRecordById(db DBExector, id string, collectionId string, filters ...FilterFunc) (*model.Record, error)
	SaveRecord(db DBExector, r *model.Record, filters ...FilterFunc) error
	DeleteRecord(db DBExector, r *model.Record) error

	FindAuthRecordByEmail(db DBExector, collectionId string, email string) (*model.Record, error)
	FindAuthRecordByToken(db DBExector, collectionId string, token string) (*model.Record, error)
	AuthRecordEmailIsUnique(db DBExector, collection string, email string) bool
	// FindRecordsByExpr(...dbx.Expression) ([]*model.Record, error)
	UpdateRecord(db DBExector, collection string, r *model.Record, filters ...FilterFunc) error

	// ExpandRecord()
}

type BaseStore struct {
	db             *dbx.DB
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

func (s *BaseStore) DB() *dbx.DB {
	return s.db
}

func (s *BaseStore) Path() string {
	return s.path
}

func (s *BaseStore) MigrationsPath() string {
	return s.migrationsPath
}
